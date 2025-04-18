package globals

import (
	"encoding/json"
	"log"
	"os"
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
	NombreSyscall string `json:"syscallname"` //no se si necesito esto
	TamProceso    int    `json:"processsize"`
	archivo       string `json:"file"`
	Pid           int    `json:"pid"`
}

//					VARIABLES GLOBALES
/*
este apartado es para poder comunicarnos entre distintos archivos .go (memoria, globals y utils) usando variables globales
*/

var ClientConfig *Config                                                                     //variable global que apunta a un struct que contiene toda la config, despues lo vamos a usar en el main
var MemoriaPrincipal []byte                                                                  //variable donde se guarda la memoria principal
var MemoriaKernel map[int]PaqueteRecibidoMemoriadeKernel                                     // memoria del kernel (donde guardo segmento de codigo basicamente)
var PaqueteInfoProceso *PaqueteRecibidoMemoriadeKernel = new(PaqueteRecibidoMemoriadeKernel) //variable global donde guardo lo que me mande el kernel (info del proceso)
var PaginasDisponibles []int

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
