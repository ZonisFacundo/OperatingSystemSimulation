package utilsCPU

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
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

	globals.ID.Pid = request.Pid
	globals.ID.Pc = request.Pc

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

func DevolverPidYPCInterrupcion(w http.ResponseWriter, r *http.Request, pc int, pid int) {
	var request PaqueteInterrupcion

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("El kernel nos interrumpio.")

	var respuesta RespuestaKernel
	respuesta.Pc = pc
	respuesta.Pid = pid

	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)
}
