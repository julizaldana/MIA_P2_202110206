import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Link, useParams } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFloppyDisk } from '@fortawesome/free-solid-svg-icons';

const Pantalla2p = () => {
    const { nombreDisco } = useParams();
    const [particiones, setParticiones] = useState([]);
    const [cargando, setCargando] = useState(true);

    useEffect(() => {
        // Hacer la solicitud al backend para obtener la lista de particiones
        axios.get('http://localhost:8080/enviarparticiones')
            .then(response => {
                setParticiones(response.data || []); // Si response.data es null, establece particiones como un array vacío
                setCargando(false);
            })
            .catch(error => {
                console.error('Error al obtener las particiones:', error);
                setCargando(false);
            });
    }, []);

    return (
        <div style={{ fontSize: '1.5em' }}>
            <h1 style={{ fontSize: '2em' }}> PARTICIONES </h1>
            <h2 style={{ fontSize: '1.5em' }}>Nombre del Disco: {nombreDisco}</h2>
            {cargando ? (
                <div>Cargando...</div>
            ) : particiones && particiones.length > 0 ? (
                <ul>
                    {particiones.map((particion, index) => (
                        <li key={index}>
                            <FontAwesomeIcon icon={faFloppyDisk} /> {particion}
                        </li>
                    ))}
                </ul>
            ) : (
                <p>No existen particiones para este disco.</p>
            )}
            <button
                type="button"
                className="button"
                style={{ borderRadius: '20px', fontSize: '1.2em' }}
            >
                <Link to="/pantalla2">Regresar</Link>
            </button>
        </div>
    );
}

export default Pantalla2p;
