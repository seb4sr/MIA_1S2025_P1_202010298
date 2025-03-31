package Comandos

import (
	"Proyecto/Herramientas"
	"Proyecto/Structs"
	"fmt"
	"os"
	"strconv"
	"strings"
)

 
 
func Mount(entrada []string) {
	var name string  
	var path string  
	paramC := true

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR MOUNT, valor desconocido de parametros ", valores[1])
			return  
		}

		 
		if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1], "\"", "")
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("ERROR MOUNT: El disco no existe")
				paramC = false
				break  
			}
			 
		} else if strings.ToLower(valores[0]) == "name" {
			 
			name = strings.ReplaceAll(valores[1], "\"", "")
			 
			name = strings.TrimSpace(name)

			 
		} else {
			fmt.Println("ERROR MOUNT: Parametro desconocido: ", valores[0])
			paramC = false
			break  
		}
	}

	if paramC {
		if path != "" {
			if name != "" {
				 
				disco, err := Herramientas.OpenFile(path)
				if err != nil {
					fmt.Println("ERROR NO SE PUEDE LEER EL DISCO ")
					return
				}

				 
				var mbr Structs.MBR
				 
				if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
					return
				}

				 
				defer disco.Close()

				montar := true  
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name {
						montar = false
						if string(mbr.Partitions[i].Type[:]) != "E" {
							if string(mbr.Partitions[i].Status[:]) != "A" {
								var id string
								var nuevaLetra byte = 'A'  
								contador := 1
								modificada := false  

								 
								for k := 0; k < len(Structs.Pmontaje); k++ {
									if Structs.Pmontaje[k].MPath == path {
										 
										Structs.Pmontaje[k].Cont = Structs.Pmontaje[k].Cont + 1
										contador = int(Structs.Pmontaje[k].Cont)
										nuevaLetra = Structs.Pmontaje[k].Letter
										modificada = true
										break
									}
								}

								if !modificada {
									if len(Structs.Pmontaje) > 0 {
										nuevaLetra = Structs.Pmontaje[len(Structs.Pmontaje)-1].Letter + 1
									}
									Structs.AddPathM(path, nuevaLetra, 1)
								}

								id = "48" + strconv.Itoa(contador) + string(nuevaLetra)  
								fmt.Println("ID:  Letra ", string(nuevaLetra), " cont ", contador)
								 
								Structs.AddMontadas(id, path)

								 
								copy(mbr.Partitions[i].Status[:], "A")
								copy(mbr.Partitions[i].Id[:], id)
								mbr.Partitions[i].Correlative = int32(contador)

								 
								if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {  
									return
								}
								fmt.Println("Particion con nombre ", name, " montada correctamente. ID: ", id)
							} else {
								fmt.Println("ERROR MOUNT. ESTA PARTICION YA FUE MONTADA PREVIAMENTE")
								return
							}
						} else {
							fmt.Println("ERROR MOUNT. No se puede montar una particion extendida")
							return
						}
					}
				}

				if montar {
					fmt.Println("ERROR MOUNT. No se pudo montar la particion ", name)
					fmt.Println("ERROR MOUNT. No se encontro la particion")
					return
				}

			} else {
				fmt.Println("ERROR: FALTA NAME  EN MOUNT")
			}
		} else {
			fmt.Println("ERROR: FALTA PARAMETRO PATH EN MOUNT")
		}
	}
}
