package utilsIO

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Handshakepaquete struct {
	Nombre string `json:"name"` //lo de name, ip, port es como va a construir el json la maquina cuando lo pasemos a json
	Ip     string `json:"ip"`   // es fundamental ponerlo
	Puerto int    `json:"port"`
}

type RespuestaHandshakeKernel struct { // aca va el formato que va a tener lo que esperas del server
	Mensaje string `json:"message"`
}

type PaqueteRecibidoIO struct {
	Mensaje string `json:"message"`
	Tiempo  int    `json:"tiempo"`
}

type PaqueteRespuestaKERNEL struct {
	Mensaje string `json:"message"`
}

/*
func ConfigurarLogger() {
	logFile, err := os.OpenFile("io.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
*/
/*
conexion entre IO (Client) con Kernel (Server)
enviamos handshake con datos del modulo y esperamos respuesta
*/

func ConfigurarLogger(ioId string) {
	logFileName := fmt.Sprintf("IO-%s.log", ioId)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	//prefija cada línea de log con el ioname:
	log.SetPrefix(fmt.Sprintf("[IO-%s] ", ioId))
}

func PeticionClienteIOServidorKERNEL(nombre string, ipKernel string, puertoKernel int, ipIO string, puertoIO int) {

	var paquete Handshakepaquete
	paquete.Nombre = nombre
	paquete.Ip = ipIO
	paquete.Puerto = puertoIO

	PaqueteFormatoJson, err := json.Marshal(paquete)
	if err != nil {
		//aca tiene que haber un logger
		log.Printf("error al convertir a json")
		return
	}
	cliente := http.Client{} //crea un "cliente"

	url := fmt.Sprintf("http://%s:%d/IO", ipKernel, puertoKernel) //url del server

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(PaqueteFormatoJson)) //genera peticion al server

	if err != nil {
		//aca tiene que haber un logger
		log.Printf("error al generar la peticion al server")
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

	//pasamos la respuesta de JSON a formato paquete que nos mando el server

	var respuesta RespuestaHandshakeKernel //para eso declaramos una variable con el struct que esperamos que nos envie el server
	err = json.Unmarshal(body, &respuesta) //pasamos de bytes al formato de nuestro paquete lo que nos mando el server
	if err != nil {
		log.Printf("Error al decodificar el JSON")
		return
	}
	log.Printf("La respuesta del server fue: %s\n", respuesta.Mensaje)

}

func RetornoClienteKERNELServidorIO(w http.ResponseWriter, r *http.Request) {

	var request PaqueteRecibidoIO
	log.Printf("llegue.\n")

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("I/O Iniciado.\n")
	IniciarSleep(request.Tiempo)

	//Leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("El kernel nos envio esto: %s\n", request.Mensaje)
	log.Printf("I/O Finalizado. \n")
	//Respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si

	var respuestaIO PaqueteRespuestaKERNEL
	respuestaIO.Mensaje = "I/O Finalizado"
	respuestaJSON, err := json.Marshal(respuestaIO)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func IniciarSleep(tiempo int) {
	time.Sleep(time.Duration(tiempo) * time.Second)

}
