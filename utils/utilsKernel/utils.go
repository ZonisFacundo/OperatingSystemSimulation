package utilsKernel

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
)

func ConfigurarLogger() {
	logFile, err := os.OpenFile("kernel.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func RecibirDatosIO(w http.ResponseWriter, r *http.Request) {

	var request HandshakepaqueteIO

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("El cliente nos mando esto: \n nombre: %s  \n puerto: %d \n IP: %s \n", request.Nombre, request.Puerto, request.Ip) //capaz que hay que sacarlo

	CrearStructIO(request.Ip, request.Puerto, request.Nombre)

	//Respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuestaIO RespuestaalIO
	respuestaIO.Mensaje = "conexion realizada con exito"
	respuestaJSON, err := json.Marshal(respuestaIO)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func FinalizarIO(w http.ResponseWriter, r *http.Request) {

	var request HandshakepaqueteFinIO

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("el IO: %s se desconecto", request.Nombre)

	ioCerrada := ObtenerIO(request.Nombre)
	enviarExitProcesosIO(ioCerrada)
	ListaIO = removerIO(&ioCerrada)

	var respuestaIO RespuestaalIO
	respuestaIO.Mensaje = "conexion realizada con exito"
	respuestaJSON, err := json.Marshal(respuestaIO)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func RecibirDatosCPU(w http.ResponseWriter, r *http.Request) {

	var request HandshakepaqueteCPU

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("Handshake recibido: Port: %d - Instance: %s - Ip: %s", request.Puerto, request.Instancia, request.Ip)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuesta RespuestaalCPU
	respuesta.Mensaje = "Conexion realizada con exito"
	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}

	crearStructCPU(request.Ip, request.Puerto, request.Instancia)

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func RecibirProceso(w http.ResponseWriter, r *http.Request) {

	var request HandshakepaqueteCPUPCB

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("contexto de devolucion del proceso: %s", request.Syscall)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuesta RespuestaalCPU

	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}
	log.Printf("Conexion establecida con exito \n")
	cpuServidor := ObtenerCpu(request.InstanciaCPU)
	cpuServidor.Disponible = true
	PCBUtilizar := ObtenerPCB(cpuServidor.Pid) // ya no hace falta porque esta en el struct
	PCBUtilizar.RafagaAnterior = float32(PCBUtilizar.TiempoEnvioExc.Sub(time.Now()))
	respuesta.Mensaje = "interrupcion"
	switch request.Syscall {
	case "I/O":
		//interrumpir
		if ExisteIO(request.Parametro2) {
			PlanificadorCortoPlazo()
			ioServidor := ObtenerIO(request.Parametro2)
			AgregarColaIO(ioServidor, PCBUtilizar.Pid, request.Parametro1)
			PasarBlocked(PCBUtilizar)
			log.Printf("## (<%d>) - Bloqueado por IO: < %s > \n", PCBUtilizar.Pid, ioServidor.Instancia)
			MandarProcesoAIO(ioServidor)
			if len(ioServidor.ColaProcesos) > 0 {
				ioServidor.ColaProcesos = ioServidor.ColaProcesos[1:]
			}
		} else {
			FinalizarProceso(PCBUtilizar)
		} //remplanificar
		log.Printf("## (<%d>) - Solicitó syscall: <IO> \n", PCBUtilizar.Pid)
	case "EXIT":
		FinalizarProceso(PCBUtilizar)
		log.Printf("## (<%d>) - Solicitó syscall: <EXIT> \n", PCBUtilizar.Pid)
	case "DUMP_MEMORY":
		PlanificadorCortoPlazo()
		DumpDelProceso(PCBUtilizar, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory) //revisar
		log.Printf("## (<%d>) - Solicitó syscall: <DUMP_MEMORY> \n", PCBUtilizar.Pid)
	case "INIT_PROC":
		respuesta.Mensaje = ""
		CrearPCB(request.Parametro1, request.Parametro2)
		log.Printf("## (<%d>) - Solicitó syscall: <INIT_PROC> \n", PCBUtilizar.Pid)
		cpuServidor.Disponible = false
		EnviarProcesoACPU(PCBUtilizar, &cpuServidor)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func UtilizarIO(ioServer IO, pid int, tiempo int) {

	var paquete PaqueteEnviadoKERNELaIO
	paquete.Pid = pid
	paquete.Tiempo = tiempo

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al convertir a json.")
		return
	}
	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/KERNELIO", ioServer.Ip, ioServer.Port) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json") //le avisa al server que manda la data en json format

	respuestaJSON, err := cliente.Do(req) //recibe la respuesta del server

	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	if respuestaJSON.StatusCode != http.StatusOK {

		log.Printf("Status de respuesta del I/0 %s no fue la esperada.\n", ioServer.Instancia)
		FinalizarProcesosIO(&ioServer)
		return
	}
	defer respuestaJSON.Body.Close() //cerramos algo supuestamente importante de cerrar pero no se que hace

	log.Printf("Conexion establecida con exito \n")
	//pasamos de JSON a formato bytes lo que nos paso el paquete
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}
	var respuesta PaqueteRecibidoDeIO
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}
	log.Printf("La respuesta del I/O %s fue: %s\n", ioServer.Instancia, respuesta.Mensaje)
	MandarProcesoAIO(ioServer)

}

