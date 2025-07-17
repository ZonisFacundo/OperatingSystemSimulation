package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsIO"
)

func main() {
	// instancias IO
	instanceID := os.Args[1]
	configSuffix := instanceID

	if strings.HasPrefix(instanceID, "IO-") {
		configSuffix = instanceID[3:] // Quita "IO-" las instancias deberian empezar con IO-
	}

	configPath := fmt.Sprintf("./io/globals/config%s.json", configSuffix)

	log.Printf("Usando config: %s para instancia: %s\n", configPath, instanceID)

	//--------------

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("No pude abrir listener dinámico: %v", err)
	}
	defer listener.Close()

	puerto := listener.Addr().(*net.TCPAddr).Port

	globals.CargarConfig(configPath, instanceID)
	nombre := os.Args[1]
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	utilsIO.ConfigurarLogger(nombre) //nombre e instance_ID son lo mismo(no lo toco para no romper nada)
	utilsIO.PeticionClienteIOServidorKERNEL(nombre, globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Ip_io, puerto)
	http.HandleFunc("/KERNELIO", utilsIO.RetornoClienteKERNELServidorIO)

	go func() {
		sig := <-sigs
		log.Printf("Señal recibida: %s", sig)
		utilsIO.NotificarFinalizacionAlKernel(nombre, globals.ClientConfig.Ip_kernel, globals.ClientConfig.Port_kernel, globals.ClientConfig.Ip_io, puerto)
		done <- true
	}()

	go func() {
		addr := listener.Addr().String()
		log.Printf("[IO-%s] Servidor HTTP corriendo en %s", nombre, addr)
		if err := http.Serve(listener, nil); err != nil &&
			!strings.Contains(err.Error(), "use of closed network connection") {
			log.Fatalf("Error al iniciar el servidor HTTP: %s", err)
		}
	}()

	<-done
	log.Println("Finalización controlada del módulo IO.")

}
