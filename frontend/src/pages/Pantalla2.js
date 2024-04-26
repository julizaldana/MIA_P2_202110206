import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCompactDisc } from '@fortawesome/free-solid-svg-icons';


const Pantalla2 = () => {
    const [nombresArchivos, setNombresArchivos] = useState([]);
    const [cargando, setCargando] = useState(true);

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
                        <FontAwesomeIcon icon={faCompactDisc} className="icono-disco" /> {nombre}
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default Pantalla2;
