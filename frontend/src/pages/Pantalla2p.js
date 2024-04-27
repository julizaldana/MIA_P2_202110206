import React from 'react';
import { Link, useParams } from 'react-router-dom';

const Pantalla2p = () => {
    // Obtener el parámetro de la URL (nombre del disco)
    const { nombreDisco } = useParams();

    return (
        <div>
            <h1> PARTICIONES </h1>
            <h2>Nombre del Disco: {nombreDisco}</h2>
            <button
                type="button"
                className="button"
                style={{ borderRadius: '20px' }}
            >
                <Link to="/pantalla2">Regresar</Link>
            </button>
        </div>
    );
}

export default Pantalla2p;
