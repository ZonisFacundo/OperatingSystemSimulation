package utilsCPU

type Proceso struct {
	Pc  int `json:"pc"`
	Pid int `json:"pid"`
	//Mensaje string `json:"messageCPU"`
}
type Interrupcion struct {
	TiempoInterrup int  `json:"interrup"`
	InterrupValida bool `json:"interrupValida"`
}
type HandshakeCPU struct {
	Ip        string `json:"ip"`
	Puerto    int    `json:"port"`
	Instancia string `json:"instance_id"`
}
type HandshakeMemory struct {
	Ip     string `json:"ip"` // es fundamental ponerlo
	Puerto int    `json:"port"`
	Pid    int    `json:"pid"`
	Pc     int    `json:"pc"`
}

type PackageFinEjecucion struct {
	Pid       int    `json:"pid"`
	Pc        int    `json:"pc"`
	Contexto  string `json:"context"`
	Instancia string `json:"instance_id"`
}

type WriteStruct struct {
	Datos     string `json:"datos"`
	Direccion int    `json:"adress"`
}
type ReadStruct struct {
	Tamaño    int `json:"datos"`
	Direccion int `json:"adress"`
}

type HandshakeKERNEL struct {
	Ip        string `json:"ip"`
	Puerto    int    `json:"port"`
	Instancia string `json:"instancia"`
}
type RespuestaalCPU struct {
	Mensaje string `json:"messageCPU"`
}
type RespuestaKernel struct {
	Mensaje string `json:"messageCPU"`
}
