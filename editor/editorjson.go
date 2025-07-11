package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func modifyJSONFile(filePath string, io, cpu, kernel, memoria string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo %s: %v", filePath, err)
	}

	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("error al deserializar JSON en %s: %v", filePath, err)
	}

	if _, exists := content["ip_io"]; exists && io != "" {
		content["ip_io"] = io
	}
	if _, exists := content["ip_cpu"]; exists && cpu != "" {
		content["ip_cpu"] = cpu
	}
	if _, exists := content["ip_kernel"]; exists && kernel != "" {
		content["ip_kernel"] = kernel
	}
	if _, exists := content["ip_memory"]; exists && memoria != "" {
		content["ip_memory"] = memoria
	}

	modifiedData, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar JSON en %s: %v", filePath, err)
	}

	if err := ioutil.WriteFile(filePath, modifiedData, 0644); err != nil {
		return fmt.Errorf("error al guardar el archivo %s: %v", filePath, err)
	}

	return nil
}

func main() {
	io := flag.String("io", "", "Valor para ip_io")
	cpu := flag.String("cpu", "", "Valor para ip_cpu")
	kernel := flag.String("kernel", "", "Valor para ip_kernel")
	memoria := flag.String("memoria", "", "Valor para ip_memory")
	flag.Parse()

	if *io == "" || *cpu == "" || *kernel == "" || *memoria == "" {
		fmt.Println("Todos los parámetros (io, cpu, kernel, memoria) son obligatorios.")
		flag.Usage()
		os.Exit(1)
	}

	dirs := []string{
		"/home/utnso/tp-2025-1c-NutriGO/cpu/globals",
		"/home/utnso/tp-2025-1c-NutriGO/io/globals",
		"/home/utnso/tp-2025-1c-NutriGO/kernel/globals",
		"/home/utnso/tp-2025-1c-NutriGO/memoria/globals",
	}

	for _, dir := range dirs {
		files, err := filepath.Glob(filepath.Join(dir, "*.json"))
		if err != nil {
			fmt.Printf("Error al buscar archivos JSON en %s: %v\n", dir, err)
			continue
		}

		for _, filePath := range files {
			fmt.Println("Modificando:", filePath)
			if err := modifyJSONFile(filePath, *io, *cpu, *kernel, *memoria); err != nil {
				fmt.Printf("Error al modificar el archivo %s: %v\n", filePath, err)
			}
		}
	}
}
