package main

//go run github.com/pressly/goose/v3/cmd/goose postgres postgres://postgres:postgres@postgresql-fencelive:5434/fencelive?sslmode=disable up

//go:generate go run github.com/99designs/gqlgen generate
//go:generate go run github.com/go-jet/jet/v2/cmd/jet -dsn=postgres://postgres:postgres@localhost:5434/fencelive?sslmode=disable -path=../internal/ports/database/gen

import (
	"aphrodite/internal/config"
	"aphrodite/internal/handlers"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}

func run() error {
	log.Println("Reding configuration...")
	configuration := config.LoadConfig()
	s, _ := json.MarshalIndent(configuration, "", "\t")
	log.Println(string(s))
	router := mux.NewRouter()

	log.Println("Calling serve")
	return serve(router, configuration)
}

func serve(mux *mux.Router, config *config.Config) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	WSHandler := handlers.NewWsHandler(config.WebSocketConfig)

	mux.Handle("/ping", WSHandler.Ping())
	mux.Handle("/echo", WSHandler.Echo())
	mux.Handle("/conn/{id}", WSHandler.HandleConnections())
	mux.Handle("/ws", WSHandler.GetConnections())
	mux.Handle("/Dio", WSHandler.DIOEndpoint())
	mux.Handle("/Yev", WSHandler.YEVEndpoint())

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := corsMiddleware.Handler(mux)
	api := http.Server{
		Addr:         "0.0.0.0:" + config.ServerConfig.Port,
		ReadTimeout:  config.ServerConfig.ReadTimeout,
		WriteTimeout: config.ServerConfig.WriteTimeout,
		Handler:      handler,
	}
	serverErrors := make(chan error, 1)
	go func() {
		sugar.Infof("Connect to http://localhost:%s/", config.ServerConfig.Port)
		if config.ServerConfig.TLSEnable {
			serverErrors <- api.ListenAndServeTLS(config.ServerConfig.TLSCertPath, config.ServerConfig.TLSKeyPath)
		} else {
			serverErrors <- api.ListenAndServe()
		}
	}()

	select {
	case err := <-serverErrors:
		return err

	case sig := <-shutdown:
		ctx, cancel := context.WithTimeout(context.Background(), config.ServerConfig.ShutdownTimeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			err = api.Close()
		}

		switch {
		case sig == syscall.SIGKILL:
			return errors.New("integrity error shuting down")

		case err != nil:
			return err
		}
		return nil
	}
}
