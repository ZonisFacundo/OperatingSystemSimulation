package utilsMemoria

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
)

//"github.com/sisoputnfrba/tp-golang/memoria/auxiliares"

func InicializarSwap() {
	file, err := os.Create(fmt.Sprintf("%s/swapfile.bin", globals.ClientConfig.Swapfile_path))

	if err != nil {
		log.Printf("error al crear el archivo swap")
		return
	}

	defer file.Close()

}

func RetornoClienteKernelServidorMemoriaSwapADisco(w http.ResponseWriter, r *http.Request) {

	var paqueteDeKernel PaqueteRecibidoMemoriadeKernel2
	err := json.NewDecoder(r.Body).Decode(&paqueteDeKernel) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	retorno := SwapADisco(paqueteDeKernel.Pid)

	var respuesta respuestaalKernel

	if retorno == -1 {
		respuesta.Mensaje = "ERROR AL SWAPPEAR \n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusNotImplemented)
		w.Write(respuestaJSON)
		return
	} else {

		respuesta.Mensaje = "listo \n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respuestaJSON)
		return
	}
}
func RetornoClienteKernelServidorMemoriaSwapAMemoria(w http.ResponseWriter, r *http.Request) {

	var paqueteDeKernel PaqueteRecibidoMemoriadeKernel2
	err := json.NewDecoder(r.Body).Decode(&paqueteDeKernel) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	retorno := EntraEnMemoria(len(globals.MemoriaKernel[paqueteDeKernel.Pid].TablaSimple)) //se fija si entra en memoria o no

	var respuesta respuestaalKernel

	if retorno == -2 {

		respuesta.Mensaje = "NO HAY SUFICIENTE ESPACIO PARA EL PROCESO EN MP \n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusInsufficientStorage)
		w.Write(respuestaJSON)
		return

	} else {
		retorno = SwapAMemoria(paqueteDeKernel.Pid)
		if retorno == -1 {

			respuesta.Mensaje = "ERROR AL SWAPPEAR A MEMORIA (error al abrir archivo) \n"

			respuestaJSON, err := json.Marshal(respuesta)
			if err != nil {
				return
			}

			w.WriteHeader(http.StatusNotImplemented)
			w.Write(respuestaJSON)
			return
		}

		respuesta.Mensaje = "listo \n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respuestaJSON)
		return
	}
}

func SwapADisco(pid int) int { //incompleta

	file, err := os.OpenFile(fmt.Sprintf("%s/swapfile.bin", globals.ClientConfig.Swapfile_path), os.O_APPEND|os.O_RDWR, 0644)

	if err != nil {
		log.Printf("error al abrir el archivo SWAP a la hora de llevarlo a disco para pid: %d\n", pid)
		return -1
	}

	buffer := make([]byte, globals.ClientConfig.Page_size)

	for i := 0; i < len(globals.MemoriaKernel[pid].TablaSimple); i++ {
		for j := 0; j < globals.ClientConfig.Page_size; j++ {
			//buffer[j] = append(buffer, globals.MemoriaPrincipal[((globals.MemoriaKernel[pid].TablaSimple[i])*globals.ClientConfig.Page_size)+j])
			buffer[j] = globals.MemoriaPrincipal[((globals.MemoriaKernel[pid].TablaSimple[i])*globals.ClientConfig.Page_size)+j]
		}
		bytesEscritos, err := file.Write(buffer)
		if err != nil {
			log.Printf("error al escribir en el swap en la pagina numero: %d del proceso de pid: %d\n", globals.MemoriaKernel[pid].TablaSimple[i], pid)
		}
		log.Printf("%d fueron escritos en la ultima iteracion (a disco) \n", bytesEscritos)

	}

	//ya fueron escritas las paginas en el swap... ahora tenemos que implementar logica para guardar los bytes exactos donde esta cada pag
	// TO DO: IMPLEMENTAR DONDE SE VAN A GUARDAR LAS PAGINAS, SI HAY ESPACIO INFINITO O NO

	CambiarAMenos1TodasLasTablas(pid)
	return 0
}

/*
func SwapAMemoria(pid int) int {

	file, err := os.OpenFile(fmt.Sprintf("%s/swapfile.bin", globals.ClientConfig.Swapfile_path), os.O_APPEND|os.O_RDWR, 0644)

	if err != nil {
		log.Printf("error al abrir el archivo SWAP a la hora de llevarlo a disco para pid: %d\n", pid)
		return -1
	}
}
*/
