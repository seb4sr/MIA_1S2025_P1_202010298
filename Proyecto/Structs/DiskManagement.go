package Structs

import (
	"Proyecto/Herramientas"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// NOTA: Recordar que los atributos de los struct deben iniciar con mayuscula
type MBR struct {
	MbrSize    int32        //mbr_tamano
	FechaC     [16]byte     //mbr_fecha_creacion
	Id         int32        //mbr_dsk_signature (random de forma unica)
	Fit        [1]byte      //dsk_fit
	Partitions [4]Partition //mbr_partitions
}

type Partition struct {
	Status      [1]byte  //part_status
	Type        [1]byte  //part_type
	Fit         [1]byte  //part_fit
	Start       int32    //part_start
	Size        int32    //part_s
	Name        [16]byte //part_name
	Correlative int32    //part_correlative
	Id          [4]byte  //part_id
}

// Setear valores de la particion
func (p *Partition) SetInfo(newType string, fit string, newStart int32, newSize int32, name string, correlativo int32) {
	p.Size = newSize
	p.Start = newStart
	p.Correlative = correlativo
	copy(p.Name[:], name)
	copy(p.Fit[:], fit)
	copy(p.Status[:], "I")
	copy(p.Type[:], newType)
}

// Metodos de Partition
func GetName(nombre string) string {
	posicionNulo := strings.IndexByte(nombre, 0)
	//Si posicionNulo retorna -1 no hay bytes nulos
	if posicionNulo != -1 {
		//guarda la cadena hasta el primer byte nulo (elimina los bytes nulos)
		nombre = nombre[:posicionNulo]
	}
	return nombre
}

func GetId(nombre string) string {
	//si existe id, no contiene bytes nulos
	posicionNulo := strings.IndexByte(nombre, 0)
	//si posicionNulo  no es -1, no existe id.
	if posicionNulo != -1 {
		nombre = "-"
	}
	return nombre
}

func (p *Partition) GetEnd() int32 {
	return p.Start + p.Size
}

type EBR struct {
	Status [1]byte //part_mount (si esta montada)
	Type   [1]byte
	Fit    [1]byte  //part_fit
	Start  int32    //part_start
	Size   int32    //part_s
	Name   [16]byte //part_name
	Next   int32    //part_next
}

func (e *EBR) SetInfo(fit string, newStart int32, newSize int32, name string, newNext int32) {
	e.Size = newSize
	e.Start = newStart
	e.Next = newNext
	copy(e.Name[:], name)
	copy(e.Fit[:], fit)
	copy(e.Status[:], "I")
	copy(e.Type[:], "L")
}

func (e *EBR) GetEnd() int32 {
	return e.Start + e.Size + int32(binary.Size(e))
}

func GetIdMount(data Mount) string {
	return data.MPath
}

// Reportes de los Structs
func PrintMBR(data MBR) {
	fmt.Println("\n     Disco")
	fmt.Printf("CreationDate: %s, fit: %s, size: %d, id: %d\n", string(data.FechaC[:]), string(data.Fit[:]), data.MbrSize, data.Id)
	for i := 0; i < 4; i++ {
		fmt.Printf("Partition %d: %s, %s, %d, %d, %s, %d\n", i, string(data.Partitions[i].Name[:]), string(data.Partitions[i].Type[:]), data.Partitions[i].Start, data.Partitions[i].Size, string(data.Partitions[i].Fit[:]), data.Partitions[i].Correlative)
	}
}

func PrintEbr(data EBR) {
	fmt.Println("part_status ", string(data.Status[:]))
	fmt.Println("part_type ", string(data.Type[:]))
	fmt.Println("part_fit: ", string(data.Fit[:]))
	fmt.Println("part_start: ", data.Start)
	fmt.Println("part_s ", data.Size)
	fmt.Println("part_name: ", string(data.Name[:]))
	fmt.Println("next_part: ", data.Next)
}

func RepGraphviz(data MBR, disco *os.File) string {
	disponible := int32(0)
	cad := ""
	inicioLibre := int32(binary.Size(data)) //Para ir guardando desde donde hay espacio libre despues de cada particion
	for i := 0; i < 4; i++ {
		if data.Partitions[i].Size > 0 {

			disponible = data.Partitions[i].Start - inicioLibre
			inicioLibre = data.Partitions[i].Start + data.Partitions[i].Size

			//reporta si hay espacio libre antes de la particion
			if disponible > 0 {
				cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#808080' COLSPAN=\"2\"> ESPACIO LIBRE <br/> %d bytes </td> \n </tr> \n", disponible)
			}
			//Reporta el contenido de la particion
			cad += " <tr>\n  <td bgcolor='DeepSkyBlue' COLSPAN=\"2\"> PARTICION </td> \n </tr> \n"
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(data.Partitions[i].Status[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_type </td> \n  <td bgcolor='LightSkyBlue'> %s </td> \n </tr> \n", string(data.Partitions[i].Type[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(data.Partitions[i].Fit[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_start </td> \n  <td bgcolor='LightSkyBlue'> %d </td> \n </tr> \n", data.Partitions[i].Start)
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", data.Partitions[i].Size)
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_name </td> \n  <td bgcolor='LightSkyBlue'> %s </td> \n </tr> \n", GetName(string(data.Partitions[i].Name[:])))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_id </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", GetId(string(data.Partitions[i].Id[:])))
			if string(data.Partitions[i].Type[:]) == "E" {
				cad += repLogicas(data.Partitions[i], disco)
			}
		}
	}

	//si hay espacio despues de la 4ta particion
	disponible = data.MbrSize - inicioLibre
	if disponible > 0 {
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#808080' COLSPAN=\"2\"> ESPACIO LIBRE <br/> %d bytes </td> \n </tr> \n", disponible)
	}

	return cad
}

func repLogicas(particion Partition, disco *os.File) string {
	cad := ""

	var actual EBR
	if err := Herramientas.ReadObject(disco, &actual, int64(particion.Start)); err != nil {
		fmt.Println("REP ERROR: No se encontro un ebr para reportar logicas")
		return ""
	}

	//Primera logica
	if actual.Size != 0 {
		cad += " <tr>\n  <td bgcolor='SteelBlue' COLSPAN=\"2\"> PARTICION LOGICA </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Status[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_next </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Next)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Fit[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_start </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Start)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", actual.Size)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_name </td> \n  <td bgcolor='SkyBlue'> %s </td> \n </tr> \n", GetName(string(actual.Name[:])))
	}

	//resto de logicas
	for actual.Next != -1 {
		if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
			fmt.Println("REP ERROR: fallo al leer particiones logicas")
			return ""
		}
		cad += " <tr>\n  <td bgcolor='SteelBlue' COLSPAN=\"2\"> PARTICION LOGICA </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Status[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_next </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Next)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Fit[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_start </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Start)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", actual.Size)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_name </td> \n  <td bgcolor='SkyBlue'> %s </td> \n </tr> \n", GetName(string(actual.Name[:])))
	}
	return cad
}
