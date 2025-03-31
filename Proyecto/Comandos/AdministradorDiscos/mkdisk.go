package Comandos

import (
	"Proyecto/Herramientas"
	"Proyecto/Structs"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

 
func Mkdisk(entrada []string) string{

	var size int				
	var pathE string		
	fit :="F"		
	unit := 1048576	
	Valido := true	
	InitSize := false	
	InitPath := false

	
	for _,parametro :=range entrada[1:]{
		 
		tmp := strings.TrimRight(parametro," ")

		 
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR MKDIS, valor desconocido de parametros ",valores[1])
			Valido = false
			return "ERROR MKDIS, valor desconocido de parametros "+valores[1]
		}
		
		 
		if strings.ToLower(valores[0])=="size"{
			
			InitSize = true
			var err error
			size, err = strconv.Atoi(valores[1])  
			 
			if err != nil {
				fmt.Println("MKDISK Error: -size debe ser un valor numerico. se leyo ", valores[1])
				Valido = false
				return "MKDISK Error: -size debe ser un valor numerico. se leyo "+ valores[1]
			} else if size <= 0 {  
				fmt.Println("MKDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", valores[1])
				Valido = false
				return "MKDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo "+ valores[1]
			}

		 
		}else if strings.ToLower(valores[0])=="fit"{
			if strings.ToLower(valores[0])=="bf"{
				fit = "B"
			}else if strings.ToLower(valores[0])=="wf"{
				fit = "W"
			}else if strings.ToLower(valores[0])!="ff"{
				fmt.Println("EEROR: PARAMETRO FIT INCORRECTO. VALORES ACEPTADO: FF, BF,WF. SE INGRESO: ",valores[1])
			}			
		
		 
		} else if strings.ToLower(valores[0]) == "unit" {
			 
			if strings.ToLower(valores[1]) == "k" {
				 
				unit = 1024
				 
			} else if strings.ToLower(valores[1]) != "m" {
				fmt.Println("MKDISK Error en -unit. Valores aceptados: k, m. ingreso: ", valores[1])
				Valido = false
				return "MKDISK Error en -unit. Valores aceptados: k, m. ingreso: "+valores[1]
			}

		 
		} else if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1],"\"","")
			InitPath = true
			
		 
		} else {
			fmt.Println("MKDISK Error: Parametro desconocido: ", valores[0])
			Valido = false
			return "MKDISK Error: Parametro desconocido: "+ valores[0]  
		}
	}
	
	if Valido{
		if InitSize{
			if InitPath{
				tam := size * unit
				 
				err := Herramientas.CrearDisco(pathE)
				if err != nil {
					fmt.Println("MKDISK Error: ", err)
					return "MKDISK Error: "+err.Error()
				}
				 
				file, err := Herramientas.OpenFile(pathE)
				if err != nil {
					return "MKDISK Error: "+err.Error()
				}

				datos := make([]byte, tam)
				newErr := Herramientas.WriteObject(file, datos, 0)
				if newErr != nil {
					fmt.Println("MKDISK Error: ", newErr)
					return "MKDISK Error: " + newErr.Error()
				}

			 
			ahora := time.Now()
			 
			 
			minutos := ahora.Minute()

			 
			rand.Seed(time.Now().Unix())
			num := rand.Intn(100)

			 
			cad := fmt.Sprintf("%02d%02d", num, minutos)
			 
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
				return "ERROR"
			}

			 
			defer file.Close()

			fmt.Println("\n Se creo el disco de forma exitosa")

			 
			var TempMBR Structs.MBR
			if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
				return "ERROR"
			}
			Structs.PrintMBR(TempMBR)

			fmt.Println("\n======End MKDISK======")

			disco := strings.Split(pathE,"/")
			return "Disco " + disco[len(disco)-1] + " creado "

			}else{
				fmt.Println("ERROR: Debe ingresar el parametro Path")
				return "ERROR: Debe ingresar el parametro Path"
			}
		}else{
			fmt.Println("ERROR: Debe ingresar el parametro Size")
			return "ERROR: Debe ingresar el parametro Size"
		}
	}

	return ""
	
}