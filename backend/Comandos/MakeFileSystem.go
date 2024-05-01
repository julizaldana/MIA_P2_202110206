package Comandos

import (
	"MIA_P2_202110206/Structs"
	"bytes"
	"encoding/binary"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"
)

//SE DEBE DE CREAR UN SISTEMA DE ARCHIVOS EXT2/EXT3 EN LA RAIZ, SE TENDRÁ UN ARCHIVO TXT, CON LOS USUARIOS, GRUPOS, INFORMACION DE LOS MISMOS.

func ValidarDatosMKFS(context []string, w http.ResponseWriter) {
	tipo := "full"
	id := ""
	fs := "" // Variable para almacenar el tipo de sistema de archivos

	for i := 0; i < len(context); i++ {
		current := context[i]
		comando := strings.Split(current, "=")

		switch comando[0] {
		case "id":
			id = comando[1]
		case "type":
			if Comparar(comando[1], "full") {
				tipo = comando[1]
			} else {
				Error("MKFS", "El comando type debe tener valores especificos")
				return
			}
		case "fs":
			fs = comando[1]
		default:
			Error("MKFS", "Parámetro no reconocido para mkfs")
			return
		}
	}

	if id == "" {
		Error("MKFS", "El comando MKFS requiere un id de particion para poder formatear.")
		return
	}

	// Determinar el tipo de sistema de archivos
	if fs == "" || fs == "2fs" {
		crearFileSystem2(id, tipo, w)
	} else if fs == "3fs" {
		crearFileSystem3(id, tipo, w)
	} else {
		Error("MKFS", "Tipo de sistema de archivos no es válido, se admite 2fs/3fs para EXT2/EXT3")
		return
	}
}

// CREAR UN SISTEMA DE ARCHIVOS EXT2 POR DEFECTO
func crearFileSystem(id string, t string, w http.ResponseWriter) {
	p := ""
	particion := GetMount("MKFS", id, &p)                                                                                                                                          //Obtener una particion con el GetMount(), que se obtiene ingresando el id de particion.                                                                                                                                        //SE LLAMA AL COMANDO PARA OBTENER LA PARTCION
	n := math.Floor(float64(particion.Part_s-int64(unsafe.Sizeof(Structs.SuperBloque{})))) / float64(4+unsafe.Sizeof(Structs.Inodos{})+3*unsafe.Sizeof(Structs.BloquesArchivos{})) //obtener n del calculo para obtener el tamaño

	//Creacion de superbloque
	spr := Structs.NewSuperBloque()
	spr.S_magic = 0xEF53
	spr.S_inode_s = int64(unsafe.Sizeof(Structs.Inodos{}))
	spr.S_block_s = int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))
	spr.S_inodes_count = int64(n)
	spr.S_free_inodes_count = int64(n)
	spr.S_blocks_count = int64(3 * n)
	spr.S_free_blocks_count = int64(3 * n)
	fecha := time.Now().String()
	copy(spr.S_mtime[:], fecha)
	spr.S_mnt_count = spr.S_mnt_count + 1
	spr.S_filesystem_type = 2
	ext2(spr, particion, int64(n), p, w)
}

//CREAR UN SISTEMA DE ARCHIVOS EXT2 con fs=2fs

func crearFileSystem2(id string, fs string, w http.ResponseWriter) {
	p := ""
	particion := GetMount("MKFS", id, &p)
	n := math.Floor(float64(particion.Part_s-int64(unsafe.Sizeof(Structs.SuperBloque{})))) / float64(4+unsafe.Sizeof(Structs.Inodos{})+3*unsafe.Sizeof(Structs.BloquesArchivos{})) //obtener n del calculo para obtener el tamaño

	//Creacion de superbloque
	spr := Structs.NewSuperBloque()
	spr.S_magic = 0xEF53
	spr.S_inode_s = int64(unsafe.Sizeof(Structs.Inodos{}))
	spr.S_block_s = int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))
	spr.S_inodes_count = int64(n)
	spr.S_free_inodes_count = int64(n)
	spr.S_blocks_count = int64(3 * n)
	spr.S_free_blocks_count = int64(3 * n)
	fecha := time.Now().String()
	copy(spr.S_mtime[:], fecha)
	spr.S_mnt_count = spr.S_mnt_count + 1
	spr.S_filesystem_type = 2
	ext2(spr, particion, int64(n), p, w)

}

