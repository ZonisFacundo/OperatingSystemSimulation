package globals

import (
	"encoding/json"
	"log"
	"os"

)

type Config struct {
	Ip_cpu            string `json:"ip_cpu"`
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
	Page_size         int    `json:"size_page"`
	Niveles           int    `json:"size_level"`
	Entradas          int    `json:"size_entries"`
}
type Instruccion struct { // instruccion obtenida de memoria
	Pc              int    `json:"pc"`
	Pid             int    `json:"pid"`
	Interrup        bool   `json:"interruption"` //Valores de la interrupción.
	DireccionLog    int    `json:"adress_log"`
	Traducida       bool   `json:"traducida"`
	DireccionFis    int    `json:"adress_fis"` //Para Read and Write -> Dirección lógica que pasa memoria.
	InstructionType string `json:"type"`       //Contexto de la ejecución, es decir, la string que entra en el execute.
	Valor           int    `json:"value"`      //Parámetro para GOTO
	Tamaño          int    `json:"size"`       //Parámetro para el READ e INIT_PROC.
	ArchiInstr      string `json:"archiInstr"`
	Tiempo          int    `json:"time"` //Parámetro para NOOP.
	Datos           string `json:"datos"`
	Syscall         string `json:"syscall"`
	Frame           int    `json:"frame"`
	Desplazamiento  int    `json:"desplazamiento"`
	Dispositivo     string `json:"dispositive"`
	NroPag          int    `json:"page_number"`
	PosicionPag     int    `json:"pos_number"`
}

type TLB struct {
	Entradas       []Entrada
	Tamanio        int
	PosDeReemplazo int
}

type Entrada struct {
	PID          int
	NroPagina    int
	Direccion    int
	UltimoAcceso int64
}

type EntradaCacheDePaginas struct {
	PID             int
	NroPag          int
	Contenido       string
	DireccionFisica int
	Modificada      bool
	BitUso          bool
}

type CacheDePaginas struct {
	Tamanio      int
	Entradas     []EntradaCacheDePaginas
	PosReemplazo int
}

var ID Instruccion
var ClientConfig *Config
var Interruption bool

var Tlb TLB
var CachePaginas CacheDePaginas
var AlgoritmoReemplazo string
var AlgoritmoReemplazoTLB string

func CargarConfig(path string, instanceID string) {

	conjuntodebytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("ATENCION: error al leer config.json. Verificá que estés ejecutando desde el directorio raíz del proyecto.\n")
		return
	}

	var configgenerica Config
	err = json.Unmarshal(conjuntodebytes, &configgenerica)
	if err != nil {
		log.Printf("Error al decodificar JSON -> struct de Go")
		return
	}

	if configgenerica.Instance_id != instanceID {
		log.Fatalf("ERROR: El instance_id en el config.json es '%s', pero ejecutaste la CPU como '%s'.", configgenerica.Instance_id, instanceID)
	}

	ClientConfig = &configgenerica
	log.Printf("Configuración cargada para la instancia '%s'", instanceID)
}

func InitCache() {
	CachePaginas = CacheDePaginas{
		Entradas:     make([]EntradaCacheDePaginas, 0, ClientConfig.Cache_entries),
		Tamanio:      ClientConfig.Cache_entries,
		PosReemplazo: 0,
	}
}

func InitTlb() {
	Tlb = TLB{
		Entradas:       make([]Entrada, 0, ClientConfig.Tlb_entries),
		Tamanio:        ClientConfig.Tlb_entries,
		PosDeReemplazo: 0,
	}
}
