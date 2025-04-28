package globals

import (
	"encoding/json"
	"log"
	"os"

	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

type Config struct {
	Port_cpu          int    `json:"port_cpu"`
	Ip_memory         string `json:"ip_memory"`
	Port_memory       int    `json:"port_memory"`
	Ip_kernel         string `json:"ip_kernel"`
	Port_kernel       int    `json:"port_kernel"`
	Tlb_entries       int    `json:"tlb_entries"`
	Tlb_replacement   string `json:"tlb_replacement"`
	Cache_entries     int    `json:"cache_entries"`
	Cache_replacement string `json:"cache_replacement"`
	Cache_delay       int    `json:"cache_delay"`
	Log_level         string `json:"log_level"`
	Instance_id       string `json:"instance_id"`   
}
type Instruccion struct { // instruccion obtenida de memoria
	ProcessValues   utilsCPU.Proceso      `json:"instruction"`  //Valores de PID y PC
	Interrup        utilsCPU.Interrupcion `json:"interruption"` //Valores de la interrupción.
	Direccion       int                   `json:"adress"`       //Para Read and Write -> Dirección lógica que pasa memoria.
	InstructionType string                `json:"message"`      //Contexto de la ejecución, es decir, la string que entra en el execute.
	Valor           *int                  `json:"value"`        //Parámetro para GOTO
	Tamaño          *int                  `json:"size"`         //Parámetro para el READ e INIT_PROC.
	Tiempo          *int                  `json:"time"`         //Parámetro para NOOP.
	Datos           *string               `json:"datos"`        //Parámetro para el WRITE.
}

var Instruction utilsCPU.Proceso
var ID Instruccion
var ClientConfig *Config

func CargarConfig(path string) {

	conjuntodebytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("ATENCION \n ATENCION\n error al recibir los datos del .json \n esto ocurre porque NO ESTAS EJECUTANDO EL PROYECTO DESDE EL DIRECTORIO CORRESPONDIENTE \n, el path que recibe el cargar config espera que ejecutes los programas desde el directorio ~/tp-2025-1c-NutriGO, seguramente lo estas haciendo desde ~/tp-2025-1c-NutriGO/nombredelmodulo\n")
		return
	}

	var configgenerica Config
	err = json.Unmarshal(conjuntodebytes, &configgenerica) //traducimos de .json a go digamosle
	if err != nil {
		log.Printf("Error al decodificar datos JSON -> GOLANG")
		return
	}

	ClientConfig = &configgenerica //hacemos que nuestro puntero (variable global) apunte a donde guardamos los datos

}
