package utilsMemoria

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	//"github.com/sisoputnfrba/tp-golang/memoria/auxiliares"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
)

//					STRUCTS

type Handshakepaquete struct {
	Instruccion string `json:"instruccion"`
}

// CAMBIO EL PAQUETE QUE RECIBO DE KERNEL
/*vamos  tener que usar este eventualmente
type PaqueteRecibidoMemoriadeKernel struct {
	NombreSyscall string `json:"syscallname"` //no se si necesito esto
	TamProceso    int    `json:"processsize"`
	archivo       string `json:"file"`
	Pid           int    `json:"pid"`
}
*/
type PaqueteRecibidoMemoriadeKernel struct {
	Pid        int    `json:"pid"`
	TamProceso int    `json:"tamanioproceso"`
	Archivo    string `json:"file"`
}

type respuestaalCPU struct {
	Mensaje string `json:"message"`
}
type respuestaalKernel struct {
	Mensaje string `json:"message"`
}

//						FUNCIONES.

func ConfigurarLogger() {
	logFile, err := os.OpenFile("memory.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func RetornoClienteCPUServidorMEMORIA(w http.ResponseWriter, r *http.Request) {

	var request Handshakepaquete

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("(CPU) El cliente nos mando esto: \n instruccion: %s.\n", request.Instruccion)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuestaCpu respuestaalCPU
	respuestaCpu.Mensaje = "Recibi de CPU"
	respuestaJSON, err := json.Marshal(respuestaCpu)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

/*
func RetornoClienteKernelServidorMEMORIA(w http.ResponseWriter, r *http.Request) {

	var request HandshakepaqueteKernel

	err := json.NewDecoder(r.Body).Decode(&request) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("(kernel) El cliente nos mando esto: \n nombre pseudocodigo: %s, tamanio proceso: %d.\n", request.NombreCodigo, request.TamanioProceso)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuestaCpu respuestaalCPU
	respuestaCpu.Mensaje = "Recibi de Kernel"
	respuestaJSON, err := json.Marshal(respuestaCpu)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}
*/
//cambio api con kernel para recibir paquete deseado
func RetornoClienteKernelServidorMEMORIA(w http.ResponseWriter, r *http.Request) {

	var DondeGuardarProceso int
	var respuestaKernel respuestaalKernel

	err := json.NewDecoder(r.Body).Decode(globals.PaqueteInfoProceso) //guarda en una variable global lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("Recibido del kernel: \n pid: %d  tam: %d  tambien recibimos un archivo\n", (*globals.PaqueteInfoProceso).Pid, (*globals.PaqueteInfoProceso).TamProceso)

	//el kernel quiere saber si podemos guardar eso en memoria, para eso vamos a consultar el espacio que tenemos
	DondeGuardarProceso = EntraEnMemoria(globals.PaqueteInfoProceso.TamProceso)
	log.Printf("lo guardamos a partir de la pagina %d \n", DondeGuardarProceso)

	if DondeGuardarProceso < 0 {
		log.Printf("NO HAY ESPACIO EN MEMORIA PARA GUARDAR EL PROCESO \n")
		respuestaKernel.Mensaje = "No hay espacio para guardar el proceso en memoria crack"
		respuestaJSON, err := json.Marshal(respuestaKernel)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusInsufficientStorage) //http tiene un mensaje de error especificamente para esto, tremendo
		w.Write(respuestaJSON)
	} else {
		//	respuesta del server al cliente

		//llevamos contenido del archivo al map

		(*globals.PaqueteInfoProceso).Archivo = "/home/utnso/archivosprueba/archi.txt" //voy a tener que recibir un archivo de kernel, esto es de prueba

		LeerArchivoYCargarMap((*globals.PaqueteInfoProceso).Archivo, (*globals.PaqueteInfoProceso).Pid)

		respuestaKernel.Mensaje = "Recibi de Kernel"
		respuestaJSON, err := json.Marshal(respuestaKernel)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respuestaJSON)
	}
}

