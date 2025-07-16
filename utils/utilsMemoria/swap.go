package utilsMemoria

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/sisoputnfrba/tp-golang/memoria/auxiliares"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
)

//"github.com/sisoputnfrba/tp-golang/memoria/auxiliares"

/*
que hace inicializarSwap?
crea un archivo con la ruta indicada en el config
*/

func InicializarSwap() {
	file, err := os.Create(globals.ClientConfig.Swapfile_path)

	if err != nil {
		log.Printf("error al crear el archivo swap (InicializarSwap)\n")
		return
	} else {
		log.Printf("Swap creado (InicializarSwap) \n")
	}

	defer file.Close()

}

func RetornoClienteKernelServidorMemoriaSwapADisco(w http.ResponseWriter, r *http.Request) {

	time.Sleep(time.Duration(globals.ClientConfig.Swap_delay) * (time.Millisecond))

	log.Printf("\n\n KERNEL SOLICITO SWAP DE MEMORIA A DISCO (suspension) (RetornoClienteKernelServidorMemoriaSwapADisco)\n\n")

	globals.Sem_Swap.Lock()

	var paqueteDeKernel PaqueteRecibidoMemoriadeKernel2
	err := json.NewDecoder(r.Body).Decode(&paqueteDeKernel) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		globals.Sem_Swap.Unlock()

		return
	}

	retorno := SwapADisco(paqueteDeKernel.Pid)

	var respuesta respuestaalKernel

	if retorno == -1 {
		respuesta.Mensaje = "ERROR AL SWAPPEAR (RetornoClienteKernelServidorMemoriaSwapADisco)\n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {
			globals.Sem_Swap.Unlock()

			return
		}

		w.WriteHeader(http.StatusNotImplemented)
		w.Write(respuestaJSON)
		globals.Sem_Swap.Unlock()

		return
	} else {

		respuesta.Mensaje = "Swappeado a disco (RetornoClienteKernelServidorMemoriaSwapADisco) \n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {
			globals.Sem_Swap.Unlock()

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respuestaJSON)
		globals.Sem_Swap.Unlock()

		return
	}

}
func RetornoClienteKernelServidorMemoriaSwapAMemoria(w http.ResponseWriter, r *http.Request) {

	time.Sleep(time.Duration(globals.ClientConfig.Swap_delay) * (time.Millisecond))
	log.Printf("\n\n KERNEL SOLICITO SWAP DE DISCO A MEMORIA (des - suspension) (RetornoClienteKernelServidorMemoriaSwapAMemoria)\n\n")

	var paqueteDeKernel PaqueteRecibidoMemoriadeKernel2
	err := json.NewDecoder(r.Body).Decode(&paqueteDeKernel) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	log.Printf(" el swap tam de este proceso es: %d", globals.MemoriaKernel[paqueteDeKernel.Pid].SwapTam)

	globals.Sem_Swap.Lock()
	globals.Sem_Bitmap.Lock()

	retorno := EntraEnMemoria(globals.MemoriaKernel[paqueteDeKernel.Pid].SwapTam) //se fija si entra en memoria o no

	var respuesta respuestaalKernel

	if retorno == -2 {

		respuesta.Mensaje = "NO HAY SUFICIENTE ESPACIO PARA EL PROCESO EN MP (RetornoClienteKernelServidorMemoriaSwapAMemoria) \n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {

			globals.Sem_Bitmap.Unlock()
			globals.Sem_Swap.Unlock() //ok

			return
		}

		w.WriteHeader(http.StatusInsufficientStorage)
		w.Write(respuestaJSON)

		globals.Sem_Bitmap.Unlock() //ok
		globals.Sem_Swap.Unlock()

		return

	} else {
		retorno = SwapAMemoria(paqueteDeKernel.Pid)
		if retorno == -1 {

			respuesta.Mensaje = "ERROR AL SWAPPEAR A MEMORIA (error al abrir archivo) (RetornoClienteKernelServidorMemoriaSwapAMemoria) \n"

			respuestaJSON, err := json.Marshal(respuesta)
			if err != nil {
				globals.Sem_Bitmap.Unlock() //ok
				globals.Sem_Swap.Unlock()

				return
			}

			w.WriteHeader(http.StatusNotImplemented)
			w.Write(respuestaJSON)

			globals.Sem_Bitmap.Unlock()
			globals.Sem_Swap.Unlock()
			return
		}

		respuesta.Mensaje = "Swappeado a Memoria (RetornoClienteKernelServidorMemoriaSwapAMemoria) \n"

		respuestaJSON, err := json.Marshal(respuesta)
		if err != nil {

			globals.Sem_Bitmap.Unlock()
			globals.Sem_Swap.Unlock()

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respuestaJSON)

		var pagi globals.Pagina
		pagi, _ = LeerPaginaCompleta(0 * globals.ClientConfig.Page_size) //ver semaforos adentro
		log.Printf("LEO PAGINA 0: \n\n\n%b\n\n\n", pagi.Info)
		auxiliares.Mostrarmemoria()

		globals.Sem_Bitmap.Unlock()
		globals.Sem_Swap.Unlock()
		return
	}

	// Todo bien con esta funcion xd
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

	file, err := os.OpenFile(globals.ClientConfig.Swapfile_path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	log.Printf("path: %s\n", globals.ClientConfig.Swapfile_path)
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
			//log.Printf("valor al que accedo en memoria: %d", ((globals.MemoriaKernel[pid].TablaSimple[i])*globals.ClientConfig.Page_size)+j)
			buffer[j] = globals.MemoriaPrincipal[((globals.MemoriaKernel[pid].TablaSimple[i])*globals.ClientConfig.Page_size)+j]
		}
		bytesEscritosRecien, _ = file.Write(buffer)

		bytesEscritos += bytesEscritosRecien
	}
	log.Printf("%d fueron escritos en disco \t(SwapADisco)\n", bytesEscritos)
	//ya fueron escritas las paginas en el swap... ahora tenemos que implementar logica para guardar los bytes exactos donde esta cada pag
	// TO DO: IMPLEMENTAR DONDE SE VAN A GUARDAR LAS PAGINAS, SI HAY ESPACIO INFINITO O NO

	auxiliares.ActualizarSwapInfo(currentPos, bytesEscritos, pid)
	log.Printf("SWAP FINALIZADO, MUESTRO DATOS DEL SWAP DEL PROCESO \n PID: %d \n Posicion en disco: %d \n tam en swap: %d\n", pid, globals.MemoriaKernel[pid].SwapStart, globals.MemoriaKernel[pid].SwapTam)

	CambiarAMenos1TodasLasTablas(pid)
	defer file.Close()

	auxiliares.InicializarSiNoLoEstaMap(globals.Instruction.Pid)
	globals.Sem_Instruccion.Lock()
	globals.Sem_Metricas.Lock()
	globals.MetricasProceso[globals.Instruction.Pid].ContadorBajadasSWAP++
	globals.Sem_Metricas.Unlock()
	globals.Sem_Instruccion.Unlock()

	log.Printf("\n\tmuestro el contenido del SWAP\t\n")

	auxiliares.MostrarArchivo(globals.ClientConfig.Swapfile_path)

	/* DEBUG 13-07
	for i := 0; i < globals.ClientConfig.Memory_size; i++ {
		globals.MemoriaPrincipal[i] = 0
	}
	*/
	return 0

	//10 bananas sobre 10 platanos, o sea todo ok
}

