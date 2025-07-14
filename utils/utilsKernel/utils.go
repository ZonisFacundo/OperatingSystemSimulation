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

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//og.Printf("El cliente nos mando esto: \n nombre: %s  \n puerto: %d \n IP: %s \n", request.Nombre, request.Puerto, request.Ip)

	CrearStructIO(request.Ip, request.Puerto, request.Nombre)

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
	ListaIO = removerIO(ioCerrada)

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

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//log.Printf("Handshake recibido: Port: %d - Instance: %s - Ip: %s", request.Puerto, request.Instancia, request.Ip)

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

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//log.Printf("contexto de devolucion del proceso: %s", request.Syscall)

	var respuesta RespuestaalCPU

	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}
	//log.Printf("Conexion establecida con exito \n")
	cpuServidor := ObtenerCpu(request.InstanciaCPU)
	PCBUtilizar := ObtenerPCB(cpuServidor.Pid)
	PCBUtilizar.Pc = request.Pc
	PCBUtilizar.RafagaAnterior = float32(PCBUtilizar.TiempoEnvioExc.Sub(time.Now()))

	switch request.Syscall {
	case "IO":
		respuesta.Mensaje = "interrupcion"
		cpuServidor.Disponible = true
		if ExisteIO(request.Parametro2) {
			SemCortoPlazo <- struct{}{}
			ioServidor := ObtenerIO(request.Parametro2)
			go PlanificadorMedianoPlazo(PCBUtilizar)
			AgregarColaIO(ioServidor, PCBUtilizar, request.Parametro1)
			PasarBlocked(PCBUtilizar)
			log.Printf("## (<%d>) - Bloqueado por IO: < %s > \n", PCBUtilizar.Pid, ioServidor.Instancia)
			MandarProcesoAIO(ioServidor)
			log.Printf("YA LO MANDE A IO PARA PID: (<%d>) - IO SERVIDOR INSTANCIA: < %s > \n", PCBUtilizar.Pid, ioServidor.Instancia)

			/*if len(ioServidor.ColaProcesos) > 0 {
				ioServidor.ColaProcesos = ioServidor.ColaProcesos[1:]
			}*/
		} else {
			FinalizarProceso(PCBUtilizar)
		}

	case "EXIT":
		log.Printf("## (<%d>) - Solicitó syscall: <EXIT> \n", PCBUtilizar.Pid)
		respuesta.Mensaje = "interrupcion"
		cpuServidor.Disponible = true
		FinalizarProceso(PCBUtilizar)

	case "DUMP_MEMORY":
		log.Printf("## (<%d>) - Solicitó syscall: <DUMP_MEMORY> \n", PCBUtilizar.Pid)
		respuesta.Mensaje = "interrupcion"
		cpuServidor.Disponible = true
		SemCortoPlazo <- struct{}{}
		DumpDelProceso(PCBUtilizar, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)

	case "INIT_PROC":
		log.Printf("## (<%d>) - Solicitó syscall: <INIT_PROC> \n", PCBUtilizar.Pid)
		respuesta.Mensaje = "NO INTERRUMPAS GIL"
		CrearPCB(request.Parametro1, request.Parametro2)
		cpuServidor.Disponible = false
		EnviarProcesoACPU(PCBUtilizar, cpuServidor)
		w.WriteHeader(http.StatusOK)

		w.Write(respuestaJSON)
		return
	}

	w.WriteHeader(http.StatusFound)
	w.Write(respuestaJSON)

}

