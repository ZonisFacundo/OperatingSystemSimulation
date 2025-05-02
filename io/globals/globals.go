package globals

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Ip_io       string `json:"ip_io"`
	Port_io     int    `json:"port_io"`
	Ip_kernel   string `json:"ip_kernel"`
	Port_kernel int    `json:"port_kernel"`
	Log_level   string `json:"log_level"`
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
