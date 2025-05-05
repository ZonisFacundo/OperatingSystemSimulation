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

/*
func ConfigurarLogger() {
	logFile, err := os.OpenFile("cpu.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
*/

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
	var request HandshakeKERNEL

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var respuesta Proceso

	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}

	log.Printf("El Kernel envio PID: %d - PC: %d", respuesta.Pid, respuesta.Pc)

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)
}

func EnvioPortKernel(ip string, puerto int, instancia string, portcpu int) {

	var paquete HandshakeCPU

	paquete.Ip = ip
	paquete.Puerto = puerto
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

func FinEjecucion(ip string, puerto int, pid int, pc int, instancia string, contexto string) {
	var paquete PackageFinEjecucion

	paquete.Pid = pid
	paquete.Pc = pc
	paquete.Contexto = contexto
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

	var respuesta RespuestaKernel

	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
	}

	log.Printf("El Kernel recibió correctamente el PID y el PC.\n")
}

func EnvioDirLogica(ip string, puerto int, dirLogica []int) {

	var paquete EnvioDirLogicaAMemoria

	// paquete.Ip = ip
	// paquete.Puerto = puerto
	paquete.DirLogica = dirLogica
	

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

	log.Printf("Enviado todo a Memoria.")

}