func ConsultarProcesoConMemoria(pcb *PCB, ip string, puerto int) {

	var paquete PaqueteEnviadoKERNELaMemoria
	paquete.Pid = pcb.Pid
	paquete.TamProceso = pcb.TamProceso
	paquete.Archivo = pcb.Archivo

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al convertir a json.")
		return
	}
	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/KERNELMEMORIA", ip, puerto) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json") //le avisa al server que manda la data en json format

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	defer respuestaJSON.Body.Close() //cerramos algo supuestamente importante de cerrar pero no se que hace

	log.Printf("Conexion establecida con exito \n")
	//pasamos de JSON a formato bytes lo que nos paso el paquete
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	//pasamos la respuesta de JSON a formato paquete que nos mando el server

	var respuesta PaqueteRecibidoDeMemoria //para eso declaramos una variable con el struct que esperamos que nos envie el server
	err = json.Unmarshal(body, &respuesta) //pasamos de bytes al formato de nuestro paquete lo que nos mando el server
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}
	log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)
	if respuestaJSON.StatusCode == http.StatusOK {
		log.Printf("Se pasa el proceso PID: %d a READY", pcb.Pid)
		PasarReady(pcb)
	}

	//en mi caso era un mensaje, por eso el struct tiene mensaje string, vos por ahi estas esperando 14 ints, no necesariamente un struct

}

func EnviarProcesoACPU(pcb *PCB, cpu *CPU) {

	var paquete PaqueteEnviadoKERNELaCPU

	paquete.PC = pcb.Pc
	paquete.Pid = pcb.Pid

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al convertir a json.")
		return
	}

	cliente := http.Client{} // Crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/KERNELCPU", cpu.Ip, cpu.Port) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {

		//aca tiene que haber un logger
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json") // Le avisa al server que manda la data en json format

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return
	}

	if respuestaJSON.StatusCode != http.StatusOK {
		log.Printf("Código de respuesta del server: %d\n", respuestaJSON.StatusCode)
		log.Printf("Status de respuesta el server no fue la esperada.\n")
		return
	}

	defer respuestaJSON.Body.Close() //cerramos algo supuestamente importante de cerrar pero no se que hace

	log.Printf("Conexion establecida con exito \n")
	//pasamos de JSON a formato bytes lo que nos paso el paquete
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	} //pasamos la respuesta de JSON a formato paquete que nos mando el server

	var respuesta PaqueteRecibido
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}
	log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)

}

