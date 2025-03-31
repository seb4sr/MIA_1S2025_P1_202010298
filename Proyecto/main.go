package main

import (
	Comandos "Proyecto/Comandos"
	DM "Proyecto/Comandos/AdministradorDiscos"  
	FS "Proyecto/Comandos/SistemaDeArchivos"    
	"encoding/json"
	"net/http"

	"bufio"
	"fmt"
	"os"
	"strings"

	 
	"github.com/rs/cors"
)

type Entrada struct {
	Text string `json:"text"`
}

type StatusResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func main() {
	 
	http.HandleFunc("/analizar", getCadenaAnalizar)

	 
	 
	c := cors.Default()

	 
	handler := c.Handler(http.DefaultServeMux)

	 
	fmt.Println("Servidor escuchando en http: //localhost:8080")
	http.ListenAndServe(":8080", handler)

}



func getCadenaAnalizar(w http.ResponseWriter, r *http.Request) {

	var respuesta string  
	 
	w.Header().Set("Content-Type", "application/json")

	var status StatusResponse
	 
	if r.Method == http.MethodPost {
		var entrada Entrada
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
			http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
			status = StatusResponse{Message: "Error al decodificar JSON", Type: "unsucces"}
			json.NewEncoder(w).Encode(status)
			return
		}

		 
		lector := bufio.NewScanner(strings.NewReader(entrada.Text))
		 
		for lector.Scan() {
			 
			if lector.Text() != "" {
				 
				linea := strings.Split(lector.Text(), "#")  
				if len(linea[0]) != 0 {
					 
					 
					respuesta += "***********************************************************************************************************\n"
					respuesta += "Comando en ejecucion: " + linea[0] + "\n"
					respuesta += analizar(linea[0]) + "\n"
				}
				 
				if len(linea) > 1 && linea[1] != "" {
					fmt.Println("#" + linea[1] + "\n")
					respuesta += "#" + linea[1] + "\n"
				}
			}

		}

		 
		w.WriteHeader(http.StatusOK)

		status = StatusResponse{Message: respuesta, Type: "succes"}
		json.NewEncoder(w).Encode(status)

	} else {
		 
		status = StatusResponse{Message: "Metodo no permitido", Type: "unsucces"}
		json.NewEncoder(w).Encode(status)
	}
}

func analizar(entrada string) string {
	 
	 

	respuesta := ""

	 
	tmp := strings.TrimRight(entrada, " ")
	parametros := strings.Split(tmp, " -")

	 

	 

	 
	if strings.ToLower(parametros[0]) == "mkdisk" {
		 
		 
		if len(parametros) > 1 {
			 
			respuesta = DM.Mkdisk(parametros)
		} else {
			fmt.Println("MKDISK ERROR: parametros no encontrados")
			respuesta = "MKDISK ERROR: parametros no encontrados"
		}

	} else if strings.ToLower(parametros[0]) == "fdisk" {
		 
		if len(parametros) > 1 {
			DM.Fdisk(parametros)
		} else {
			fmt.Println("FDISK ERROR: parametros no encontrados")
		}
	} else if strings.ToLower(parametros[0]) == "mount" {
		 
		if len(parametros) > 1 {
			DM.Mount(parametros)
		} else {
			fmt.Println("FDISK ERROR: parametros no encontrados")
		}

		 
	} else if strings.ToLower(parametros[0]) == "mkfs" {
		 
		if len(parametros) > 1 {
			FS.Mkfs(parametros)
		} else {
			fmt.Println("MKFS ERROR: parametros no encontrados")
		}
		 
	} else if strings.ToLower(parametros[0]) == "rep" {
		 
		if len(parametros) > 1 {
			Comandos.Rep(parametros)
			 
		} else {
			fmt.Println("REP ERROR: parametros no encontrados")
		}

	} else if strings.ToLower(parametros[0]) == "exit" {
		fmt.Println("Salida exitosa")
		os.Exit(0)

	} else if strings.ToLower(parametros[0]) == "" {
		 
	} else {
		fmt.Println("Comando no reconocible")
	}

	return respuesta
}