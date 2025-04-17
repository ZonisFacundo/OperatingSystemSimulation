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

type HandshakepaqueteCPU struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"port"`
}

type HandshakepaqueteMemoria struct {
	Instruccion string `json:"instruccion"`
	Ip          string `json:"ip"` // es fundamental ponerlo
	Puerto      int    `json:"port"`
}

type RespuestaHandshakeKernel struct { // aca va el formato que va a tener lo que esperas del server
	Mensaje string `json:"message"`
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("cpu.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func PeticionCLienteCPUServidorMEMORIA(instruccion string, ip string, puerto int) {

	var paquete HandshakepaqueteMemoria

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("Error al convertir a json.\n")
		return
	}
	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/CPUMEMORIA", globals.ClientConfig.Ip_memory, globals.ClientConfig.Port_memory)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {
		log.Printf("Error al generar la peticion al server.\n")
		return
	}

	req.Header.Set("Content-Type", "application/json") //le avisa al server que manda la data en json format

	respuestaJSON, err := cliente.Do(req) //recibe la respuesta del server
	/* que tipo de dato tiene respuestaJSON?
		type respuestaJSON struct {
	    Status     string
	    StatusCode int
	    Header     Header
	    Body       io.ReadCloser  // ← This is what you're accessing
	    // ... other fields ...
	*/

	//ya definido por go de esa forma

	if err != nil {
		log.Printf("Error al recibir respuesta.\n")
		return

	}

	if respuestaJSON.StatusCode != http.StatusOK {

		log.Printf("Status de respuesta el server no fue la esperada.\n")
		return
	}
	defer respuestaJSON.Body.Close()

	fmt.Printf("Conexion establecida con exito.\n")
	//pasamos de JSON a formato bytes lo que nos paso el paquete
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	//pasamos la respuesta de JSON a formato paquete que nos mando el server

	var respuesta RespuestaHandshakeKernel //para eso declaramos una variable con el struct que esperamos que nos envie el server
	err = json.Unmarshal(body, &respuesta) //pasamos de bytes al formato de nuestro paquete lo que nos mando el server
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
		return
	}
	log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)
	//en mi caso era un mensaje, por eso el struct tiene mensaje string, vos por ahi estas esperando 14 ints, no necesariamente un struct

}

func PeticionClienteCPUServidorKERNEL(ip string, puerto int) {

	var paquete HandshakepaqueteCPU

	paquete.Ip = ip
	paquete.Puerto = puerto

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

	req.Header.Set("Content-Type", "application/json") //le avisa al server que manda la data en json format

	respuestaJSON, err := cliente.Do(req) //recibe la respuesta del server
	/* que tipo de dato tiene respuestaJSON?
		type respuestaJSON struct {
	    Status     string
	    StatusCode int
	    Header     Header
	    Body       io.ReadCloser  // ← This is what you're accessing
	    // ... other fields ...
	*/

	//ya definido por go de esa forma

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
	//pasamos de JSON a formato bytes lo que nos paso el paquete
	body, err := io.ReadAll(respuestaJSON.Body)

	if err != nil {
		return
	}

	//pasamos la respuesta de JSON a formato paquete que nos mando el server

	var respuesta Instruccion //para eso declaramos una variable con el struct que esperamos que nos envie el server

	err = json.Unmarshal(body, &respuesta) //pasamos de bytes al formato de nuestro paquete lo que nos mando el server
	if err != nil {
		log.Printf("Error al decodificar el JSON.\n")
	}

	log.Printf("Recibido del Kernel el PID: %d y el PC: %d.\n", respuesta.Pid, respuesta.Pc) //en mi caso era un mensaje, por eso el struct tiene mensaje string, vos por ahi estas esperando 14 ints, no necesariamente un struct

}