func UtilizarIO(ioServer *IO, pcb *PCB, tiempo int) {

	log.Printf("\n\n\n CHUPETE EN EL ORTO INSIIIIIDE usado por PID: %d \n\n\n", pcb.Pid)
	var paquete PaqueteEnviadoKERNELaIO
	paquete.Pid = pcb.Pid
	paquete.Tiempo = tiempo

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {

		log.Printf("Error al convertir a json.")
		return
	}
	cliente := http.Client{}

	url := fmt.Sprintf("http://%s:%d/KERNELIO", ioServer.Ip, ioServer.Port)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))
	if err != nil {

		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	respuestaJSON, err := cliente.Do(req)

	if err != nil {
		FinalizarProceso(pcb)
		return

	}

	defer respuestaJSON.Body.Close()

	//log.Printf("Conexion establecida con exito \n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}
	var respuesta PaqueteRecibidoDeIO
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		//log.Printf("Error al decodificar el JSON.\n")
		return
	}

	if respuestaJSON.StatusCode == http.StatusOK {

		//log.Printf("La respuesta del I/O %s fue: %s\n", ioServer.Instancia, respuesta.Mensaje)

		log.Printf("## (<%d>) finalizó IO y pasa a READY \n", pcb.Pid)

		if EstaEnColaBlock(pcb) {
			log.Printf("pasa de blocked a ready el pid %d\n", pcb.Pid)
			PasarReady(pcb, ColaBlock)
			//RemoverDeColaProcesoIO(ioServer)

		} else if EstaEnColaSuspBlock(pcb) {
			log.Printf("pasa desde %s a susp ready el pid %d\n", pcb.EstadoActual, pcb.Pid)

			PasarSuspReady(pcb)
			//RemoverDeColaProcesoIO(ioServer)

		} else {
			log.Printf("mira flaco, este pcb no esta ni en blocked ni en susp blocked \n")
		}
		ioServer.Disponible = true
		RemoverDeColaProcesoIO(ioServer)
		MandarProcesoAIO(ioServer)

	}
}

func EstaEnColaSuspBlock(pcbChequear *PCB) bool {
	for _, pcb := range ColaSuspBlock {
		if pcb == pcbChequear {
			return true
		}
	}
	return false
}

func ConsultarProcesoConMemoria(pcb *PCB, ip string, puerto int, cola []*PCB) {

	var paquete PaqueteEnviadoKERNELaMemoria
	paquete.Pid = pcb.Pid
	paquete.TamProceso = pcb.TamProceso
	paquete.Archivo = pcb.Archivo

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {

		log.Printf("Error al convertir a json.")
		return
	}
	cliente := http.Client{}

	url := fmt.Sprintf("http://%s:%d/KERNELMEMORIA", ip, puerto)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {

		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	defer respuestaJSON.Body.Close()

	//log.Printf("Pregunto si puedo pasar a ready un proceso \n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibidoDeMemoria
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}
	log.Printf("La respuesta del server memoria fue: %s\n", respuesta.Mensaje)
	if respuestaJSON.StatusCode == http.StatusOK {
		log.Printf("Se pasa el proceso PID: %d a READY", pcb.Pid)
		PasarReady(pcb, cola)
	} else {
		log.Printf("no se puede pasar a ready al PID: %d porque memoria basicamnete nos dijo que hay quilombo", pcb.Pid)

	}

}

func EnviarProcesoACPU(pcb *PCB, cpu *CPU) {

	var paquete PaqueteEnviadoKERNELaCPU

	paquete.PC = pcb.Pc
	paquete.Pid = pcb.Pid

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.")
		return
	}

	cliente := http.Client{}

	url := fmt.Sprintf("http://%s:%d/KERNELCPU", cpu.Ip, cpu.Port)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {

		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return
	}

	if respuestaJSON.StatusCode != http.StatusOK {
		//log.Printf("Código de respuesta del server: %d\n", respuestaJSON.StatusCode)
		//log.Printf("Status de respuesta el server no fue la esperada.\n")
		return
	}

	defer respuestaJSON.Body.Close()
	//log.Printf("Conexion establecida con exito \n")

	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibido
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		//log.Printf("Error al decodificar el JSON.\n")
		return
	}
	//log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)

}

func InterrumpirCPU(cpu *CPU) {

	var paquete PaqueteInterrupcion

	paquete.Mensaje = "Interrupcion del proceso"

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.")
		return
	}

	cliente := http.Client{}

	url := fmt.Sprintf("http://%s:%d/INTERRUPCIONCPU", cpu.Ip, cpu.Port)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {

		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return
	}

	if respuestaJSON.StatusCode != http.StatusOK {
		//log.Printf("Código de respuesta del server: %d\n", respuestaJSON.StatusCode)
		//log.Printf("Status de respuesta el server no fue la esperada.\n")
		return
	}

	defer respuestaJSON.Body.Close()

	//	log.Printf("Conexion establecida con exito \n")

	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibidoDeCPU
	pcb := ObtenerPCB(respuesta.Pid)
	pcb.Pc = respuesta.Pc
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		//log.Printf("Error al decodificar el JSON.\n")
		return
	}
	//log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)
	//log.Printf("PID: %d PC: %d", respuesta.Pid, respuesta.Pc)
}

