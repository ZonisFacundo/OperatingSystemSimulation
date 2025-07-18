package utilsMemoria

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/sisoputnfrba/tp-golang/memoria/auxiliares"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
)

// STRUCTS
type PaqueteRecibidoMemoriadeCPU struct {
	Pid int `json:"pid"`
	Pc  int `json:"pc"`
}

type PaqueteRecibidoMemoriadeKernel struct {
	Pid        int    `json:"pid"`
	TamProceso int    `json:"tamanioproceso"`
	Archivo    string `json:"file"`
}

// Hice este struct para la respuesta generica (Santi)
type PaqueteRecibidoMemoriadeKernel2 struct {
	Pid     int    `json:"pid"`
	Mensaje string `json:"message"`
}

type respuestaalKernel struct {
	Mensaje string `json:"message"`
}
type respuestaalCPU struct {
	Mensaje string `json:"message"`
}

type PaqueteCPUHandshake struct {
	Entradas int `json:"ent"`
	Niveles  int `json:"niv"`
	TamPag   int `json:"tam"`
}

// http codigo

// handshake para pasar datos del globals de memoria nada mas... ignorar
func HandshakeACpu(w http.ResponseWriter, r *http.Request) {

	var recibo respuestaalCPU
	err := json.NewDecoder(r.Body).Decode(&recibo) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var datos PaqueteCPUHandshake

	datos.Entradas = globals.ClientConfig.Entries_per_page
	datos.Niveles = globals.ClientConfig.Number_of_levels
	datos.TamPag = globals.ClientConfig.Page_size
	log.Printf("entradas: %d, niveles: %d, tam pagina: %d, (Envio valores del config a cpu) (Handshake) \n\n", datos.Entradas, datos.Niveles, datos.TamPag)

	respuestaJSON, err := json.Marshal(datos)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func RetornoClienteCPUServidorMEMORIA(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(globals.ClientConfig.Memory_delay) * time.Millisecond)

	var InstruccionLocal globals.Instru
	// globals.Sem_Instruccion.Lock()
	err := json.NewDecoder(r.Body).Decode(&InstruccionLocal) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		// globals.Sem_Instruccion.Unlock()
		return
	}
	globals.Sem_MemoriaKernel.Lock()
	// globals.Sem_Instruccion.Lock()
	auxiliares.InicializarSiNoLoEstaMap(InstruccionLocal.Pid)
	globals.MetricasProceso[InstruccionLocal.Pid].ContadorInstruccionesSolicitadas++
	//globals.Sem_Instruccion.Unlock()
	globals.Sem_MemoriaKernel.Unlock()
	//log.Printf("## PID: <%d> - Obtener instrucción: <%d> - Instrucción: %s\n", globals.Instruction.Pid, globals.Instruction.Pid, globals.MemoriaKernel[globals.Instruction.Pid].Instrucciones[globals.Instruction.Pc])

	log.Printf("## PID: <%d>- Obtener instrucción: <%d> - Instrucción: %s\n", InstruccionLocal.Pid, InstruccionLocal.Pc, globals.MemoriaKernel[InstruccionLocal.Pid].Instrucciones[InstruccionLocal.Pc])
	// globals.Sem_Instruccion.Unlock()

	//	respuesta del server al cliente, no hace falta en este modulo pero en el que estas trabajando seguro que si
	var respuestaCpu respuestaalCPU

	//log.Printf("\nla longitud del archivo de instrucciones es: %d\n\n", len(globals.MemoriaKernel[globals.Instruction.Pid].Instrucciones))

	//log.Printf("estamos mandandole a CPU, del pid: %d la instrucion del pc: %d la cual es %s \n\n", globals.Instruction.Pid, globals.Instruction.Pc, globals.MemoriaKernel[globals.Instruction.Pid].Instrucciones[globals.Instruction.Pc])
	globals.Sem_MemoriaKernel.Lock()
	// globals.Sem_Instruccion.Lock()
	respuestaCpu.Mensaje = globals.MemoriaKernel[InstruccionLocal.Pid].Instrucciones[InstruccionLocal.Pc]
	// globals.Sem_Instruccion.Unlock()
	globals.Sem_MemoriaKernel.Unlock()

	respuestaJSON, err := json.Marshal(respuestaCpu)

	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)
}

