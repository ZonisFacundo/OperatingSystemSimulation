package utilsMemoria

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/auxiliares"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
)

//"github.com/sisoputnfrba/tp-golang/memoria/auxiliares"

/*
que hace inicializarSwap?
crea un archivo con la ruta indicada en el config
*/

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
		respuesta.Mensaje = "ERROR AL SWAPPEAR (RetornoClienteKernelServidorMemoriaSwapADisco)\n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusNotImplemented)
		w.Write(respuestaJSON)
		return
	} else {

		respuesta.Mensaje = "listo (RetornoClienteKernelServidorMemoriaSwapADisco) \n"

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

		respuesta.Mensaje = "NO HAY SUFICIENTE ESPACIO PARA EL PROCESO EN MP (RetornoClienteKernelServidorMemoriaSwapAMemoria) \n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusInsufficientStorage)
		w.Write(respuestaJSON)
		return

	} else {
		//retorno = SwapAMemoria(paqueteDeKernel.Pid) comento esto para que no de error
		if retorno == -1 {

			respuesta.Mensaje = "ERROR AL SWAPPEAR A MEMORIA (error al abrir archivo) (RetornoClienteKernelServidorMemoriaSwapAMemoria) \n"

			respuestaJSON, err := json.Marshal(respuesta)
			if err != nil {
				return
			}

			w.WriteHeader(http.StatusNotImplemented)
			w.Write(respuestaJSON)
			return
		}

		respuesta.Mensaje = "listo (RetornoClienteKernelServidorMemoriaSwapAMemoria) \n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respuestaJSON)
		return
	}
}

/*
Que hace swapadisco?

copia el contenido en memoria del proceso y lo pega en un arhcivo, ademas, en MemoriaKernel, vamos a almacenar la posicion donde fue almacenado este
y el tamano hasta donde fue almacenado para luego poder recuperarlo en otra funcion
tambien se cambia a -1 todas las reservas de paginas que tenia asociadas en las distintas estructuras de memoria para liberar la memoria para otro proceso
*/

func SwapADisco(pid int) int { //incompleta

	var bytesEscritos int = 0
	var bytesEscritosRecien int = 0

	file, err := os.OpenFile(fmt.Sprintf("%s/swapfile.bin", globals.ClientConfig.Swapfile_path), os.O_APPEND|os.O_RDWR, 0644)

	if err != nil {
		log.Printf("error al abir el archivo (SwapADisco)\n")
	}

	currentPos, err := file.Seek(0, io.SeekCurrent) //guardamos el numero de byte en el que arrancamos a escribir

	if err != nil {
		log.Printf("error al abrir el archivo SWAP a la hora de llevarlo a disco para pid: %d \t (SwapADisco)\n", pid)
		return -1
	}

	buffer := make([]byte, globals.ClientConfig.Page_size)

	for i := 0; i < len(globals.MemoriaKernel[pid].TablaSimple); i++ {
		for j := 0; j < globals.ClientConfig.Page_size; j++ {
			//buffer[j] = append(buffer, globals.MemoriaPrincipal[((globals.MemoriaKernel[pid].TablaSimple[i])*globals.ClientConfig.Page_size)+j])
			buffer[j] = globals.MemoriaPrincipal[((globals.MemoriaKernel[pid].TablaSimple[i])*globals.ClientConfig.Page_size)+j]
		}
		bytesEscritosRecien, _ = file.Write(buffer)

		bytesEscritos += bytesEscritosRecien
	}
	log.Printf("%d fueron escritos en disco \t(SwapADisco)\n", bytesEscritos)
	//ya fueron escritas las paginas en el swap... ahora tenemos que implementar logica para guardar los bytes exactos donde esta cada pag
	// TO DO: IMPLEMENTAR DONDE SE VAN A GUARDAR LAS PAGINAS, SI HAY ESPACIO INFINITO O NO

	auxiliares.ActualizarSwapInfo(currentPos, bytesEscritos, pid)

	CambiarAMenos1TodasLasTablas(pid)
	defer file.Close()
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
