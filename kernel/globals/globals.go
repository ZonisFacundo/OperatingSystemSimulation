package globals

type Config struct {
	Ip_kernel   string `json:"ip_kernel"`
	Port_io     int    `json:"port_io"`
	Port_kernel int    `json:"port_kernel"`
	Log_level   string `json:"log_level"`
}

var ClientConfig *Config
