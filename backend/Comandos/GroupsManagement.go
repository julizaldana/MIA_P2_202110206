package Comandos

import (
	"MIA_P2_202110206/Structs"
	"bytes"
	"encoding/binary"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func ValidarDatosGrupos(context []string, action string, w http.ResponseWriter) {
	name := ""
	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "name") {
			name = tk[1]
		}
	}
	if name == "" {
		Error(action+"GRP", "No se encontró el parámetro name en el comando.")
		return
	}
	if Comparar(action, "MK") {
		mkgrp(name, w)
	} else if Comparar(action, "RM") {
		rmgrp(name, w)
	} else {
		Error(action+"GRP", "No se reconoce este comando.")
		return
	}
}

// CREAR GRUPO
func mkgrp(n string, w http.ResponseWriter) {
	if !Comparar(Logged.User, "root") {
		Error("MKGRP", "Solo el usuario \"root\" puede acceder a estos comandos.")
		MandarError("MKGRP", "Solo el usuario \"root\" puede acceder a estos comandos.", w)
		return
	}

	var path string
	partition := GetMount("MKGRP", Logged.Id, &path)
	if string(partition.Part_status) == "0" {
		Error("MKGRP", "No se encontró la partición montada con el id: "+Logged.Id)
		MandarError("MKGRP", "No se encontró la partición montada con el id: "+Logged.Id, w)
		return
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKGRP", "No se ha encontrado el disco.")
		return
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("MKGRP", "Error al leer el archivo")
		return
	}
	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = leerBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("MKGRP", "Error al leer el archivo")
		return
	}

	var fb Structs.BloquesArchivos
	txt := ""
	for bloque := 1; bloque < 16; bloque++ {
		if inode.I_block[bloque-1] == -1 {
			break
		}
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(bloque-1), 0)

		data = leerBytes(file, int(unsafe.Sizeof(Structs.BloquesArchivos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)

		if err_ != nil {
			Error("MKGRP", "Error al leer el archivo")
			return
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	c := 0
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if linea[2] == 'G' || linea[2] == 'g' {
			c++
			in := strings.Split(linea, ",")
			if in[2] == n {
				if linea[0] != '0' {
					Error("MKGRP", "EL nombre "+n+", ya está en uso.")
					MandarError("MKGRP", "EL nombre "+n+", ya está en uso.", w)
					return
				}
			}
		}
	}
	txt += strconv.Itoa(c+1) + ",G," + n + "\n"

	tam := len(txt)
	var cadenasS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenasS = append(cadenasS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenasS = append(cadenasS, txt)
		}
	} else {
		cadenasS = append(cadenasS, txt)
	}
	if len(cadenasS) > 16 {
		Error("MKGRP", "Se ha llenado la cantidad de archivos posibles y no se pueden generar más.")
		MandarError("MKGRP", "Se ha llenado la cantidad de archivos posibles y no se pueden generar más.", w)
		return
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("MKGRP", "No se ha encontrado el disco.")
		MandarError("MKGRP", "No se ha encontrado el disco.", w)
		return
	}

	for i := 0; i < len(cadenasS); i++ {

		var fbAux Structs.BloquesArchivos
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			EscribirBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenasS[i])

		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
		var bin6 bytes.Buffer
		binary.Write(&bin6, binary.BigEndian, fbAux)
		EscribirBytes(file, bin6.Bytes())

	}
	for i := 0; i < len(cadenasS); i++ {
		inode.I_block[i] = int64(0)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	EscribirBytes(file, inodos.Bytes())

	Mensaje("MKGRP", "Grupo "+n+", creado correctamente!")
	MandarMensaje("MKGRP", "Grupo "+n+", creado correctamente!", w)
	file.Close()
}

//ELIMINAR GRUPO

