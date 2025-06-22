package utilsCPU

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func ConfigurarLogger(cpuId string) {
	logFileName := fmt.Sprintf("CPU-%s.log", cpuId)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	//prefija cada línea de log con el cpuId:
	log.SetPrefix(fmt.Sprintf("[CPU-%s] ", cpuId))
}

func RecibirPCyPID(w http.ResponseWriter, r *http.Request) {
	//var request HandshakeKERNEL
	var request Proceso

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("El Kernel envio PID: %d - PC: %d", request.Pid, request.Pc)

	var respuesta RespuestaKernel
	respuesta.Mensaje = "PC y PID recbidos correctamente"
	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)
}

func EnvioPortKernel(ip string, puerto int, instancia string, portcpu int, ipcpu string) {

	var paquete HandshakeCPU

	paquete.Ip = ipcpu
	paquete.Puerto = portcpu
	paquete.Instancia = instancia

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}

	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/handshake", ip, puerto) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

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

		log.Printf("Status de respuesta el server no fue la esperada.\n")
		return
	}
	defer respuestaJSON.Body.Close()

	log.Printf("Conexion establecida con exito.\n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta RespuestaalCPU
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
	}

	log.Printf("Conexión realizada con exito con el Kernel.")

}

func FinEjecucion(ip string, puerto int, pid int, pc int, instancia string, syscall string, parametro1 int, parametro2 string) { // si no reciben parametros que sean  0 y "" que nosostros ahi no los usamos
	var paquete PackageFinEjecucion

	paquete.Pid = pid
	paquete.Pc = pc
	paquete.Syscall = syscall
	paquete.InstanciaCPU = instancia
	paquete.Parametro1 = parametro1
	paquete.Parametro2 = parametro2

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}

	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/PCB", ip, puerto) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

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

		log.Printf("Status de respuesta el server no fue la esperada.\n")
		return
	}

	defer respuestaJSON.Body.Close()

	log.Printf("Conexion establecida con exito.\n")
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	var respuesta RespuestaKernel

	if respuesta.Mensaje == "interrupcion" {
		//interrumpir
		//finalizar ejecución de ahora y esperar el próximo PID y PC, como? no tengo la mas minima idea.
	}
	
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
	}

	log.Printf("El Kernel recibió correctamente el PID y el PC.\n")
}
