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

func GOTO(pcInstr int, valor int) int {
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
		log.Printf("entro aca y rompo")
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
