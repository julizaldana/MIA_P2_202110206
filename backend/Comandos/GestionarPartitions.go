package Comandos

import (
	"MIA_P2_202110206/Structs"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

// Estructura para representar una partición
type Particion struct {
	Name string `json:"name"`
}

// Función para obtener las particiones asociadas con un disco
/*func MandarParticiones(disco string) []Particion {
	// Crear un slice para almacenar las particiones
	var particiones []Particion

	// Construir la ruta del archivo basado en el nombre del disco
	rutaBase := "./MIA/Discos/"
	nombreDisco := disco
	rutaDisco := filepath.Join(rutaBase, nombreDisco)

	// Abrir el archivo del disco
	file, err := os.Open(strings.ReplaceAll(rutaDisco, "\"", ""))
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
	particionesMBR := GetParticiones(disk)

	// Recorrer las particiones y obtener sus nombres
	for _, particion := range particionesMBR {
		// Solo graficar la partición si su part_start es diferente de -1
		if particion.Part_start != -1 {
			partName := strings.TrimRight(string(particion.Part_name[:]), "\x00")
			particion := Particion{
				Name: partName}

			particiones = append(particiones, particion)
		}
	}
	log.Println(particiones)
	return particiones
	// Ahora tienes todas las particiones en el slice 'particiones'
	// Puedes enviar este slice como respuesta
}
*/

func MostrarParticiones(disco string) {
	// Construir la ruta del archivo basado en el nombre del disco
	rutaBase := "./MIA/Discos/"
	nombreDisco := disco
	rutaDisco := filepath.Join(rutaBase, nombreDisco)

	//file
	file, err := os.Open(strings.ReplaceAll(rutaDisco, "\"", ""))

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

	partitions := GetParticiones(disk)
	var extended Structs.Particion
	ext := false
	for i := 0; i < 4; i++ {
		if partitions[i].Part_status == '1' {
			if partitions[i].Part_type == "E"[0] || partitions[i].Part_type == "e"[0] {
				ext = true
				extended = partitions[i]
			}
		}
	}

	//nombreparticiones := ""
	content := ""
	var positions [5]int64
	var positionsii [5]int64
	positions[0] = disk.Mbr_partition_1.Part_start - (1 + int64(unsafe.Sizeof(Structs.MBR{})))
	positions[1] = disk.Mbr_partition_2.Part_start - disk.Mbr_partition_1.Part_start + disk.Mbr_partition_1.Part_s
	positions[2] = disk.Mbr_partition_3.Part_start - disk.Mbr_partition_2.Part_start + disk.Mbr_partition_2.Part_s
	positions[3] = disk.Mbr_partition_4.Part_start - disk.Mbr_partition_3.Part_start + disk.Mbr_partition_3.Part_s
	positions[4] = disk.Mbr_tamano + 1 - disk.Mbr_partition_4.Part_s + disk.Mbr_partition_4.Part_s
	copy(positionsii[:], positions[:])

	//logic := 0
	tmplogic := ""
	if ext {
		auxEbr := Structs.NewEBR()

		file, err = os.Open(strings.ReplaceAll(rutaDisco, "\n", ""))

		if err != nil {
			Error("REP", "No se ha encontrado el disco")
			return
		}

		file.Seek(extended.Part_start, 0)
		data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &auxEbr)
		if err_ != nil {
			Error("REP", "Error al leer el archivo")
			return
		}
		file.Close()
		//var tamGen int64 = 0
		for auxEbr.Part_next != -1 {
			//tamGen += auxEbr.Part_s
			//res := float64(auxEbr.Part_s) / float64(disk.Mbr_tamano)
			//res = res * 100
			//s := fmt.Sprintf("%.2f", res)
			nombrePart := strings.TrimRight(string(auxEbr.Part_name[:]), "\x00")
			tmplogic += "Particion Logica " + nombrePart + "\n"

			//resta := float64(auxEbr.Part_next) - (float64(auxEbr.Part_start) + float64(auxEbr.Part_s))
			//resta = resta / float64(disk.Mbr_tamano)
			//resta = resta * 10000.00
			//resta = math.Round(resta) / 100.00 //PARA OBTENER LOS PORCENTAJES
			//f resta != 0 {
			//s = fmt.Sprintf("%f", resta)
			//tmplogic += "<td>\"Logica\n " + nombrePart + s + "% libre de la partición extendida\"</td>\n"
			//logic++
			//}
			//logic += 2 //Son los id, para los nodos de graphviz
			file, err = os.Open(strings.ReplaceAll(rutaDisco, "\"", ""))

			if err != nil {
				Error("REP", "No se ha encontrado el disco")
				return
			}

			file.Seek(auxEbr.Part_next, 0)
			data = leerBytes(file, int(unsafe.Sizeof(Structs.EBR{})))
			buffer = bytes.NewBuffer(data)
			err_ = binary.Read(buffer, binary.BigEndian, &auxEbr)
			if err_ != nil {
				Error("REP", "Error al leer el archivo")
				return
			}
			file.Close()
		}
		//resta := float64(extended.Part_s) - float64(tamGen)
		//resta = resta / float64(disk.Mbr_tamano)
		//resta = math.Round(resta * 100)
		//if resta != 0 {
		//s := fmt.Sprintf("%.2f", resta)
		//nombrePart := strings.TrimRight(string(auxEbr.Part_name[:]), "\x00")
		//tmplogic += "<td>\"Libre \n" + nombrePart + s + "% de la partición extendida. \"</td>\n"
		//logic++
		//}
		//tmplogic += "</tr>\n\n"
		//logic += 2
	}
	//var tamPrim int64
	for i := 0; i < 4; i++ {
		if partitions[i].Part_type == 'E' {
			//tamPrim += partitions[i].Part_s
			//res := float64(partitions[i].Part_s) / float64(disk.Mbr_tamano)
			//res = math.Round(res*10000.00) / 100.00
			//s := fmt.Sprintf("%.2f", res)
			partName := strings.TrimRight(string(partitions[i].Part_name[:]), "\x00")
			content += "Particion Extendida " + partName + "\n"
		} else if partitions[i].Part_start != -1 {
			//tamPrim += partitions[i].Part_s
			//res := float64(partitions[i].Part_s) / float64(disk.Mbr_tamano)
			//res = math.Round(res*10000.00) / 100.00
			//s := fmt.Sprintf("%.2f", res)
			partName := strings.TrimRight(string(partitions[i].Part_name[:]), "\x00")
			content += "Particion Primaria " + partName + "\n"
		}
	}

	//if tamPrim != 0 {
	//ibre := disk.Mbr_tamano - tamPrim
	//res := float64(libre) / float64(disk.Mbr_tamano)
	//res = math.Round(res * 100)
	//s := fmt.Sprintf("%.2f", res)
	//content += "<td ROWSPAN='2'> Libre \n" + s + "% del disco </td>"

	//}
	content += tmplogic

	//fmt.Println(content)
	//log.Println("https://quickchart.io/graphviz?graph=" + content)
	//se crean reportes en /MIA/Reportes/
	// Definir la ruta de la imagen
	rutaB := "./MIA/Particiones/" // Ruta base donde se guardarán los reportes

	// Escribir el contenido en un archivo .txt
	pd := filepath.Join(rutaB, "part.txt")
	b := []byte(content)
	err_ = ioutil.WriteFile(pd, b, 0644)
	if err_ != nil {
		log.Fatal(err_)
	}

}

// Función para recibir el nombre del disco desde el frontend y obtener las particiones asociadas
func RecibirNombreDisco(w http.ResponseWriter, r *http.Request) {
	// Decodificar el JSON que viene del frontend con el nombre del disco
	var requestData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener el nombre del disco del mapa
	nombreDisco := requestData["nombreDisco"]

	log.Println(nombreDisco)
	MostrarParticiones(nombreDisco)
}

/*
// Función para enviar la lista de particiones al frontend
func EnviarListaParticiones(w http.ResponseWriter, r *http.Request) {
	particiones := MandarParticiones()
	// Estructurar la lista de particiones en formato JSON
	jsonData, err := json.Marshal(particiones)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar el encabezado de respuesta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Enviar la lista de particiones como respuesta
	w.Write(jsonData)
}

// Función adaptadora para enviar la lista de particiones al frontend
func AdaptadorEnviarListaParticiones(particiones []Particion) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Llamar a la función para enviar la lista de particiones
		EnviarListaParticiones(particiones, w)
	}
}

*/
