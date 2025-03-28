package Comandos

import (
	"Proyecto/Herramientas"
	"Proyecto/Structs"
	"fmt"
	"path/filepath"
	"strings"
)

//id -> 48 (dos ultimos digitos de su carnet) + correlativo particion + letra correlativa del disco
//EJ: 481A, 482A, 483A, 484A -> se obtiene en el mount

func Rep(parametros []string) {
	fmt.Println("REP")
	var name string //obligatorio Nombre del tipo de reporte a generar
	var path string //obligatorio Nombre que tendrá el reporte
	var id string   //obligatorio sera el del disco o el de la particion
	//var ruta string //opcional para file y ls
	paramC := true //valida que todos los parametros sean correctos

	for _, parametro := range parametros[1:] {
		//quito los espacios en blano despues de cada parametro
		tmp2 := strings.TrimRight(parametro, " ")
		//divido cada parametro entre nombre del parametro y su valor # -size=25 -> -size, 25
		tmp := strings.Split(tmp2, "=")

		//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
		if len(tmp) != 2 {
			fmt.Println("REP Error: Valor desconocido del parametro ", tmp[0])
			paramC = false
			break //para finalizar el ciclo for con el error y no ejecutar lo que haga falta
		}

		if strings.ToLower(tmp[0]) == "name" {
			name = strings.ToLower(tmp[1])
		} else if strings.ToLower(tmp[0]) == "path" {
			// Eliminar comillas
			name = strings.ReplaceAll(tmp[1], "\"", "")
			path = name
		} else if strings.ToLower(tmp[0]) == "id" {
			id = strings.ToUpper(tmp[1]) //Mayusculas para tratarlo como case insensitive
		} else if strings.ToLower(tmp[0]) == "ruta" {
			//ruta = strings.ToLower(tmp[1])
		} else {
			fmt.Println("REP Error: Parametro desconocido: ", tmp[0])
			paramC = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if paramC {
		if name != "" && id != "" && path != "" {
			switch name {
			case "mbr":
				fmt.Println("reporte mbr")
				mbr(path, id)
			case "disk":
				fmt.Println("reporte disk")
			default:
				fmt.Println("REP Error: Reporte ", name, " desconocido")
			}
		} else {
			fmt.Println("REP Error: Faltan parametros")
		}
	}
}

func mbr(path string, id string) {
	var pathDico string
	existe := false

	//BUsca en struck de particiones montadas el id ingresado
	for _, montado := range Structs.Montadas {
		if montado.Id == id {
			pathDico = montado.PathM
			existe = true
			break
		}
	}

	//if true { //para probar los reporte hayan o no particiones montadas
	if existe {
		//Reporte
		tmp := strings.Split(path, "/") // /dir1/dir2/reporte
		nombreReporte := strings.Split(tmp[len(tmp)-1], ".")[0]

		//Disco a reportar
		tmp = strings.Split(pathDico, "/")
		disco := strings.Split(tmp[len(tmp)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return
		}

		// Close bin file
		defer file.Close()

		//reporte graphviz (cad es el contenido del reporte)
		//mbr
		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='SlateBlue' COLSPAN=\"2\"> Reporte MBR </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_tamano </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", mbr.MbrSize)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#AFA1D1'> mbr_fecha_creacion </td> \n  <td bgcolor='#AFA1D1'> %s </td> \n </tr> \n", string(mbr.FechaC[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_disk_signature </td> \n  <td bgcolor='Azure'> %d </td> \n </tr>  \n", mbr.Id)
		cad += Structs.RepGraphviz(mbr, file)
		cad += "</table> > ]\n}"

		//reporte requerido
		carpeta := filepath.Dir(path)
		rutaReporte := "." + carpeta + "/" + nombreReporte + ".dot"

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombreReporte)
		fmt.Println(" Reporte MBR del disco " + disco + " creado exitosamente")
	} else {
		fmt.Println("REP Error: Id no existe")
	}
}
