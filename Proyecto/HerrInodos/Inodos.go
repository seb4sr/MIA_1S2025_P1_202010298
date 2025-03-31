package toolsinodos

import (
	"Proyecto/Herramientas"
	"Proyecto/Structs"
	"encoding/binary"
	"os"
	"strings"
	"time"
)

func BuscarInodo(idInodo int32, path string, superBloque Structs.Superblock, file *os.File) int32 {
	 
	stepsPath := strings.Split(path, "/")
	 
	tmpPath := stepsPath[1:]
	 

	 
	var Inode0 Structs.Inode
	Herramientas.ReadObject(file, &Inode0, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
	 
	var folderBlock Structs.Folderblock
	for i := 0; i < 12; i++ {
		idBloque := Inode0.I_block[i]
		if idBloque != -1 {
			Herramientas.ReadObject(file, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
			 
			for j := 2; j < 4; j++ {
				 
				apuntador := folderBlock.B_content[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GetB_name(string(folderBlock.B_content[j].B_name[:]))
					if tmpPath[0] == pathActual {
						 
						if len(tmpPath) > 1 {
							return buscarIrecursivo(apuntador, tmpPath[1:], superBloque.S_inode_start, superBloque.S_block_start, file)
						} else {
							return apuntador
						}
					}
				}
			}
		}
	}
	 
	 
	 
	return idInodo
}

 
func buscarIrecursivo(idInodo int32, path []string, iStart int32, bStart int32, file *os.File) int32 {
	 
	var inodo Structs.Inode
	Herramientas.ReadObject(file, &inodo, int64(iStart+(idInodo*int32(binary.Size(Structs.Inode{})))))

	 
	 
	var folderBlock Structs.Folderblock
	for i := 0; i < 12; i++ {
		idBloque := inodo.I_block[i]
		if idBloque != -1 {
			Herramientas.ReadObject(file, &folderBlock, int64(bStart+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
			 
			for j := 2; j < 4; j++ {
				apuntador := folderBlock.B_content[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GetB_name(string(folderBlock.B_content[j].B_name[:]))					
					if path[0] == pathActual {
						if len(path) > 1 {
							 
							return buscarIrecursivo(apuntador, path[1:], iStart, bStart, file)
						} else {
							 
							return apuntador
						}
					}
				}
			}
		}
	}
	 
	 
	return -1
}

func CreaCarpeta(idInode int32, carpeta string, initSuperBloque int64, disco *os.File)int32{
	var superBloque Structs.Superblock
	Herramientas.ReadObject(disco, &superBloque, initSuperBloque)

	var inodo Structs.Inode
	Herramientas.ReadObject(disco, &inodo, int64(superBloque.S_inode_start+(idInode*int32(binary.Size(Structs.Inode{})))))

	 
	for i := 0; i < 12; i++ {
		idBloque := inodo.I_block[i]
		if idBloque != -1{
			 
			var folderBlock Structs.Folderblock
			Herramientas.ReadObject(disco, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

			 
			for j := 2; j < 4; j++{
				apuntador := folderBlock.B_content[j].B_inodo
				 
				if apuntador == -1 {
					 
					copy(folderBlock.B_content[j].B_name[:], carpeta)
					ino := superBloque.S_first_ino  
					folderBlock.B_content[j].B_inodo = ino
					 
					Herramientas.WriteObject(disco, folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

					 
					var newInodo Structs.Inode
					newInodo.I_uid = Structs.UsuarioActual.IdUsr
					newInodo.I_gid = Structs.UsuarioActual.IdGrp
					newInodo.I_size = 0  
					 
					ahora := time.Now()
					date := ahora.Format("02/01/2006 15:04")
					copy(newInodo.I_atime[:], date)
					copy(newInodo.I_ctime[:], date)
					copy(newInodo.I_mtime[:], date)
					copy(newInodo.I_type[:], "0")  
					copy(newInodo.I_mtime[:], "664")

					 
					for i := int32(0); i < 15; i++ {
						newInodo.I_block[i] = -1
					}
					 
					block := superBloque.S_first_blo
					newInodo.I_block[0] = block
					 
					Herramientas.WriteObject(disco, newInodo, int64(superBloque.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

					 
					var newFolderBlock Structs.Folderblock
					newFolderBlock.B_content[0].B_inodo = ino  
					copy(newFolderBlock.B_content[0].B_name[:], ".")
					newFolderBlock.B_content[1].B_inodo = folderBlock.B_content[0].B_inodo  
					copy(newFolderBlock.B_content[1].B_name[:], "..")
					newFolderBlock.B_content[2].B_inodo = -1
					newFolderBlock.B_content[3].B_inodo = -1
					 
					Herramientas.WriteObject(disco, newFolderBlock, int64(superBloque.S_block_start+(block*int32(binary.Size(Structs.Folderblock{})))))

					 
					superBloque.S_free_inodes_count -= 1
					superBloque.S_free_blocks_count -= 1
					superBloque.S_first_blo += 1
					superBloque.S_first_ino += 1
					 
					Herramientas.WriteObject(disco, superBloque, initSuperBloque)

					 
					Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_block_start+block))

					 
					Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_inode_start+ino))
					 
					return ino
				}
			} 
		 
		}else{
			 
			 
			block := superBloque.S_first_blo  
			inodo.I_block[i] = block
			 
			Herramientas.WriteObject(disco, &inodo, int64(superBloque.S_inode_start+(idInode*int32(binary.Size(Structs.Inode{})))))

			 
			var folderBlock Structs.Folderblock
			bloque := inodo.I_block[0]  
			Herramientas.ReadObject(disco, &folderBlock, int64(superBloque.S_block_start+(bloque*int32(binary.Size(Structs.Folderblock{})))))

			 
			var newFolderBlock1 Structs.Folderblock
			newFolderBlock1.B_content[0].B_inodo = folderBlock.B_content[0].B_inodo  
			copy(newFolderBlock1.B_content[0].B_name[:], ".")
			newFolderBlock1.B_content[1].B_inodo = folderBlock.B_content[1].B_inodo  
			copy(newFolderBlock1.B_content[1].B_name[:], "..")
			ino := superBloque.S_first_ino                         
			newFolderBlock1.B_content[2].B_inodo = ino             
			copy(newFolderBlock1.B_content[2].B_name[:], carpeta)  
			newFolderBlock1.B_content[3].B_inodo = -1
			 
			Herramientas.WriteObject(disco, newFolderBlock1, int64(superBloque.S_block_start+(block*int32(binary.Size(Structs.Folderblock{})))))

			 
			var newInodo Structs.Inode
			newInodo.I_uid = Structs.UsuarioActual.IdUsr
			newInodo.I_gid = Structs.UsuarioActual.IdGrp
			newInodo.I_size = 0  
			 
			ahora := time.Now()
			date := ahora.Format("02/01/2006 15:04")
			copy(newInodo.I_atime[:], date)
			copy(newInodo.I_ctime[:], date)
			copy(newInodo.I_mtime[:], date)
			copy(newInodo.I_type[:], "0")  
			copy(newInodo.I_mtime[:], "664")

			 
			for i := int32(0); i < 15; i++ {
				newInodo.I_block[i] = -1
			}
			 
			block2 := superBloque.S_first_blo + 1
			newInodo.I_block[0] = block2
			 
			Herramientas.WriteObject(disco, newInodo, int64(superBloque.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

			 
			var newFolderBlock2 Structs.Folderblock
			newFolderBlock2.B_content[0].B_inodo = ino  
			copy(newFolderBlock2.B_content[0].B_name[:], ".")
			newFolderBlock2.B_content[1].B_inodo = newFolderBlock1.B_content[0].B_inodo  
			copy(newFolderBlock2.B_content[1].B_name[:], "..")
			newFolderBlock2.B_content[2].B_inodo = -1
			newFolderBlock2.B_content[3].B_inodo = -1
			 
			Herramientas.WriteObject(disco, newFolderBlock2, int64(superBloque.S_block_start+(block2*int32(binary.Size(Structs.Folderblock{})))))

			 
			superBloque.S_free_inodes_count -= 1
			superBloque.S_free_blocks_count -= 2
			superBloque.S_first_blo += 2
			superBloque.S_first_ino += 1
			Herramientas.WriteObject(disco, superBloque, initSuperBloque)

			 
			Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_block_start+block))
			Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_block_start+block2))

			 
			Herramientas.WriteObject(disco, byte(1), int64(superBloque.S_bm_inode_start+ino))
			return ino
		}
	} 
	return 0
}			