package globals

type Config struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
}

var ClientConfig *Config
