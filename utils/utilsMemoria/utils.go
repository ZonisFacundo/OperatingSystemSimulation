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

/*
	type respuestaalKernel struct {
		Mensaje string `json:"message"`
	}
*/
type respuestaalKernel struct {
	Mensaje string `json:"message"`
	Exito   bool   `json:"exito"`
}
type respuestaalCPU struct {
	Mensaje string `json:"message"`
}

// http codigo

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
	log.Printf("Recibido del kernel: \n pid: %d  tam: %d  tambien recibimos un archivo con esta ruta: %s \n", (PaqueteInfoProceso).Pid, (PaqueteInfoProceso).TamProceso, (PaqueteInfoProceso.Archivo))

	//el kernel quiere saber si podemos guardar eso en memoria, para eso vamos a consultar el espacio que tenemos
	DondeGuardarProceso = EntraEnMemoria(PaqueteInfoProceso.TamProceso, PaqueteInfoProceso.Pid) //devuelve menor a 0 si no entra en memoria el proceso

	if DondeGuardarProceso < 0 {
		log.Printf("NO HAY ESPACIO EN MEMORIA PARA GUARDAR EL PROCESO \n")
		respuestaKernel.Mensaje = "No hay espacio para guardar el proceso en memoria crack"
		respuestaKernel.Exito = false //estamos probando esto como respuesta ademas del mesanje
		respuestaJSON, err := json.Marshal(respuestaKernel)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusInsufficientStorage) //http tiene un mensaje de error especificamente para esto, tremendo
		w.Write(respuestaJSON)
	} else {
		CrearProceso(PaqueteInfoProceso)

		respuestaKernel.Mensaje = "Recibi de Kernel"
		respuestaKernel.Exito = true //estamos probando esto como respuesta ademas del mesanje
		respuestaJSON, err := json.Marshal(respuestaKernel)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respuestaJSON)
	}
}

