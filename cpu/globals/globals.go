package globals

type Config struct {
	Port_cpu          int    `json:"port_cpu"`
	Ip_memory         string `json:"ip_memory"`
	Port_memory       int    `json:"port_memory"`
	Ip_kernel         string `json:"ip_kernel"`
	Port_kernel       int    `json:"port_kernel"`
	Tlb_entries       int    `json:"tlb_entries"`
	Tlb_replacement   string `json:"tlb_replacement"`
	Cache_entries     int    `json:"cache_entries"`
	Cache_replacement string `json:"cache_replacement"`
	Cache_delay       int    `json:"cache_delay"`
	Log_level         string `json:"log_level"`
}

var ClientConfig *Config
