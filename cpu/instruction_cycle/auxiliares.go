package instruction_cycle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

/*func NOOP(Tiempo int) int {
	return Tiempo
}*/

type CPUMMU struct {
	Entradas int `json:"ent"`
	Niveles  int `json:"niv"`
	TamPag   int `json:"tam"`
}

func GOTO(pcInstr int, valor int) int {
	pcInstr--
	return pcInstr + valor
}

func EnvioDirLogica(ip string, puerto int, dirLogica []int) {

	var paquete utilsCPU.EnvioDirLogicaAMemoria

	paquete.Ip = ip
	paquete.Puerto = puerto
	paquete.DirLogica = dirLogica

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}

	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/TRADUCCIONLOGICAAFISICA", ip, puerto) //url del server

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

	var frame utilsCPU.MarcoDeMemoria
	err = json.Unmarshal(body, &frame)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
	}

	log.Printf("Recibido de memoria el frame: %d", frame.Frame)

	globals.ID.Frame = frame.Frame

}

func RecibirPCyPID(w http.ResponseWriter, r *http.Request) {
	//var request HandshakeKERNEL
	var request utilsCPU.Proceso

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("El Kernel envio PID: %d - PC: %d", request.Pid, request.Pc)

	globals.ID.ProcessValues.Pid = request.Pid
	globals.ID.ProcessValues.Pc = request.Pc

	var respuesta utilsCPU.RespuestaKernel
	respuesta.Mensaje = "PC y PID recbidos correctamente"
	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)
}

func RecibirDatosMMU(ip string, puerto int) {
	var paquete utilsCPU.ReciboMMU

	paquete.Ip = ip
	paquete.Puerto = puerto
	paquete.Mensaje = "mensa"

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}

	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/HANDSHAKE", ip, puerto) //url del server

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

	var respuesta CPUMMU
	err = json.Unmarshal(body, &respuesta)
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
	}

	globals.ClientConfig.Entradas = respuesta.Entradas
	globals.ClientConfig.Page_size = respuesta.TamPag
	globals.ClientConfig.Niveles = respuesta.Niveles
	log.Printf("imprimo el tam de pagina: %d", globals.ClientConfig.Page_size)

	log.Printf("Conexión realizada con exito con el Kernel.")

}
