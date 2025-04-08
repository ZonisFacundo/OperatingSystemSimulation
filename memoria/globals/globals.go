package globals

type Config struct {
	Port_memory      int    `json:"port_memory"`
	Memory_size      int    `json:"memory_size"`
	Page_size        int    `json:"page_size"`
	Entries_per_page int    `json:"entries_per_page"`
	Number_of_levels int    `json:"number_of_levels"`
	Memory_delay     int    `json:"memory_delay"`
	Swapfile_path    string `json:"swapfile_path"`
	Swap_delay       int    `json:"swap_delay"`
	Log_level        string `json:"log_level"`
	Dump_path        string `json:"dump_path"`
}

var ClientConfig *Config
