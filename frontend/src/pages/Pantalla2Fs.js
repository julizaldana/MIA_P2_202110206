import React, { useState } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const Pantalla2Fs = () => {
    // Obtener el parámetro de la URL 
    const { user } = useParams();
    const [logout, setLogout] = useState(false); // Estado para controlar el logout

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

    return (
        <div>
            <h1> SISTEMA DE ARCHIVOS </h1>
            <h2>Usuario Logueado: {user}</h2>
            <ul>
                <FontAwesomeIcon icon="fa-solid fa-folder" />
                <FontAwesomeIcon icon="fa-regular fa-file" />
            </ul>
            {!logout && ( // Renderizar solo si no se ha hecho logout
                <button
                    type="button"
                    className="btn btn-danger"
                    onClick={handleSubmit}
                >
                    LOGOUT
                </button>
            )}
        </div>
    );
}

export default Pantalla2Fs;
