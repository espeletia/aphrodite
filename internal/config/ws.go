package config

type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	Dio             string
	Yev             string
}

func LoadWebSocketConfig() WebSocketConfig {
	webSocketConfig := &WebSocketConfig{}
	v := configViper("ws")
	v.BindEnv("Dio", "DIO")
	v.BindEnv("Yev", "YEV")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	v.Unmarshal(webSocketConfig)
	return *webSocketConfig
}
