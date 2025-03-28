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

func Fdisk(parametros []string) {
	fmt.Println("FDISK")
	//PARAMETROS: -size -unit -path -type -fit -name
	var size int    //obligatorio si es creacion
	var path string //obligatorio (es el "path", es una letra nombre de la particion, path ya esta fijado)
	var name string //obligatorio Nombre de la particion
	unit := 1024    //opcional /valor por defecto en KB por eso es 1024
	typee := "P"    //opcional Valores: P, E, L
	fit := "W"      //opcional valores para fit: f, w, b

	var opcion int        // 0 -> crear; 1 -> add; 2 -> delete (por defecto es 0 = CREAR)
	paramC := true        //Para validar que los parametros cumplen con los requisitos
	sizeInit := false     //Sirve para saber si se inicializo size (por si no viniera el parametro por ser opcional) false -> no inicializado
	var sizeValErr string //Para reportar el error si no se pudo convertir a entero el size

	//mismo proceso que el mkdisk para manejar parametros
	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		tmp := strings.Split(tmp2, "=")

		//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
		if len(tmp) != 2 {
			fmt.Println("FDISK Error: Valor desconocido del parametro ", tmp[0])
			paramC = false
			break
		}

		//SIZE
		if strings.ToLower(tmp[0]) == "size" {
			sizeInit = true
			var err error
			size, err = strconv.Atoi(tmp[1]) //se convierte el valor en un entero
			if err != nil {
				sizeValErr = tmp[1] //guarda para el reporte del error si es necesario validar size
			}

			//PATH
		} else if strings.ToLower(tmp[0]) == "path" {
			//homonimo al path
			path = tmp[1]
			nombreDisco := strings.Split(path, "/")
			disco := nombreDisco[len(nombreDisco)-1]
			//Se valida si existe el disco ingresado
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("FDISK Error: El disco ", disco, " no existe")
				paramC = false
				break // Terminar el bucle porque encontramos un nombre único
			}

			//NAME
		} else if strings.ToLower(tmp[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(tmp[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)

			//UNIT
		} else if strings.ToLower(tmp[0]) == "unit" {
			//k ya esta predeterminado
			if strings.ToLower(tmp[1]) == "b" {
				//asigno el valor del parametro en su respectiva variable
				unit = 1
			} else if strings.ToLower(tmp[1]) == "m" {
				unit = 1048576 //1024*1024
			} else if strings.ToLower(tmp[1]) != "k" {
				fmt.Println("FDISK Error en -unit. Valores aceptados: b, k, m. ingreso: ", tmp[1])
				paramC = false
				break
			}

			//TYPE
		} else if strings.ToLower(tmp[0]) == "type" {
			//p esta predeterminado
			if strings.ToLower(tmp[1]) == "e" {
				typee = "E"
			} else if strings.ToLower(tmp[1]) == "l" {
				typee = "L"
			} else if strings.ToLower(tmp[1]) != "p" {
				fmt.Println("FDISK Error en -type. Valores aceptados: e, l, p. ingreso: ", tmp[1])
				paramC = false
				break
			}

			//FIT
		} else if strings.ToLower(tmp[0]) == "fit" {
			//Si el ajuste es BF (best fit)
			if strings.ToLower(tmp[1]) == "bf" {
				//asigno el valor del parametro en su respectiva variable
				fit = "B"
				//Si el ajuste es WF (worst fit)
			} else if strings.ToLower(tmp[1]) == "ff" {
				//asigno el valor del parametro en su respectiva variable
				fit = "F"
				//Si el ajuste es ff ya esta definido por lo que si es distinto es un error
			} else if strings.ToLower(tmp[1]) != "wf" {
				fmt.Println("FDISK Error en -fit. Valores aceptados: BF, FF o WF. ingreso: ", tmp[1])
				paramC = false
				break
			}

			//ERROR EN LOS PARAMETROS LEIDOS
		} else {
			fmt.Println("FDISK Error: Parametro desconocido ", tmp[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	//Si va a crear una particion verificar el size
	if opcion == 0 && paramC {
		if sizeInit { //Si viene el parametro size
			if sizeValErr == "" { //Si es un numero (si es numero la variable sizeValErr sera una cadena vacia)
				if size <= 0 { //se valida que sea mayor a 0 (positivo)
					fmt.Println("FDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", size)
					paramC = false
				}
			} else { //Si sizeValErr es una cadena (por lo que no se pudo dar valor a size)
				fmt.Println("FDISK Error: -size debe ser un valor numerico. se leyo ", sizeValErr)
				paramC = false
			}
		} else { //Si no viene el parametro size
			fmt.Println("FDISK Error: No se encuentra el parametro -size")
			paramC = false
		}
	}

	//si todos los parametros son correctos
	if paramC {
		if path != "" && name != "" {
			// Abrir y cargar el disco
			filepath := path
			disco, err := Herramientas.OpenFile(filepath)
			if err != nil {
				fmt.Println("FDisk Error: No se pudo leer el disco")
				return
			}

			//Se crea un mbr para cargar el mbr del disco
			var mbr Structs.MBR
			//Guardo el mbr leido
			if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
				return
			}

			//CREAR (opcion: 0 -> crear; 1 -> add; 2 -> delete)
			if opcion == 0 {

				//Si la particion es tipo extendida validar que no exista alguna extendida
				isPartExtend := false //Indica si se puede usar la particion extendida
				isName := true        //Valida si el nombre no se repite (true no se repite)
				if typee == "E" {
					for i := 0; i < 4; i++ {
						tipo := string(mbr.Partitions[i].Type[:])
						//fmt.Println("tipo ", tipo)
						if tipo != "E" {
							isPartExtend = true
						} else {
							isPartExtend = false
							isName = false //Para que ya no evalue el nombre ni intente hacer nada mas
							fmt.Println("FDISK Error. Ya existe una particion extendida")
							fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
							break
						}
					}
				}

				//verificar si  el nombre existe en las particiones primarias o extendida
				if isName {
					for i := 0; i < 4; i++ {
						nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
						if nombre == name {
							isName = false
							fmt.Println("FDISK Error. Ya existe la particion : ", name)
							fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
							break
						}
					}
				}

				//verificar si existe en las logicas

				//INGRESO DE PARTICIONES PRIMARIAS Y/O EXTENDIDA (SIN LOGICAS)
				sizeNewPart := size * unit //Tamaño de la nueva particion (tamaño * unidades)
				guardar := false           //Indica si se debe guardar la particion, es decir, escribir en el disco
				var newPart Structs.Partition
				if (typee == "P" || isPartExtend) && isName { //para que  isPartExtend sea true, typee tendra que ser "E"
					sizeMBR := int32(binary.Size(mbr)) //obtener el tamaño del mbr (el que ocupa fisicamente: 165)
					//Para manejar los demas ajustes hacer un if del fit para llamar a la funcion adecuada
					//F = primer ajuste; B = mejor ajuste; else -> peor ajuste

					//INSERTAR PARTICION (Primer ajuste)
					mbr, newPart = primerAjuste(mbr, typee, sizeMBR, int32(sizeNewPart), name, fit) //int32(sizeNewPart) es para castear el int a int32 que es el tipo que tiene el atributo en el struct Partition
					guardar = newPart.Size != 0

					//escribimos el MBR en el archivo. Lo que no se llegue a escribir en el archivo (aqui) se pierde, es decir, los cambios no se guardan
					if guardar {
						//sobreescribir el mbr
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {
							return
						}

						//Se agrega el ebr de la particion extendida en el disco
						if isPartExtend {
							var ebr Structs.EBR
							ebr.Start = newPart.Start
							ebr.Next = -1
							if err := Herramientas.WriteObject(disco, ebr, int64(ebr.Start)); err != nil {
								return
							}
						}
						//para verificar que lo guardo
						var TempMBR2 Structs.MBR
						// Read object from bin file
						if err := Herramientas.ReadObject(disco, &TempMBR2, 0); err != nil {
							return
						}
						fmt.Println("\nParticion con nombre " + name + " creada exitosamente")
						Structs.PrintMBR(TempMBR2)
					} else {
						//Lo podría eliminar pero tendria que modificar en el metodo del ajuste todos los errores para que aparezca el nombre que se intento ingresar como nueva particion
						fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
					}

				} //else if para ingreso de particiones logicas
				//a esta altura sigue abierto el archivo

				//------------------------------ADD---------------------

				//--------------------- Eliminar particiones -----------------------------------------------------

			} else {
				//Probablemente nunca entre aqui (se podría quitar)
				fmt.Println("FDISK Error. Operación desconocida (operaciones aceptadas: crear, modificar o eliminar)")
			}
			//Fin operaciones crear, modificar (add) y eliminar

			// Cierro el disco
			defer disco.Close()
			fmt.Println("======End FDISK======")
		} else {
			fmt.Println("FDISK Error. No se encontro parametro letter y/o name")
		}
	} //Fin if paramC
} //Fin FDisk

func primerAjuste(mbr Structs.MBR, typee string, sizeMBR int32, sizeNewPart int32, name string, fit string) (Structs.MBR, Structs.Partition) {
	var newPart Structs.Partition
	var noPart Structs.Partition //para revertir el set info (simula volverla null)

	//PARTICION 1 (libre) - (size = 0 no se ha creado)
	if mbr.Partitions[0].Size == 0 {
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if mbr.Partitions[1].Size == 0 {
			if mbr.Partitions[2].Size == 0 {
				//caso particion 4 (no existe)
				if mbr.Partitions[3].Size == 0 {
					//859 <= 1024 - 165
					if sizeNewPart <= mbr.MbrSize-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
					}
				}
			}
		}
		//Fin de 1 no existe

		//PARTICION 2 (no existe)
	} else if mbr.Partitions[1].Size == 0 {
		//Si no hay espacio antes de particion 1
		newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2) //el nuevo inicio es donde termina 1
		if mbr.Partitions[2].Size == 0 {
			if mbr.Partitions[3].Size == 0 {
				if sizeNewPart <= mbr.MbrSize-newPart.Start {
					mbr.Partitions[1] = newPart
				} else {
					newPart = noPart
					fmt.Println("FDISK Error. Espacio insuficiente")
				}
			}
		}
		//Fin particion 2 no existe

		//PARTICION 3
	} else if mbr.Partitions[2].Size == 0 {
		//despues de 2
		newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
		if mbr.Partitions[3].Size == 0 {
			if sizeNewPart <= mbr.MbrSize-newPart.Start {
				mbr.Partitions[2] = newPart
			} else {
				newPart = noPart
				fmt.Println("FDISK Error. Espacio insuficiente")
			}
		}
		//Fin particion 3

		//PARTICION 4
	} else if mbr.Partitions[3].Size == 0 {
		if sizeNewPart <= mbr.MbrSize-mbr.Partitions[2].GetEnd() {
			//despues de 3
			newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
			mbr.Partitions[3] = newPart
		} else {
			newPart = noPart
			fmt.Println("FDISK Error. Espacio insuficiente")
		}
		//Fin particion 4
	} else {
		newPart = noPart
		fmt.Println("FDISK Error. Particiones primarias y/o extendidas ya no disponibles")
	}

	return mbr, newPart
}
