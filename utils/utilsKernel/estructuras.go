package utilsKernel

type Estado string

type PCB struct {
	Pid            int              `json:"pid"`
	PC             int              `json:"PC"`
	EstadoActual   Estado           `json:"estadoActual"`
	TamProceso     int              `json:"tamanioProceso"`
	MetricaEstados map[Estado]int   `json:"metricaEstados"` //falta verlo
	TiempoEstados  map[Estado]int64 `json:"tiempoEstados"`  // falta verlo
}

type CPU struct {
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	Instancia  string `json:"instancia"`
	Disponible bool   `json:"disponible"`
}

type HandshakepaqueteIO struct {
	Nombre string `json:"name"`
	Ip     string `json:"ip"`
	Puerto int    `json:"port"`
}

type HandshakepaqueteCPU struct {
	Ip        string `json:"ip"`
	Puerto    int    `json:"port"`
	Instancia string `json:"instancia"`
}

type HandshakepaqueteCPUPCB struct {
	Pid       string `json:"pid"`
	Pc        int    `json:"pc"`
	Contexto  string `json:"contexto"`
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

type PaqueteEnviadoKERNELaCPU struct {
	Pid int `json:"pid"`
	PC  int `json:"pc"`
}
type RespuestaalIO struct {
	Mensaje string `json:"message"`
}

type Respuesta struct {
	Mensaje string `json:"message"`
}
type PaqueteRecibidoDeMemoria struct {
	Mensaje string `json:"message"`
}

type PaqueteRecibidoDeIO struct {
	Mensaje string `json:"message"`
}

type PaqueteRecibido struct {
	Mensaje string `json:"message"`
}

type PaqueteRecibidoDeCPU struct {
	Mensaje string `json:"message"`
	Pid     int    `json:"pid"`
}

var ColaNew []PCB
var ColaReady []PCB
var ContadorPCB int = 0
var ListaCPU []CPU
