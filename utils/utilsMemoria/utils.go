package utilsMemoria

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/auxiliares"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
)

//					STRUCTS

type PaqueteRecibidoMemoriadeCPU struct {
	Pid int `json:"pid"`
	Pc  int `json:"pc"`
}

type PaqueteRecibidoMemoriadeKernel struct {
	Pid        int    `json:"pid"`
	TamProceso int    `json:"tamanioproceso"`
	Archivo    string `json:"file"`
}
type respuestaalKernel struct {
	Mensaje string `json:"message"`
}
type respuestaalCPU struct {
	Mensaje string `json:"message"`
}

// FUNCIONES.
func ConfigurarLogger() {
	logFile, err := os.OpenFile("memory.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func RetornoClienteCPUServidorMEMORIA(w http.ResponseWriter, r *http.Request) {

	err := json.NewDecoder(r.Body).Decode(&globals.Instruction) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Cliente envio: \n pid: %d \n pc: %d", globals.Instruction.Pid, globals.Instruction.Pc)

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuestaCpu respuestaalCPU
	respuestaCpu.Mensaje = globals.MemoriaKernel[globals.Instruction.Pid].Instrucciones[globals.Instruction.Pc]
	respuestaJSON, err := json.Marshal(respuestaCpu)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func RetornoClienteKernelServidorMEMORIA(w http.ResponseWriter, r *http.Request) {

	var DondeGuardarProceso int
	var respuestaKernel respuestaalKernel
	var PaqueteInfoProceso PaqueteRecibidoMemoriadeKernel //variable global donde guardo lo que me mande el kernel (info del proceso)

	err := json.NewDecoder(r.Body).Decode(&PaqueteInfoProceso) //guarda en una variable global lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//leo lo que nos mando el cliente, en este caso un struct de dos strings y un int
	log.Printf("Recibido del kernel: \n pid: %d  tam: %d  tambien recibimos un archivo\n", (PaqueteInfoProceso).Pid, (PaqueteInfoProceso).TamProceso)
	(PaqueteInfoProceso).Archivo = "/home/utnso/archivosprueba/archi.txt" //voy a tener que recibir un archivo de kernel, esto es de prueba

	//el kernel quiere saber si podemos guardar eso en memoria, para eso vamos a consultar el espacio que tenemos
	DondeGuardarProceso = EntraEnMemoria(PaqueteInfoProceso.TamProceso, PaqueteInfoProceso.Pid) //devuelve menor a 0 si no entra en memoria el proceso

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
		CrearProceso(PaqueteInfoProceso)

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
		globals.MemoriaPrincipal[200] = 1S
	*/
}

func InicializarPaginasDisponibles() {

	globals.PaginasDisponibles = make([]int, (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size))

	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ {
		globals.PaginasDisponibles[i] = 0
	}

}

/*
	DEPRECADO

func EscanearMemoria() {

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
*/

/*
QUE HACE RESERVAR MEMORIA?

reservar memoria basicamente recibe informacion sobre un proceso que quiere iniciar kernel y guarda en el map que tenemos con informacion basica de proceso las paginas que este tiene reservada en memoria
*/
func ReservarMemoria(tam int, pid int) int {

	var PaginasNecesarias float64 = math.Ceil(float64(tam) / float64(globals.ClientConfig.Page_size)) //redondea para arriba para saber cuantas paginas ocupa

	var frames globals.ProcesoEnMemoria
	frames.TablaSimple = make([]int, 0) //inicializa el slice donde vamos a guardar la tabla de paginas simple para el proceso

	var PaginasEncontradas int = 0
	if EntraEnMemoria(tam, pid) >= 0 {
		for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ { //recorremos array de paginas disponibles a ver si encontramos la cantidad que necesitamos contiguas en memoria

			if globals.PaginasDisponibles[i] == 0 {
				PaginasEncontradas++
				frames.TablaSimple = append(frames.TablaSimple, i)
				globals.PaginasDisponibles[i] = 1 //reservamos la pagina (podemos hacerlo ya que se llamo a EntraEnMemoria anteriormente)

				if PaginasEncontradas == int(PaginasNecesarias) {
					auxiliares.ActualizarTablaSimple(frames, pid)

					auxiliares.MostrarProceso(pid)

					return 1 //devuelvo numero positivo para indicar que fue un exito, asignamos todas las paginas al proceso
				}
			}
		}
	}
	return -1
}

func EntraEnMemoria(tam int, pid int) int {

	var PaginasNecesarias float64 = math.Ceil(float64(tam) / float64(globals.ClientConfig.Page_size)) //redondea para arriba para saber cuantas paginas ocupa
	log.Printf("necesitamos %f paginas para guardar este proceso, dejame ver si tenemos", PaginasNecesarias)

	var PaginasEncontradas int = 0

	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ { //recorremos array de paginas disponibles para ver si entran todas las paginas del proceso

		if globals.PaginasDisponibles[i] == 0 {
			PaginasEncontradas++

			if PaginasEncontradas == int(PaginasNecesarias) {

				return 1 //devuelvo numero positivo para indicar que fue entra
			}
		}
	}
	return -1 //no entra en memoria
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

		if buffer[i] == 10 { //ASCII para \n
			Contenido.Instrucciones = append(Contenido.Instrucciones, Line) //agrega la instruccion al slice de strings (donde cada elemento (cada string) es una instruccion)

			Line = ""
		}
		Line += string(buffer[i]) //va armando un string caracter a caracter hasta formar una instruccion (cuando lee \n)

	}
	//	globals.MemoriaKernel[Pid].Instrucciones = Contenido.Instrucciones    esto no anda, hay que hacerlo con una copia //carga instrucciones al map global, lo que verdaderamente importa

	//creo una funcion para hacerlo porque sino rompe
	auxiliares.ActualizarInstrucciones(Contenido, Pid)
	//lo muestro a ver si funco
	for j := 0; j < len(globals.MemoriaKernel[Pid].Instrucciones); j++ {
		fmt.Printf("%s", globals.MemoriaKernel[Pid].Instrucciones[j])
	}

}
func CrearProceso(paquete PaqueteRecibidoMemoriadeKernel) {
	if ReservarMemoria(paquete.TamProceso, paquete.Pid) < 0 { //ReservarMemoria devuelve <0 si hubo un error, si no hubieron errores actualiza el map y reserva la memoria para el proceso

		log.Printf("error al reservar memoria para el proceso de pid: %d", (paquete).Pid)
		return
	}

	//llevamos contenido del archivo al map
	LeerArchivoYCargarMap((paquete).Archivo, (paquete).Pid)
	log.Printf("## PID: %d - Proceso Creado - Tamaño: %d \n", paquete.Pid, paquete.TamProceso)

}
func CrearTablaDePaginas() {

}