func InterrumpirCPU(cpu *CPU) {

	var paquete PaqueteInterrupcion

	paquete.mensaje = "Interrupcion del proceso"

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al convertir a json.")
		return
	}

	cliente := http.Client{} // Crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/InterrupcionCPU", cpu.Ip, cpu.Port) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {

		//aca tiene que haber un logger
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json") // Le avisa al server que manda la data en json format

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return
	}

	if respuestaJSON.StatusCode != http.StatusOK {
		log.Printf("Código de respuesta del server: %d\n", respuestaJSON.StatusCode)
		log.Printf("Status de respuesta el server no fue la esperada.\n")
		return
	}

	defer respuestaJSON.Body.Close() //cerramos algo supuestamente importante de cerrar pero no se que hace

	log.Printf("Conexion establecida con exito \n")
	//pasamos de JSON a formato bytes lo que nos paso el paquete
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	} //pasamos la respuesta de JSON a formato paquete que nos mando el server

	var respuesta PaqueteRecibido
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}
	log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)

} //falta la respuesta de CPU

func InformarMemoriaFinProceso(pcb *PCB, ip string, puerto int) {

	var paquete PaqueteEnviadoKERNELaMemoria2
	paquete.Pid = pcb.Pid
	paquete.Mensaje = fmt.Sprintf("El proceso PID: %d  termino su ejecucion y se paso a EXIT", pcb.Pid)

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al convertir a json.")
		return
	}
	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/FinProceso", ip, puerto) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json") //le avisa al server que manda la data en json format

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	defer respuestaJSON.Body.Close() //cerramos algo supuestamente importante de cerrar pero no se que hace

	log.Printf("Conexion establecida con exito \n")
	//pasamos de JSON a formato bytes lo que nos paso el paquete
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	//pasamos la respuesta de JSON a formato paquete que nos mando el server

	var respuesta PaqueteRecibidoDeMemoria //para eso declaramos una variable con el struct que esperamos que nos envie el server
	err = json.Unmarshal(body, &respuesta) //pasamos de bytes al formato de nuestro paquete lo que nos mando el server
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}
	log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)
	PlanificadorLargoPlazo()
	PlanificadorCortoPlazo()

}

func CrearPCB(tamanio int, archivo string) { //pid unico arranca de 0
	pcbUsar := &PCB{
		Pid:                ContadorPCB,
		Pc:                 0,
		EstadoActual:       "NEW",
		TamProceso:         tamanio,
		MetricaEstados:     make(map[Estado]int),
		TiempoLlegada:      make(map[Estado]time.Time),
		TiempoEstados:      make(map[Estado]int64),
		Archivo:            archivo,
		TiempoEnvioExc:     time.Now(),
		RafagaAnterior:     0,
		EstimacionAnterior: globals.ClientConfig.Initial_estimate,
	}
	ColaNew = append(ColaNew, pcbUsar)

	log.Printf("## (<%d>) Se crea el proceso - Estado: NEW \n", pcbUsar.Pid)
	pcbUsar.MetricaEstados["NEW"]++
	pcbUsar.TiempoLlegada["NEW"] = time.Now()
	ContadorPCB++
	PlanificadorLargoPlazo()
}

func LeerConsola() string {
	// Leer de la consola
	reader := bufio.NewReader(os.Stdin)
	log.Println("Presione enter para inciar el planificador")
	text, _ := reader.ReadString('\n')
	//log.Print(text)
	return text
}

func IniciarPlanifcador(tamanio int, archivo string) {
	for true {
		text := LeerConsola()
		if text == "\n" {
			log.Printf("Planificador de largo plazo ejecutando")
			CrearPCB(tamanio, archivo)
			break
		}
	}
}

