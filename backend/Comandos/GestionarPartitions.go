package Comandos

import (
	"MIA_P2_202110206/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strings"
	"unsafe"
)

//Se le mete por paramaetro el disco del frontend para poder buscar en la ruta quemada del proyecto
func mandarParticiones(disco string) {
	// Construir la ruta del archivo basado en el driveLetter
	rutaBase := "./MIA/Discos/"
	nombreDisco := disco
	path := rutaBase + nombreDisco

	//file
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))

	if err != nil {
		Error("REP", "No se ha encontrado el disco.")
		return
	}
	var disk Structs.MBR
	file.Seek(0, 0)

	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &disk)
	if err_ != nil {
		Error("REP", "Error al leer el archivo")
		return
	}
	file.Close()

	contenidoParticiones := graficarPart(disk)

	file, err = os.Open(strings.ReplaceAll(path, "\n", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco")
		return
	}

	partitions := GetParticiones(disk)
	ext := false
	var extended Structs.Particion
	for i := 0; i < 4; i++ {
		if partitions[i].Part_status == '1' {
			if partitions[i].Part_type == "E"[0] || partitions[i].Part_type == "e"[0] {
				ext = true
				extended = partitions[i]
			}
		}
	}

	if ext {
		auxEbr := Structs.NewEBR()
		file.Seek(extended.Part_start, 0)
		data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &auxEbr)
		if err_ != nil {
			Error("REP", "Error al leer el archivo")
			return
		}

		for auxEbr.Part_next != -1 {
			partName := strings.TrimRight(string(auxEbr.Part_name[:]), "\x00")

			// Mover el puntero al siguiente registro de arranque extendido
			file.Seek(auxEbr.Part_next, 0)

			// Leer el siguiente EBR
			data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &auxEbr)
			if err_ != nil {
				Error("REP", "Error al leer el archivo")
				return
			}
		}

		file.Close() // Cerrar el archivo después de terminar de leer todos los EBRs
	}
}

func graficarPart(disk Structs.MBR) string {
	contenido := ""
	particiones := GetParticiones(disk)
	for i := 0; i < len(particiones); i++ {
		// Solo graficar la partición si su part_start es diferente de -1
		if particiones[i].Part_start != -1 {
			contenido += generarTablaParticion(particiones[i])
		}
	}
	return contenido
}

func generarTablaPart(particion Structs.Particion) string {
	partName := strings.TrimRight(string(particion.Part_name[:]), "\x00")

	tabla := ""
	tabla += partName + ""
	return tabla
}
