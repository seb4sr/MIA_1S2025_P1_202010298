package Comandos

import (
	"Proyecto/Herramientas"
	"Proyecto/Structs"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// mkdisk -Size=3000 -unit=K -path=/home/user/Disco1.mia​
func Mkdisk(parametros []string) string {
	salida := ""
	fmt.Println("MKDISK")
	
	var size int      
	var path string   
	fit := "F"        
	unit := 1048576   
	paramC := true    
	sizeInit := false 
	pathInit := false 

	
	for _, parametro := range parametros[1:] {

		tmp2 := strings.TrimRight(parametro, " ")

		
		tmp := strings.Split(tmp2, "=")

		
		if len(tmp) != 2 {
			fmt.Println("MKDISK Error: Valor desconocido del parametro ", tmp[0])
			paramC = false
			
			return "Valor desconocido del parametro"

		}
		
		
		if strings.ToLower(tmp[0]) == "size" {
			sizeInit = true
			var err error
			size, err = strconv.Atoi(tmp[1]) 
			
			if err != nil {
				fmt.Println("MKDISK Error: -size debe ser un valor numerico. se leyo ", tmp[1])
				paramC = false
				break
			} else if size <= 0 { 
				fmt.Println("MKDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", tmp[1])
				paramC = false
				break
			}
			
		} else if strings.ToLower(tmp[0]) == "fit" {

			if strings.ToLower(tmp[1]) == "bf" {

				fit = "B"

				} else if strings.ToLower(tmp[1]) == "wf" {

					fit = "W"

					} else if strings.ToLower(tmp[1]) != "ff" {
				fmt.Println("MKDISK Error en -fit. Valores aceptados: BF, FF o WF. ingreso: ", tmp[1])
				paramC = false
				break
			}

			} else if strings.ToLower(tmp[0]) == "unit" {

				if strings.ToLower(tmp[1]) == "k" {

					unit = 1024

			} else if strings.ToLower(tmp[1]) != "m" {
				fmt.Println("MKDISK Error en -unit. Valores aceptados: k, m. ingreso: ", tmp[1])
				paramC = false
				break
			}

			} else if strings.ToLower(tmp[0]) == "path" {
			pathInit = true
			path = tmp[1]

			} else {
			fmt.Println("MKDISK Error: Parametro desconocido: ", tmp[0])
			paramC = false
			break 
		}
	}

	if paramC {

		if sizeInit && pathInit {
			tam := size * unit


			nombreDisco := strings.Split(path, "/")
			disco := nombreDisco[len(nombreDisco)-1]

			err := Herramientas.CrearDisco(path)
			if err != nil {
				fmt.Println("MKDISK Error:: ", err)
			}

			file, err := Herramientas.OpenFile(path)
			if err != nil {
				return "error"
			}

			datos := make([]byte, tam)
			newErr := Herramientas.WriteObject(file, datos, 0)
			if newErr != nil {
				fmt.Println("MKDISK Error: ", newErr)
				return "error"
			}

			ahora := time.Now()
			segundos := ahora.Second()
			minutos := ahora.Minute()
			cad := fmt.Sprintf("%02d%02d", segundos, minutos)
			idTmp, err := strconv.Atoi(cad)
			if err != nil {
				fmt.Println("MKDISK Error: no converti fecha en entero para id")
			}
			
			var newMBR Structs.MBR
			newMBR.MbrSize = int32(tam)
			newMBR.Id = int32(idTmp)
			copy(newMBR.Fit[:], fit)
			copy(newMBR.FechaC[:], ahora.Format("02/01/2006 15:04"))
			
			if err := Herramientas.WriteObject(file, newMBR, 0); err != nil {
				return "Error al escribir el MBR"
			}

			
			defer file.Close()

			fmt.Println("\n Se creo el disco ", disco, " de forma exitosa")

			//imprimir el disco creado para validar que todo este correcto
			var TempMBR Structs.MBR
			if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
				return "error"
			}
			Structs.PrintMBR(TempMBR)

			fmt.Println("\n======End MKDISK======")

		} else {
			fmt.Println("MKDISK Error: Falta parametro -size")
		}
	}
	return salida
}
