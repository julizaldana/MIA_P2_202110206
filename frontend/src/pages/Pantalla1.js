import React, { useState} from 'react';
import axios from 'axios'; // Se importa Axios para realizar la solicitud HTTP
import "../css/Pantalla1.css";

const Pantalla1 = () => {
    const [comando, setComando] = useState('');
    const [notificacion, setNotificacion] = useState({ tipo: '', mensaje: '' });

    // Función para manejar el envío del formulario
    const handleSubmit = async () => {
        const data = {
            comandos: [comando]
        };

        try {
            const response = await axios.post('http://localhost:8080/analizador', data);
            console.log(response.data);
            console.log("Se manda a backend los comandos");

            //obtener una respuesta del backend respecto al comando enviado
            obtenerMensajes();


        } catch (error) {
            console.error('Error al enviar los datos:', error);
            mostrarNotificacion('danger', 'Error al procesar los comandos');
        }
    };

    // Función para obtener mensajes del backend
    const obtenerMensajes = async () => {
        try {
            const response = await axios.get('http://localhost:8080/notificacion');
            console.log(response.data);

            // Mostrar mensajes recibidos del backend
            if (response.data && response.data.length > 0) {
                response.data.forEach((mensaje) => {
                    if (mensaje.operacion === "ERROR") {
                        mostrarNotificacion('danger', mensaje.mensaje);
                    } else {
                        mostrarNotificacion('success', mensaje.mensaje);
                    }
                });
            }
        } catch (error) {
            console.error('Error al obtener los mensajes:', error);
        }
    };


    const mostrarNotificacion = (tipo, mensaje) => {
        setNotificacion({ tipo, mensaje });
        setTimeout(() => {
            setNotificacion({ tipo: '', mensaje: '' });
        }, 5000); // Ocultar la notificación después de 5 segundos
    };

    return (
        <div style={{ textAlign: 'center' }}>
            <h1>COMANDOS</h1>
            {/* Muestra la notificación si existe */}
            {notificacion.mensaje && (
                <div className={`alert alert-${notificacion.tipo}`} role="alert">
                    {notificacion.mensaje}
                </div>
            )}
            <div className="input-container">
                <p className="bash-text">
                    <span className="user">user</span><span className="vm">@mia-go-file</span>:
                    <span className="char">~</span>$
                </p>
                <textarea
                    className="input textarea"
                    placeholder="Ingresar comando..."
                    value={comando}
                    onChange={(e) => setComando(e.target.value)}
                ></textarea>
            </div>
            <br></br>
            <button
                type="button"
                className="button2"
                style={{ borderRadius: '20px' }}
                onClick={handleSubmit}
            >
                Enviar
            </button>
        </div>
    )
}

export default Pantalla1;
