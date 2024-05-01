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

    const mostrarParticionesMontadas = () => {
        axios.get('http://localhost:8080/obtenerparticionesmontadas')
            .then(response => {
                setParticionesMontadas(response.data || [] );
            })
            .catch(error => {
                console.error('Error al obtener las particiones montadas:', error);
                setParticionesMontadas([]);
            });
    };

    if (cargando) {
        return <div>Cargando...</div>;
    }

    if (!nombresArchivos || nombresArchivos.length === 0) {
        return <div><h1>No hay Discos Creados</h1></div>;
    }

    const mandarnombredisco = async (nombreDisco) => {
        const data = {
            nombreDisco: nombreDisco
        };
    
        try {
            const response = await axios.post('http://localhost:8080/mandarnombredisco', data);
            console.log(response.data);
            console.log("Se manda a backend el disco");
    
        } catch (error) {
            console.error('Error al enviar los datos:', error);
        }
    };
    

    return (
        <div style={{ fontSize: '1.5em' }}>
            <h1 style={{ fontSize: '2em' }}>DISCOS, PARTICIONES Y LOGIN</h1>
            <h2 style={{ fontSize: '1.5em' }}>Nombres de Discos:</h2>
            <ul>
                {nombresArchivos.map(nombre => (
                    <li key={nombre}>
                        <Link 
                            to={`/particiones/${nombre}`}
                            onClick={() => mandarnombredisco(nombre)} 
                        >
                            <FontAwesomeIcon icon={faCompactDisc} className="icono-disco" style={{ fontSize: '2em' }} /> {nombre}
                        </Link>
                    </li>
                ))}
            </ul>
            <br></br>

            <div>
                <h2 style={{ fontSize: '1.5em' }}>Particiones Montadas:</h2>
                {particionesMontadas.length > 0 ? (
                    <ul>
                        {particionesMontadas.map((particion, index) => (
                            <li key={index}>
                                <Link to={`/login/${particion.id}`}>
                                    <span style={{ fontSize: '1.2em' }}>id: {particion.id} || nombre: {particion.nombre}</span>
                                </Link>
                            </li>
                        ))}
                    </ul>
                ) : (
                    <p>No hay particiones montadas</p>
                )}
            </div>
            <button
                type="button"
                className="button"
                style={{ borderRadius: '20px', fontSize: '1.2em' }}
                onClick={mostrarParticionesMontadas}
            >
                Mostrar Particiones Montadas
            </button>  
         
        </div>
    );
}

export default Pantalla2;
