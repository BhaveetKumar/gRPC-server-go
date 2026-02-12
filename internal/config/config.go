package config

type ServerConfig struct {
	Host string
	Port int
}

type ClientConfig struct {
	ServerAddress  string
	TimeoutSeconds int
}

type LogConfig struct {
	EnableRequestID bool
}

type AppConfig struct {
	Environment string
	Server      ServerConfig
	Client      ClientConfig
	Log         LogConfig
}
