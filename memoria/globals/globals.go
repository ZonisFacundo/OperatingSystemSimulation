package globals

import (
	"encoding/json"
	"log"
	"os"

	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

// STRUCTS
type Config struct {
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
	Instrucciones []string `json:"instructions"`
	TablaSimple   []int    `json:"tablasimple"` //basicamente la tabla de paginas simple para el proceso...
}

type Nodo struct {
	Siguiente []*Nodo `json:"node"`
	Marco     []int   `json:"frame"`
}
type DireccionFisica struct {
	Marco          int `json:"frame"`
	Desplazamiento int `json:"offset"`
	Direccion      int `json:"address"`
}
type PaqueteWrite struct {
	Direccion int  `json:"address"`
	Contenido byte `json:"content"`
}
type DireccionLogica struct {
	Ip        string `json:"ip"`
	Puerto    int    `json:"port"`
	DirLogica []int  `json:"dir_logica"`
}

type BytePaquete struct {
	Info byte `json:"info"`
}
type Pagina struct {
	Info []byte `json:"info"`
}

var Instruction *utilsCPU.Proceso

//					VARIABLES GLOBALES
/*
este apartado es para poder comunicarnos entre distintos archivos .go (memoria, globals y utils) usando variables globales
*/

var ClientConfig *Config                                                    //variable global que apunta a un struct que contiene toda la config, despues lo vamos a usar en el main
var MemoriaPrincipal []byte                                                 //variable donde se guarda la memoria principal
var MemoriaKernel map[int]ProcesoEnMemoria = make(map[int]ProcesoEnMemoria) // memoria del kernel (donde guardo segmento de codigo basicamente) y paginas reservadas para cada proceso con pid como key
var PaginasDisponibles []int                                                //nos indica el estado de cada pagina, ocupada o libre
var PunteroBase *Nodo = nil

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
