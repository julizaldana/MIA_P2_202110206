import React, { useState } from 'react';
import axios from 'axios'; // Se importa Axios para realizar la solicitud HTTP
import "../css/Pantalla1.css";

//PANTALLA 1 - PARA INGRESAR COMANDOS

const Pantalla1 = () => {
    // Estado para almacenar el valor del textarea
    const [comando, setComando] = useState('');

    // Función para manejar el envío del formulario
    const handleSubmit = async () => {
        // Estructura los datos en el formato requerido por el backend
        const data = {
            comandos: [comando]
        };

        try {
            // Realiza la solicitud POST al backend
            const response = await axios.post('http://localhost:8080/analizador', data);
            console.log(response.data); // Imprime la respuesta del backend en la consola
        } catch (error) {
            console.error('Error al enviar los datos:', error);
        }
    };

    return (
        <div style={{ textAlign: 'center' }}>
            <h1>COMANDOS</h1>
            <div className="input-container">
                <p className="bash-text">
                    <span className="user">user</span><span className="vm">@mia-go-file</span>:
                    <span className="char">~</span>$
                </p>
                <textarea 
                    className="input textarea" 
                    placeholder="Ingresar comando..."
                    value={comando} // Asigna el valor del textarea al estado
                    onChange={(e) => setComando(e.target.value)} // Actualiza el estado al cambiar el valor del textarea
                ></textarea>
            </div>
            <br></br>
            <button 
                type="button" 
                className="button2" 
                style={{ borderRadius: '20px' }}
                onClick={handleSubmit} // Maneja el envío del formulario al hacer clic en el botón
            >
                Enviar
            </button>
        </div>
    )
}

export default Pantalla1;
