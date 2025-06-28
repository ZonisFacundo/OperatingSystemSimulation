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
	PCBUtilizar := ObtenerPCB(cpuServidor.Pid) // ya no hace falta porque esta en el struct
	PCBUtilizar.Pc = request.Pc
	PCBUtilizar.RafagaAnterior = float32(PCBUtilizar.TiempoEnvioExc.Sub(time.Now()))

	switch request.Syscall {
	case "IO":
		log.Printf("## (<%d>) - Solicitó syscall: <IO> \n", PCBUtilizar.Pid)
		respuesta.Mensaje = "interrupcion"
		cpuServidor.Disponible = true
		if ExisteIO(request.Parametro2) {
			SemCortoPlazo <- struct{}{}
			ioServidor := ObtenerIO(request.Parametro2)
			PCBUtilizar.TiempoEnvioBlock = time.Now()
			go PlanificadorMedianoPlazo(PCBUtilizar)
			AgregarColaIO(ioServidor, PCBUtilizar, request.Parametro1)
			PasarBlocked(PCBUtilizar)
			log.Printf("## (<%d>) - Bloqueado por IO: < %s > \n", PCBUtilizar.Pid, ioServidor.Instancia)
			MandarProcesoAIO(ioServidor)
			if len(ioServidor.ColaProcesos) > 0 {
				ioServidor.ColaProcesos = ioServidor.ColaProcesos[1:]
			}
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
		respuesta.Mensaje = "NO INTERRUMPAS GIL"
		log.Printf("## (<%d>) - Solicitó syscall: <INIT_PROC> \n", PCBUtilizar.Pid)
		CrearPCB(request.Parametro1, request.Parametro2)
		cpuServidor.Disponible = false
		EnviarProcesoACPU(PCBUtilizar, &cpuServidor)
		w.WriteHeader(http.StatusOK)

		w.Write(respuestaJSON)
		return
	}
	log.Printf("PID: %d PC: %d", request.Pid, request.Pc)
	w.WriteHeader(http.StatusFound)
	w.Write(respuestaJSON)

}

func UtilizarIO(ioServer IO, pcb *PCB, tiempo int) {

	var paquete PaqueteEnviadoKERNELaIO
	paquete.Pid = pcb.Pid
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

	if respuestaJSON.StatusCode == http.StatusOK {

		log.Printf("La respuesta del I/O %s fue: %s\n", ioServer.Instancia, respuesta.Mensaje)

		if EstaEnColaBlock(pcb) {
			PasarReady(pcb)
		} else {
			PasarSuspReady(pcb)
		}

		MandarProcesoAIO(ioServer)
	}
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

	paquete.Mensaje = "Interrupcion del proceso"

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		//aca tiene que haber un logger
		log.Printf("Error al convertir a json.")
		return
	}

	cliente := http.Client{} // Crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/INTERRUPCIONCPU", cpu.Ip, cpu.Port) //url del server

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

	var respuesta PaqueteRecibidoDeCPU
	pcb := ObtenerPCB(respuesta.Pid)
	pcb.Pc = respuesta.Pc
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}
	//log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)
	log.Printf("PID: %d PC: %d", respuesta.Pid, respuesta.Pc)
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
	SemLargoPlazo <- struct{}{}
	SemCortoPlazo <- struct{}{}

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
	for true {
		<-SemLargoPlazo //wait()
		if len(ColaSuspReady) != 0 {
			MutexColaNew.Lock()
			pcbChequear := CriterioColaNew(ColaSuspReady)
			MutexColaNew.Unlock()
			ConsultarProcesoConMemoria(pcbChequear, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)

		} else if len(ColaNew) != 0 {
			MutexColaNew.Lock()
			pcbChequear := CriterioColaNew(ColaNew)
			MutexColaNew.Unlock()
			ConsultarProcesoConMemoria(pcbChequear, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)

		} else {
			SemLargoPlazo <- struct{}{} //signal()
			time.Sleep(1 * time.Second)

		}
	}
}

