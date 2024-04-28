import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom'; 
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCompactDisc } from '@fortawesome/free-solid-svg-icons';


const Pantalla2 = () => {
    const [nombresArchivos, setNombresArchivos] = useState([]);
    const [cargando, setCargando] = useState(true);
    const [particionesMontadas, setParticionesMontadas] = useState([]);

    useEffect(() => {
        // Hacer la solicitud al backend para obtener los nombres de los archivos (discos)
        axios.get('http://localhost:8080/obtenerdiscos')
            .then(response => {
                setNombresArchivos(response.data);
                setCargando(false);
            })
            .catch(error => {
                console.error('Error al obtener los nombres de archivos:', error);
                setCargando(false);
            });
    }, []);

    // Función para manejar el clic del botón para mostrar particiones montadas
    const mostrarParticionesMontadas = () => {
        axios.get('http://localhost:8080/obtenerparticionesmontadas')
            .then(response => {
                setParticionesMontadas(response.data);
            })
            .catch(error => {
                console.error('Error al obtener las particiones montadas:', error);
                // Limpiar la lista de particiones montadas en caso de error
                setParticionesMontadas([]);
            });
    };

    // Mostrar un mensaje de carga mientras se obtienen los datos
    if (cargando) {
        return <div>Cargando...</div>;
    }

    // Verificar si nombresArchivos es nulo o está vacío
    if (!nombresArchivos || nombresArchivos.length === 0) {
        return <div><h1>No hay Discos Creados</h1></div>;
    }

    return (
        <div>
            <h1>DISCOS, PARTICIONES Y LOGIN</h1>
            <h2>Nombres de Archivos:</h2>
            <ul>
                {nombresArchivos.map(nombre => (
                    <li key={nombre}>
                        {/* Utilizar Link para redirigir al usuario a una nueva pantalla */}
                        <Link to={`/particiones/${nombre}`}>
                            <FontAwesomeIcon icon={faCompactDisc} className="icono-disco" /> {nombre}
                        </Link>
                    </li>
                ))}
            </ul>
            <br></br>
            <button
                type="button"
                className="button"
                style={{ borderRadius: '20px' }}
                onClick={mostrarParticionesMontadas}
            >
                Mostrar Particiones Montadas
            </button>  

            {/* Mostrar la lista de particiones montadas */}
            <div>
                <h2>Particiones Montadas:</h2>
                {particionesMontadas.length > 0 ? (
                    <ul>
                        {particionesMontadas.map((particion, index) => (
                            <li key={index}>
                                <Link to={`/login/${particion.id}`}>{/* Agregar enlace */}
                                <span>id: {particion.id} || nombre: {particion.nombre}</span>
                                </Link>
                            </li>
                        ))}
                    </ul>
                ) : (
                    <p>No hay particiones montadas</p>
                )}
            </div>         
        </div>
    );
}

export default Pantalla2;