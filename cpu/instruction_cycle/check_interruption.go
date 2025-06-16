package instruction_cycle

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/utilsCPU"
)

func CheckInterruption() bool {
	
	// Crear un canal para saber cuándo llegó una interrupción
	interruption := make(chan bool)

	// Definir handler con cierre sobre el canal
	mux := http.NewServeMux()
	mux.HandleFunc("/KERNELINTERRUPTION", func(w http.ResponseWriter, r *http.Request) {
		RecibirInterrupcion(w, r)
		interruption <- true // Señal de interrupción recibida
	})

	// Crear el servidor con configuración
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", globals.ClientConfig.Port_cpu),
		Handler: mux,
	}

	// Correr servidor en una goroutine
	go func() {
		log.Println("¿Hay interrupción? Esperando...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error en servidor: %v", err)
		}
	}()

	// Esperar a que llegue una interrupción
	<-interruption

	// Apagar el servidor suavemente
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	// Evaluar si hubo interrupción
	if globals.Interruption {
		globals.Interruption = false
		return true
	}

	return false
}

func RecibirInterrupcion(w http.ResponseWriter, r *http.Request) {

	var request utilsCPU.Interrupcion

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	globals.Interruption = request.Interrup // Señalamos que hay interrupción

	respuesta := utilsCPU.RespuestaKernel{
		Mensaje: "Interrupción recibida.",
	}
	respuestaJSON, err := json.Marshal(respuesta)
	if err != nil {
		http.Error(w, "Error generando respuesta", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respuestaJSON)
}