//CREAR UN SISTEMA DE ARCHIVOS EXT3 con fs=3fs

func crearFileSystem3(id string, fs string, w http.ResponseWriter) {
	p := ""
	particion := GetMount("MKFS", id, &p)
	n := math.Floor(float64(particion.Part_s-int64(unsafe.Sizeof(Structs.SuperBloque{})))) / float64(4+unsafe.Sizeof(Structs.Journaling{})+unsafe.Sizeof(Structs.Inodos{})+3*unsafe.Sizeof(Structs.BloquesArchivos{})) //obtener n del calculo para obtener el tamaño

	//Creacion de superbloque
	spr := Structs.NewSuperBloque()
	spr.S_magic = 0xEF53
	spr.S_inode_s = int64(unsafe.Sizeof(Structs.Inodos{}))
	spr.S_block_s = int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))
	spr.S_inodes_count = int64(n)
	spr.S_free_inodes_count = int64(n)
	spr.S_blocks_count = int64(3 * n)
	spr.S_free_blocks_count = int64(3 * n)
	fecha := time.Now().String()
	copy(spr.S_mtime[:], fecha)
	spr.S_mnt_count = spr.S_mnt_count + 1
	spr.S_filesystem_type = 3

	//inicializar journal
	journal := Structs.NewJournaling()
	copy(journal.J_operacion[:], "mkdir")
	copy(journal.J_path[:], "/")
	copy(journal.J_content[:], "-")
	copy(journal.J_fecha[:], fecha)

	ext3(spr, journal, particion, int64(n), p, w)

}