func InformarMemoriaFinProceso(pcb *PCB, ip string, puerto int) {

	var paquete PaqueteEnviadoKERNELaMemoria2
	paquete.Pid = pcb.Pid
	paquete.Mensaje = fmt.Sprintf("El proceso PID: %d  termino su ejecucion y se paso a EXIT", pcb.Pid)

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.")
		return
	}
	cliente := http.Client{}

	url := fmt.Sprintf("http://%s:%d/FinProceso", ip, puerto)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	defer respuestaJSON.Body.Close()
	//log.Printf("Conexion establecida con exito \n")

	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibidoDeMemoria
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		//log.Printf("Error al decodificar el JSON.\n")
		return
	}
	log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)
	SemLargoPlazo <- struct{}{}
	SemCortoPlazo <- struct{}{}

}

func CrearPCB(tamanio int, archivo string) {
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
	MutexColaNew.Lock()
	ColaNew = append(ColaNew, pcbUsar)
	MutexColaNew.Unlock()

	log.Printf("## (<%d>) Se crea el proceso - Estado: NEW \n", pcbUsar.Pid)
	pcbUsar.MetricaEstados["NEW"]++
	pcbUsar.TiempoLlegada["NEW"] = time.Now()
	ContadorPCB++
	SemLargoPlazo <- struct{}{}
}

/*
func CrearPCBPrueba(tamanio int, archivo string) { //pid unico arranca de 0

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
			EstimacionAnterior: 5,
		}
		MutexColaNew.Lock()
		ColaNew = append(ColaNew, pcbUsar)
		MutexColaNew.Unlock()

		log.Printf("## (<%d>) Se crea el proceso - Estado: NEW \n", pcbUsar.Pid)
		pcbUsar.MetricaEstados["NEW"]++
		pcbUsar.TiempoLlegada["NEW"] = time.Now()
		ContadorPCB++
		SemLargoPlazo <- struct{}{}
	}
*/

func LeerConsola() string {
	reader := bufio.NewReader(os.Stdin)
	log.Println("Presione enter para inciar el planificador")
	text, _ := reader.ReadString('\n')
	return text
}

func IniciarPlanifcador(tamanio int, archivo string) {
	for true {
		text := LeerConsola()
		if text == "\n" {
			//log.Printf("Planificador de largo plazo ejecutando")
			CrearPCB(tamanio, archivo)
			break
		}
	}
}

// Dejo nuestro planificador porque es el que deberia de funcar mas mejor (aproposito)

func PlanificadorLargoPlazo() {
	for true {
		<-SemLargoPlazo
		if len(ColaSuspReady) != 0 {
			MutexColaNew.Lock()
			pcbChequear := CriterioColaNew(ColaSuspReady)
			MutexColaNew.Unlock()
			ConsultarProcesoConMemoria(pcbChequear, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, ColaSuspReady)

			/*
				MutexColaSuspReady.Lock()                      // <<< CAMBIO: lock de ColaSuspReady (antes MutexColaNew)
				pcb := CriterioColaNew(ColaSuspReady)          // <<< igual criterio
				ColaSuspReady = removerPCB(ColaSuspReady, pcb) // <<< CAMBIO: quito manualmente de SuspReady
				MutexColaSuspReady.Unlock()

				PasarReady(pcb, ColaSuspReady) // <<< CAMBIO: paso directo a READY
				continue                       // <<< CAMBIO: sigo al siguiente ciclo
			*/

		} else if len(ColaNew) != 0 {
			MutexColaNew.Lock()
			pcbChequear := CriterioColaNew(ColaNew)
			MutexColaNew.Unlock()
			ConsultarProcesoConMemoria(pcbChequear, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory, ColaNew)

		} //else {
		//	SemLargoPlazo <- struct{}{}
		//time.Sleep(1 * time.Second)

		//}

	}
}

// CODIGO PRE CAMBIO DE FACU KERNEL