func RetornoClienteKernelServidorMEMORIA(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(globals.ClientConfig.Memory_delay) * time.Millisecond)

	time.Sleep(time.Duration(globals.ClientConfig.Memory_delay) * time.Millisecond)
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
	globals.Sem_Bitmap.Lock()                                           //TODO
	DondeGuardarProceso = EntraEnMemoria(PaqueteInfoProceso.TamProceso) //devuelve menor a 0 si no entra en memoria el proceso
	globals.Sem_Bitmap.Unlock()

	if DondeGuardarProceso == -1 {
		log.Printf("NO HAY ESPACIO EN MEMORIA PARA GUARDAR EL PROCESO \n")
		respuestaKernel.Mensaje = "No hay espacio para guardar el proceso en memoria crack"
		respuestaJSON, err := json.Marshal(respuestaKernel)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusInsufficientStorage) //http tiene un mensaje de error especificamente para esto, tremendo
		w.Write(respuestaJSON)
		return
	} else {

		/*
			verificamos si el archivo que nos enviaron es valido
		*/

		_, err := os.ReadFile(PaqueteInfoProceso.Archivo)

		if err != nil {
			log.Printf("error al abrir el archivo de instrucciones enviado por kernel, pid: %d\n", PaqueteInfoProceso.Pid)
			log.Printf("TENEMOS ESPACIO, EL PROBLEMA ES EN EL ARCHIVO, MANDO ESTE STATUS CODE PORQUE ES MAS PRACTICO POR EL CODIGO DE KERNEL\n")
			w.WriteHeader(http.StatusInsufficientStorage)
			return
		}

		CrearProceso(PaqueteInfoProceso)

		respuestaKernel.Mensaje = "Recibi de Kernel"
		respuestaJSON, err := json.Marshal(respuestaKernel)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(respuestaJSON)
		return
	}
}

