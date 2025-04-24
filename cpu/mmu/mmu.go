package mmu

/*import (
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

type MMU struct {
	ProcesoActual       *utilsCPU.Instruccion
	Niveles             int
	TamPagina           int
	cant_entradas_tabla int
	TablasPaginas       map[int]int
}

1) Definir las estructuras de datos para la tabla de paginas
2) Traducir direccion logica a fisica
3)Gestionar TLB "se implementará una TLB para agilizar la traducción de las direcciones lógicas a direcciones físicas"
La TLB contará con la siguiente estructura base: [ página | marco ]


func TraducirDireccion(mmu *MMU, direccionLogica int) {
	nroPagina := floor(direccionLogica / mmu.TamPagina)
	entrada_nivel_X = floor(nro_página  / cant_entradas_tabla ^ (mmu.Niveles - X)) % cant_entradas_tabla   //X??
	desplazamiento := direccionLogica % mmu.TamPagina
	...
	}
*/
