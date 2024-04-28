import React from 'react';
import axios from 'axios';
import { Link, useParams } from 'react-router-dom';

const Pantalla2Fs = () => {
    // Obtener el parámetro de la URL 
    const { user } = useParams();

    const handleSubmit = async () => {
        const comando = `logout`;
        const data = {
            comandos: [comando]
        };

        try {
            const response = await axios.post('http://localhost:8080/logout', data);
            console.log(response.data);
            console.log("Se mandó al backend el comando de logout");
        } catch (error) {
            console.error('Error al enviar los datos:', error);
            // Manejo del error
        }
    };

    return (
        <div>
            <h1> SISTEMA DE ARCHIVOS </h1>
            <h2>Usuario Logueado: {user}</h2>
            <button
                type="button"
                className="btn btn-danger"
                style={{ borderRadius: '20px' }}
                onClick={handleSubmit}
            >
                <Link to="/pantalla2">LOGOUT</Link>
            </button>
        </div>
    );
}

export default Pantalla2Fs;