func ext2(spr Structs.SuperBloque, p Structs.Particion, n int64, path string, w http.ResponseWriter) {
	spr.S_bm_inode_start = p.Part_s + int64(unsafe.Sizeof(Structs.SuperBloque{}))
	spr.S_bm_block_start = spr.S_bm_inode_start + n
	spr.S_inode_start = spr.S_bm_block_start + (3 * n)
	spr.S_block_start = spr.S_bm_inode_start + (n * int64(unsafe.Sizeof(Structs.Inodos{})))

	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco.")
		return
	}

	//AQUI SE ESCRIBE EL SUPERBLOQUE YA EN EL ARCHIVO BINARIO, EN EL INICIO DE LA PARTICION.
	file.Seek(p.Part_start, 0)
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, spr)
	EscribirBytes(file, binario2.Bytes())

	//Se escribe el bitmap inodos despues del superbloque
	zero := '0'
	file.Seek(spr.S_bm_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioZero bytes.Buffer
		binary.Write(&binarioZero, binary.BigEndian, zero)
		EscribirBytes(file, binarioZero.Bytes())
	}

	//Se inicializa bitmap bloques despues de bitmap de inodos
	file.Seek(spr.S_bm_block_start, 0)
	for i := 0; i < 3*int(n); i++ {
		var binarioZero bytes.Buffer
		binary.Write(&binarioZero, binary.BigEndian, zero)
		EscribirBytes(file, binarioZero.Bytes())
	}
	inode := Structs.NewInodos()
	//INICIALIZANDO EL INODO 0
	inode.I_uid = -1
	inode.I_gid = -1
	inode.I_s = -1
	for i := 0; i < len(inode.I_block); i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = -1
	inode.I_perm = -1

	//Se escribe en el inicio de los inodos, se escribe el inodo 0, que contendrá la informacion del directorio raiz
	file.Seek(spr.S_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioInodos bytes.Buffer
		binary.Write(&binarioInodos, binary.BigEndian, inode)
		EscribirBytes(file, binarioInodos.Bytes())
	}
	//INICIALIZANDO EL BLOQUE 0, que se escribirá en el espacio de informacion de los bloques
	folder := Structs.NewBloquesCarpetas()

	for i := 0; i < len(folder.B_content); i++ {
		folder.B_content[i].B_inodo = -1
	}

	file.Seek(spr.S_block_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioFolder bytes.Buffer
		binary.Write(&binarioFolder, binary.BigEndian, folder)
		EscribirBytes(file, binarioFolder.Bytes())
	}
	file.Close()

	recuperado := Structs.NewSuperBloque()
	//ABRIR ARCHIVO
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)

	file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco.")
		return
	}
	//Lee los bytes del struct superbloque del archivo
	file.Seek(p.Part_start, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &recuperado)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo")
		return
	}
	file.Close()

	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_s = 0
	fecha := time.Now().String()
	copy(inode.I_atime[:], fecha)
	copy(inode.I_ctime[:], fecha)
	copy(inode.I_mtime[:], fecha)
	inode.I_type = 0
	inode.I_perm = 664
	inode.I_block[0] = 0

	fb := Structs.NewBloquesCarpetas()
	copy(fb.B_content[0].B_name[:], ".")
	fb.B_content[0].B_inodo = 0
	copy(fb.B_content[1].B_name[:], "..")
	fb.B_content[1].B_inodo = 0
	copy(fb.B_content[2].B_name[:], "users.txt")
	fb.B_content[2].B_inodo = 1
	//copy(fb.B_content[3].B_name[:], "--")
	//fb.B_content[3].B_inodo = 0

	dataArchivo := "1,G,root\n1,U,root,root,123\n"
	inodetmp := Structs.NewInodos()
	inodetmp.I_uid = 1
	inodetmp.I_gid = 1
	inodetmp.I_s = int64(unsafe.Sizeof(dataArchivo) + unsafe.Sizeof(Structs.BloquesCarpetas{}))

	copy(inodetmp.I_atime[:], fecha)
	copy(inodetmp.I_ctime[:], fecha)
	copy(inodetmp.I_mtime[:], fecha)
	inodetmp.I_type = 1
	inodetmp.I_perm = 664
	inodetmp.I_block[0] = 1

	inode.I_s = inodetmp.I_s + int64(unsafe.Sizeof(Structs.BloquesCarpetas{})) + int64(unsafe.Sizeof(Structs.Inodos{}))

	var fileb Structs.BloquesArchivos
	copy(fileb.B_content[:], dataArchivo)

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco.")
		return
	}
	file.Seek(spr.S_bm_inode_start, 0)
	caracter := '1'

	var bin1 bytes.Buffer
	binary.Write(&bin1, binary.BigEndian, caracter)
	EscribirBytes(file, bin1.Bytes())
	EscribirBytes(file, bin1.Bytes())

	file.Seek(spr.S_bm_block_start, 0)
	var bin2 bytes.Buffer
	binary.Write(&bin2, binary.BigEndian, caracter)
	EscribirBytes(file, bin2.Bytes())
	EscribirBytes(file, bin1.Bytes())

	file.Seek(spr.S_inode_start, 0)

	var bin3 bytes.Buffer
	binary.Write(&bin3, binary.BigEndian, inode)
	EscribirBytes(file, bin3.Bytes())

	file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var bin4 bytes.Buffer
	binary.Write(&bin4, binary.BigEndian, inodetmp)
	EscribirBytes(file, bin4.Bytes())

	file.Seek(spr.S_block_start, 0)

	var bin5 bytes.Buffer
	binary.Write(&bin5, binary.BigEndian, fb)
	EscribirBytes(file, bin5.Bytes())

	//fmt.Println(spr.S_block_start + int64(unsafe.Sizeof(Structs.BloquesCarpetas{})))

	file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{})), 0)
	var bin6 bytes.Buffer
	binary.Write(&bin6, binary.BigEndian, fileb)
	EscribirBytes(file, bin6.Bytes())

	file.Close()

	nombreParticion := ""
	for i := 0; i < len(p.Part_name); i++ {
		if p.Part_name[i] != 0 {
			nombreParticion += string(p.Part_name[i])
		}
	}
	Mensaje("MKFS", "Se ha formateado un sistema EXT2 en la partición "+nombreParticion+" de manera correcta.")
	MandarMensaje("MKFS", "Se ha formateado un sistema EXT2 en la partición "+nombreParticion+" de manera correcta.", w)
}

