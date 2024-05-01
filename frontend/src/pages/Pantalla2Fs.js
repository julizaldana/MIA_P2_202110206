import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Link, useParams } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFile } from '@fortawesome/free-solid-svg-icons';


const Pantalla2Fs = () => {
    // Obtener el parámetro de la URL 
    const { user } = useParams();
    const { idParticion } = useParams();
    const [logout, setLogout] = useState(false); // Estado para controlar el logout
    const [archivos, setArchivos] = useState([]);
    const [cargando, setCargando] = useState(true);

    useEffect(() => {
        // Hacer la solicitud al backend para obtener la lista de particiones
        axios.get('http://localhost:8080/enviararchivos', )
            .then(response => {
                setArchivos(response.data || []); // Si response.data es null, establece particiones como un array vacío
                setCargando(false);
            })
            .catch(error => {
                console.error('Error al obtener los archivos:', error);
                setCargando(false);
            });
    }, []);


    const handleSubmit = async () => {
        const comando = `logout`;
        const data = {
            comandos: [comando]
        };

        try {
            const response = await axios.post('http://localhost:8080/logout', data);
            console.log(response.data);
            console.log("Se mandó al backend el comando de logout");
            setLogout(true); // Cambiar el estado a true al hacer logout
        } catch (error) {
            console.error('Error al enviar los datos:', error);
            // Manejo del error
        }
    };


    const mandaridparticion = async (idParticion) => {
        const data = {
            idParticion: idParticion
        };
    
        try {
            const response = await axios.post('http://localhost:8080/mandaridparticion', data);
            console.log(response.data);
            console.log("Se manda a backend el disco");
    
        } catch (error) {
            console.error('Error al enviar los datos:', error);
        }
    };
    

    return (
        <div style={{ fontSize: '1.5em' }}>
            <h1> SISTEMA DE ARCHIVOS </h1>
            <h2>Usuario Logueado: {user}</h2>
            {!logout && ( // Renderizar solo si no se ha hecho logout
                <button
                    type="button"
                    className="btn btn-danger"
                    onClick={handleSubmit}
                    style={{ fontSize: '1em' }}
                >
                <Link to="/pantalla2">LOGOUT</Link>
                </button>
            )}
            <br></br>
            <br></br>
            <button
                    type="button"
                    className="btn btn-primary"
                    onClick={() => { mandaridparticion(idParticion)}}
                    style={{ fontSize: '1em' }}
                > RECARGAR SISTEMA
                </button>
            {cargando ? (
                <div>Cargando...</div>
            ) : archivos && archivos.length > 0 ? ( // Verifica si particiones no es null y tiene al menos un elemento
                <ul>
                    <br></br>
                    {archivos.map((archivo, index) => (
                        <li key={index}>
                            <FontAwesomeIcon icon={faFile} /> {archivo}
                        </li>
                    ))}
                </ul>
            ) : (
                <p>No existen archivos en el sistema de archivos.</p>
            )}

        </div>
    );
}

export default Pantalla2Fs;
