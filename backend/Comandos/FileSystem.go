package Comandos

import (
	"MIA_P2_202110206/Structs"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

//LEER DEL ARCHIVO BINARIO COMO EL REPORTE TREE, ARCHIVOS Y CARPETAS, HACER DOS ENDPOINTS UNO PARA ARCHIVOS Y OTRO PARA CARPETAS, Y SEGUN LOS ENDPOINTS METODOS GET
//OBTENER LA INFORMAICON Y USAR ICONOS DE FOLDER Y FILE PARA EL SISTEMA DE ARCHIVOS.

func MostrarArchivos(id string) {
	var pth string
	spr := Structs.NewSuperBloque()
	inode := Structs.NewInodos()
	partition := GetMount("REP", id, &pth)

	if partition.Part_start == -1 {
		return
	}

	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))

	if err != nil {
		Error("REP", "No se ha encontrado el disco.")
		return
	}

	file.Seek(partition.Part_start, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.SuperBloque{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &spr)
	if err_ != nil {
		Error("REP", "Error al leer el archivo")
		return
	}

	file.Seek(spr.S_inode_start, 0)
	data = leerBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
	buffer = bytes.NewBuffer(data)
	err_ = binary.Read(buffer, binary.BigEndian, &inode)
	if err_ != nil {
		Error("REP", "Error al leer el archivo")
		return
	}

	freeI := GetFree(spr, pth, "BI")

	content := ""
	for i := 0; i < int(freeI); i++ {
		//atime := arregloString(inode.I_atime)
		//ctime := arregloString(inode.I_ctime)
		//mtime := arregloString(inode.I_mtime)
		//content += "inode" + strconv.Itoa(i) + "  [label = <<table>\n" +
		//"<tr><td COLSPAN = '2' BGCOLOR=\"#000080\">" +
		//"<font color=\"white\">INODO " + strconv.Itoa(i) + "</font>" +
		//"</td></tr>\n " +
		//"<tr><td BGCOLOR=\"#87CEFA\">NOMBRE</td><td BGCOLOR=\"#87CEFA\" >VALOR</td></tr>\n" +
		//"<tr><td>i_uid</td><td>" + strconv.Itoa(int(inode.I_uid)) + "</td></tr>\n" +
		//"<tr><td>i_gid</td><td>" + strconv.Itoa(int(inode.I_gid)) + "</td></tr>\n" +
		//"<tr><td>i_size</td><td>" + strconv.Itoa(int(inode.I_s)) + "</td></tr>\n" +
		//"<tr><td>i_atime</td><td>" + atime + "</td></tr>\n" +
		//"<tr><td>i_ctime</td><td>" + ctime + "</td></tr>\n" +
		//"<tr><td>i_mtime</td><td>" + mtime + "</td></tr>\n"
		for j := 0; j < 16; j++ {
			content += ""
		}
		content += ""

		if inode.I_type == 0 {

			for j := 0; j < 16; j++ {
				if inode.I_block[j] != -1 {
					bloquesUsados = append(bloquesUsados, inode.I_block[j])
					contadorBloques++
					if existeEnArreglo(bloquesUsados, inode.I_block[j]) == 1 {
						foldertmp := Structs.BloquesCarpetas{}

						file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*inode.I_block[j]+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*32*inode.I_block[j], 0)
						data = leerBytes(file, int(unsafe.Sizeof(Structs.BloquesCarpetas{})))
						buffer = bytes.NewBuffer(data)
						err_ = binary.Read(buffer, binary.BigEndian, &foldertmp)
						if err_ != nil {
							Error("REP", "Error al leer el archivo")
							return
						}

						if foldertmp.B_content[2].B_inodo == -1 {
							continue
						}

						content += ""

						for k := 0; k < 4; k++ {
							ctmp := ""
							for name := 0; name < len(foldertmp.B_content[k].B_name); name++ {
								if foldertmp.B_content[k].B_name[name] != 0 {
									ctmp += string(foldertmp.B_content[k].B_name[name])
								}
								//content += ctmp + "\n"
							}
							//content += ctmp + "\n" //DEVUELVE ARCHIVOS
						}

						for b := 0; b < 4; b++ {
							if foldertmp.B_content[b].B_inodo != -1 {
								ctmp := ""
								temp := ""
								for name := 0; name < len(foldertmp.B_content[b].B_name); name++ {
									if foldertmp.B_content[b].B_name[name] != 0 {
										ctmp += string(foldertmp.B_content[b].B_name[name])
									}
								}
								if ctmp != "." && ctmp != ".." {
									temp += ctmp + "\n"
								}
								content += ctmp + "\n"
							}
						}
					}
				}
			}

		} else {
			for j := 0; j < 16; j++ {
				if inode.I_block[j] != -1 {
					if j < 16 {
						var contador int64 = 0
						var posicion int
						for {
							foldertmp := Structs.NewBloquesCarpetas()
							file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*contador+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*32*contador, 0)
							data = leerBytes(file, int(unsafe.Sizeof(Structs.BloquesCarpetas{})))
							buffer = bytes.NewBuffer(data)
							err_ = binary.Read(buffer, binary.BigEndian, &foldertmp)

							if err_ != nil {
								Error("REP", "Error al leer el archivo")
								return
							}
							salir := false
							for l := 0; l < 4; l++ {
								if foldertmp.B_content[l].B_inodo == inode.I_block[0] {
									posicion = l
									salir = true
									break
								}
							}
							if salir {
								break
							}
							contador++
						}
						if posicion == 2 || posicion == 0 {
							file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*contador+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*32*contador+int64(unsafe.Sizeof(Structs.BloquesCarpetas{})), 0)
							for k := 0; k < 16; k++ {
								contadorArchivos++
								if inode.I_block[k] == -1 {
									break
								}
								content += ""
								filetmp := Structs.BloquesArchivos{}
								data = leerBytes(file, int(unsafe.Sizeof(Structs.BloquesArchivos{})))
								buffer = bytes.NewBuffer(data)
								err_ = binary.Read(buffer, binary.BigEndian, &filetmp)
								if err_ != nil {
									Error("REP", "Error al leer el archivo")
									return
								}

								contenido := ""
								for arch := 0; arch < len(filetmp.B_content); arch++ {
									if filetmp.B_content[arch] != 0 {
										contenido += string(filetmp.B_content[arch])
									}
								}
								content += ""
							}
						} else if posicion == 3 {
							file.Seek(spr.S_block_start+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))*contador+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*32*contador+int64(unsafe.Sizeof(Structs.BloquesCarpetas{}))+int64(unsafe.Sizeof(Structs.BloquesArchivos{}))*int64(16), 0)
							for k := 0; k < 16; k++ {
								contadorArchivos++
								if inode.I_block[k] == -1 {
									break
								}
								content += ""
								filetmp := Structs.BloquesArchivos{}
								data = leerBytes(file, int(unsafe.Sizeof(Structs.BloquesArchivos{})))
								buffer = bytes.NewBuffer(data)
								err_ = binary.Read(buffer, binary.BigEndian, &filetmp)
								if err_ != nil {
									Error("REP", "Error al leer el archivo")
									return
								}

								contenido := ""
								for arch := 0; arch < len(filetmp.B_content); arch++ {
									if filetmp.B_content[arch] != 0 {
										contenido += string(filetmp.B_content[arch])
									}
								}
								content += ""
							}
						}
						break
					}
				} else {
					break
				}
			}
		}
		file.Seek(spr.S_inode_start+int64(unsafe.Sizeof(Structs.Inodos{}))*int64(i+1), 0)
		data = leerBytes(file, int(unsafe.Sizeof(Structs.Inodos{})))
		buffer = bytes.NewBuffer(data)
		err_ = binary.Read(buffer, binary.BigEndian, &inode)
		if err_ != nil {
			Error("REP", "Error al leer el archivo")
			return
		}
	}
	file.Close()
	content += ""

	// Definir la ruta de la imagen
	rutaBase := "./MIA/Almacenamiento/" // Ruta base donde se guardarán los nombres de archivos

	// Escribir el contenido en un archivo de texto.
	pd := filepath.Join(rutaBase, "archivos.txt")
	// Abrir el archivo en modo de escritura con la bandera os.O_APPEND
	file, err = os.OpenFile(pd, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close() // Asegúrate de cerrar el archivo al final de la función

	// Escribir el contenido en el archivo
	if _, err := file.WriteString(content); err != nil {
		log.Fatal(err)
	}

}

// Función para recibir el id de la particion desde el frontend y obtener los archivos y carpetas asociadas
func RecibirIdParticion(w http.ResponseWriter, r *http.Request) {
	// Decodificar el JSON que viene del frontend con el nombre
	var requestData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener el nombre de la particion
	idparticion := requestData["idParticion"]

	log.Println(idparticion)
	MostrarArchivos(idparticion)
}
