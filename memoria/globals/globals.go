package globals

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

// STRUCTS
type Config struct {
	Ip_memory        string `json:"ip_memory"`
	Port_memory      int    `json:"port_memory"`
	Memory_size      int    `json:"memory_size"`
	Page_size        int    `json:"page_size"`
	Entries_per_page int    `json:"entries_per_page"`
	Number_of_levels int    `json:"number_of_levels"`
	Memory_delay     int    `json:"memory_delay"`
	Swapfile_path    string `json:"swapfile_path"`
	Swap_delay       int    `json:"swap_delay"`
	Log_level        string `json:"log_level"`
	Dump_path        string `json:"dump_path"`
}

type PaqueteRecibidoMemoriadeKernel struct {
	Pid        int    `json:"pid"`
	TamProceso int    `json:"tamanioProceso"`
	Archivo    string `json:"file"`
}

type ProcesoEnMemoria struct {
	Instrucciones          []string `json:"instructions"`
	TablaSimple            []int    `json:"tablasimple"` //basicamente la tabla de paginas simple para el proceso...
	PunteroATablaDePaginas *Nodo    `json:"tabladepaginas"`
	SwapStart              int64    `json:"swapstart"` //POSIBILIDAD DE ERROR
	SwapTam                int      `json:"swaptam"`
}

// le damos mas bits al int porque puede ser largo el n de byte
type Nodo struct {
	Siguiente []*Nodo `json:"node"`
	Marco     []int   `json:"frame"`
}
type Marco struct {
	Frame int `json:"frame"`
}
type DireccionFisica struct {
	Direccion int `json:"adress"`
	Tamaño    int `json:"value"`
}

type PaqueteWrite struct {
	Direccion int    `json:"adress"`
	Contenido string `json:"content"`
}
type DireccionLogica struct {
	Ip        string `json:"ip"`
	Puerto    int    `json:"port"`
	DirLogica []int  `json:"dir_logica"`
}

type BytePaquete struct {
	Info []byte `json:"info"`
}
type Pagina struct {
	Info []byte `json:"info"`
}

var Instruction *utilsCPU.Proceso

type Metricas struct {
	Pid                              utilsCPU.Proceso `json:"pidparametricas"`
	ContadorAccesosTablaPaginas      int              `json:"accesos"`
	ContadorInstruccionesSolicitadas int              `json:"totalinstr"`
	ContadorBajadasSWAP              int              `json:"bajadasswap"`
	ContadorSubidasAMemoria          int              `json:"subidasmemoria"`
	ContadorReadMemoria              int              `json:"readmemory"`
	ContadorWriteMemoria             int              `json:"writememory"`
}

//					VARIABLES GLOBALES
/*
este apartado es para poder comunicarnos entre distintos archivos .go (memoria, globals y utils) usando variables globales
*/

var ClientConfig *Config                                                    //variable global que apunta a un struct que contiene toda la config, despues lo vamos a usar en el main
var MemoriaPrincipal []byte                                                 //variable donde se guarda la memoria principal
var MemoriaKernel map[int]ProcesoEnMemoria = make(map[int]ProcesoEnMemoria) // memoria del kernel (donde guardo segmento de codigo basicamente) y paginas reservadas para cada proceso con pid como key
var PaginasDisponibles []int                                                //nos indica el estado de cada pagina, ocupada o libre
var PunteroBase *Nodo = nil
var PaginasSwap []int = make([]int, 0)                          //la idea es que se vea algo asi por ejemplo: PaginasSwap = [2, 2, 5] --> significa que el primer marco de pagina del swap y el segundo los ocupan paginas del proceso 2, la tercera del proceso 5, si se deswappea el proceso 2 quedaria --> [-1, -1, 5]
var Contador int = 0                                            //solo para debug, ni proteger porque no va a servir de nada el dia de la entrega
var ContadorTabla int = 0                                       //lo uso para contar donde estamos parados en la tabla de paginas global (la del map del proceso)
var MetricasProceso map[int]*Metricas = make(map[int]*Metricas) //Uso éste map para guardar las metricas por proceso, ¿debería inicializarlo en 1?

//SEMAFOROS

var Sem_Swap sync.Mutex
var Sem_Mem sync.Mutex
var Sem_Bitmap sync.Mutex
var Sem_Contador sync.Mutex
var Sem_Instruccion sync.Mutex
var Sem_MemoriaKernel sync.Mutex
var Sem_Metricas sync.Mutex

// FUNCIONES
func CargarConfig(path string) {

	conjuntodebytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("ATENCION \n ATENCION\n error al recibir los datos del .json \n esto ocurre porque NO ESTAS EJECUTANDO EL PROYECTO DESDE EL DIRECTORIO CORRESPONDIENTE \n, el path que recibe el cargar config espera que ejecutes los programas desde el directorio ~/tp-2025-1c-NutriGO, seguramente lo estas haciendo desde ~/tp-2025-1c-NutriGO/nombredelmodulo\n")
		return
	}

	var configgenerica Config
	err = json.Unmarshal(conjuntodebytes, &configgenerica) //traducimos de .json a go digamosle
	if err != nil {
		log.Printf("error al decodificar datos json -> golang")
		return
	}

	ClientConfig = &configgenerica //hacemos que nuestro puntero (variable global) apunte a donde guardamos los datos

}
