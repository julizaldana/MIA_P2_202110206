package Comandos

import (
	"MIA_P2_202110206/Structs"
	"bytes"
	"encoding/binary"
	"net/http"
	"os"
	"strings"
	"unsafe"
)

func MandarParticiones(disco string, w http.ResponseWriter) {

	// Construir la ruta del archivo basado en el nombre del disco
	rutaBase := "./MIA/Discos/"
	nombreDisco := disco
	path := rutaBase + nombreDisco

	// Abrir el archivo del disco
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	if err != nil {
		Error("REP", "No se ha encontrado el disco.")
		return
	}
	defer file.Close()

	// Leer el MBR del disco
	var disk Structs.MBR
	file.Seek(0, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &disk)
	if err != nil {
		Error("REP", "Error al leer el archivo")
		return
	}

	// Obtener las particiones del MBR
	particiones := GetParticiones(disk)

	// Crear una lista para almacenar los nombres de las particiones
	nombresParticiones := []string{}

	// Recorrer las particiones y obtener sus nombres
	for _, particion := range particiones {
		// Solo graficar la partición si su part_start es diferente de -1
		if particion.Part_start != -1 {
			partName := strings.TrimRight(string(particion.Part_name[:]), "\x00")
			nombresParticiones = append(nombresParticiones, partName)
		}
	}

}

/*
func enviarListaParticiones(w http.ResponseWriter, r *http.Request) {
	// Convertir la lista de nombres de particiones a formato JSON
	//jsonNombres, err := json.Marshal(nombresParticiones)
	if err != nil {
		Error("REP", "Error al convertir a JSON")
		return
	}

	// Escribir la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonNombres)
}
*/
