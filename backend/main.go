package main

import (
	"MIA_P2_202110206/Comandos"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//PROYECTO 2 - MANEJO E IMPLEMENTACIÓN DE ARCHIVOS
//NOMBRE: JULIO ALEJANDRO ZALDAÑA RÍOS 		CARNET: 202110206

var logued = false //variable booleana para verificar si un usuario estará logueado en su sesión

//gorilla mux sirve para levantar un servidor con golang - go get -u github.com/gorilla/mux
//librería cors permite cualquier ingreso de peticiones desde cualquier puerto. go get github.com/rs/cors

//Estructura para recibir datos del front, para comandos.
type DatosEntrada struct {
	Comandos []string `json:"comandos"`
}

//Con el main se declara el servidor.

func main() {
	router := mux.NewRouter() // declarar router

	// Ruta para el endpoint "/analizador"
	router.HandleFunc("/analizador", analizador).Methods("POST")
	// Ruta para el endpoint "/obtenerdiscos" para obtener los datos de los discos
	router.HandleFunc("/obtenerdiscos", obtenerNombresArchivos).Methods("GET")

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
	Exec("./prueba.script")
	fmt.Fprintf(w, "Script ejecutado exitosamente")

	// Devuelve una respuesta exitosa
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(" - Datos recibidos correctamente en Backend"))
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

func funciones(token string, tks []string) {
	if token != "" {
		if Comandos.Comparar(token, "EXECUTE") {
			fmt.Println("=*=*=*=*=*=*= FUNCION EXECUTE =*=*=*=*=*=*=")
			FuncionExec(tks)
		} else if Comandos.Comparar(token, "MKDISK") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKDISK =*=*=*=*=*=*=")
			Comandos.ValidarDatosMKDISK(tks)
		} else if Comandos.Comparar(token, "RMDISK") {
			fmt.Println("=*=*=*=*=*=*= FUNCION RMDISK =*=*=*=*=*=*=*=")
			Comandos.RMDISK(tks)
		} else if Comandos.Comparar(token, "FDISK") {
			fmt.Println("=*=*=*=*=*=*= FUNCION FDISK =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosFDISK(tks)
		} else if Comandos.Comparar(token, "REP") {
			fmt.Println("=*=*=*=*=*=*= FUNCION REP =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosREP(tks)
		} else if Comandos.Comparar(token, "MOUNT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MOUNT =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosMOUNT(tks)
		} else if Comandos.Comparar(token, "UNMOUNT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION UNMOUNT =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosUNMOUNT(tks)
		} else if Comandos.Comparar(token, "MKFS") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKFS =*=*=*=*=*=*=*=")
			Comandos.ValidarDatosMKFS(tks)
		} else if Comandos.Comparar(token, "LOGIN") {
			fmt.Println("=*=*=*=*=*=*= FUNCION LOGIN =*=*=*=*=*=*=*=")
			if logued {
				Comandos.Error("LOGIN", "Ya hay un usuario en línea.")
				return
			} else {
				logued = Comandos.ValidarDatosLOGIN(tks)
			}
		} else if Comandos.Comparar(token, "LOGOUT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION LOGOUT =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("LOGOUT", "Aún no se ha iniciado sesión.")
				return
			} else {
				logued = Comandos.CerrarSesion()
			}
		} else if Comandos.Comparar(token, "MKGRP") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKGRP =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("MKGRP", "Aún no se ha iniciado sesión.")
				return
			} else {
				Comandos.ValidarDatosGrupos(tks, "MK")
			}
		} else if Comandos.Comparar(token, "RMGRP") {
			fmt.Println("=*=*=*=*=*=*= FUNCION RMGRP =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RMGRP", "Aún no se ha iniciado sesión.")
				return
			} else {
				Comandos.ValidarDatosGrupos(tks, "RM")
			}
		} else if Comandos.Comparar(token, "MKUSR") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKUSR =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("MKUSR", "Aún no se ha iniciado sesión.")
				return
			} else {
				Comandos.ValidarDatosUsers(tks, "MK")
			}
		} else if Comandos.Comparar(token, "RMUSR") {
			fmt.Println("=*=*=*=*=*=*= FUNCION RMUSR =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RMUSR", "Aún no se ha iniciado sesión.")
				return
			} else {
				Comandos.ValidarDatosUsers(tks, "RM")
			}
		} else if Comandos.Comparar(token, "RMUSR") {
			fmt.Println("=*=*=*=*=*=*= FUNCION RMUSR =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RMUSR", "Aún no se ha iniciado sesión.")
				return
			} else {
				Comandos.ValidarDatosUsers(tks, "RM")
			}
		} else if Comandos.Comparar(token, "MKDIR") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKDIR =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("MKDIR", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				var p string
				particion := Comandos.GetMount("MKDIR", Comandos.Logged.Id, &p)
				Comandos.ValidarDatosMKDIR(tks, particion, p)
			}
		} else if Comandos.Comparar(token, "MKFILE") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MKFILE =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("MKFILE", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				Comandos.ValidarDatosMKFILE(tks)
			}
		} else if Comandos.Comparar(token, "CAT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION CAT =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("CAT", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				//Comandos.ValidarDatosCAT(tks)
			}
		} else if Comandos.Comparar(token, "REMOVE") {
			fmt.Println("=*=*=*=*=*=*= FUNCION REMOVE =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("REMOVE", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				//Comandos.ValidarDatosREMOVE(tks)
			}
		} else if Comandos.Comparar(token, "EDIT") {
			fmt.Println("=*=*=*=*=*=*= FUNCION EDIT =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("EDIT", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				//Comandos.ValidarDatosEDIT(tks)
			}
		} else if Comandos.Comparar(token, "RENAME") {
			fmt.Println("=*=*=*=*=*=*= FUNCION RENAME =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RENAME", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				//Comandos.ValidarDatosRENAME(tks)
			}
		} else if Comandos.Comparar(token, "COPY") {
			fmt.Println("=*=*=*=*=*=*= FUNCION COPY =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RENAME", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				//Comandos.ValidarDatosCOPY(tks)
			}
		} else if Comandos.Comparar(token, "MOVE") {
			fmt.Println("=*=*=*=*=*=*= FUNCION MOVE =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RENAME", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				//Comandos.ValidarDatosMOVE(tks)
			}
		} else if Comandos.Comparar(token, "FIND") {
			fmt.Println("=*=*=*=*=*=*= FUNCION FIND =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RENAME", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				//Comandos.ValidarDatosFIND(tks)
			}
		} else if Comandos.Comparar(token, "CHMOD") {
			fmt.Println("=*=*=*=*=*=*= FUNCION CHMOD =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RENAME", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				//Comandos.ValidarDatosCHMOD(tks)
			}
		} else if Comandos.Comparar(token, "CHGRP") {
			fmt.Println("=*=*=*=*=*=*= FUNCION CHGRP =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("CHGRP", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				Comandos.ValidarDatosCHGRP(tks)
			}
		} else if Comandos.Comparar(token, "CHOWN") {
			fmt.Println("=*=*=*=*=*=*= FUNCION CHOWN =*=*=*=*=*=*=*=")
			if !logued {
				Comandos.Error("RENAME", "Aún no se ha iniciado sesión con ningún usuario.")
				return
			} else {
				//Comandos.ValidarDatosCHOWN(tks)
			}
		} else {
			Comandos.Error("ANALIZADOR", "No se reconoce el comando \""+token+"\"")
		}
	}
}

func FuncionExec(tokens []string) {
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
	Exec(path)
}

func Exec(path string) {
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
			funciones(tk, tokens)
		}
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error al leer el archivo: %s", err)
	}
	file.Close()
}