func RetornoClienteCPUServidorMEMORIATraduccionLogicaAFisica(w http.ResponseWriter, r *http.Request) {

	var DireccionLogica []int

	err := json.NewDecoder(r.Body).Decode(&DireccionLogica) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := 0; i < globals.ClientConfig.Number_of_levels+1; i++ {
		log.Printf("entranda nivel %d: %d\n", i, DireccionLogica[i])

	}

	log.Printf("desplazamiento %d: \n", DireccionLogica[globals.ClientConfig.Number_of_levels+1])

	var Traduccion globals.DireccionFisica = TraducirLogicaAFisica(DireccionLogica, globals.PunteroBase)

	respuestaJSON, err := json.Marshal(Traduccion)
	if err != nil {
		return
	}

	if Traduccion.Direccion == -1 {
		log.Printf("ERROR, envio una entrada mayor a la cantidad de entradas posibles en la configuracion actual \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if Traduccion.Direccion == -2 {
		log.Printf("ERROR, envio un desplazamiento (%d) mayor al tam de pagina de la configuracion actual (%d)  \n", DireccionLogica[0], globals.ClientConfig.Page_size)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("DIRECCION FISICA HALLADA:  %d: \n", Traduccion.Direccion)
	log.Printf("MARCO:  %d: \n", Traduccion.Marco)
	log.Printf("DESPLAZAMIENTO:  %d: \n", Traduccion.Desplazamiento)

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

// lee y devuelve a CPU lo que quiere de memoria principal
func RetornoClienteCPUServidorMEMORIARead(w http.ResponseWriter, r *http.Request) {

	var PaqueteDireccion globals.DFisica
	err := json.NewDecoder(r.Body).Decode(&PaqueteDireccion) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ContenidoDireccion globals.BytePaquete
	ContenidoDireccion.Info = globals.MemoriaPrincipal[PaqueteDireccion.DireccionFisica]

	respuestaJSON, err := json.Marshal(ContenidoDireccion)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func RetornoClienteCPUServidorMEMORIAWrite(w http.ResponseWriter, r *http.Request) {

	var PaqueteInfoWrite globals.PaqueteWrite

	err := json.NewDecoder(r.Body).Decode(&PaqueteInfoWrite) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	PaqueteInfoWrite.Contenido = globals.MemoriaPrincipal[PaqueteInfoWrite.Direccion]

	var rta respuestaalCPU
	rta.Mensaje = "OK\n"

	respuestaJSON, err := json.Marshal(rta)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

/*
/
/
/
/
	FUNCIONES NO HTTP
/
/
/
/

*/

func ConfigurarLogger() {
	logFile, err := os.OpenFile("memory.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
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

	//creo una funcion para hacerlo porque sino rompeutilsMemoria
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

func CrearEInicializarTablaDePaginas(PunteroANodo *globals.Nodo, nivel int) {

	if nivel < 0 {
		log.Printf("No puede haber una estructura con niveles negativos...")
		return
	}
	if nivel == globals.ClientConfig.Number_of_levels {

		(*PunteroANodo).Marco = make([]int, globals.ClientConfig.Entries_per_page)
		for j := 0; j < globals.ClientConfig.Entries_per_page; j++ {
			(*PunteroANodo).Marco[j] = -1 //lo dejo en -1 porque si los dejo en 0 podria significar una pagina valida
		}
		return
	}
	(*PunteroANodo).Siguiente = make([]*globals.Nodo, globals.ClientConfig.Entries_per_page) //inicializa el globals.Nodo -> sgte

	for entrada := 0; entrada < globals.ClientConfig.Entries_per_page; entrada++ {

		(*PunteroANodo).Siguiente[entrada] = new(globals.Nodo)
		CrearEInicializarTablaDePaginas((*PunteroANodo).Siguiente[entrada], nivel+1)

	}

}

/*
que hace traducirlogicaafisica?

recibe el slice de cpu
DireccionLogica[0] = desplazamiento
DireccionLogica[1] = entrada nivel 1
...
DireccionLogica[n] = entrada nivel n

a partir de estos datos accede a la tabla de paginas y devuelve el marco asociado a tal direccion logica, ademas el desplazamiento en otro valor aparte.
Tambien devuelvo la direccion en forma de bytes para que CPU use la que mas le guste
*/
func TraducirLogicaAFisica(DireccionLogica []int, PunteroNodo *globals.Nodo) globals.DireccionFisica {

	var DireccionFisica globals.DireccionFisica

	//VERIFICO SI LOS DATOS QUE MANDO CPU TIENEN SENTIDO (O SEA, NO HAY VALORES MAYORES A LOS DE LA CANTIDAD DE NIVELES/ENTRADAS/TAMDEPAGINA QUE TENEMOS DEFINIDOS)
	for i := 1; i <= globals.ClientConfig.Number_of_levels; i++ { //arrancamos desde 1 porque en 0 esta el desplazamiento, nos fijamos si la entrada nivel n es mayor a la cantidad de entradas por tabla
		if DireccionLogica[i] >= globals.ClientConfig.Entries_per_page {

			DireccionFisica.Desplazamiento = -1
			DireccionFisica.Marco = -1
			DireccionFisica.Direccion = -1 //para marcar error
			return DireccionFisica
		}
	}

	if DireccionLogica[0] >= globals.ClientConfig.Page_size { //nos envio un desplazamiento dentro de la pagina mayor al tam de la pagina
		DireccionFisica.Desplazamiento = -2
		DireccionFisica.Marco = -2
		DireccionFisica.Direccion = -2 //para marcar error
		return DireccionFisica
	}

	//SI LLEGAMOS ACA, LO QUE ENVIO CPU TIENE SENTIDO

	marco := AccedeAEntrada(DireccionLogica, 1, PunteroNodo)

	DireccionFisica.Marco = marco
	DireccionFisica.Desplazamiento = DireccionLogica[0]                                       //desplazamiento
	DireccionFisica.Direccion = (marco * globals.ClientConfig.Page_size) + DireccionLogica[0] //tam de pagina * numero de pagina + desplazamiento dentro de pagina

	return DireccionFisica

}

/*
Direccion logica tiene la forma de
DireccionLogica[0] = desplazamiento
DireccionLogica[1] = entrada nivel 1
...
DireccionLogica[n] = entrada nivel n

Para aumentar expresividad en el codigo (no estar agregando i - 1 en los loops por ejemplo)
*/

func AccedeAEntrada(DireccionLogica []int, nivel int, PunteroNodo *globals.Nodo) int {

	if nivel == globals.ClientConfig.Number_of_levels { //significa que ya estamos parados en el nivel que contiene los marcos
		return ((*PunteroNodo).Marco[nivel])

	} else {
		return AccedeAEntrada(DireccionLogica, nivel+1, (*PunteroNodo).Siguiente[DireccionLogica[nivel]])

	}
}

/*
Leer Página completa
Se deberá devolver el contenido correspondiente de la página a partir del byte enviado como dirección física dentro de la Memoria de Usuario, que deberá coincidir con la posición del byte 0 de la página.
Actualizar página completa
Se escribirá la página completa a partir del byte 0 que igual será enviado como dirección física, esta operación se realizará dentro de la Memoria de Usuario y se responderá como OK.
*/
/*
func LeerPaginaCompleta (direccion int) {

	if direccion % globals.ClientConfig.Page_size != 0 {

		log.Printf("")
	}
}*/
