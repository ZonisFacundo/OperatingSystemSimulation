package instruction_cycle

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

func CheckInterruption() bool {

	//Hacer conexión http con Kernel para que nos mande las interrupciones y verificar ésto.

	if globals.Interruption {

		globals.Interruption = false
		return true
	}

	return false
}
