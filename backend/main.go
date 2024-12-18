package main

import (
	"MIA_P2_202110206/Comandos"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//PROYECTO 2 - MANEJO E IMPLEMENTACIÓN DE ARCHIVOS
//NOMBRE: JULIO ALEJANDRO ZALDAÑA RÍOS 		CARNET: 202110206

var logued = false //variable booleana para verificar si un usuario estará logueado en su sesión

//gorilla mux sirve para levantar un servidor con golang - go get -u github.com/gorilla/mux
//librería cors permite cualquier ingreso de peticiones desde cualquier puerto. go get github.com/rs/cors

// Estructura para recibir datos del front, para comandos.
type DatosEntrada struct {
	Comandos []string `json:"comandos"`
}

// Estructura para recibir datos del front, para comandos.
type NombreDisco struct {
	NombreDsk []string `json:"nombredsk"`
}

//Con el main se declara el servidor.

func main() {
	router := mux.NewRouter() // declarar router

	// Ruta para el endpoint "/analizador"
	router.HandleFunc("/analizador", analizador).Methods("POST")
	// Ruta para el endpoint "/obtenerdiscos" para obtener los datos de los discos
	router.HandleFunc("/obtenerdiscos", obtenerNombresArchivos).Methods("GET")
	//Ruta para notificaciones "/notificacion"
	router.HandleFunc("/notificacion", Comandos.ObtenerMensajes).Methods("GET")
	//Ruta para mandar los datos de las particiones
	router.HandleFunc("/mandarnombredisco", Comandos.RecibirNombreDisco).Methods("POST")
	router.HandleFunc("/enviarparticiones", sendPartitions).Methods("GET")
	// Ruta para obtener la lista de particiones montadas
	//router.HandleFunc("/obtenerparticionesmontadas", Comandos.ObtenerParticionesMontadas).Methods("GET")
	router.HandleFunc("/obtenerreportes", obtenerReportes).Methods("GET")
	router.HandleFunc("/mandaridparticion", Comandos.RecibirIdParticion).Methods("POST")
	router.HandleFunc("/enviararchivos", sendFiledata).Methods("GET")
	router.HandleFunc("/iniciarsesion", analizador).Methods("POST")
	router.HandleFunc("/logout", analizador).Methods("POST")
	router.HandleFunc("/obtenerparticionesmontadas", func(w http.ResponseWriter, r *http.Request) {
		// Obtener la lista de particiones montadas
		particionesMontadas := Comandos.ListaPartMount()

		// Imprimir en la consola la lista de particiones montadas en formato JSON
		particionesMontadasJSON, err := json.Marshal(particionesMontadas)
		if err != nil {
			log.Println("Error al convertir la lista de particiones montadas a JSON:", err)
		} else {
			log.Println(string(particionesMontadasJSON))
		}

		// Configurar el encabezado de respuesta como JSON
		w.Header().Set("Content-Type", "application/json")

		// Enviar la lista de particiones montadas como JSON en la respuesta
		if err := json.NewEncoder(w).Encode(particionesMontadas); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// Manejador CORS
	handler := cors.Default().Handler(router)

	log.Println("Server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

/*
func allowCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		handler.ServeHTTP(w, r)
	})
}

func inicial(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>¡Hola Desde el servidor!</h1>")
}
*/
// Función para enviar la lista de informacion del sistema al frontend
func sendFiledata(w http.ResponseWriter, r *http.Request) {
	// Obtener la lista de archivos en la carpeta
	archivos, err := ioutil.ReadDir("./MIA/Almacenamiento")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Crear una estructura para almacenar los nombres de las particiones
	var datafile []string

	// Leer cada archivo en la carpeta
	for _, archivo := range archivos {
		// Ignorar los directorios
		if !archivo.IsDir() {
			// Leer el contenido del archivo línea por línea
			contenido, err := ioutil.ReadFile("./MIA/Almacenamiento/" + archivo.Name())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Dividir el contenido en líneas
			lineas := strings.Split(string(contenido), "\n")

			// Agregar cada línea no vacía a la lista de particiones
			for _, linea := range lineas {
				linea = strings.TrimSpace(linea) // Eliminar espacios en blanco al inicio y al final de la línea
				if linea != "" && linea != "." && linea != ".." {
					datafile = append(datafile, linea)
				}
			}
		}
	}

	// Estructurar la lista de particiones en formato JSON
	jsonData, err := json.Marshal(datafile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar el encabezado de respuesta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Enviar la lista de particiones como respuesta
	w.Write(jsonData)
}

// Función para enviar la lista de particiones al frontend
func sendPartitions(w http.ResponseWriter, r *http.Request) {
	// Obtener la lista de archivos en la carpeta
	archivos, err := ioutil.ReadDir("./MIA/Particiones")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Crear una estructura para almacenar los nombres de las particiones
	var particiones []string

	// Leer cada archivo en la carpeta
	for _, archivo := range archivos {
		// Ignorar los directorios
		if !archivo.IsDir() {
			// Leer el contenido del archivo línea por línea
			contenido, err := ioutil.ReadFile("./MIA/Particiones/" + archivo.Name())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Dividir el contenido en líneas
			lineas := strings.Split(string(contenido), "\n")

			// Agregar cada línea no vacía a la lista de particiones
			for _, linea := range lineas {
				linea = strings.TrimSpace(linea) // Eliminar espacios en blanco al inicio y al final de la línea
				if linea != "" {
					particiones = append(particiones, linea)
				}
			}
		}
	}

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

func analizador(w http.ResponseWriter, r *http.Request) {
	var datos DatosEntrada                        //se realiza una variable con estructura
	err := json.NewDecoder(r.Body).Decode(&datos) //se mete el body del json a esa variable
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = guardarDatos("./prueba.script", datos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ejecutar el archivo de script con el comando Exec como lo hacia en el proyecto 1
	Exec("./prueba.script", w)
	//fmt.Fprintf(w, "Script ejecutado exitosamente")

	// Devuelve una respuesta exitosa
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte(" - Datos recibidos correctamente en Backend"))
	//Comandos.MandarMensaje("ANALIZADOR", "Comandos procesados correctamente", w)
}

/*
func verParticiones(w http.ResponseWriter, r *http.Request) {
	var disco NombreDisco
	err :=
	err := json.NewDecoder(r.Body).Decode(&disco) //se mete el body del json a esa variable
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Comandos.MandarParticiones(disco, w)

	w.WriteHeader(http.StatusOK)

}
*/

// FUNCION PARA OBTENER REPORTES DE BACKEND
func obtenerReportes(w http.ResponseWriter, r *http.Request) {
	// Obtener la lista de archivos en la carpeta
	archivos, err := ioutil.ReadDir("./MIA/Reportes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Crear una estructura para almacenar los datos de los reportes
	type Reporte struct {
		Nombre    string `json:"nombre"`
		Contenido string `json:"contenido"`
	}
	var reportes []Reporte

	// Iterar sobre cada archivo y obtener su contenido codificado en base64
	for _, archivo := range archivos {
		if !archivo.IsDir() && (strings.HasSuffix(archivo.Name(), ".jpg") || strings.HasSuffix(archivo.Name(), ".png")) {
			contenido, err := ioutil.ReadFile(filepath.Join("./MIA/Reportes", archivo.Name()))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reporte := Reporte{
				Nombre:    archivo.Name(),
				Contenido: base64.StdEncoding.EncodeToString(contenido),
			}
			reportes = append(reportes, reporte)
		}
	}

	// Convertir la lista de reportes a JSON y enviarla como respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reportes)
}

func guardarDatos(archivo string, datos DatosEntrada) error {
	// Abrir o crear el archivo
	file, err := os.Create(archivo)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribir los comandos en el archivo
	for _, comando := range datos.Comandos {
		_, err := file.WriteString(strings.TrimSpace(comando) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func obtenerNombresArchivos(w http.ResponseWriter, r *http.Request) {
	archivos, err := ioutil.ReadDir("./MIA/Discos")
	if err != nil {
		// Manejar el error
		http.Error(w, "Error al leer los archivos", http.StatusInternalServerError)
		return
	}

	var nombres []string
	for _, archivo := range archivos {
		nombres = append(nombres, archivo.Name())
	}

	// Convertir la lista de nombres de los discos a JSON y enviarla como respuesta
	json.NewEncoder(w).Encode(nombres)
}

/*
func enviarListaParticiones(w http.ResponseWriter, r *http.Request) {
	// Convertir la lista de nombres de particiones a formato JSON
	jsonNombres, err := json.Marshal(nombresParticiones)
	if err != nil {
		Error("REP", "Error al convertir a JSON")
		return
	}

	var nombreDisco NombreDisco

	disco :=

		Comandos.MandarMensaje()

	// Escribir la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonNombres)
}
*/

func Comando(text string) string {
	var tkn string
	terminar := false
	for i := 0; i < len(text); i++ {
		if terminar {
			if string(text[i]) == " " || string(text[i]) == "-" {
				break
			}
			tkn += string(text[i])
		} else if string(text[i]) != " " && !terminar {
			if string(text[i]) == "#" {
				tkn = text
			} else {
				tkn += string(text[i])
				terminar = true
			}
		}
	}
	return tkn
}

func SepararTokens(texto string) []string {
	var tokens []string
	if texto == "" {
		return tokens
	}
	texto += " "
	var token string
	estado := 0
	for i := 0; i < len(texto); i++ {
		c := string(texto[i])
		if estado == 0 && c == "-" {
			estado = 1
		} else if estado == 0 && c == "#" {
			continue
		} else if estado != 0 {
			if estado == 1 {
				if c == "=" {
					estado = 2
				} else if c == " " {
					continue
				} else if (c == "P" || c == "p") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
					continue
				} else if (c == "R" || c == "r") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
					continue
				}
			} else if estado == 2 {
				if c == " " {
					continue
				}
				if c == "\"" {
					estado = 3
					continue
				} else {
					estado = 4
				}
			} else if estado == 3 {
				if c == "\"" {
					estado = 4
					continue
				}
			} else if estado == 4 && c == "\"" {
				tokens = []string{}
				continue
			} else if estado == 4 && c == " " {
				estado = 0
				tokens = append(tokens, token)
				token = ""
				continue
			}
			token += c
		}
	}
	return tokens
}

func funciones(token string, tks []string, w http.ResponseWriter) {
	if token != "" {
		if Comandos.Comparar(token, "EXECUTE") {
			fmt.Println("=*=*=*=*=*=*= FUNCION EXECUTE =*=*=*=*=*=*=")
			FuncionExec(tks, w)
		} else if Comandos.Comparar(token, "MKDISK") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKDISK =*=*=*=*=*=*=")
			Comandos.ValidarDatosMKDISK(tks, w)
		} else if Comandos.Comparar(token, "RMDISK") {
			fmt.Println("=*=*=*=*=*=*= FUNCION RMDISK =*=*=*=*=*=*=*=")
			Comandos.RMDISK(tks, w)
		} else if Comandos.Comparar(token, "FDISK") {
			fmt.Println("=*=*=*=*=*=*= FUNCION FDISK =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosFDISK(tks, w)
		} else if Comandos.Comparar(token, "REP") {
			fmt.Println("=*=*=*=*=*=*= FUNCION REP =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosREP(tks, w)
		} else if Comandos.Comparar(token, "MOUNT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MOUNT =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosMOUNT(tks, w)
		} else if Comandos.Comparar(token, "UNMOUNT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION UNMOUNT =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosUNMOUNT(tks, w)
		} else if Comandos.Comparar(token, "MKFS") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKFS =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosMKFS(tks, w)
		} else if Comandos.Comparar(token, "LOGIN") {
			fmt.Println("=*=*=*=*=*=*= FUNCION LOGIN =*=*=*=*=*=*=*=")
			if logued {
				Comandos.Error("LOGIN", "Ya hay un usuario en línea.")
				Comandos.MandarError("LOGIN", "No se puede hacer LOGIN porque ya hay un usuario en línea.", w)
				return
			} else {
				logued = Comandos.ValidarDatosLOGIN(tks, w)
			}
		} else if Comandos.Comparar(token, "LOGOUT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION LOGOUT =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("LOGOUT", "Aún no se ha iniciado sesión.")
				Comandos.MandarError("LOGOUT", "No se puede hacer un LOGOUT sin una sesión iniciada.", w)
				return
			} else {
				logued = Comandos.CerrarSesion(w)
			}
		} else if Comandos.Comparar(token, "MKGRP") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKGRP =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("MKGRP", "Aún no se ha iniciado sesión.")
				Comandos.MandarError("MKGRP", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				Comandos.ValidarDatosGrupos(tks, "MK", w)
			}
		} else if Comandos.Comparar(token, "RMGRP") {
			fmt.Println("=*=*=*=*=*=*= FUNCION RMGRP =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RMGRP", "Aún no se ha iniciado sesión.")
				Comandos.MandarError("RMGRP", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				Comandos.ValidarDatosGrupos(tks, "RM", w)
			}
		} else if Comandos.Comparar(token, "MKUSR") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKUSR =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("MKUSR", "Aún no se ha iniciado sesión.")
				Comandos.MandarError("MKUSR", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				Comandos.ValidarDatosUsers(tks, "MK", w)
			}
		} else if Comandos.Comparar(token, "RMUSR") {
			fmt.Println("=*=*=*=*=*=*= FUNCION RMUSR =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RMUSR", "Aún no se ha iniciado sesión.")
				Comandos.MandarError("RMUSR", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				Comandos.ValidarDatosUsers(tks, "RM", w)
			}
		} else if Comandos.Comparar(token, "MKDIR") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKDIR =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("MKDIR", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("MKDIR", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				var p string
				particion := Comandos.GetMount("MKDIR", Comandos.Logged.Id, &p)
				Comandos.ValidarDatosMKDIR(tks, particion, p, w)
			}
		} else if Comandos.Comparar(token, "MKFILE") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKFILE =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("MKFILE", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("MKFILE", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				var p string
				particion := Comandos.GetMount("MKFILE", Comandos.Logged.Id, &p)
				Comandos.ValidarDatosMKFILE(tks, particion, p, w)
			}
		} else if Comandos.Comparar(token, "CAT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION CAT =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("CAT", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("CAT", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				//Comandos.ValidarDatosCAT(tks)
			}
		} else if Comandos.Comparar(token, "REMOVE") {
			fmt.Println("=*=*=*=*=*=*= FUNCION REMOVE =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("REMOVE", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("REMOVE", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				//Comandos.ValidarDatosREMOVE(tks)
			}
		} else if Comandos.Comparar(token, "EDIT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION EDIT =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("EDIT", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("EDIT", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				//Comandos.ValidarDatosEDIT(tks)
			}
		} else if Comandos.Comparar(token, "RENAME") {
			fmt.Println("=*=*=*=*=*=*= FUNCION RENAME =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RENAME", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("RENAME", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				//Comandos.ValidarDatosRENAME(tks)
			}
		} else if Comandos.Comparar(token, "COPY") {
			fmt.Println("=*=*=*=*=*=*= FUNCION COPY =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("COPY", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("COPY", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				//Comandos.ValidarDatosCOPY(tks)
			}
		} else if Comandos.Comparar(token, "MOVE") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MOVE =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("MOVE", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("CHOWN", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				//Comandos.ValidarDatosMOVE(tks)
			}
		} else if Comandos.Comparar(token, "FIND") {
			fmt.Println("=*=*=*=*=*=*= FUNCION FIND =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("FIND", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("CHOWN", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				//Comandos.ValidarDatosFIND(tks)
			}
		} else if Comandos.Comparar(token, "CHMOD") {
			fmt.Println("=*=*=*=*=*=*= FUNCION CHMOD =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("CHMOD", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("CHOWN", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				//Comandos.ValidarDatosCHMOD(tks)
			}
		} else if Comandos.Comparar(token, "CHGRP") {
			fmt.Println("=*=*=*=*=*=*= FUNCION CHGRP =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("CHGRP", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("CHOWN", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				Comandos.ValidarDatosCHGRP(tks)
			}
		} else if Comandos.Comparar(token, "CHOWN") {
			fmt.Println("=*=*=*=*=*=*= FUNCION CHOWN =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("CHOWN", "Aún no se ha iniciado sesión con ningún usuario.")
				Comandos.MandarError("CHOWN", "No se ha iniciado sesión para ejecutar este comando.", w)
				return
			} else {
				//Comandos.ValidarDatosCHOWN(tks)
			}
		} else {
			Comandos.Error("ANALIZADOR", "No se reconoce el comando \""+token+"\"")
			Comandos.MandarError("ANALIZADOR", "No se reconoce el comando \""+token+"\"", w)
		}
	}
}

func FuncionExec(tokens []string, w http.ResponseWriter) {
	path := ""
	for i := 0; i < len(tokens); i++ {
		datos := strings.Split(tokens[i], "=")
		if Comandos.Comparar(datos[0], "path") {
			path = datos[1]
		}
	}
	if path == "" {
		Comandos.Error("EXECUTE", "Se requiere el parámetro \"path\" para este comando")
		return
	}
	Exec(path, w)
}

func Exec(path string, w http.ResponseWriter) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		texto := fileScanner.Text()
		texto = strings.TrimSpace(texto)
		tk := Comando(texto)
		if texto != "" {
			if Comandos.Comparar(tk, "pause") {
				fmt.Println("************************************** FUNCIÓN PAUSE **************************************")
				var pause string
				Comandos.Mensaje("PAUSE", "Presione \"enter\" para continuar...")
				fmt.Scanln(&pause)
				continue
			} else if string(texto[0]) == "#" {
				//fmt.Println("************************************** COMENTARIO **************************************")
				//Comandos.Mensaje("COMENTARIO", texto)
				continue
			}
			texto = strings.TrimLeft(texto, tk)
			tokens := SepararTokens(texto)
			funciones(tk, tokens, w)
		}
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error al leer el archivo: %s", err)
	}
	file.Close()
}
