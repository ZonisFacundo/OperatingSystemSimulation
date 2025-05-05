package utilsKernel

type Estado string

type PCB struct {
	Pid            int              `json:"pid"`
	Pc             int              `json:"pc"`
	EstadoActual   Estado           `json:"estadoActual"`
	TamProceso     int              `json:"tamanioProceso"`
	MetricaEstados map[Estado]int   `json:"metricaEstados"` //falta verlo
	TiempoEstados  map[Estado]int64 `json:"tiempoEstados"`  // falta verlo un abrazo
	Archivo        string           `json:"file"`
}

/*
type PaqueteRecibidoMemoriadeKernel struct {
	Pid        int    `json:"pid"`
	TamProceso int    `json:"tamanioproceso"`
	Archivo    string `json:"file"`
}
*/

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
	Instancia string `json:"instance_id"`
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
	PC  int `json:"pc"`
	Pid int `json:"pid"`
}
type RespuestaalIO struct {
	Mensaje string `json:"message"`
}

type PaqueteEnviadoKERNELaIO struct {
	Pid    int `json:"pid"`
	Tiempo int `json:"tiempo"`
}

type RespuestaalCPU struct {
	Mensaje string `json:"messageCPU"`
}
type PaqueteRecibidoDeMemoria struct {
	Mensaje string `json:"message"`
}

type PaqueteRecibidoDeIO struct {
	Mensaje string `json:"message"` //Tiene que ser igual de ambos lados.
}

type PaqueteRecibido struct {
	Mensaje string `json:"messageCPU"`
}

type PaqueteRecibidoDeCPU struct {
	Mensaje string `json:"messageCPU"`
	Pid     int    `json:"pid"`
	Pc      int    `json:"pc"`
}

var ColaNew []PCB
var ColaReady []PCB
var ContadorPCB int = 0
var ListaCPU []CPU