func PlanificadorLargoPlazo() {
	if len(ColaSuspReady) != 0 {
		pcbChequear := CriterioColaNew(ColaSuspReady)
		ConsultarProcesoConMemoria(pcbChequear, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	} else if len(ColaNew) != 0 {
		pcbChequear := CriterioColaNew(ColaNew)
		ConsultarProcesoConMemoria(pcbChequear, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	}
}

func PlanificadorCortoPlazo() {
	if len(ColaReady) != 0 {
		pcbChequear, hayDesalojo := CriterioColaReady()
		CPUDisponible, noEsVacio := TraqueoCPU() //drakukeo en su defecto
		if noEsVacio {
			log.Printf("se pasa el proceso PID: %d a EXECUTE", pcbChequear.Pid) //solo para saber que esta funcionando
			PasarExec(pcbChequear)
			CPUDisponible.Disponible = false
			CPUDisponible.Pid = pcbChequear.Pid //le asigno el pid al cpu que lo va a ejecutar
			EnviarProcesoACPU(pcbChequear, CPUDisponible)

		} else if hayDesalojo {
			pcbDesalojar, cpuDesalojar := RafagaMasLargaDeLosCPU()
			if calcularRafagaEstimada(pcbChequear) <= CalcularTiempoRestanteEjecucion(pcbDesalojar) {
				InterrumpirCPU(cpuDesalojar)
				PasarReady(pcbDesalojar)
				PasarExec(pcbChequear)
				cpuDesalojar.Pid = pcbChequear.Pid
			}
		}
	}
}

func FIFO(cola []*PCB) *PCB {
	if len(cola) == 0 {
		return &PCB{}
	}
	pcb := cola[0]
	return pcb
}

func ProcesoMasChicoPrimero(cola []*PCB) *PCB {
	if len(cola) == 0 {
		return &PCB{}
	}
	pcbTamanioMinimo := cola[0]
	for _, pcb := range cola {
		if pcb.TamProceso <= pcbTamanioMinimo.TamProceso {
			pcbTamanioMinimo = pcb
		}
	}
	return pcbTamanioMinimo
}

func Sjf() *PCB {
	if len(ColaReady) == 0 {
		return &PCB{}
	}
	pcbEstimacionMinima := ColaReady[0]
	for _, pcb := range ColaReady {
		if calcularRafagaEstimada(pcb) <= calcularRafagaEstimada(pcbEstimacionMinima) {
			pcbEstimacionMinima = pcb
		}
	}
	pcbEstimacionMinima.EstimacionAnterior = calcularRafagaEstimada(pcbEstimacionMinima)
	return pcbEstimacionMinima
}

func RafagaMasLargaDeLosCPU() (*PCB, *CPU) {
	if len(ListaExec) == 0 {
		return &PCB{}, &CPU{}
	}
	pcbEstimacionMasLarga := ListaExec[0]
	for _, pcb := range ListaExec {
		if CalcularTiempoRestanteEjecucion(pcb) >= CalcularTiempoRestanteEjecucion(pcbEstimacionMasLarga) { //estamiacionAnterior - calcularTiempoEnExec
			pcbEstimacionMasLarga = pcb
		}
	}
	//pcbEstimacionMinima.EstimacionAnterior = calcularRafagaEstimada(pcbEstimacionMinima)
	cpuConLaRafagaLarga := ObtenerCpuEnFuncionDelPid(pcbEstimacionMasLarga.Pid)

	return pcbEstimacionMasLarga, &cpuConLaRafagaLarga
}

func CalcularTiempoRestanteEjecucion(pcb *PCB) float32 {
	return pcb.EstimacionAnterior - float32(pcb.TiempoEnvioExc.Sub(time.Now()))
}

func calcularRafagaEstimada(pcb *PCB) float32 {
	return globals.ClientConfig.Alpha*pcb.RafagaAnterior + (1-globals.ClientConfig.Alpha)*pcb.EstimacionAnterior
}

func PasarReady(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado READY  \n", pcb.Pid, pcb.EstadoActual)
	ColaReady = append(ColaReady, pcb)
	ColaNew = removerPCB(ColaNew, pcb)
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "READY"
	pcb.MetricaEstados["READY"]++
	pcb.TiempoLlegada["READY"] = time.Now()

	PlanificadorCortoPlazo()
}

func PasarExec(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado EXECUTE \n", pcb.Pid, pcb.EstadoActual)
	ListaExec = append(ListaExec, pcb)
	ColaReady = removerPCB(ColaReady, pcb)
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "EXECUTE"
	pcb.TiempoLlegada["EXECUTE"] = time.Now()
	pcb.MetricaEstados["EXECUTE"]++
	pcb.TiempoEnvioExc = time.Now()

}

func PasarBlocked(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado BLOCKED \n", pcb.Pid, pcb.EstadoActual)
	ColaBlock = append(ColaBlock, pcb)
	ListaExec = removerPCB(ListaExec, pcb)
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "BLOCKED"
	pcb.MetricaEstados["BLOCKED"]++
	pcb.TiempoLlegada["BLOCKED"] = time.Now()

	PlanificadorCortoPlazo()
}

func removerPCB(cola []*PCB, pcb *PCB) []*PCB {
	for i, item := range cola {
		if item.Pid == pcb.Pid {
			return append(cola[:i], cola[1+i:]...)
		}
	}
	return cola
}

func CriterioColaNew(cola []*PCB) *PCB {
	if globals.ClientConfig.Ready_ingress_algorithm == "FIFO" {
		return FIFO(cola)
	} else {
		return ProcesoMasChicoPrimero(cola)
	}

}

func CriterioColaReady() (*PCB, bool) {
	if globals.ClientConfig.Scheduler_algorithm == "FIFO" {
		return FIFO(ColaReady), false
	} else if globals.ClientConfig.Scheduler_algorithm == "SJF" {
		return Sjf(), false
	} else {
		return Sjf(), true
	}
}

func TraqueoCPU() (*CPU, bool) {
	for i := range ListaCPU {
		if ListaCPU[i].Disponible {
			return &ListaCPU[i], true
		}
	}
	return nil, false
} //Esto busca un cpu disponible

func crearStructCPU(ip string, puerto int, instancia string) {
	ListaCPU = append(ListaCPU, CPU{
		Ip:         ip,
		Port:       puerto,
		Disponible: true,
		Instancia:  instancia,
	})
}

func ObtenerCpu(instancia string) CPU {
	for _, cpu := range ListaCPU {
		if cpu.Instancia == instancia {
			return cpu
		}
	}
	return CPU{}
} //Nos dice que instancia de CPU es

func ObtenerCpuEnFuncionDelPid(pid int) CPU {
	for _, cpu := range ListaCPU {
		if cpu.Pid == pid {
			return cpu
		}
	}
	return CPU{}
} //Nos dice que instancia de CPU es

func FinalizarProceso(pcb *PCB) {
	log.Printf("El proceso PID: %d termino su ejecucion y se paso a EXIT \n", pcb.Pid)
	pcb.EstadoActual = "EXIT"
	pcb.MetricaEstados["EXIT"]++
	ColaExit = append(ColaExit, pcb) //es un esquema de como podria finalizar el proceso, puede cambiarse esto
	InformarMemoriaFinProceso(pcb, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	log.Printf("## (<%d>) - Finaliza el proceso \n", pcb.Pid)
	//Métricas de Estado: “## (<PID>) - Métricas de estado: NEW (NEW_COUNT) (NEW_TIME), READY (READY_COUNT) (READY_TIME), …”
	log.Printf("## (<%d>) - Métricas de estado: NEW NEW_COUNT: %d NEW_TIME: %d, READY READY_COUNT: %d READY_TIME: %d, EXECUTE EXECUTE_COUNT: %d EXECUTE_TIME: %d, BLOCKED BLOCKED_COUNT: %d BLOCKED_TIME: %d", pcb.Pid, pcb.MetricaEstados["NEW"], pcb.TiempoEstados["NEW"], pcb.MetricaEstados["READY"], pcb.TiempoEstados["READY"], pcb.MetricaEstados["EXECUTE"], pcb.TiempoEstados["EXECUTE"], pcb.MetricaEstados["BLOCKED"], pcb.TiempoEstados["BLOCKED"])

}

func CrearStructIO(ip string, puerto int, instancia string) {
	ListaIO = append(ListaIO, IO{
		Ip:           ip,
		Port:         puerto,
		Instancia:    instancia,
		ColaProcesos: []PCBIO{},
		Disponible:   true,
	})
}

func ObtenerIO(instancia string) IO {
	for _, io := range ListaIO {
		if io.Instancia == instancia {
			return io
		}
	}
	return IO{}
}

func ExisteIO(instancia string) bool {
	for _, io := range ListaIO {
		if io.Instancia == instancia {
			return true
		}
	}
	return false
}

func AgregarColaIO(io IO, pid int, tiempo int) {
	io.ColaProcesos = append(io.ColaProcesos, PCBIO{
		Pid:    pid,
		Tiempo: tiempo,
	})
}

func ObtenerPCB(pid int) *PCB {
	for _, pcb := range ListaExec {
		if pcb.Pid == pid {
			return pcb
		}
	}
	return &PCB{}
}

func MandarProcesoAIO(io IO) {
	if io.Disponible {
		io.Disponible = false
		go UtilizarIO(io, io.ColaProcesos[0].Pid, io.ColaProcesos[0].Tiempo)

	}
}

func DumpDelProceso(pcb *PCB, ip string, puerto int) {

	var paquete PaqueteEnviadoKERNELaMemoria2
	paquete.Pid = pcb.Pid
	paquete.Mensaje = fmt.Sprintf("El proceso PID: %d  requiere que se haga un DUMP del mismo", pcb.Pid)

	PasarBlocked(pcb)

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al convertir a json.")
		return
	}
	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/KERNELMEMORIADUMP", ip, puerto) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json") //le avisa al server que manda la data en json format

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	defer respuestaJSON.Body.Close() //cerramos algo supuestamente importante de cerrar pero no se que hace

	log.Printf("Conexion establecida con exito \n")
	//pasamos de JSON a formato bytes lo que nos paso el paquete
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	//pasamos la respuesta de JSON a formato paquete que nos mando el server

	var respuesta PaqueteRecibidoDeMemoria //para eso declaramos una variable con el struct que esperamos que nos envie el server
	err = json.Unmarshal(body, &respuesta) //pasamos de bytes al formato de nuestro paquete lo que nos mando el server
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}

	log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)

	if respuestaJSON.StatusCode == http.StatusOK {
		log.Printf("Se pudo hacer el DUMP del proceso con el PID: %d ", pcb.Pid)
		PasarReady(pcb)
	} else {
		log.Printf("No se pudo hacer el DUMP del proceso con el PID: %d ", pcb.Pid)
		FinalizarProceso(pcb) //Mando a exit al proceso
	}

}

func FinalizarProcesosIO(io *IO) {
	for _, proceso := range io.ColaProcesos {
		pcb := ObtenerPCB(proceso.Pid)
		if pcb != nil {
			log.Printf("El proceso PID: %d  no pudo ser atendido por el I/O %s y se pasa a EXIT", pcb.Pid, io.Instancia)
			FinalizarProceso(pcb)
		}
	}
	ListaIO = removerIO(io) //alto gil
}

func removerIO(io *IO) []IO {
	for i, item := range ListaIO {
		if item.Instancia == io.Instancia {
			return append(ListaIO[:i], ListaIO[1+i:]...)
		}
	}
	return ListaIO
}

func enviarExitProcesosIO(io IO) {
	for _, proceso := range io.ColaProcesos {
		pcb := ObtenerPCB(proceso.Pid)
		if pcb != nil {
			log.Printf("El proceso PID: %d  se pasa a EXIT por desconexion del I/O %s", pcb.Pid, io.Instancia)
			FinalizarProceso(pcb)
		}
	}
}
