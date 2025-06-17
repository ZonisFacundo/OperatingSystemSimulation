package utilsKernel

import (
	"time"
)

type Estado string

type PCB struct {
	Pid                int                  `json:"pid"`
	Pc                 int                  `json:"pc"`
	EstadoActual       Estado               `json:"estadoActual"`
	TamProceso         int                  `json:"tamanioProceso"`
	MetricaEstados     map[Estado]int       `json:"metricaEstados"`
	TiempoLlegada      map[Estado]time.Time `json:"tiempoLLegada"`
	TiempoEstados      map[Estado]int64     `json:"tiempoEstados"`
	Archivo            string               `json:"file"`
	RafagaAnterior     float32              `json:"rafagaAnterior"` //capaz dsp lo cambiamos a time xd
	EstimacionAnterior float32              `json:"estimacionAnterior"`
	TiempoEnvioExc     time.Time            `json:"tiempoEnvioExc"` //sirve para calcular el timpo de ejecucion
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
	Pid        int    `json:"pid"` //Lo agregue porque no sabemos q cpu ejecuta q proceso
}

type IO struct {
	Ip           string  `json:"ip"`
	Port         int     `json:"port"`
	Instancia    string  `json:"instancia"`
	Disponible   bool    `json:"disponible"`
	ColaProcesos []PCBIO `json:"colaprocesos"`
}

type PCBIO struct {
	Pid    int `json:"pid"`
	Tiempo int `json:"tiempo"`
}

type HandshakepaqueteIO struct {
	Nombre string `json:"name"`
	Ip     string `json:"ip"`
	Puerto int    `json:"port"`
}

type HandshakepaqueteFinIO struct {
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
	Pid          int    `json:"pid"`
	Pc           int    `json:"pc"`
	Syscall      string `json:"syscall"`
	Parametro1   int    `json:"parametro1"`
	Parametro2   string `json:"parametro2"`
	InstanciaCPU string `json:"instanciaCPU"`
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

type PaqueteEnviadoKERNELaMemoria2 struct {
	Pid     int    `json:"pid"`
	Mensaje string `json:"message"`
}

type PaqueteEnviadoKERNELaCPU struct {
	PC  int `json:"pc"`
	Pid int `json:"pid"`
}

type PaqueteInterrupcion struct {
	mensaje string `json:"message"`
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

var ColaNew []*PCB
var ColaReady []*PCB
var ListaExec []*PCB
var ColaBlock []*PCB
var ColaSuspBlock []*PCB
var ColaSuspReady []*PCB
var ColaExit []*PCB
var ContadorPCB int = 0
var ListaCPU []CPU
var ListaIO []IO
