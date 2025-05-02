package main

import (
	"log"
	"net/http"
	"time"

	"os"

	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsIO"
)

func main() {
	globals.CargarConfig("./io/globals/config.json")
	nombre := os.Args[1]

	utilsIO.ConfigurarLogger(nombre)
	println(nombre)
	go func() {
		time.Sleep(4 * time.Second)

		utilsIO.PeticionClienteIOServidorKERNEL(nombre, globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Ip_io, globals.ClientConfig.Port_io)
	}()

	http.HandleFunc("/KERNELIO", utilsIO.RetornoClienteKERNELServidorIO)
	log.Printf("Servidor IO corriendo.\n")
	http.ListenAndServe(":8003", nil)

}