/*
func PlanificadorLargoPlazo() {
	for {
		<-SemLargoPlazo

		// 1) Si hay procesos en SUSP.READY, pasar a READY:
		if len(ColaSuspReady) != 0 {
			MutexColaSuspReady.Lock()
			pcb := CriterioColaNew(ColaSuspReady)
			// Lo saco de la cola de suspendidos listos
			ColaSuspReady = removerPCB(ColaSuspReady, pcb)
			MutexColaSuspReady.Unlock()

			// Y lo paso a READY (dispara SemCortoPlazo internamente)
			PasarReady(pcb, ColaSuspReady)

			continue // vuelvo al for, ya consumí esta señal
		}

		// 2) Si hay procesos en SUSP.BLOCKED, hacer swap‑in y pasarlos a SUSP.READY:
		if len(ColaSuspBlock) != 0 {
			MutexColaSuspBlock.Lock()
			pcb := CriterioColaNew(ColaSuspBlock)
			// Lo quito de SUSP.BLOCKED
			ColaSuspBlock = removerPCB(ColaSuspBlock, pcb)
			MutexColaSuspBlock.Unlock()

			// Swap‑in en memoria
			SwapInProceso(pcb)

			// Lo paso a SUSP.READY (no dispara semáforo por defecto)
			PasarSuspReady(pcb)

			continue
		}

		// 3) Por último, los procesos NEW entran a memoria:
		if len(ColaNew) != 0 {
			MutexColaNew.Lock()
			pcb := CriterioColaNew(ColaNew)
			ColaNew = removerPCB(ColaNew, pcb)
			MutexColaNew.Unlock()

			ConsultarProcesoConMemoria(
				pcb,
				globals.ClientConfig.Ip_memory,
				globals.ClientConfig.Port_memory,
				ColaNew,
			)
		}
	}
}
*/

func PlanificadorCortoPlazo() {
	for true {
		<-SemCortoPlazo

		if len(ColaReady) != 0 {
			MutexColaReady.Lock()
			pcbChequear, hayDesalojo := CriterioColaReady()
			MutexColaReady.Unlock()
			CPUDisponible, noEsVacio := TraqueoCPU()
			if noEsVacio {
				//log.Printf("se pasa el proceso PID: %d a EXECUTE", pcbChequear.Pid)
				PasarExec(pcbChequear)
				CPUDisponible.Disponible = false
				CPUDisponible.Pid = pcbChequear.Pid
				EnviarProcesoACPU(pcbChequear, CPUDisponible)

			} else if hayDesalojo {
				pcbDesalojar, cpuDesalojar := RafagaMasLargaDeLosCPU()
				if calcularRafagaEstimada(pcbChequear) < CalcularTiempoRestanteEjecucion(pcbDesalojar) {
					InterrumpirCPU(cpuDesalojar)
					log.Printf("## (<%d>) - Desalojado por algoritmo SJF/SRT \n", cpuDesalojar.Pid)
					PasarReady(pcbDesalojar, ListaExec)
					PasarExec(pcbChequear)
					cpuDesalojar.Pid = pcbChequear.Pid
				}
			} else {
				SemCortoPlazo <- struct{}{}
				time.Sleep(1 * time.Second)
			}
		} // else {
		//SemCortoPlazo <- struct{}{}
		//time.Sleep(1 * time.Second)

		//}
	}
}

