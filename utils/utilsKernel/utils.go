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

func RetornoClienteIOServidorKERNEL(w http.ResponseWriter, r *http.Request) {

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
	respuestaIO.Mensaje = "Se envio un string al Kernel."
	respuestaJSON, err := json.Marshal(respuestaIO)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}
func RetornoClienteCPUServidorKERNEL(w http.ResponseWriter, r *http.Request) {

	var request HandshakepaqueteCPU

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("Handshake recibido de la instancia: %s", request.Instancia)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuesta RespuestaalCPU
	respuesta.Mensaje = "conexion realizada con exito"
	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}

	crearStructCPU(request.Ip, request.Puerto, request.Instancia)

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func RetornoClienteCPUServidorKERNEL2(w http.ResponseWriter, r *http.Request) {

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
	respuesta.Mensaje = "conexion realizada con exito"
	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}

	cpuServidor := ObtenerCpu(request.Instancia)
	cpuServidor.Disponible = true

	if request.Contexto == "RUNNING" {
		//hacer algo a chequear
	} else {
		//cambiar estado de pcb

	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

} // conexion kernel --> IO lado del cliente (kernel)
func PeticionClienteKERNELServidorIO(ip string, puerto int) {

	var paquete RespuestaalIO
	paquete.Mensaje = "mensaje enviado a kernel desde io"

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

		log.Printf("Status de respuesta el server no fue la esperada.\n")
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
	log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)

}
func PeticionClienteKERNELServidorMemoria(pcb PCB, ip string, puerto int) {

	var paquete PaqueteEnviadoKERNELaMemoria
	paquete.Pid = pcb.Pid
	paquete.TamProceso = pcb.TamProceso

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
		println("Se pasa el proceso a READY")
		PasarReady(pcb)
	}

	//en mi caso era un mensaje, por eso el struct tiene mensaje string, vos por ahi estas esperando 14 ints, no necesariamente un struct

}

func PeticionClienteKERNELServidorCPU(pcb PCB, cpu CPU) {

	var paquete PaqueteEnviadoKERNELaCPU

	paquete.Pid = pcb.Pid
	paquete.PC = pcb.PC

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

func CrearPCB(tamanio int) { //pid unico arranca de 0
	ColaNew = append(ColaNew, PCB{
		Pid:            ContadorPCB,
		PC:             0,
		EstadoActual:   "NEW",
		TamProceso:     tamanio,
		MetricaEstados: make(map[Estado]int),
		TiempoEstados:  make(map[Estado]int64),
	})
	ContadorPCB++
	PlanificadorLargoPlazo()
}

func LeerConsola() string {
	// Leer de la consola
	reader := bufio.NewReader(os.Stdin)
	log.Println("Precione enter para inciar el planificador")
	text, _ := reader.ReadString('\n')
	//log.Print(text)
	return text
}

func IniciarPlanifcador() {
	for true {
		text := LeerConsola()
		for text == "\n" {

			println("Planificador de largo plazo ejecutando") //solo para saber que esta funcionando
			PlanificadorLargoPlazo()
			PlanificadorCortoPlazo()
			time.Sleep(5 * time.Second)
		}
	}
}

func PlanificadorLargoPlazo() {
	if len(ColaNew) != 0 {
		pcbChequear := CriterioColaNew()
		PeticionClienteKERNELServidorMemoria(pcbChequear, globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)
	} else {
		println("No hay procesos en cola NEW")
	}
}

func PlanificadorCortoPlazo() {
	if len(ColaReady) != 0 {
		pcbChequear := CriterioColaReady()
		CPUDisponible, noEsVacio := TraqueoCPU() //drakukeo en su defecto
		if noEsVacio {
			println("Proceso ejecutandose") //solo para saber que esta funcionando
			pcbChequear.EstadoActual = "EXEC"
			CPUDisponible.Disponible = false
			PeticionClienteKERNELServidorCPU(pcbChequear, CPUDisponible)

		}
	} else {
		println("No hay procesos en cola READY")
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