func PlanificadorCortoPlazo() {
	for true {
		<-SemCortoPlazo

		if len(ColaReady) != 0 {
			//log.Printf("Planificador de corto plazo ejecutando")
			//time.Sleep(10 * time.Second)
			MutexColaReady.Lock()
			log.Printf("entreg al mutex")
			pcbChequear, hayDesalojo := CriterioColaReady()
			MutexColaReady.Unlock()
			log.Printf("sali del mutex")
			CPUDisponible, noEsVacio := TraqueoCPU()
			if noEsVacio {
				log.Printf("se pasa el proceso PID: %d a EXECUTE", pcbChequear.Pid) //solo para saber que esta funcionando
				PasarExec(pcbChequear)
				CPUDisponible.Disponible = false
				CPUDisponible.Pid = pcbChequear.Pid //le asigno el pid al cpu que lo va a ejecutar
				EnviarProcesoACPU(pcbChequear, CPUDisponible)

			} else if hayDesalojo {
				pcbDesalojar, cpuDesalojar := RafagaMasLargaDeLosCPU()
				if calcularRafagaEstimada(pcbChequear) < CalcularTiempoRestanteEjecucion(pcbDesalojar) {
					InterrumpirCPU(cpuDesalojar)
					PasarReady(pcbDesalojar)
					PasarExec(pcbChequear)
					cpuDesalojar.Pid = pcbChequear.Pid
				}
			} else {
				SemCortoPlazo <- struct{}{}
				time.Sleep(1 * time.Second)
			}
		} else {
			SemCortoPlazo <- struct{}{}
			time.Sleep(1 * time.Second)

		}
	}
}

func PlanificadorMedianoPlazo(pcb *PCB) {
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
	MutexColaReady.Lock()
	ColaReady = append(ColaReady, pcb)
	MutexColaReady.Unlock()
	MutexColaNew.Lock()
	ColaNew = removerPCB(ColaNew, pcb)
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

func PasarSuspBlock(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado SUSP.BLOCKED \n", pcb.Pid, pcb.EstadoActual)
	//HacerSwap(pcb, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	MutexColaSuspBlock.Lock()
	ColaBlock = append(ColaSuspBlock, pcb)
	MutexColaSuspBlock.Unlock()
	MutexColaBlock.Lock()
	ColaBlock = removerPCB(ColaBlock, pcb)
	MutexColaBlock.Unlock()
	pcb.TiempoEstados[pcb.EstadoActual] = +time.Since(pcb.TiempoLlegada[pcb.EstadoActual]).Milliseconds()
	pcb.EstadoActual = "SUSP.BLOCKED"
	pcb.MetricaEstados["SUSP.BLOCKED"]++
	pcb.TiempoLlegada["SUSP.BLOCKED"] = time.Now()
}

func PasarSuspReady(pcb *PCB) {
	log.Printf("## (<%d>) Pasa del estado %s al estado SUSP.READY \n", pcb.Pid, pcb.EstadoActual)
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

func CrearStructCPU2(ip string, puerto int, instancia string) CPU {
	return (CPU{
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

func AgregarColaIO(io IO, pcb *PCB, tiempo int) {
	io.ColaProcesos = append(io.ColaProcesos, PCBIO{
		Pcb:    pcb,
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

func EstaEnColaBlock(pcbChequear *PCB) bool {
	for _, pcb := range ColaBlock {
		if pcb == pcbChequear {
			return true
		}
	}
	return false
}

func MandarProcesoAIO(io IO) {
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
		if proceso.Pcb != nil {
			log.Printf("El proceso PID: %d  se pasa a EXIT por desconexion del I/O %s", proceso.Pcb.Pid, io.Instancia)
			FinalizarProceso(proceso.Pcb)
		}
	}
	removerIO(&io)
	log.Printf("Se desconecto el I/O %s", io.Instancia)
}

func InicializarSemaforos() {
	SemLargoPlazo = make(chan struct{}, 100)
	SemCortoPlazo = make(chan struct{}, 100)
}