func PlanificadorMedianoPlazo(pcb *PCB) {
	pcb.TiempoEnvioBlock = time.Now()
	for true {
		if EstaEnColaBlock(pcb) {
			if time.Since(pcb.TiempoEnvioBlock) >= time.Duration(globals.ClientConfig.Suspension_time)*time.Millisecond {
				PasarSuspBlock(pcb)
				break
			}
		} else {
			break
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
		if calcularRafagaEstimada(pcb) < calcularRafagaEstimada(pcbEstimacionMinima) {
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
		if CalcularTiempoRestanteEjecucion(pcb) >= CalcularTiempoRestanteEjecucion(pcbEstimacionMasLarga) {
			pcbEstimacionMasLarga = pcb
		}
	}
	//pcbEstimacionMinima.EstimacionAnterior = calcularRafagaEstimada(pcbEstimacionMinima)
	cpuConLaRafagaLarga := ObtenerCpuEnFuncionDelPid(pcbEstimacionMasLarga.Pid)

	return pcbEstimacionMasLarga, cpuConLaRafagaLarga
}

func CalcularTiempoRestanteEjecucion(pcb *PCB) float32 {
	return pcb.EstimacionAnterior - float32(pcb.TiempoEnvioExc.Sub(time.Now()))
}

func calcularRafagaEstimada(pcb *PCB) float32 {
	return globals.ClientConfig.Alpha*pcb.RafagaAnterior + (1-globals.ClientConfig.Alpha)*pcb.EstimacionAnterior
}

func PasarReady(pcb *PCB, colaSacar []*PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado READY  \n", pcb.Pid, pcb.EstadoActual)
	MutexColaReady.Lock()
	ColaReady = append(ColaReady, pcb)
	MutexColaReady.Unlock()
	MutexColaNew.Lock()
	ColaNew = removerPCB(ColaNew, pcb)
	ColaBlock = removerPCB(ColaBlock, pcb)
	ColaSuspBlock = removerPCB(ColaSuspBlock, pcb)
	ListaExec = removerPCB(ListaExec, pcb)
	MutexColaNew.Unlock()
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "READY"
	pcb.MetricaEstados["READY"]++
	pcb.TiempoLlegada["READY"] = time.Now()

	SemCortoPlazo <- struct{}{}
}

func PasarExec(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado EXECUTE \n", pcb.Pid, pcb.EstadoActual)
	MutexListaExec.Lock()
	ListaExec = append(ListaExec, pcb)
	MutexListaExec.Unlock()
	MutexColaReady.Lock()
	ColaReady = removerPCB(ColaReady, pcb)
	MutexColaReady.Unlock()
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "EXECUTE"
	pcb.TiempoLlegada["EXECUTE"] = time.Now()
	pcb.MetricaEstados["EXECUTE"]++
	pcb.TiempoEnvioExc = time.Now()

}

func PasarBlocked(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado BLOCKED \n", pcb.Pid, pcb.EstadoActual)
	MutexColaBlock.Lock()
	ColaBlock = append(ColaBlock, pcb)
	MutexColaBlock.Unlock()
	MutexListaExec.Lock()
	ListaExec = removerPCB(ListaExec, pcb)
	MutexListaExec.Unlock()
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "BLOCKED"
	pcb.MetricaEstados["BLOCKED"]++
	pcb.TiempoLlegada["BLOCKED"] = time.Now()

	SemCortoPlazo <- struct{}{}

}

//CODIGO PREVIO A TOQUETEO FACU

/*
func PasarSuspBlock(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado SUSP.BLOCKED \n", pcb.Pid, pcb.EstadoActual)
	SwapDelProceso(pcb, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	MutexColaSuspBlock.Lock()
	ColaSuspBlock = append(ColaSuspBlock, pcb)
	MutexColaSuspBlock.Unlock()
	MutexColaBlock.Lock()
	ColaBlock = removerPCB(ColaBlock, pcb)
	MutexColaBlock.Unlock()
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "SUSP.BLOCKED"
	pcb.MetricaEstados["SUSP.BLOCKED"]++
	pcb.TiempoLlegada["SUSP.BLOCKED"] = time.Now()
}
*/

func PasarSuspBlock(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado SUSP.BLOCKED \n", pcb.Pid, pcb.EstadoActual)

	// <<< CAMBIO: añadir a ColaSuspBlock, no a ColaBlock
	log.Printf("\n \n\n## (<%d>) ESTOY ARRIBA DEL MUTEX \n", pcb.Pid)
	MutexColaSuspBlock.Lock()
	log.Printf("\n \n\n## (<%d>) ESTOY MUY ABAJO DEL MUTEX \n", pcb.Pid)
	ColaSuspBlock = append(ColaSuspBlock, pcb)
	log.Printf("\n \n\n## (<%d>) ESTOY EN SUSPENDE BLOCK FORRO \n", pcb.Pid)
	MutexColaSuspBlock.Unlock()

	// quitar de la cola BLOCK normal
	MutexColaBlock.Lock()
	ColaBlock = removerPCB(ColaBlock, pcb)
	MutexColaBlock.Unlock()

	// actualizar métricas y estado
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "SUSP.BLOCKED"
	pcb.MetricaEstados["SUSP.BLOCKED"]++
	pcb.TiempoLlegada["SUSP.BLOCKED"] = time.Now()

	SwapDelProceso(pcb, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
}

func PasarSuspReady(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado SUSP.READY \n", pcb.Pid, pcb.EstadoActual)
	//pedimos a la memoria q traiga devuelta el proceso del swap
	//SwapInProceso(pcb)
	MutexColaSuspReady.Lock()
	ColaSuspReady = append(ColaSuspReady, pcb)
	MutexColaSuspReady.Unlock()
	MutexColaSuspBlock.Lock()
	ColaSuspBlock = removerPCB(ColaSuspBlock, pcb)
	MutexColaSuspBlock.Unlock()
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "SUSP.READY"
	pcb.MetricaEstados["SUSP.READY"]++
	pcb.TiempoLlegada["SUSP.READY"] = time.Now()

	SwapInProceso(pcb)
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
	} else if globals.ClientConfig.Ready_ingress_algorithm == "PMCP" {
		return ProcesoMasChicoPrimero(cola)
	}
	return &PCB{}

}

func CriterioColaReady() (*PCB, bool) {
	if globals.ClientConfig.Scheduler_algorithm == "FIFO" {
		return FIFO(ColaReady), false
	} else if globals.ClientConfig.Scheduler_algorithm == "SJF" {
		return Sjf(), false
	} else if globals.ClientConfig.Scheduler_algorithm == "SJF/SRT" {
		return Sjf(), true
	}
	return &PCB{}, false
}

func TraqueoCPU() (*CPU, bool) {
	for i := range ListaCPU {
		if ListaCPU[i].Disponible {
			return &ListaCPU[i], true
		}
	}
	return nil, false
}

func crearStructCPU(ip string, puerto int, instancia string) {
	ListaCPU = append(ListaCPU, CPU{
		Ip:         ip,
		Port:       puerto,
		Disponible: true,
		Instancia:  instancia,
	})
}

func CrearStructCPU2(ip string, puerto int, instancia string) CPU {
	return (CPU{
		Ip:         ip,
		Port:       puerto,
		Disponible: true,
		Instancia:  instancia,
	})

}

func ObtenerCpu(instancia string) *CPU {
	for i := range ListaCPU {
		if ListaCPU[i].Instancia == instancia {
			return &ListaCPU[i]
		}
	}
	return nil
}

func ObtenerCpuEnFuncionDelPid(pid int) *CPU {
	for i := range ListaCPU {
		if ListaCPU[i].Pid == pid {
			return &ListaCPU[i]
		}
	}
	return nil
}
func FinalizarProceso(pcb *PCB) {
	//log.Printf("El proceso PID: %d termino su ejecucion y se paso a EXIT \n", pcb.Pid)
	pcb.EstadoActual = "EXIT"
	log.Printf("\n\nEl proceso PID: %d esta tratande de exitearrrrrr \n", pcb.Pid)
	pcb.MetricaEstados["EXIT"]++
	log.Printf("El proceso PID: %d paso por el map con las manos arriba tomando tequila \n\n\n", pcb.Pid)

	ColaExit = append(ColaExit, pcb)
	InformarMemoriaFinProceso(pcb, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	log.Printf("## (<%d>) - Finaliza el proceso \n", pcb.Pid)
	log.Printf("## (<%d>) - Métricas de estado: NEW NEW_COUNT: %d NEW_TIME: %d, READY READY_COUNT: %d READY_TIME: %d, EXECUTE EXECUTE_COUNT: %d EXECUTE_TIME: %d, BLOCKED BLOCKED_COUNT: %d BLOCKED_TIME: %d, SUSP.BLOCKED  SUSP.BLOCKED_COUNT: %d SUSP.BLOCKED_TIME: %d, SUSP.READY  SUSP.READY_COUNT: %d SUSP.READY_TIME: %d \n", pcb.Pid, pcb.MetricaEstados["NEW"], pcb.TiempoEstados["NEW"], pcb.MetricaEstados["READY"], pcb.TiempoEstados["READY"], pcb.MetricaEstados["EXECUTE"], pcb.TiempoEstados["EXECUTE"], pcb.MetricaEstados["BLOCKED"], pcb.TiempoEstados["BLOCKED"], pcb.MetricaEstados["SUSP.BLOCKED"], pcb.TiempoEstados["SUSP.BLOCKED"], pcb.MetricaEstados["SUSP.READY"], pcb.TiempoEstados["SUSP.READY"])
}

/*
FACU facu funcion original (no modificada por facu)

func FinalizarProceso(pcb *PCB) {
	//log.Printf("El proceso PID: %d termino su ejecucion y se paso a EXIT \n", pcb.Pid)
	pcb.EstadoActual = "EXIT"
	pcb.MetricaEstados["EXIT"]++
	ColaExit = append(ColaExit, pcb)
	InformarMemoriaFinProceso(pcb, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	log.Printf("## (<%d>) - Finaliza el proceso \n", pcb.Pid)
	log.Printf("## (<%d>) - Métricas de estado: NEW NEW_COUNT: %d NEW_TIME: %d, READY READY_COUNT: %d READY_TIME: %d, EXECUTE EXECUTE_COUNT: %d EXECUTE_TIME: %d, BLOCKED BLOCKED_COUNT: %d BLOCKED_TIME: %d, SUSP.BLOCKED  SUSP.BLOCKED_COUNT: %d SUSP.BLOCKED_TIME: %d, SUSP.READY  SUSP.READY_COUNT: %d SUSP.READY_TIME: %d \n", pcb.Pid, pcb.MetricaEstados["NEW"], pcb.TiempoEstados["NEW"], pcb.MetricaEstados["READY"], pcb.TiempoEstados["READY"], pcb.MetricaEstados["EXECUTE"], pcb.TiempoEstados["EXECUTE"], pcb.MetricaEstados["BLOCKED"], pcb.TiempoEstados["BLOCKED"], pcb.MetricaEstados["SUSP.BLOCKED"], pcb.TiempoEstados["SUSP.BLOCKED"], pcb.MetricaEstados["SUSP.READY"], pcb.TiempoEstados["SUSP.READY"])

}
*/

func CrearStructIO(ip string, puerto int, instancia string) {
	ListaIO = append(ListaIO, IO{
		Ip:           ip,
		Port:         puerto,
		Instancia:    instancia,
		ColaProcesos: []PCBIO{},
		Disponible:   true,
	})
}

func ObtenerIO(instancia string) *IO {
	for i := range ListaIO {
		if ListaIO[i].Instancia == instancia {
			return &ListaIO[i]
		}
	}
	return nil
}

func ExisteIO(instancia string) bool {
	for _, io := range ListaIO {
		if io.Instancia == instancia {
			return true
		}
	}
	return false
}

func AgregarColaIO(io *IO, pcb *PCB, tiempo int) {
	io.ColaProcesos = append(io.ColaProcesos, PCBIO{
		Pcb:    pcb,
		Tiempo: tiempo,
	})
}

/*
	func RemoverDeColaProcesoIO(io *IO) []PCBIO {
		//return append(io.ColaProcesos[:0], io.ColaProcesos[1+0:]...)
	}
*/
func RemoverDeColaProcesoIO(io *IO) {
	if len(io.ColaProcesos) == 0 {
		return
	}
	// Reconstruye la slice sin el primer elemento
	io.ColaProcesos = io.ColaProcesos[1:]
}

func ObtenerPCB(pid int) *PCB {

	for _, pcb := range ListaExec {
		if pcb.Pid == pid {
			return pcb
		}
	}
	//FACU Facu facu
	//si llega aca es porque no encontro en la listaEXECUTE a dicho pid (ocurre en EXIT que lo saca de esa lista antes de ejecutar esto, entonces rompe todo cuando se quiere acceder a un PCB nulo como el que devuelven)
	//dejo esto como esta y voy a finalizar proceso, saco el acceso a EXIT a ese PCB y lo muevo a justo antes de que se lo quite de la lista execute
	for i := 0; i < 20; i++ {
		log.Printf("voy a devolver uno nulo pibeeeeee dice el proceso de pid: %d\n", pid)
	}
	return &PCB{}
}

func EstaEnColaBlock(pcbChequear *PCB) bool {
	for _, pcb := range ColaBlock {
		if pcb == pcbChequear {
			return true
		}
	}
	return false
}

func MandarProcesoAIO(io *IO) {

	if io.Disponible && len(io.ColaProcesos) > 0 {
		io.Disponible = false
		go UtilizarIO(io, io.ColaProcesos[0].Pcb, io.ColaProcesos[0].Tiempo)

	}
}

func DumpDelProceso(pcb *PCB, ip string, puerto int) {

	var paquete PaqueteEnviadoKERNELaMemoria2
	paquete.Pid = pcb.Pid
	paquete.Mensaje = fmt.Sprintf("El proceso PID: %d  requiere que se haga un DUMP del mismo", pcb.Pid)

	PasarBlocked(pcb)

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.")
		return
	}
	cliente := http.Client{}

	url := fmt.Sprintf("http://%s:%d/KERNELMEMORIADUMP", ip, puerto)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	defer respuestaJSON.Body.Close()
	//log.Printf("Conexion establecida con exito \n")

	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibidoDeMemoria
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		//log.Printf("Error al decodificar el JSON.\n")
		return
	}

	//log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)

	if respuestaJSON.StatusCode == http.StatusOK {
		log.Printf("Se pudo hacer el DUMP del proceso con el PID: %d ", pcb.Pid)
		PasarReady(pcb, ColaBlock)
	} else {
		log.Printf("No se pudo hacer el DUMP del proceso con el PID: %d ", pcb.Pid)
		FinalizarProceso(pcb)
	}

}

func SwapDelProceso(pcb *PCB, ip string, puerto int) {

	var paquete PaqueteEnviadoKERNELaMemoria2
	paquete.Pid = pcb.Pid
	paquete.Mensaje = fmt.Sprintf("El proceso PID: %d  requiere que se haga un SWAP del mismo", pcb.Pid)

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.")
		return
	}
	cliente := http.Client{}

	url := fmt.Sprintf("http://%s:%d/SWAPADISCO", ip, puerto)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson))

	if err != nil {
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	respuestaJSON, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	defer respuestaJSON.Body.Close()

	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta PaqueteRecibidoDeMemoria
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		//log.Printf("Error al decodificar el JSON.\n")
		return
	}

	//log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)

	if respuestaJSON.StatusCode != http.StatusOK {
		log.Printf("Error al hacer el SWAP")
	}

	SemLargoPlazo <- struct{}{}

}

