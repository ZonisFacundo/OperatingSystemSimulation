package utilsCPU

type Proceso struct {
	Pc  int `json:"pc"`
	Pid int `json:"pid"`
}
type Interrupcion struct {
	TiempoInterrup int  `json:"interrup"`
	InterrupValida bool `json:"interrupValida"`
}
type HandshakeCPU struct {
	Ip        string `json:"ip"`
	Puerto    int    `json:"port"`
	Instancia string `json:"instance"`
}
type HandshakeMemory struct {
	Ip     string `json:"ip"` // es fundamental ponerlo
	Puerto int    `json:"port"`
	Pid    int    `json:"pid"`
	Pc     int    `json:"pc"`
}

type WriteStruct struct {
	Datos     string `json:"datos"`
	Direccion int    `json:"adress"`
}
type ReadStruct struct {
	Tamaño    int `json:"datos"`
	Direccion int `json:"adress"`
}
