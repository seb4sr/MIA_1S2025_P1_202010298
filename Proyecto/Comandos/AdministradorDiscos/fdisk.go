package Comandos

import (
	"Proyecto/Herramientas"
	"Proyecto/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Fdisk(entrada []string) string{
	var respuesta string
	 
	 
	unit:=1024 	 
	tipe:="P"	 
	fit :="W"	 
	var size int			 
	var pathE string		 
	var name string			 
	Valido := true         
	InitSize := false      
	InitPath :=false
	var sizeValErr string  
	

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			fmt.Println("ERROR FDISK, valor desconocido de parametros ",valores[1])
			respuesta += "ERROR FDISK, valor desconocido de parametros " + valores[1]+ "\n"
			Valido = false
			 
			return respuesta
		}

		 
		if strings.ToLower(valores[0])=="size"{
			InitSize = true
			var err error
			 
			 
			size, err = strconv.Atoi(valores[1])  
			if err != nil {
				sizeValErr = valores[1]  
			}

		 
		} else if strings.ToLower(valores[0]) == "unit" {
			 
			if strings.ToLower(valores[1]) == "b" {
				unit = 1
				 
			} else if strings.ToLower(valores[1]) == "m" {
				unit = 1048576  
			} else if strings.ToLower(valores[1]) != "k" {
				fmt.Println("FDISK Error en -unit. Valores aceptados: b, k, m. ingreso: ", valores[1])
				Valido = false
				respuesta += "FDISK Error en -unit. Valores aceptados: b, k, m. ingreso: " + valores[1]+ "\n"
				return respuesta
			}

		 
		} else if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1],"\"","")
			InitPath = true	

			_, err := os.Stat(pathE)
			if os.IsNotExist(err) {
				fmt.Println("FDISK Error: El disco no existe")
				respuesta +=  "FDISK Error: El disco no existe"+ "\n"
				Valido = false
				return respuesta  
			}
		
		 
		} else if strings.ToLower(valores[0]) == "type" {
			 
			if strings.ToLower(valores[1]) == "e" {
				tipe = "E"
			} else if strings.ToLower(valores[1]) == "l" {
				tipe = "L"
			} else if strings.ToLower(valores[1]) != "p" {
				fmt.Println("FDISK Error en -type. Valores aceptados: e, l, p. ingreso: ", valores[1])
				respuesta += "FDISK Error en -type. Valores aceptados: e, l, p. ingreso: " + valores[1]+ "\n"
				Valido = false
				return respuesta
			}

		 
		}else if strings.ToLower(valores[0])=="fit"{
			if strings.ToLower(strings.TrimSpace(valores[1]))=="bf"{
				fit = "B"
			}else if strings.ToLower(valores[1])=="ff"{
				fit = "F"
			}else if strings.ToLower(valores[1])!="wf"{
				fmt.Println("EEROR: PARAMETRO FIT INCORRECTO. VALORES ACEPTADO: FF, BF,WF. SE INGRESO:",valores[1])
				respuesta += "EEROR: PARAMETRO FIT INCORRECTO. VALORES ACEPTADO: FF, BF,WF. SE INGRESO:"+valores[1]+ "\n"
				return respuesta
			}
			
			
		 
		} else if strings.ToLower(valores[0]) == "name" {
			 
			name = strings.ReplaceAll(valores[1], "\"", "")
			 
			name = strings.TrimSpace(name)
		
		 
		} else {
			fmt.Println("FDISK Error: Parametro desconocido: ", valores[0])
			respuesta += "FDISK Error: Parametro desconocido: "+ valores[0]+ "\n"
			return respuesta  
		}
		
	}

	

	if InitPath{
		if InitSize{
			if sizeValErr == "" {  
				if size <= 0 {  
					fmt.Println("FDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", size)
					respuesta += "FDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo " + strconv.Itoa(size)+ "\n"
					Valido = false
					return respuesta
				}
			} else {  
				fmt.Println("FDISK Error: -size debe ser un valor numerico. se leyo ", sizeValErr)
				respuesta +="FDISK Error: -size debe ser un valor numerico. se leyo " + sizeValErr+ "\n"
				Valido = false
				return respuesta
			}
		}else{
			fmt.Println("ERROR: FALTO PARAMETRO SIZE")
			respuesta += "ERROR: FALTO PARAMETRO SIZE"+ "\n"
			Valido =false
		}
	}else{
		fmt.Println("ERROR: FALTO PARAMETRO PATH")
		respuesta += "ERROR: FALTO PARAMETRO PATH"+ "\n"
		Valido =false
	}

	if Valido{
		if name != "" {
			 
			disco, err := Herramientas.OpenFile(pathE)
			if err != nil {
				fmt.Println("FDisk Error: No se pudo leer el disco")
				respuesta += "FDisk Error: No se pudo leer el disco"+ "\n"
				return  respuesta
			}

			 
			var mbr Structs.MBR
			 
			if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
				respuesta += "Error Read " + err.Error()+ "\n"
				return  respuesta
			}

			 
			isPartExtend := false  
			isName := true         
			if tipe == "E" {
				for i := 0; i < 4; i++ {
					tipo := string(mbr.Partitions[i].Type[:])
					
					if tipo != "E" {
						isPartExtend = true
					} else {
						isPartExtend = false
						isName = false  
						fmt.Println("FDISK Error. Ya existe una particion extendida")
						fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
						respuesta += "FDISK Error. Ya existe una particion extendida \nFDISK Error. No se puede crear la nueva particion con nombre:  " + name+ "\n"
						return respuesta
					}
				}
			}

			 
			if isName {
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name {
						isName = false
						fmt.Println("FDISK Error. Ya existe la particion : ", name)
						fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
						respuesta += "FDISK Error. Ya existe la particion : " + name + "\nFDISK Error. No se puede crear la nueva particion con nombre: " + name+ "\n"
						return respuesta
					}
				}
			}

			if isName {
				 
				var partExtendida Structs.Partition
				 
				if string(mbr.Partitions[0].Type[:]) == "E" {
					partExtendida = mbr.Partitions[0]
				} else if string(mbr.Partitions[1].Type[:]) == "E" {
					partExtendida = mbr.Partitions[1]
				} else if string(mbr.Partitions[2].Type[:]) == "E" {
					partExtendida = mbr.Partitions[2]
				} else if string(mbr.Partitions[3].Type[:]) == "E" {
					partExtendida = mbr.Partitions[3]
				}

				if partExtendida.Size != 0 {
					var actual Structs.EBR
					if err := Herramientas.ReadObject(disco, &actual, int64(partExtendida.Start)); err != nil {
						respuesta += "Error Read " + err.Error()+ "\n"
						return respuesta
					}

					 
					if Structs.GetName(string(actual.Name[:])) == name {
						isName = false
					} else {
						for actual.Next != -1 {
							 
							if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
								respuesta += "Error Read " + err.Error()+ "\n"
								return respuesta
							}
							if Structs.GetName(string(actual.Name[:])) == name {
								isName = false
								break
							}
						}
					}

					if !isName {
						fmt.Println("FDISK Error. Ya existe la particion : ", name)
						fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
						respuesta += "FDISK Error. Ya existe la particion : " + name
						respuesta += "\nFDISK Error. No se puede crear la nueva particion con nombre: " + name+ "\n"
						return respuesta 
					}
				}
			}

			 
			sizeNewPart := size * unit  
			guardar := false            
			var newPart Structs.Partition
			if (tipe == "P" || isPartExtend) && isName {  
				sizeMBR := int32(binary.Size(mbr))  
				 
				 

				 
				var resTem string
				mbr, newPart, resTem = primerAjuste(mbr, tipe, sizeMBR, int32(sizeNewPart), name, fit)  
				respuesta += resTem
				guardar = newPart.Size != 0

				 
				if guardar {
					 
					if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {
						respuesta += "Error Write " +err.Error()+ "\n"
						return  respuesta
					}

					 
					if isPartExtend {
						var ebr Structs.EBR
						ebr.Start = newPart.Start
						ebr.Next = -1
						if err := Herramientas.WriteObject(disco, ebr, int64(ebr.Start)); err != nil {
							respuesta += "Error Write " +err.Error()+ "\n"
							return  respuesta
						}
					}
					 
					var TempMBR2 Structs.MBR
					 
					if err := Herramientas.ReadObject(disco, &TempMBR2, 0); err != nil {
						respuesta += "Error Read " + err.Error()+ "\n"
						return  respuesta
					}
					fmt.Println("\nParticion con nombre " + name + " creada exitosamente")
					respuesta += "\nParticion con nombre " + name + " creada exitosamente"+ "\n"
					Structs.PrintMBR(TempMBR2)
				} else {
					 
					fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
					respuesta += "FDISK Error. No se puede crear la nueva particion con nombre: "+ name
					return respuesta
				}
			
			 
			}else if tipe == "L" && isName {
				var partExtend Structs.Partition
				if string(mbr.Partitions[0].Type[:]) == "E" {
					partExtend = mbr.Partitions[0]
				} else if string(mbr.Partitions[1].Type[:]) == "E" {
					partExtend = mbr.Partitions[1]
				} else if string(mbr.Partitions[2].Type[:]) == "E" {
					partExtend = mbr.Partitions[2]
				} else if string(mbr.Partitions[3].Type[:]) == "E" {
					partExtend = mbr.Partitions[3]
				} else {
					fmt.Println("FDISK Error. No existe una particion extendida en la cual crear un particion logica")
					respuesta += "FDISK Error. No existe una particion extendida en la cual crear un particion logica"+ "\n"
					return respuesta
				}

				 
				if partExtend.Size != 0 {
					 
					respuesta += primerAjusteLogicas(disco, partExtend, int32(sizeNewPart), name, fit) + "\n" 
					 
				}
			}
			return respuesta
			
		}else{
			respuesta := "ERROR: FALTA PARAMETRO NAME"+ "\n"
			fmt.Println("ERROR: FALTA PARAMETRO NAME")
			return respuesta
			
		}
	}
	return respuesta
}

 
func primerAjuste(mbr Structs.MBR, typee string, sizeMBR int32, sizeNewPart int32, name string, fit string) (Structs.MBR, Structs.Partition, string) {
	var respuesta string
	var newPart Structs.Partition
	var noPart Structs.Partition  

	 
	if mbr.Partitions[0].Size == 0 {
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if mbr.Partitions[1].Size == 0 {
			if mbr.Partitions[2].Size == 0 {
				 
				if mbr.Partitions[3].Size == 0 {
					 
					if sizeNewPart <= mbr.MbrSize-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
					}
				} else {
					 
					 
					if sizeNewPart <= mbr.Partitions[3].Start-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
						}
					}
				}
				 
			} else {
				 
				 
				if sizeNewPart <= mbr.Partitions[2].Start-sizeMBR {
					mbr.Partitions[0] = newPart
				} else {
					 
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
						}
					} else {
						 
						 
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							 
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3  
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							 
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
						}
					}  
				}  
			}  
		} else {
			 
			 
			if sizeNewPart <= mbr.Partitions[1].Start-sizeMBR {
				mbr.Partitions[0] = newPart
			} else {
				 
				 
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[2].Size == 0 {
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
						}
					} else {
						 
						 
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							 
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							if sizeNewPart <= mbr.MbrSize-newPart.Start {  
								mbr.Partitions[2] = mbr.Partitions[3]
								mbr.Partitions[3] = newPart
								 
								mbr.Partitions[2].Correlative = 3
							} else {
								newPart = noPart
								fmt.Println("FDISK Error. Espacio insuficiente")
								respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
							}
						}  
					}  
				} else {
					 
					 
					if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
						mbr.Partitions[0] = mbr.Partitions[1]
						mbr.Partitions[1] = newPart
						 
						mbr.Partitions[0].Correlative = 1
						mbr.Partitions[1].Correlative = 2
					} else if mbr.Partitions[3].Size == 0 {
						 
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
						}
					} else {
						 
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							 
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							 
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
						}
					}  
				}  
			}  
		}  
		 

		 
	} else if mbr.Partitions[1].Size == 0 {
		 
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {  
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			 
			mbr.Partitions[1].Correlative = 2
		} else {
			 
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)  
			if mbr.Partitions[2].Size == 0 {
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[1] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta +="FDISK Error. Espacio insuficiente"+ "\n"
					}
				} else {
					 
					 
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[1] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						 
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
					}
				}  
			} else {
				 
				 
				if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
					mbr.Partitions[1] = newPart
				} else {
					 
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[3].Correlative = 4
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
						}
					} else {
						 
						 
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							 
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							 
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
						}
					}  
				}  
			}  
		}  
		 

		 
	} else if mbr.Partitions[2].Size == 0 {
		 
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			 
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			 
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				 
				mbr.Partitions[2].Correlative = 3
			} else {
				 
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[2] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
					}
				} else {
					 
					 
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[2] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						 
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
					}
				}  
			}  
		}  
		 

		 
	} else if mbr.Partitions[3].Size == 0 {
		 
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[3] = mbr.Partitions[2]
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			 
			mbr.Partitions[3].Correlative = 4
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			 
			 
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				 
				mbr.Partitions[3].Correlative = 4
				mbr.Partitions[2].Correlative = 3
			} else if sizeNewPart <= mbr.Partitions[2].Start-mbr.Partitions[1].GetEnd() {
				 
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = newPart
				 
				mbr.Partitions[3].Correlative = 4
			} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[2].GetEnd() {
				 
				newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
				mbr.Partitions[3] = newPart
			} else {
				newPart = noPart
				fmt.Println("FDISK Error. Espacio insuficiente")
				respuesta += "FDISK Error. Espacio insuficiente"+ "\n"
			}
		}  
		 
	} else {
		newPart = noPart
		fmt.Println("FDISK Error. Particiones primarias y/o extendidas ya no disponibles")
		respuesta += "FDISK Error. Particiones primarias y/o extendidas ya no disponibles"+ "\n"
	}

	return mbr, newPart, respuesta
}