func InicializarMemoria() {

	globals.MemoriaPrincipal = make([]byte, globals.ClientConfig.Memory_size) //inicializa la memoria segun lo que decia el enunciado

	//Descomentar si tenes ganas de ver si anda
	/*
		globals.MemoriaPrincipal[22] = 1
		globals.MemoriaPrincipal[80] = 1
		globals.MemoriaPrincipal[200] = 1
	*/
}

func InicializarPaginasDisponibles() {

	globals.PaginasDisponibles = make([]int, (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size))

	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ {
		globals.PaginasDisponibles[i] = 0
	}

}

func ActualizaPaginasDisponibles() {

	//buscamos espacio contiguo en memoria, la memoria esta dividida en paginas
	//primer for recorre de a paginas, segundo for recorre cada pagina buscando ver si esta libre o no
	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ {
		//	fmt.Printf("entre al ciclo i \n")
		for j := 0; j < globals.ClientConfig.Page_size; j++ {
			//		fmt.Printf("entre al ciclo j \n")

			if globals.MemoriaPrincipal[(i*globals.ClientConfig.Page_size)+j] != 0 {
				//			fmt.Printf(" \n \n \n DIJE QUE ESTA OCUPADAAA \n \n")

				globals.PaginasDisponibles[i] = 1   //marcamos que esta ocupada
				j += globals.ClientConfig.Page_size //salimos de la pagina si sabemos que esta ocupada
			} else if j == globals.ClientConfig.Page_size-1 {
				globals.PaginasDisponibles[i] = 0 //marcamos que esta desocupada
			}

		}

	}

}
func EntraEnMemoria(tam int) int {

	var PaginasNecesarias float64 = math.Ceil(float64(tam) / float64(globals.ClientConfig.Page_size)) //redondea para arriba para saber cuantas paginas ocupa
	log.Printf("necesitamos %f paginas para guardar este proceso, dejame ver si tenemos", PaginasNecesarias)

	var PaginasContiguasEncontradas int = 0

	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ { //recorremos array de paginas disponibles a ver si encontramos la cantidad que necesitamos contiguas en memoria

		if globals.PaginasDisponibles[i] == 0 {
			PaginasContiguasEncontradas++
			if PaginasContiguasEncontradas == int(PaginasNecesarias) {
				return (i - int(PaginasNecesarias) + 1) //devuelvo el indice del primer marco de pagina que vamos a usar para guardar el proceso
			}
		} else {
			PaginasContiguasEncontradas = 0
		}
	}
	return -1
}

func LeerArchivoYCargarMap(FilePath string, Pid int) {

	var buffer []byte
	var err error
	var Contenido globals.ProcesoEnMemoria //guardo lo que voy viendo del archivo organizadito para pasarselo a MemoriaKernel
	Contenido.Instrucciones = make([]string, 0)
	var Line string = ""
	buffer, err = os.ReadFile(FilePath)

	if err != nil {
		log.Printf("Error al leer el archivo enviado por Kernel Pid: %d", Pid)
	}

	for i := 0; i < (len(buffer)); i++ {

		Line += string(buffer[i]) //va armando un string caracter a caracter hasta formar una instruccion (cuando lee \n)

		if buffer[i] == 10 { //ASCII para \n
			Contenido.Instrucciones = append(Contenido.Instrucciones, Line) //agrega la instruccion al slice de strings (donde cada elemento (cada string) es una instruccion)

			Line = ""
		}

	}

	globals.MemoriaKernel[Pid] = Contenido

	//lo muestro a ver si funco
	for i := 0; i < len(globals.MemoriaKernel); i++ {
		for j := 0; j < len(globals.MemoriaKernel[i].Instrucciones); j++ {
			fmt.Printf("%s", globals.MemoriaKernel[i].Instrucciones[j])
		}

	}

}