//guardo esto por las dudas
//	file, err := os.OpenFile(fmt.Sprintf("%s/swapfile.bin", globals.ClientConfig.Swapfile_path), os.O_APPEND|os.O_RDWR, 0644) //no deberia |os.O_CREATE nunca

func SwapAMemoria(pid int) int {
	var bytestotales int = 0

	log.Printf("ENTRE A SWAPpppp ")
	file, err := os.OpenFile(globals.ClientConfig.Swapfile_path, os.O_RDWR, 0644) //no deberia |os.O_CREATE nunca

	if err != nil {
		log.Printf("error al abrir el archivo SWAP a la hora de llevarlo a disco para pid: %d    (SwapAMemoria)\n", pid)
		return -1
	}
	if EntraEnMemoria(globals.MemoriaKernel[pid].SwapTam) < 0 {

		log.Printf("No hay espacio para cargar el proceso en MP, pid: %d  tam proceso: %d bytes (SwapAMemoria)\n", pid, globals.MemoriaKernel[pid].SwapTam)
		return -1
	}

	//ya esta protegido de race conditions esta funcion
	ReservarMemoriaSwapeado(pid, globals.MemoriaKernel[pid].SwapTam)

	ActualizarTodasLasTablasEnBaseATablaSimple(pid) //actualiza tabla de paginas multinivel y paginas disponibles

	pos, _ := file.Seek(globals.MemoriaKernel[pid].SwapStart, io.SeekStart)
	log.Printf("posicion donde estamos parados, listos para copiar de disco: %d, posicion donde arranca el proceso %d en disco: %d", pos, pid, globals.MemoriaKernel[pid].SwapStart)
	//ADMINISTRATIVAMENTE, TODO TERMINADO, AHORA NOS QUEDA ACTUALIZAR LA IMAGEN EN MP

	var bufferPagina []byte = make([]byte, globals.ClientConfig.Page_size)

	for i := 0; i < (globals.MemoriaKernel[pid].SwapTam / globals.ClientConfig.Page_size); i++ {

		//bufferPagina, _ = os.ReadFile(fmt.Sprintf("%s/swapfile.bin", globals.ClientConfig.Swapfile_path))

		bytesleidos, _ := file.Read(bufferPagina)

		for j := 0; j < bytesleidos; j++ {
			globals.MemoriaPrincipal[globals.MemoriaKernel[pid].TablaSimple[i]*globals.ClientConfig.Page_size+j] = bufferPagina[j]
		}
		bytestotales += bytesleidos
	}
	defer file.Close()
	log.Printf("bytes leidos de disco y copiados en memoria: %d 	(SwapAMemoria)\n", bytestotales)

	auxiliares.InicializarSiNoLoEstaMap(globals.Instruction.Pid)

	globals.Sem_Instruccion.Lock()
	globals.Sem_Metricas.Lock()
	globals.MetricasProceso[globals.Instruction.Pid].ContadorSubidasAMemoria++
	globals.Sem_Metricas.Unlock()
	globals.Sem_Instruccion.Unlock()

	//debug 13-07
	log.Printf("\n\tmuestro el contenido del SWAP PARA DEBUGEAR MAS QUE NADA\t\n")

	auxiliares.MostrarArchivo(globals.ClientConfig.Swapfile_path)

	return 1
}