func primerAjusteLogicas(disco *os.File, partExtend Structs.Partition, sizeNewPart int32, name string, fit string) string{
	var respuesta string
	 
	save := true  
	var actual Structs.EBR
	sizeEBR := int32(binary.Size(actual))  
	 

	 
	if err := Herramientas.ReadObject(disco, &actual, int64(partExtend.Start)); err != nil {
		respuesta += "Error Read " + err.Error()+ "\n"
		return respuesta
	}

	 
	 
	 

	 
	if actual.Size == 0 {
		if actual.Next == -1 {
			 
			if sizeNewPart+sizeEBR <= partExtend.Size {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, -1)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " +err.Error()+ "\n"
					return respuesta
				}
				save = false  
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre "+ name+ " creada correctamente"+ "\n"
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas")
				respuesta += "FDISK Error. Espacio insuficiente logicas"+ "\n"
			}
		} else {
			 
			 
			disponible := actual.Next - partExtend.Start  
			if sizeNewPart+sizeEBR <= disponible {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, actual.Next)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " +err.Error()+ "\n"
					return respuesta
				}
				save = false  
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre " + name+ " creada correctamente"+ "\n"
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas 2")
				respuesta += "FDISK Error. Espacio insuficiente logicas"+ "\n"
			}
		}
		 
	}

	if save {
		 
		for actual.Next != -1 {
			 
			if sizeNewPart+sizeEBR <= actual.Next-actual.GetEnd() {
				break
			}
			 
			if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
				respuesta += "Error Read " + err.Error()+ "\n"
				return respuesta
			}

		}

		 
		if actual.Next == -1 {
			 
			if sizeNewPart+sizeEBR <= (partExtend.GetEnd() - actual.GetEnd()) {
				 
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " +err.Error()+ "\n"
					return respuesta
				}

				 
				newStart := actual.GetEnd()                           
				actual.SetInfo(fit, newStart, sizeNewPart, name, -1)  
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " +err.Error()+ "\n"
					return respuesta
				}
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre "+ name+" creada correctamente"+ "\n"
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas 3")
				respuesta += "FDISK Error. Espacio insuficiente logicas"+ "\n"
			}
		} else {
			 
			if sizeNewPart+sizeEBR <= (actual.Next - actual.GetEnd()) {
				siguiente := actual.Next  
				 
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " +err.Error()+ "\n"
					return respuesta
				}

				 
				newStart := actual.GetEnd()                                  
				actual.SetInfo(fit, newStart, sizeNewPart, name, siguiente)  
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " +err.Error()+ "\n"
					return respuesta
				}
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre "+ name +" creada correctamente"+ "\n"
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas 4")
				respuesta+="FDISK Error. Espacio insuficiente logicas "+ "\n"
			}
		}
	}
	return respuesta
}