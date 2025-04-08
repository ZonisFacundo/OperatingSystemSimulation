package globals

type Config struct {
	Ip_memory           string `json:"ip_memory"`
	Port_memory         int    `json:"port_memory"`
	Port_kernel         string `json:"port_kernel"`
	Scheduler_algorithm int    `json:"scheduler_algorithm"`
	Suspension_time     int    `json:"suspension_time"`
	Log_level           string `json:"log_level"`
}

var ClientConfig *Config
