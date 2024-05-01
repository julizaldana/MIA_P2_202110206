import React, { useState} from 'react';
import axios from 'axios'; // Se importa Axios para realizar la solicitud HTTP
import "../css/Pantalla1.css";

const Pantalla1 = () => {
    const [comando, setComando] = useState('');
    const [notificacion, setNotificacion] = useState({ tipo: '', mensaje: '' });
    const [notificacionesRecibidas, setNotificacionesRecibidas] = useState('');

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
                let mensajes = '';
                response.data.forEach((mensaje) => {
                    mensajes += `${mensaje.operacion}: ${mensaje.mensaje}\n`;
                    if (mensaje.operacion === "ERROR") {
                        mostrarNotificacion('danger', mensaje.mensaje);
                    } else {
                        mostrarNotificacion('success', mensaje.mensaje);
                    }
                });
                setNotificacionesRecibidas(mensajes);
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
            <br></br>
            <div className="input-container">
                <p className="bash-text" style={{ fontSize: '1.5em' }}>
                    <span className="user">u s e r</span><span className="vm">@mia-go-files</span>:
                    <span className="char">~</span>$
                </p>
                <textarea 
                    className="input textarea"
                    placeholder="Ingresar comando..."
                    value={comando}
                    onChange={(e) => setComando(e.target.value)}
                    style={{ height: '200px', resize: 'vertical', fontSize: '1.5em' }} // Ajuste del tamaño y permitir redimensionamiento vertical
                ></textarea>
            </div>
            <br></br>
            <button
                type="button"
                className="button2"
                style={{ borderRadius: '20px', fontSize: '1.5em' }}
                onClick={handleSubmit}
            >
                Enviar
            </button>
            <br></br>
            <br></br>
            <div className="input-container2">
                <p className="bash-text" style={{ fontSize: '1.5em' }}>
                    <span className="user">console</span><span className="vm">@mia-go-files</span>:
                    <span className="char">~</span>$
                </p>
                <textarea
                    className="input textarea"
                    value={notificacionesRecibidas}
                    onChange={(e) => setComando(e.target.value)}
                    style={{ height: '200px', resize: 'vertical', fontSize: '1em' }} // Ajuste del tamaño y permitir redimensionamiento vertical
                ></textarea>
            </div>
        </div>
    )
}

export default Pantalla1;
