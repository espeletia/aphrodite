package config

type WebSocketConfig struct {
	ReadBufferSize  int 
	WriteBufferSize int
}

func LoadWebSocketConfig() WebSocketConfig {
	webSocketConfig := &WebSocketConfig{}
	v := configViper("ws")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	v.Unmarshal(webSocketConfig)
	return *webSocketConfig
}