func RetornoClienteCPUServidorMEMORIATraduccionLogicaAFisica(w http.ResponseWriter, r *http.Request) {

	var Paquete globals.DireccionLogica
	Paquete.DirLogica = make([]int, 0) //inicializamos el slice

	err := json.NewDecoder(r.Body).Decode(&Paquete) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//time.Sleep(time.Duration(globals.ClientConfig.Memory_delay*(len(Paquete.DirLogica)-1)) * time.Millisecond)

	auxiliares.InicializarSiNoLoEstaMap(Paquete.DirLogica[0])

	for i := 0; i < (len(Paquete.DirLogica) - 1); i++ {
		time.Sleep(time.Duration(globals.ClientConfig.Memory_delay) * time.Millisecond)
		globals.MetricasProceso[Paquete.DirLogica[0]].ContadorAccesosTablaPaginas++
	}

	globals.Sem_MemoriaKernel.Lock()
	globals.PunteroBase = globals.MemoriaKernel[Paquete.DirLogica[0]].PunteroATablaDePaginas
	globals.Sem_MemoriaKernel.Unlock()

	var Traduccion globals.Marco = TraducirLogicaAFisica(Paquete.DirLogica, globals.PunteroBase)

	respuestaJSON, err := json.Marshal(Traduccion)
	if err != nil {
		return
	}

	if Traduccion.Frame == -1 {
		log.Printf("ERROR, envio una entrada mayor a la cantidad de entradas posibles en la configuracion actual \n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("MARCO:  %d: \n", Traduccion.Frame)
	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

// lee y devuelve a CPU lo que quiere de memoria principal
func RetornoClienteCPUServidorMEMORIARead(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(globals.ClientConfig.Memory_delay) * time.Millisecond)

	var PaqueteDireccion globals.PaqueteRead
	err := json.NewDecoder(r.Body).Decode(&PaqueteDireccion) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if PaqueteDireccion.Tamaño < 1 {
		log.Printf("\n\nTAM A LEER ES < 0, ERROR (HTTP Read)\n\n")
		return
	}
	auxiliares.InicializarSiNoLoEstaMap(PaqueteDireccion.Pid)

	globals.Sem_Metricas.Lock()
	globals.MetricasProceso[PaqueteDireccion.Pid].ContadorReadMemoria++
	globals.Sem_Metricas.Unlock()

	log.Printf("## PID: %d - <Escritura/Lectura> - Dir. Física: %d  - Tamaño: %d\n", PaqueteDireccion.Pid, PaqueteDireccion.Direccion, PaqueteDireccion.Tamaño)

	var ContenidoDireccion globals.BytePaquete
	ContenidoDireccion.Info = make([]byte, PaqueteDireccion.Tamaño)
	globals.Sem_Mem.Lock()

	ContenidoDireccion.PaginaCompleta, _ = LeerPaginaCompleta(int(math.Floor(float64(PaqueteDireccion.Direccion / globals.ClientConfig.Page_size))))

	for i := 0; i < PaqueteDireccion.Tamaño; i++ {

		ContenidoDireccion.Info[i] = globals.MemoriaPrincipal[PaqueteDireccion.Direccion]
	}
	globals.Sem_Mem.Unlock()
	respuestaJSON, err := json.Marshal(ContenidoDireccion)
	if err != nil {
		return
	}

	log.Printf("\n\nMUESTRO LO QUE LE MANDO A CPU COMO LEIDO (HTTP Read)\n\n")
	log.Print("array de bytes: \n", ContenidoDireccion.Info)

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func RetornoClienteCPUServidorMEMORIAWrite(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(globals.ClientConfig.Memory_delay) * time.Millisecond)

	var PaqueteInfoWrite globals.PaqueteWrite

	err := json.NewDecoder(r.Body).Decode(&PaqueteInfoWrite) //guarda en request lo que nos mando el cliente
	if err != nil {
		log.Printf("returneo porque no pude decodear (RetornoClienteCPUServidorMEMORIAWrite) \n")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//registramos en metrica que funco el write
	auxiliares.InicializarSiNoLoEstaMap(PaqueteInfoWrite.Pid)

	globals.Sem_Metricas.Lock()
	globals.MetricasProceso[PaqueteInfoWrite.Pid].ContadorWriteMemoria++
	globals.Sem_Metricas.Unlock()

	log.Printf("direccion recibida: %d\n", PaqueteInfoWrite.Direccion)

	bytardos := []byte(PaqueteInfoWrite.Contenido)

	log.Printf("## PID: %d - <Escritura> - Dir. Física: %d  - Tamaño: %d\n", PaqueteInfoWrite.Pid, PaqueteInfoWrite.Direccion, len(PaqueteInfoWrite.Contenido))

	globals.Sem_Mem.Lock()
	for i := 0; i < len(PaqueteInfoWrite.Contenido); i++ {
		//log.Printf("%b", bytardos[i])
		globals.MemoriaPrincipal[PaqueteInfoWrite.Direccion+i] = bytardos[i]
	}
	globals.Sem_Mem.Unlock()
	var rta respuestaalCPU
	rta.Mensaje = "OK\n"

	respuestaJSON, err := json.Marshal(rta)
	if err != nil {
		return
	}

	//log.Printf("\n\nMUESTRO LA MEMORIA DONDE SE ESCRIBIO LO QUE NOS PIDIO CPU \n\n")
	//auxiliares.Mostrarmemoria()
	//log.Printf("\n\n")

	/*
		MUCHISIMO CUIDADO, GLOBALS,CONTAOR Y MOSTRAR TABLA O DESCOMENTAS AMBAS O COMENTAS AMBAS, PORQUE HAY RACE CONDITION SINO, EL PUNTO ES QUE ES MERAMENTE PARA DEBUG EL MOSTRAR ASI QUE NO VALE LA PENA HACER ALGO AL RESPECTO
		//log.Printf("\n\nMUESTRO LA tabla de paginas multinivel de pid 0\n\n")

		//globals.ContadorTabla = 0
		//descomentar si queres ver la tabla de paginas
		//var PunteritoAux *globals.Nodo = globals.MemoriaKernel[0].PunteroATablaDePaginas
		//MostrarTablaMultinivel(0, 0, PunteritoAux)
	*/
	// auxiliares.Mostrarmemoria()

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

}

func RetornoClienteKernelServidorMemoriaDumpDelProceso(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(globals.ClientConfig.Memory_delay) * time.Millisecond)
	//Este paquete lo unico q recibe es el pid para hacerle el dump junto a un mensaje
	var paqueteDeKernel PaqueteRecibidoMemoriadeKernel2
	err := json.NewDecoder(r.Body).Decode(&paqueteDeKernel) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("## PID: <%d> - Memory Dump solicitado\n", paqueteDeKernel.Pid)

	MemoryDump(paqueteDeKernel.Pid)

	var respuesta respuestaalKernel
	respuesta.Mensaje = "DUMP REALIZADO CON EXITO \n"

	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)

	//para testear
	log.Printf("\n----------------MUESTRO TODO ASI DEBUGGEAS---------------------\n\n")
	log.Printf("\n----------------MUESTRO MEMORIA---------------------\n\n")

	// auxiliares.Mostrarmemoria()

	log.Printf("\n\n-------------ahora muestro el swap------------------\n\n")
	// auxiliares.MostrarArchivo(globals.ClientConfig.Swapfile_path)

	log.Printf("\n\n-------------ahora muestro el dump------------------\n\n")
	// auxiliares.MostrarArchivo(fmt.Sprintf("%s", globals.ClientConfig.Dump_path))

}
func RetornoClienteKernelServidorMemoriaFinProceso(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(globals.ClientConfig.Memory_delay) * time.Millisecond)
	//Este paquete lo unico q recibe es el pid para hacerle el dump junto a un mensaje
	var paqueteDeKernel PaqueteRecibidoMemoriadeKernel2
	err := json.NewDecoder(r.Body).Decode(&paqueteDeKernel) //guarda en request lo que nos mando el cliente
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	FinalizarProceso(paqueteDeKernel.Pid)

	var respuesta respuestaalKernel
	respuesta.Mensaje = "PROCESO FINALIZADO CON EXITO \n"

	respuestaJSON, err := json.Marshal(respuesta)
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

	globals.Sem_Bitmap.Lock()
	globals.PaginasDisponibles = make([]int, (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size))

	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ {
		globals.PaginasDisponibles[i] = 0
	}
	globals.Sem_Bitmap.Unlock()
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
¿QUE HACE RESERVAR MEMORIA?

reservar memoria basicamente recibe informacion sobre un proceso que quiere iniciar kernel y guarda en el map que tenemos con informacion basica de proceso
las paginas que este tiene reservada en memoria.

*/

func ReservarMemoria(tam int, pid int) int {

	var PaginasNecesarias float64 = math.Ceil(float64(tam) / float64(globals.ClientConfig.Page_size)) //redondea para arriba para saber cuantas paginas ocupa

	var frames globals.ProcesoEnMemoria
	frames.TablaSimple = make([]int, 0) //inicializa el slice donde vamos a guardar la tabla de paginas simple para el proceso

	if tam == 0 {
		return 1
	}

	var PaginasEncontradas int = 0
	globals.Sem_Bitmap.Lock()

	if EntraEnMemoria(tam) >= 0 {

		for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ { //recorremos array de paginas disponibles a ver si encontramos la cantidad que necesitamos contiguas en memoria

			if globals.PaginasDisponibles[i] == 0 {
				PaginasEncontradas++
				frames.TablaSimple = append(frames.TablaSimple, i)
				globals.PaginasDisponibles[i] = 1 //reservamos la pagina (podemos hacerlo ya que se llamo a EntraEnMemoriaYVerificaSiYaExiste anteriormente)

				if PaginasEncontradas == int(PaginasNecesarias) {
					auxiliares.ActualizarTablaSimple(frames, pid)

					//	auxiliares.MostrarProceso(pid)
					globals.Sem_Bitmap.Unlock()

					return 1 //devuelvo numero positivo para indicar que fue un exito, asignamos todas las paginas al proceso
				}
			}
		}

	}
	globals.Sem_Bitmap.Unlock()

	return -1
}

/*
SOLO USAR SI EN ALGUN TEST TRATAN DE CREAR UN PROCESO CON UN PID QUE YA FUE USADO, EN TEORIA NO
func EntraEnMemoriaYVerificaSiYaExiste(tam int, pid int) int {

	for key, _ := range globals.MemoriaKernel {
		if key == pid {
			log.Printf("UN PROCESO CON PID: %d ya existe, no se puede crear el proceso", pid)
			return -2
		}
	}
	var PaginasNecesarias float64 = math.Ceil(float64(tam) / float64(globals.ClientConfig.Page_size)) //redondea para arriba para saber cuantas paginas ocupa
	log.Printf("necesitamos %f paginas para guardar este proceso, dejame ver si tenemos", PaginasNecesarias)

	var PaginasEncontradas int = 0

	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ { //recorremos array de paginas disponibles para ver si entran todas las paginas del proceso

		if globals.PaginasDisponibles[i] == 0 {
			PaginasEncontradas++

			if PaginasEncontradas == int(PaginasNecesarias) {

				return 1 //devuelvo numero positivo para indicar que  entra
			}
		}
	}
	return -1 //no entra en memoria
}
*/
/*
	QUE HACE ENTRAENMEMORIA?

lo mismo que el anterior pero no verifica si ya existe proceso con ese pid
se fija si se encuentra disponible el tam necesario
*/
func EntraEnMemoria(tam int) int {

	if tam == 0 {
		return 1
	}
	var PaginasNecesarias float64 = math.Ceil(float64(tam) / float64(globals.ClientConfig.Page_size)) //redondea para arriba para saber cuantas paginas ocupa
	log.Printf("necesitamos %f paginas para guardar este proceso, dejame ver si tenemos	(EntraEnMemoria) \n", PaginasNecesarias)

	var PaginasEncontradas int = 0

	for i := 0; i < (globals.ClientConfig.Memory_size / globals.ClientConfig.Page_size); i++ { //recorremos array de paginas disponibles para ver si entran todas las paginas del proceso

		if globals.PaginasDisponibles[i] == 0 {
			PaginasEncontradas++

			if PaginasEncontradas == int(PaginasNecesarias) {
				return 1 //devuelvo numero positivo para indicar que  entra
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
		return
	}
	for i := 0; i < (len(buffer)); i++ {

		if buffer[i] == 10 { //ASCII para \n
			Contenido.Instrucciones = append(Contenido.Instrucciones, Line) //agrega la instruccion al slice de strings (donde cada elemento (cada string) es una instruccion)
			Line = ""
		}
		Line += string(buffer[i]) //va armando un string caracter a caracter hasta formar una instruccion (cuando lee \n)
	}
	Contenido.Instrucciones = append(Contenido.Instrucciones, Line) //Lo vuelvo a agregar porque le falta el ultimo
	//	globals.MemoriaKernel[Pid].Instrucciones = Contenido.Instrucciones    esto no anda, hay que hacerlo con una copia //carga instrucciones al map global, lo que verdaderamente importa
	auxiliares.ActualizarInstrucciones(Contenido, Pid) //esta funcion es la que hace que ande copiar el contenido en memoria

	//creo una funcion para hacerlo porque sino rompeutilsMemoria
	globals.Sem_MemoriaKernel.Lock()
	for j := 0; j < len(globals.MemoriaKernel[Pid].Instrucciones); j++ {
		fmt.Printf("%s", globals.MemoriaKernel[Pid].Instrucciones[j])
	}
	globals.Sem_MemoriaKernel.Unlock()

}

// todo lo necesario para crear un proceso nuevo
func CrearProceso(paquete PaqueteRecibidoMemoriadeKernel) {

	if ReservarMemoria(paquete.TamProceso, paquete.Pid) < 0 { //ReservarMemoria devuelve <0 si hubo un error, si no hubieron errores actualiza el map y reserva la memoria para el proceso

		log.Printf("error al reservar memoria para el proceso de pid: %d,  (CrearProceso)", (paquete).Pid)
		return
	}

	//llevamos contenido del archivo al map
	LeerArchivoYCargarMap((paquete).Archivo, (paquete).Pid)

	//ojo aca, tengo que hacer estas maniobras porque no me deja asignarle de una un solo campo del struct al map... (ignorar)
	//creamos un puntero que apunte a la base de la tabla de paginas del proceso (un nodo), luego inicializamos la tabla de paginas
	globals.Sem_MemoriaKernel.Lock()
	aux := globals.MemoriaKernel[paquete.Pid]
	aux.PunteroATablaDePaginas = new(globals.Nodo)
	globals.MemoriaKernel[paquete.Pid] = aux
	globals.Sem_MemoriaKernel.Unlock()

	globals.Sem_MemoriaKernel.Lock()
	CrearEInicializarTablaDePaginas(globals.MemoriaKernel[paquete.Pid].PunteroATablaDePaginas, 0)
	globals.Sem_MemoriaKernel.Unlock()

	//ahora nos queda asignarle los marcos correspondientes al proceso segun la tabla de paginas simple que ya tenemos creada
	globals.Sem_MemoriaKernel.Lock()
	var PunteroAux *globals.Nodo = globals.MemoriaKernel[paquete.Pid].PunteroATablaDePaginas //es necesario enviar un puntero auxiliar por parametro en esta funcion
	globals.Sem_MemoriaKernel.Unlock()

	globals.Sem_MemoriaKernel.Lock()
	// globals.Sem_Contador.Lock()
	globals.ContadorTabla = 0
	AsignarValoresATablaDePaginas(paquete.Pid, 0, PunteroAux)
	globals.ContadorTabla = 0 //lo reinicio para que cuando otro proceso quiera usarlo este bien seteado en 0 y no en algun valor tipo 14 como lo dejo el proceso anterior (es la unica varialbe global de utils)
	// globals.Sem_Contador.Unlock()
	globals.Sem_MemoriaKernel.Unlock()

	auxiliares.InicializarSiNoLoEstaMap(paquete.Pid)

	ActualizarPaginasDisponibles() //actualiza que paginas estan disponibles en este momento

	log.Printf("## PID: %d - Proceso Creado - Tamaño: %d  (CrearProceso) \n", paquete.Pid, paquete.TamProceso)
}

func CrearEInicializarTablaDePaginas(PunteroANodo *globals.Nodo, nivel int) {

	if nivel < 0 {
		log.Printf("No puede haber una estructura con niveles negativos...")
		return
	}
	if nivel == globals.ClientConfig.Number_of_levels-1 {

		(*PunteroANodo).Marco = make([]int, globals.ClientConfig.Entries_per_page)
		for j := 0; j < globals.ClientConfig.Entries_per_page; j++ {
			(*PunteroANodo).Marco[j] = -10 //lo dejo en -10 porque si los dejo en 0 podria significar una pagina valida
			//descomentar si queres ver llenado de tabla de paginas multinivel
			//globals.Contador++
			//log.Printf("%d \t contador: %d\n", (*PunteroANodo).Marco[j], globals.Contador)

		}
		return
	} else {
		(*PunteroANodo).Siguiente = make([]*globals.Nodo, globals.ClientConfig.Entries_per_page) //inicializa el globals.Nodo -> sgte

		for entrada := 0; entrada < globals.ClientConfig.Entries_per_page; entrada++ {

			(*PunteroANodo).Siguiente[entrada] = new(globals.Nodo)
			CrearEInicializarTablaDePaginas((*PunteroANodo).Siguiente[entrada], nivel+1)

		}
	}
}

/*
que hace traducirlogicaafisica?

recibe el slice de cpu
DireccionLogica[0] = pid
DireccionLogica[1] = entrada nivel 1
...
DireccionLogica[n] = entrada nivel n

a partir de estos datos accede a la tabla de paginas y devuelve el marco asociado a tal direccion logica
*/
func TraducirLogicaAFisica(DireccionLogica []int, PunteroNodo *globals.Nodo) globals.Marco {

	var MarcoAurelio globals.Marco

	//VERIFICO SI LOS DATOS QUE MANDO CPU TIENEN SENTIDO (O SEA, NO HAY VALORES MAYORES A LOS DE LA CANTIDAD DE NIVELES/ENTRADAS/TAMDEPAGINA QUE TENEMOS DEFINIDOS)
	for i := 1; i <= len(DireccionLogica)-1; i++ { //arrancamos desde 1 porque en 0 esta el desplazamiento, nos fijamos si la entrada nivel n es mayor a la cantidad de entradas por tabla
		log.Printf("DE CPU RECIBO: %d y el numero de entradas total es: %d	(TraducirLogicaAFisica)\n", DireccionLogica[i], globals.ClientConfig.Entries_per_page)

		if DireccionLogica[i] >= globals.ClientConfig.Entries_per_page {
			log.Printf("voy a devolver -1 maestro, no tenemos tantas entradas	(TraducirLogicaAFisica)\n")
			MarcoAurelio.Frame = -1
			return MarcoAurelio
		}
	}

	//SI LLEGAMOS ACA, LO QUE ENVIO CPU TIENE SENTIDO
	var marco globals.Marco

	ActualizarTodasLasTablasEnBaseATablaSimple(DireccionLogica[0])

	//log.Printf("MUESTRO EL PROCESO \n")

	//auxiliares.MostrarProceso(DireccionLogica[0])
	marco.Frame = AccedeAEntrada(DireccionLogica, 0, PunteroNodo)

	MarcoAurelio = marco
	if marco.Frame < 0 {
		log.Printf("como el frame es: %d voy a finalizar el proceso(TraducirLogicaAFisica)\n", MarcoAurelio.Frame)
		FinalizarProceso(DireccionLogica[0])
	}
	log.Printf("frame que devuelvo ante traduccion solicitada: %d	(TraducirLogicaAFisica)\n", MarcoAurelio.Frame)
	return MarcoAurelio

}

/*
Direccion logica tiene la forma de
DireccionLogica[0] = pid
DireccionLogica[1] = entrada nivel 1
...
DireccionLogica[n] = entrada nivel n

Para aumentar expresividad en el codigo (no estar agregando i - 1 en los loops por ejemplo)
*/

func AccedeAEntrada(DireccionLogica []int, nivel int, PunteroNodo *globals.Nodo) int {

	if nivel == globals.ClientConfig.Number_of_levels-1 { //significa que ya estamos parados en el nivel que contiene los marcos
		log.Printf("\n este es el valor que returneo: %d\n", (*PunteroNodo).Marco[DireccionLogica[nivel+1]])
		return ((*PunteroNodo).Marco[DireccionLogica[nivel+1]])

	} else {
		return AccedeAEntrada(DireccionLogica, nivel+1, (*PunteroNodo).Siguiente[DireccionLogica[nivel+1]]) //DESPUES DE 5 HORAS DEBUGGEANDO, AHI IBA UN +1, NO TE OLVIDES NUNCA MAS POR FAVOR, EL PRIMER ELEMENTO DEL SLICE ES EL PID, ESTABAS TOMANDO EL PID COMO LA PRIMERA ENTRADA EN ESTA FUNCION

	}
}

/*
Leer Página completa
Se deberá devolver el contenido correspondiente de la página a partir del byte enviado como dirección física dentro de la Memoria de Usuario, que deberá coincidir con la posición del byte 0 de la página.
Actualizar página completa
Se escribirá la página completa a partir del byte 0 que igual será enviado como dirección física, esta operación se realizará dentro de la Memoria de Usuario y se responderá como OK.
*/

/*
que hace LeerPaginaCompleta?
recibe una direccion fisica (tiene que ser una donde inicie una pagina) y devuelve el contenido de esa pagina en un slice de bytes
*/
func LeerPaginaCompleta(direccion int) ([]byte, error) {

	var pagina []byte

	pagina = make([]byte, globals.ClientConfig.Page_size)

	//reviso varias posibilidades... si el usuario envia una direccion distinta a un inicio de pagina es el obvio pero tambien reviso si de casualidad nos pusieron un tam de memoria no divisible por page size, tambien controlo que no nos envien la ultima direccion de la memoria
	if direccion%globals.ClientConfig.Page_size != 0 || direccion+globals.ClientConfig.Page_size >= globals.ClientConfig.Memory_size {

		log.Printf("ERROR, LA DIRECCION RECIBIDA NO CORRESPONDE A LA DE UN INICIO DE PAGINA \n")
		return pagina, fmt.Errorf("error")

	} else {
		globals.Sem_Mem.Lock()
		for i := 0; i < globals.ClientConfig.Page_size; i++ {
			pagina[i] = globals.MemoriaPrincipal[direccion+i] //vamos recorriendo la pagina en memoria y se la asignamos a la variable que vamos a devolver

		}
		globals.Sem_Mem.Unlock()
		return pagina, nil
	} //esa es la forma de go de devolver errores, no la uso en otras partes porque puedo arreglarme con valores negativos o cosas asi que siento que dejan el codigo mas expresivo, al menos para mi, devuelve dos cosas esta funcion.

}

/*
Actualizar página completa
Se escribirá la página completa a partir del byte 0 que igual será enviado como dirección física,
esta operación se realizará dentro de la Memoria de Usuario y se responderá como OK.
*/

func ActualizarPaginaCompleta(PaginaNueva globals.Pagina, direccion int) {

	//reviso varias posibilidades... si el usuario envia una direccion distinta a un inicio de pagina es el obvio pero tambien reviso si de casualidad nos pusieron un tam de memoria no divisible por page size, tambien controlo que no nos envien la ultima direccion de la memoria
	if direccion%globals.ClientConfig.Page_size != 0 || direccion+globals.ClientConfig.Page_size > globals.ClientConfig.Memory_size {

		log.Printf("ERROR, LA DIRECCION RECIBIDA NO CORRESPONDE A LA DE UN INICIO DE PAGINA \n")

	} else {
		globals.Sem_Mem.Lock()
		for i := 0; i < globals.ClientConfig.Page_size; i++ {
			globals.MemoriaPrincipal[direccion+i] = PaginaNueva.Info[i]
		}
		globals.Sem_Mem.Unlock()

	}
}

/*
que hace AsignarValoresATablaDePaginas

mira los valores que hay en la tabla de paginas simple del proceso y actualiza los de la tabla de paginas de verdad
*/

func AsignarValoresATablaDePaginas(pid int, nivel int, PunteroAux *globals.Nodo) {

	if nivel == globals.ClientConfig.Number_of_levels-1 { //significa que ya estamos parados en el nivel que contiene los marcos
		for j := 0; j < globals.ClientConfig.Entries_per_page; j++ {

			if globals.ContadorTabla < len(globals.MemoriaKernel[pid].TablaSimple) {
				(*PunteroAux).Marco[j] = globals.MemoriaKernel[pid].TablaSimple[globals.ContadorTabla]
				//	log.Printf("llene este valor     %d       , es una de las paginas que tiene, una de la tabla que printie arriba /n", (*PunteroAux).Marco[j])
				globals.ContadorTabla++
			} else {
				return
			}

		}

	} else {

		for i := 0; i < globals.ClientConfig.Entries_per_page; i++ {
			AsignarValoresATablaDePaginas(pid, nivel+1, (*PunteroAux).Siguiente[i])

		}

	}
}

func ActualizarPaginasDisponibles() {
	//globals.Sem_Bitmap.Lock()
	//recorro el map de memoria kernel (donte tenemos la tabla simple de cada proceso, basicamente la posita de que proceso tiene cada pagina sale de ahi)
	globals.Sem_MemoriaKernel.Lock()
	for _, value := range globals.MemoriaKernel { //que hace range? es literalmente un for, en cada iteración, cambia el valor de key y value, los cuales vas a usar para laburar dentro del range tal como si fuera un for con sintaxis media rara.

		for j := 0; j < len(value.TablaSimple); j++ {
			if value.TablaSimple[j] > 0 { //solo entramos si el valor es positivo, si es negativo significa que el proceso finalizo o esta en disco, porque la estructura no la borramos por si vuelve de disco a MP
				globals.PaginasDisponibles[value.TablaSimple[j]] = 1 //ej: si el proceso 4 tiene reservada las paginas [12, 23]  vamos a entrar a las posiciones del array de paginas disponibles y cambiamos la posicion 12 con un 1 y la del 23 con un 1 tambien, dando a entender que estan ocupadas. repetimos eso con todos los procesos que estan en este momento en memoria
			}
		}
	}
	globals.Sem_MemoriaKernel.Unlock()

	//globals.Sem_Bitmap.Unlock()
}

/*
//cambiar 			if contador <= len(globals.MemoriaKernel[pid].TablaSimple) {
//a menor solo
//cambiar las llamadas de las funciones actualizar einicializar de 0 a 1

/*
Que hace LIberarTablASimple?
cambia el valor a -1 de todas las paginas asociadas al proceso la tabla simple del proceso que le pases por parametro
*/
func LiberarTablaSimpleYPagsDisponibles(pid int) {

	for i := 0; i < len(globals.MemoriaKernel[pid].TablaSimple); i++ {
		globals.PaginasDisponibles[globals.MemoriaKernel[pid].TablaSimple[i]] = 0
		globals.MemoriaKernel[pid].TablaSimple[i] = -1
	}
}

/*
Que hace CambiarAMenos1TodasLasTablas

cambia a -1 toda la data relacionada a paginas del proceso, o sea, lo borras/mandas a swap, entonces tenes que llamar a esta funcion poruqe sino queda como si siguiera en memoria
*/
func CambiarAMenos1TodasLasTablas(pid int) {
	globals.Sem_Bitmap.Lock()        //bitmap
	globals.Sem_MemoriaKernel.Lock() //memkernel1

	LiberarTablaSimpleYPagsDisponibles(pid)

	//ActualizarPaginasDisponibles()

	var PunteroAux *globals.Nodo = globals.MemoriaKernel[pid].PunteroATablaDePaginas //es necesario enviar un puntero auxiliar por parametro en esta funcion
	globals.Sem_MemoriaKernel.Unlock()                                               //memkernel1

	// globals.Sem_Contador.Lock()
	globals.ContadorTabla = 0
	AsignarValoresATablaDePaginas(pid, 0, PunteroAux)
	globals.ContadorTabla = 0
	// globals.Sem_Contador.Unlock()
	auxiliares.InicializarSiNoLoEstaMap(pid)

	globals.Sem_Bitmap.Unlock() //bitmap
}

func ActualizarTodasLasTablasEnBaseATablaSimple(pid int) { //no sirve para procesos que fueron quitados de MP
	ActualizarPaginasDisponibles()
	var PunteroAux *globals.Nodo = globals.MemoriaKernel[pid].PunteroATablaDePaginas //es necesario enviar un puntero auxiliar por parametro en esta funcion
	// globals.Sem_Contador.Lock()
	globals.ContadorTabla = 0
	AsignarValoresATablaDePaginas(pid, 0, PunteroAux)
	// globals.Sem_Contador.Unlock()
	globals.ContadorTabla = 0
	auxiliares.InicializarSiNoLoEstaMap(pid)

}

/*
type Metricas struct {
	Pid                              utilsCPU.Proceso `json:"pidparametricas"`
	ContadorAccesosTablaPaginas      int              `json:"accesos"`
	ContadorInstruccionesSolicitadas int              `json:"totalinstr"`
	ContadorBajadasSWAP              int              `json:"bajadasswap"`
	ContadorSubidasAMemoria          int              `json:"subidasmemoria"`
	ContadorReadMemoria              int              `json:"readmemory"`
	ContadorWriteMemoria             int              `json:"writememory"`
}*/

/*
cambia a -1 la info del proceso a finalizar y printea las metricas del proceso
*/
func FinalizarProceso(pid int) {

	if len(globals.MemoriaKernel[pid].TablaSimple) > 0 && globals.MemoriaKernel[pid].TablaSimple[0] != -1 { //si esta en swap o es de tam 0 no entra al if basicamente porque no hace falta
		CambiarAMenos1TodasLasTablas(pid)

	}
	// globals.Sem_Instruccion.Lock()
	globals.Sem_Metricas.Lock()
	log.Printf("## PID: <%d> - Proceso Destruido - Métricas - Acc. T. Pag: <%d>; Inst.Sol.: <%d>; SWAP:<%d>; Mem.Prin.:<%d>; Lec.Mem.: <%d>, Esc.Mem.: <%d>",

		pid,
		globals.MetricasProceso[pid].ContadorAccesosTablaPaginas,
		globals.MetricasProceso[pid].ContadorInstruccionesSolicitadas,
		globals.MetricasProceso[pid].ContadorBajadasSWAP,
		globals.MetricasProceso[pid].ContadorSubidasAMemoria,
		globals.MetricasProceso[pid].ContadorReadMemoria,
		globals.MetricasProceso[pid].ContadorWriteMemoria)
	globals.Sem_Metricas.Unlock()
	// globals.Sem_Instruccion.Unlock()

	//printear metricas
	//llamar una funcion que reinicie las metricas del proceso a 0 por si se crea un proceso con ese pid
}

/*
Que hace MemoryDump?
copia el contenido de todas las paginas del proceso y las pega en un archivo
*/

func MemoryDump(pid int) {

	var bytestotales int = 0
	log.Printf("\n\n INICIANDO MEMORY DUMP (memorydump) \n\n")
	//time stamp ---> timestamp := time.Now().Unix() // Ej: 1686835231     ?????

	var path string = fmt.Sprintf("%s%d-<TIMESTAMP>.dmp", globals.ClientConfig.Dump_path, pid)
	file, err := os.Create(path) //crea archivo para el dump

	if err != nil {
		log.Printf("error al crear el archivo para el dump de pid %d \n", pid)
	}

	buffer := make([]byte, globals.ClientConfig.Page_size) //contiene el contenido de una pagina entera
	globals.Sem_MemoriaKernel.Lock()

	for i := 0; i < len(globals.MemoriaKernel[pid].TablaSimple); i++ {
		for j := 0; j < globals.ClientConfig.Page_size; j++ {
			//buffer[j] = append(buffer, globals.MemoriaPrincipal[((globals.MemoriaKernel[pid].TablaSimple[i])*globals.ClientConfig.Page_size)+j])
			buffer[j] = globals.MemoriaPrincipal[((globals.MemoriaKernel[pid].TablaSimple[i])*globals.ClientConfig.Page_size)+j]
		}
		bytesEscritos, err := file.Write(buffer)
		if err != nil {
			log.Printf("error al escribir en el archivo\n")
		}
		bytestotales += bytesEscritos
	}
	globals.Sem_MemoriaKernel.Unlock()

	log.Printf("%d bytes fueron escritos en el archivo gracias a la syscall de dump \n", bytestotales)

	//auxiliares.MostrarArchivo(path)
	defer file.Close()
}

// TODO COMENTAR ANTES DE TERMINAR PARA EVITAR RACE CONDITION EN CONTADORES
func MostrarTablaMultinivel(pid int, nivel int, PunteroAux *globals.Nodo) {

	if nivel == globals.ClientConfig.Number_of_levels-1 { //significa que ya estamos parados en el nivel que contiene los marcos
		for j := 0; j < globals.ClientConfig.Entries_per_page; j++ {

			if globals.ContadorTabla < len(globals.MemoriaKernel[pid].TablaSimple) {
				log.Printf("%d \t contador: %d\n", (*PunteroAux).Marco[j], globals.Contador)
				globals.ContadorTabla++
				globals.Contador++
			} else {
				return
			}

		}

	} else {

		for i := 0; i < globals.ClientConfig.Entries_per_page; i++ {
			MostrarTablaMultinivel(pid, nivel+1, (*PunteroAux).Siguiente[i])

		}

	}
}