/*
Que hace ReservarMemoriaSwapeado?
lo mismo que reservarmemoria pero para procesos que ya existen
verifica si el proceso entra en memoria y en tal caso asigna las paginas al proceso en la tabla simple, ademas actualiza las paginasdisponibles del globals con las nuevas del proceso
*/

func ReservarMemoriaSwapeado(pid int, tam int) {

	var PaginasNecesarias float64 = math.Ceil(float64(tam) / float64(globals.ClientConfig.Page_size)) //redondea para arriba para saber cuantas paginas ocupa

	var frames globals.ProcesoEnMemoria
	frames.TablaSimple = make([]int, 0) //inicializa el slice donde vamos a guardar la tabla de paginas simple para el proceso

	var PaginasEncontradas int = 0

	if EntraEnMemoria(tam) >= 0 {
		log.Printf("printeo tabla de paginas disponibles\n")

		for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ { //recorremos array de paginas disponibles a ver si encontramos la cantidad que necesitamos contiguas en memoria
			log.Printf("%d\n", globals.PaginasDisponibles[i])

			if globals.PaginasDisponibles[i] == 0 {
				PaginasEncontradas++
				frames.TablaSimple = append(frames.TablaSimple, i)
				globals.PaginasDisponibles[i] = 1 //reservamos la pagina (podemos hacerlo ya que se llamo a EntraEnMemoriaYVerificaSiYaExiste anteriormente)

				if PaginasEncontradas == int(PaginasNecesarias) {
					auxiliares.ActualizarTablaSimple(frames, pid) //ver semaforos

					auxiliares.MostrarProceso(pid)

					return
				}
			}
		}
	}

}