func removerIO(io *IO) []IO {
	for i, item := range ListaIO {
		if item.Instancia == io.Instancia {
			return append(ListaIO[:i], ListaIO[1+i:]...)
		}
	}
	return ListaIO
}

func enviarExitProcesosIO(io *IO) {
	for _, proceso := range io.ColaProcesos {
		if proceso.Pcb != nil {
			//log.Printf("El proceso PID: %d  se pasa a EXIT por desconexion del I/O %s", proceso.Pcb.Pid, io.Instancia)
			FinalizarProceso(proceso.Pcb)
		}
	}
	removerIO(io)
	log.Printf("Se desconecto el I/O %s", io.Instancia)
}

func InicializarSemaforos() {
	SemLargoPlazo = make(chan struct{}, 100)
	SemCortoPlazo = make(chan struct{}, 100)
}

/*
NUEVA CONEXION FACU
EN TEORIA, CUANDO QUERIAN DESWAPPEAR UN PROCESO, LO QUE HACIAN ERA CREAR UNO NUEVO CON ESE MISMO PID
AHORA NO
ARMO HTTP PARA DESWAPEAR
*/

// ATENCION, SI LE VAN A CAMBIAR EL NOMBRE, TIENEN QUE IR A CAMBIARLO TAMBIEN EN SU INVOCACION (PLANI LARGO PLAZO)
// <<< NUEVO >>> Trae un proceso swappeado devuelta a memoria
func SwapInProceso(pcb *PCB) {
	var paquete PaqueteEnviadoKERNELaMemoria2
	paquete.Pid = pcb.Pid
	paquete.Mensaje = "SWAP_IN" // o el texto que tu memoria espere

	data, _ := json.Marshal(paquete)
	cliente := http.Client{}
	url := fmt.Sprintf("http://%s:%d/SWAPAMEMORIA", globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error generando petición SWAP_IN: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cliente.Do(req)
	if err != nil {
		log.Printf("Error en SWAP_IN del PID %d: %v", pcb.Pid, err)
		return
	}
	defer resp.Body.Close()
}