func ext3(spr Structs.SuperBloque, journal Structs.Journaling, p Structs.Particion, n int64, path string, w http.ResponseWriter) {
	spr.S_bm_inode_start = p.Part_s + int64(unsafe.Sizeof(Structs.SuperBloque{})) + int64(unsafe.Sizeof(Structs.Journaling{}))
	spr.S_bm_block_start = spr.S_bm_inode_start + n
	spr.S_inode_start = spr.S_bm_block_start + (3 * n)
	spr.S_block_start = spr.S_bm_inode_start + (n * int64(unsafe.Sizeof(Structs.Inodos{})))

	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco.")
		MandarError("MKFS", "No se ha encontrado el disco.", w)
		return
	}

	//AQUI SE ESCRIBE EL SUPERBLOQUE YA EN EL ARCHIVO BINARIO, EN EL INICIO DE LA PARTICION.
	file.Seek(p.Part_start, 0)
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, spr)
	EscribirBytes(file, binario2.Bytes())

	//Aqui se deberia de escribir el journaling despues del superbloque
	file.Seek(p.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), 0)
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, journal)
	EscribirBytes(file, binario3.Bytes())

	//Se escribe el bitmap inodos despues del superbloque y journaling
	zero := '0'
	file.Seek(spr.S_bm_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioZero bytes.Buffer
		binary.Write(&binarioZero, binary.BigEndian, zero)
		EscribirBytes(file, binarioZero.Bytes())
	}

	//Se inicializa bitmap bloques despues de bitmap de inodos
	file.Seek(spr.S_bm_block_start, 0)
	for i := 0; i < 3*int(n); i++ {
		var binarioZero bytes.Buffer
		binary.Write(&binarioZero, binary.BigEndian, zero)
		EscribirBytes(file, binarioZero.Bytes())
	}
	inode := Structs.NewInodos()
	//INICIALIZANDO EL INODO 0
	inode.I_uid = -1
	inode.I_gid = -1
	inode.I_s = -1
	for i := 0; i < len(inode.I_block); i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = -1
	inode.I_perm = -1

	//Se escribe en el inicio de los inodos, se escribe el inodo 0, que contendrá la informacion del directorio raiz
	file.Seek(spr.S_inode_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioInodos bytes.Buffer
		binary.Write(&binarioInodos, binary.BigEndian, inode)
		EscribirBytes(file, binarioInodos.Bytes())
	}
	//INICIALIZANDO EL BLOQUE 0, que se escribirá en el espacio de informacion de los bloques
	folder := Structs.NewBloquesCarpetas()

	for i := 0; i < len(folder.B_content); i++ {
		folder.B_content[i].B_inodo = -1
	}

	file.Seek(spr.S_block_start, 0)
	for i := 0; i < int(n); i++ {
		var binarioFolder bytes.Buffer
		binary.Write(&binarioFolder, binary.BigEndian, folder)
		EscribirBytes(file, binarioFolder.Bytes())
	}
	file.Close()

	recuperado := Structs.NewSuperBloque()
	journalrec := Structs.NewJournaling()
	//ABRIR ARCHIVO
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)

	file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco.")
		MandarError("MKFS", "No se ha encontrado el disco.", w)
		return
	}
	//Lee los bytes del struct superbloque del archivo
	file.Seek(p.Part_start, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &recuperado)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo")
		return
	}

	file.Seek(p.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{})), 0)
	data = leerBytes(file, int(unsafe.Sizeof(Structs.Journaling{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &journalrec)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo")
		return
	}
	file.Close()

	inode.I_uid = 1
	inode.I_gid = 1
	inode.I_s = 0
	fecha := time.Now().String()
	copy(inode.I_atime[:], fecha)
	copy(inode.I_ctime[:], fecha)
	copy(inode.I_mtime[:], fecha)
	inode.I_type = 0
	inode.I_perm = 664
	inode.I_block[0] = 0

	fb := Structs.NewBloquesCarpetas()
	copy(fb.B_content[0].B_name[:], ".")
	fb.B_content[0].B_inodo = 0
	copy(fb.B_content[1].B_name[:], "..")
	fb.B_content[1].B_inodo = 0
	copy(fb.B_content[2].B_name[:], "users.txt")
	fb.B_content[2].B_inodo = 1
	//copy(fb.B_content[3].B_name[:], "--")
	//fb.B_content[3].B_inodo = 0

	dataArchivo := "1,G,root\n1,U,root,root,123\n"
	inodetmp := Structs.NewInodos()
	inodetmp.I_uid = 1
	inodetmp.I_gid = 1
	inodetmp.I_s = int64(unsafe.Sizeof(dataArchivo) + unsafe.Sizeof(Structs.BloquesCarpetas{}))

	copy(inodetmp.I_atime[:], fecha)
	copy(inodetmp.I_ctime[:], fecha)
	copy(inodetmp.I_mtime[:], fecha)
	inodetmp.I_type = 1
	inodetmp.I_perm = 664
	inodetmp.I_block[0] = 1

	inode.I_s = inodetmp.I_s + int64(unsafe.Sizeof(Structs.BloquesCarpetas{})) + int64(unsafe.Sizeof(Structs.Inodos{}))

	newjournal := Structs.NewJournaling()
	copy(newjournal.J_operacion[:], "mkfile")
	copy(newjournal.J_path[:], "/users.txt")
	copy(newjournal.J_content[:], dataArchivo)
	copy(newjournal.J_fecha[:], fecha)

	var fileb Structs.BloquesArchivos
	copy(fileb.B_content[:], dataArchivo)

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err = os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKFS", "No se ha encontrado el disco.")
		return
	}
	//Se escribe nuevo objeto journaling
	file.Seek(p.Part_start+int64(unsafe.Sizeof(Structs.SuperBloque{}))+int64(unsafe.Sizeof(Structs.Journaling{})), 0)
	var bins1 bytes.Buffer
	binary.Write(&bins1, binary.BigEndian, newjournal)
	EscribirBytes(file, bins1.Bytes())

	file.Seek(spr.S_bm_inode_start, 0)
	caracter := '1'

	var bin1 bytes.Buffer
	binary.Write(&bin1, binary.BigEndian, caracter)
	EscribirBytes(file, bin1.Bytes())
	EscribirBytes(file, bin1.Bytes())

	file.Seek(spr.S_bm_block_start, 0)
	var bin2 bytes.Buffer
	binary.Write(&bin2, binary.BigEndian, caracter)
	EscribirBytes(file, bin2.Bytes())
	EscribirBytes(file, bin1.Bytes())

	file.Seek(spr.S_inode_start, 0)

	var bin3 bytes.Buffer
	binary.Write(&bin3, binary.BigEndian, inode)
	EscribirBytes(file, bin3.Bytes())

	file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var bin4 bytes.Buffer
	binary.Write(&bin4, binary.BigEndian, inodetmp)
	EscribirBytes(file, bin4.Bytes())

	file.Seek(spr.S_block_start, 0)

	var bin5 bytes.Buffer
	binary.Write(&bin5, binary.BigEndian, fb)
	EscribirBytes(file, bin5.Bytes())

	//fmt.Println(spr.S_block_start + int64(unsafe.Sizeof(Structs.BloquesCarpetas{})))

	file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{})), 0)
	var bin6 bytes.Buffer
	binary.Write(&bin6, binary.BigEndian, fileb)
	EscribirBytes(file, bin6.Bytes())

	file.Close()

	nombreParticion := ""
	for i := 0; i < len(p.Part_name); i++ {
		if p.Part_name[i] != 0 {
			nombreParticion += string(p.Part_name[i])
		}
	}
	Mensaje("MKFS", "Se ha formateado un sistema EXT3 en la partición "+nombreParticion+" de manera correcta.")
	MandarMensaje("MKFS", "Se ha formateado un sistema EXT3 en la partición "+nombreParticion+" de manera correcta.", w)
}
