package utilsCPU

type Proceso struct {
	Pc  int `json:"pc"`
	Pid int `json:"pid"`
	//Mensaje string `json:"messageCPU"`
}

type PaqueteInterrupcion struct {
	Mensaje string `json:"message"`
}
type Interrupcion struct {
	Interrup bool `json:"interrup"`
}
type HandshakeCPU struct {
	Ip        string `json:"ip"`
	Puerto    int    `json:"port"`
	Instancia string `json:"instance_id"`
	Port_cpu  int    `json:"port_cpu"`
}
type HandshakeMemory struct {
	Ip     string `json:"ip"` // es fundamental ponerlo
	Puerto int    `json:"port"`
	Pid    int    `json:"pid"`
	Pc     int    `json:"pc"`
}

type PackageFinEjecucion struct {
	Pid          int    `json:"pid"`
	Pc           int    `json:"pc"`
	Syscall      string `json:"syscall"`
	Parametro1   int    `json:"parametro1"`
	Parametro2   string `json:"parametro2"`
	InstanciaCPU string `json:"instanciaCPU"`
}

type WriteStruct struct {
	Direccion int    `json:"adress"`
	Contenido string `json:"content"`
}
type ReadStruct struct {
	Direccion int `json:"adress"`
	Tamaño    int `json:"value"`
}

type HandshakeKERNEL struct {
	Ip        string `json:"ip"`
	Puerto    int    `json:"port"`
	Instancia string `json:"instancia"`
}
type RespuestaalCPU struct {
	Direccion int `json:"dir_logica"`
}
type RespuestaKernel struct {
	Mensaje string `json:"messageCPU"`
	Pid     int    `json:"pid"`
	Pc      int    `json:"pc"`
}
type EnvioDirLogicaAMemoria struct {
	Ip        string `json:"ip"`
	Puerto    int    `json:"port"`
	DirLogica []int  `json:"dir_logica"`
}

type ReciboMMU struct {
	Mensaje string `json:"message"`
	Ip      string `json:"ip"`
	Puerto  int    `json:"port"`
}

type MarcoDeMemoria struct {
	Frame int `json:"frame"`
}

var Pid int
var Pc int

