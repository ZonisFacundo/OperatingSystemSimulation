package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsIO"
)

func main() {
	globals.CargarConfig("./io/globals/config.json")
	nombre := os.Args[1]
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	utilsIO.ConfigurarLogger(nombre)
	utilsIO.PeticionClienteIOServidorKERNEL(nombre, globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Ip_io, globals.ClientConfig.Port_io)
	http.HandleFunc("/KERNELIO", utilsIO.RetornoClienteKERNELServidorIO)

	go func() {
		sig := <-sigs
		log.Printf("Señal recibida: %s", sig)
		utilsIO.NotificarFinalizacionAlKernel(nombre, globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Ip_io, globals.ClientConfig.Port_io)
		done <- true
	}()

	go func() {
		//log.Printf("Servidor IO corriendo en puerto 8003.")
		if err := http.ListenAndServe(":8003", nil); err != nil {
			log.Fatalf("Error al iniciar el servidor HTTP: %s", err)
		}
	}()

	<-done
	log.Println("Finalización controlada del módulo IO.")

}