func rmgrp(n string, w http.ResponseWriter) {
	if !Comparar(Logged.User, "root") {
		Error("RMGRP", "Solo el usuario \"root\" puede acceder a estos comandos.")
		MandarError("RMGRP", "Solo el usuario \"root\" puede acceder a estos comandos.", w)
		return
	}

	var path string
	partition := GetMount("RMGRP", Logged.Id, &path)
	if string(partition.Part_status) == "0" {
		Error("RMGRP", "No se encontró la partición montada con el id: "+Logged.Id)
		MandarError("RMGRP", "No se encontró la partición montada con el id: "+Logged.Id, w)
		return
	}
	//file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("RMGRP", "No se ha encontrado el disco.")
		return
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &super)
	if err_ != nil {
		Error("RMGRP", "Error al leer el archivo")
		return
	}
	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = leerBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("RMGRP", "Error al leer el archivo")
		return
	}

	var fb Structs.BloquesArchivos
	txt := ""
	for bloque := 1; bloque < 16; bloque++ {
		if inode.I_block[bloque-1] == -1 {
			break
		}
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(bloque-1), 0)

		data = leerBytes(file, int(unsafe.Sizeof(Structs.BloquesArchivos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &fb)

		if err_ != nil {
			Error("RMGRP", "Error al leer el archivo")
			return
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	aux := ""

	vctr := strings.Split(txt, "\n")
	existe := false
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if (linea[2] == 'G' || linea[2] == 'g') && linea[0] != '0' {
			in := strings.Split(linea, ",")
			if in[2] == n {
				existe = true
				aux += strconv.Itoa(0) + ",G," + in[2] + "\n"
				continue
			}
		}
		aux += linea + "\n"
	}
	if !existe {
		Error("RMGRP", "No se encontró el grupo \""+n+"\".")
		MandarError("RMGRP", "No se encontró el grupo \""+n+"\".", w)
		return
	}
	txt = aux

	tam := len(txt)
	var cadenasS []string
	if tam > 64 {
		for tam > 64 {
			aux := ""
			for i := 0; i < 64; i++ {
				aux += string(txt[i])
			}
			cadenasS = append(cadenasS, aux)
			txt = strings.ReplaceAll(txt, aux, "")
			tam = len(txt)
		}
		if tam < 64 && tam != 0 {
			cadenasS = append(cadenasS, txt)
		}
	} else {
		cadenasS = append(cadenasS, txt)
	}
	if len(cadenasS) > 16 {
		Error("RMGRP", "Se ha llenado la cantidad de archivos posibles y no se pueden generar más.")
		MandarError("RMGRP", "Se ha llenado la cantidad de archivos posibles y no se pueden generar más.", w)
		return
	}
	file.Close()

	file, err = os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	//file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("RMGRP", "No se ha encontrado el disco.")
		MandarError("RMGRP", "No se ha encontrado el disco.", w)
		return
	}
	for i := 0; i < len(cadenasS); i++ {

		var fbAux Structs.BloquesArchivos
		if inode.I_block[i] == -1 {
			file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
			var binAux bytes.Buffer
			binary.Write(&binAux, binary.BigEndian, fbAux)
			EscribirBytes(file, binAux.Bytes())
		} else {
			fbAux = fb
		}

		copy(fbAux.B_content[:], cadenasS[i])

		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(i), 0)
		var bin6 bytes.Buffer
		binary.Write(&bin6, binary.BigEndian, fbAux)
		EscribirBytes(file, bin6.Bytes())

	}
	for i := 0; i < len(cadenasS); i++ {
		inode.I_block[i] = int64(0)
	}
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	var inodos bytes.Buffer
	binary.Write(&inodos, binary.BigEndian, inode)
	EscribirBytes(file, inodos.Bytes())

	Mensaje("RMGRP", "Grupo "+n+", eliminado correctamente!")
	MandarMensaje("RMGRP", "Grupo "+n+", eliminado correctamente!", w)

	file.Close()
}

func ValidarDatosCHGRP(context []string) {
	user := ""
	grp := ""

	for i := 0; i < len(context); i++ {
		token := context[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "user") {
			user = tk[1]
		} else if Comparar(tk[0], "grp") {
			grp = tk[1]
		}
	}
	if user == "" || grp == "" {
		Error("CHGRP", "Se necesita el user y el nombre de grupo para hacer el cambio")
	} else {
		chgrp(user, grp)
	}

}

func chgrp(user string, grp string) {
	if !Comparar(Logged.User, "root") {
		Error("CHGRP", "Solo el usuario \"root\" puede acceder a estos comandos.")
		return
	}

	var path string
	partition := GetMount("CHGRP", Logged.Id, &path)
	if string(partition.Part_status) == "0" {
		Error("CHGRP", "No se encontró la partición montada con el id: "+Logged.Id)
		return
	}

	file, err := os.OpenFile(strings.ReplaceAll(path, "\"", ""), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("CHGRP", "No se ha encontrado el disco.")
		return
	}

	super := Structs.NewSuperBloque()
	file.Seek(partition.Part_start, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &super)
	if err != nil {
		Error("CHGRP", "Error al leer el archivo")
		return
	}

	inode := Structs.NewInodos()
	file.Seek(super.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{})), 0)
	data = leerBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &inode)
	if err != nil {
		Error("CHGRP", "Error al leer el archivo")
		return
	}

	var fb Structs.BloquesArchivos
	txt := ""
	for bloque := 1; bloque < 16; bloque++ {
		if inode.I_block[bloque-1] == -1 {
			break
		}
		file.Seek(super.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(bloque-1), 0)

		data = leerBytes(file, int(unsafe.Sizeof(Structs.BloquesArchivos{})))
		buffer = bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &fb)

		if err != nil {
			Error("CHGRP", "Error al leer el archivo")
			return
		}

		for i := 0; i < len(fb.B_content); i++ {
			if fb.B_content[i] != 0 {
				txt += string(fb.B_content[i])
			}
		}
	}

	vctr := strings.Split(txt, "\n")
	existeUsuario := false
	existeGrupo := false
	for i := 0; i < len(vctr)-1; i++ {
		linea := vctr[i]
		if (linea[2] == 'U' || linea[2] == 'u') && linea[0] != '0' {
			in := strings.Split(linea, ",")
			if in[3] == user {
				existeUsuario = true
				in[2] = grp // Cambiar el grupo
				vctr[i] = strings.Join(in, ",")
			}
		}
		if (linea[2] == 'G' || linea[2] == 'g') && linea[0] != '0' {
			in := strings.Split(linea, ",")
			if in[2] == grp {
				existeGrupo = true
			}
		}
	}

	if !existeUsuario {
		Error("CHGRP", "El usuario \""+user+"\" no existe.")
		return
	}

	if !existeGrupo {
		Error("CHGRP", "El grupo \""+grp+"\" no existe.")
		return
	}

	// Escribir los cambios de nuevo en el archivo
	txt = strings.Join(vctr, "\n")

	// Truncar el archivo antes de escribir los cambios
	err = file.Truncate(0)
	if err != nil {
		Error("CHGRP", "Error al truncar el archivo")
		return
	}

	// Escribir el texto modificado en el archivo
	_, err = file.WriteAt([]byte(txt), 0)
	if err != nil {
		Error("CHGRP", "Error al escribir en el archivo")
		return
	}

	Mensaje("CHGRP", "Grupo de usuario "+user+" cambiado a "+grp+" correctamente.")

	file.Close()
}
