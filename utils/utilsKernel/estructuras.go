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

type RespuestaalCPU struct {
	Mensaje string `json:"messageCPU"`
}
type PaqueteRecibidoDeMemoria struct {
	Mensaje string `json:"message"`
	Exito   bool   `json:"exito"`
}

type PaqueteRecibidoDeIO struct {
	Mensaje string `json:"message"` //Tiene que ser igual de ambos lados.
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
