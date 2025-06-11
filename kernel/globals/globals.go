package globals

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Ip_memory               string  `json:"ip_memory"`
	Port_memory             int     `json:"port_memory"`
	Ip_Kernel               string  `json:"ip_Kernel"`
	Port_kernel             int     `json:"port_kernel"`
	Scheduler_algorithm     string  `json:"scheduler_algorithm"`
	Ready_ingress_algorithm string  `json:"Ready_ingress_algorithm"`
	Alpha                   float32 `json:"alpha"`
	Initial_estimate        float32 `json:"initial_estimate"`
	Suspension_time         int     `json:"suspension_time"`
	Log_level               string  `json:"log_level"`
}

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
		log.Printf("Error al decodificar datos json -> golang")
		return
	}

	ClientConfig = &configgenerica //hacemos que nuestro puntero (variable global) apunte a donde guardamos los datos

}
