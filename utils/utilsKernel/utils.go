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
	log.Printf("El cliente nos mando esto: \n nombre: %s  \n puerto: %d \n IP: %s \n", request.Nombre, request.Puerto, request.Ip)

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
	log.Printf("contexto de devolucion del proceso: %s", request.Contexto)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuesta RespuestaalCPU
	respuesta.Mensaje = "Conexion realizada con exito"
	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}
	log.Printf("Conexion establecida con exito \n")
	cpuServidor := ObtenerCpu(request.Instancia)
	cpuServidor.Disponible = true

	if request.Contexto == "RUNNING" {
		//hacer algo a chequear
	} else {
		//cambiar estado de pcb

	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}
func UtilizarIO(ip string, puerto int, pid int, tiempo int, nombre string) {

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

	url := fmt.Sprintf("http://%s:%d/KERNELIO", ip, puerto) //url del server

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

		log.Printf("Status de respuesta del I/0 %s no fue la esperada.\n", nombre)
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
	log.Printf("La respuesta del I/O %s fue: %s\n", nombre, respuesta.Mensaje)

}
func ConsultarProcesoConMemoria(pcb PCB, ip string, puerto int) {

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

func EnviarProcesoACPU(pcb PCB, cpu CPU) {

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

func InformarMemoriaFinProceso(pcb PCB, ip string, puerto int) {

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

}

func CrearPCB(tamanio int, archivo string) { //pid unico arranca de 0
	ColaNew = append(ColaNew, PCB{
		Pid:            ContadorPCB,
		Pc:             0,
		EstadoActual:   "NEW",
		TamProceso:     tamanio,
		MetricaEstados: make(map[Estado]int),
		TiempoEstados:  make(map[Estado]int64),
		Archivo:        archivo,
	})
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
		log.Printf("Planificador de largo plazo ejecutando") //solo para saber que esta funcionando
		CrearPCB(tamanio, archivo)
		for text == "\n" {
			//log.Print("Planificador de largo plazo ejecutando, pero dentro de un for")
			//PlanificadorLargoPlazo()
			PlanificadorCortoPlazo()
			time.Sleep(5 * time.Second)
		}
	}
}

func PlanificadorLargoPlazo() {
	if len(ColaSuspReady) != 0 {
		pcbChequear := CriterioColaSuspReady()
		ConsultarProcesoConMemoria(pcbChequear, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	} else if len(ColaNew) != 0 {
		pcbChequear := CriterioColaNew()
		ConsultarProcesoConMemoria(pcbChequear, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	}
}

func PlanificadorCortoPlazo() {
	if len(ColaReady) != 0 {
		pcbChequear := CriterioColaReady()
		CPUDisponible, noEsVacio := TraqueoCPU() //drakukeo en su defecto
		if noEsVacio {
			log.Printf("se pasa el proceso PID: %d a EXECUTE", pcbChequear.Pid) //solo para saber que esta funcionando
			PasarExec(pcbChequear)
			CPUDisponible.Disponible = false
			EnviarProcesoACPU(pcbChequear, CPUDisponible)

		}
	}
}

func FIFO(cola []PCB) PCB {
	if len(cola) == 0 {
		return PCB{}
	}
	pcb := cola[0]
	return pcb
}

func PasarReady(pcb PCB) {
	ColaReady = append(ColaReady, pcb)
	ColaNew = removerPCB(ColaNew, pcb)
	pcb.EstadoActual = "READY"
}

func PasarExec(pcb PCB) {
	ColaReady = removerPCB(ColaReady, pcb)
	pcb.EstadoActual = "EXECUTE"
}

func removerPCB(cola []PCB, pcb PCB) []PCB {
	for i, item := range cola {
		if item.Pid == pcb.Pid {
			return append(cola[:i], cola[1+i:]...)
		}
	}
	return cola
}

func CriterioColaNew() PCB {
	if globals.ClientConfig.Ready_ingress_algorithm == "FIFO" {
		return FIFO(ColaNew)
	}
	return FIFO(ColaNew) //esto no va asi pero es para que no de error
}

func CriterioColaSuspReady() PCB {
	if globals.ClientConfig.Ready_ingress_algorithm == "FIFO" {
		return FIFO(ColaSuspReady)
	}
	return FIFO(ColaSuspReady) //esto no va asi pero es para que no de error
}

func CriterioColaReady() PCB {
	if globals.ClientConfig.Scheduler_algorithm == "FIFO" {
		return FIFO(ColaNew)
	}
	return FIFO(ColaNew) //esto no va asi pero es para que no de error
}

func TraqueoCPU() (CPU, bool) {
	for _, CPU := range ListaCPU {
		if CPU.Disponible {
			return CPU, true
		}
	}
	return CPU{}, false
}

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
}

func FinalizarProceso(pcb PCB) {
	log.Printf("El proceso PID: %d termino su ejecucion y se paso a EXIT", pcb.Pid)
	pcb.EstadoActual = "EXIT"
	ColaExit = append(ColaExit, pcb) //es un esquema de como podria finalizar el proceso, puede cambiarse esto
	InformarMemoriaFinProceso(pcb, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	PlanificadorLargoPlazo() // esto seria porque se libera el espacio de memoria y capaz se podria ejecutar otro proceso
}
