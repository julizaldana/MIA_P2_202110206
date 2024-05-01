package Comandos

import (
	"MIA_P2_202110206/Structs"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"unsafe"
)

//FUNCIONAMIENTOS IMPORTANTES Y NECESARIOS PARA REALIZAR TODOS LOS PROCESOS

type Mensajes struct {
	Operacion string `json:"operacion"`
	Mensaje   string `json:"mensaje"`
}

var mensajes []Mensajes

// Función para enviar mensajes de éxito o error
func MandarMensaje(op string, mensaje string, w http.ResponseWriter) {
	respuesta := Mensajes{
		Operacion: op,
		Mensaje:   mensaje,
	}

	mensajes = append(mensajes, respuesta)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mensajes)
}

func MandarError(op string, mensaje string, w http.ResponseWriter) {
	// Construye el mensaje de error
	mensajeError := ("\tERROR " + op + "\n\tTIPO: " + mensaje)

	// Envía el mensaje de error al frontend
	MandarMensaje("ERROR", mensajeError, w)
}

func ObtenerMensajes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mensajes)
}

func Comparar(a string, b string) bool {
	if strings.ToUpper(a) == strings.ToUpper(b) {
		return true
	}
	return false
}

func Error(op string, mensaje string) {
	fmt.Println(Format(RED, "\tERROR "+op+"\n\tTIPO: "+mensaje))
}

func Mensaje(op string, mensaje string) {
	fmt.Println(Format(BLUE, "\tCOMANDO: "+op+"\n\tTIPO: "+mensaje))
}

func Confirmar(mensaje string) bool {
	fmt.Println(mensaje + "(y/n)")
	var respuesta string
	fmt.Scanln(&respuesta)
	if Comparar(respuesta, "y") {
		return true
	}
	return false
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func EscribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func leerDisco(path string) *Structs.MBR {
	m := Structs.MBR{}
	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
		return nil
	}
	file.Seek(0, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &m)
	if err_ != nil {
		Error("FDISK", "Error al leer el archivo")
		return nil
	}
	var mDir *Structs.MBR = &m
	return mDir
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

const escape = "\x1b"

const (
	NONE = iota
	RED
	GREEN
	YELLOW
	BLUE
	PURPLE
)

func color(c int) string {
	if c == NONE {
		return fmt.Sprintf("%s[%dm", escape, c)
	}

	return fmt.Sprintf("%s[3%dm", escape, c)
}

func Format(c int, text string) string {
	return color(c) + text + color(NONE)
}
