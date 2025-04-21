package utilsKernel

type Estado string

type PCB struct {
	Pid            int              `json:"pid"`
	PC             int              `json:"PC"`
	EstadoActual   Estado           `json:"estadoActual"`
	TamProceso     int              `json:"tamanioProceso"`
	MetricaEstados map[Estado]int   `json:"metricaEstados"`
	TiempoEstados  map[Estado]int64 `json:"tiempoEstados"`
}

type HandshakepaqueteIO struct {
	Nombre string `json:"name"`
	Ip     string `json:"ip"`
	Puerto int    `json:"port"`
}

type HandshakepaqueteCPU struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"port"`
	Instancia string `json:"instancia"`
}

type HandshakepaqueteKERNEL struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"port"`
}

type PaqueteEnviadoKERNELaMemoria struct {
	Pid        int    `json:"pid"`
	TamProceso int    `json:"tamanioProceso"`
	Archivo    string `json:"file"`
}
type RespuestaalIO struct {
	Mensaje string `json:"message"`
}

type RespuestaalCPU struct {
	Mensaje string `json:"message"`
}
type PaqueteRecibidoKERNEL struct {
	Mensaje string `json:"message"`
	Exito   bool   `json:"exito"`
}

var ColaNew []PCB
var ColaReady []PCB
var ContadorPCB int = 0
