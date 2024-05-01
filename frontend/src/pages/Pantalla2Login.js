import React, { useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import axios from 'axios';
import "../css/Pantalla2Login.css";

const Pantalla2Login = () => {
    // Obtener el parámetro de la URL (nombre del disco)
    const { idParticion } = useParams();
    const [usuario, setUsuario] = useState('');
    const [contraseña, setContraseña] = useState('');
    const [idParticionInput, setIdParticionInput] = useState('');

    // Función para manejar el envío del formulario
    const handleSubmit = async () => {
        const comando = `login -user=${usuario} -pass=${contraseña} -id=${idParticionInput}`;
        const data = {
            comandos: [comando]
        };

        try {
            const response = await axios.post('http://localhost:8080/iniciarsesion', data);
            console.log(response.data);
            console.log("Se mandó al backend el comando de inicio de sesión");
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
            <h1> LOGIN </h1>
            <h2>Id Particion Montada: {idParticion}</h2>
            <button
                type="button"
                className="button"
                style={{ borderRadius: '20px' }}
            >
                <Link to="/pantalla2">Regresar</Link>
            </button>
            <div className="form-box">
                <form className="form">
                    <span className="title" style={{ fontSize: '1.5em' }}>Inicia Sesión!</span>
                    <span className="subtitle" style={{ fontSize: '1.5em' }}>Inicia sesión con tu cuenta de </span>
                    <span className="subtitle" style={{ fontSize: '1.5em' }}> MIA Go File</span>
                    <div className="form-container">
                        <input type="text" className="input" placeholder="usuario" value={usuario} onChange={(e) => setUsuario(e.target.value)} />
                        <input type="password" className="input" placeholder="contraseña" value={contraseña} onChange={(e) => setContraseña(e.target.value)} />
                        <input type="text" className="input" placeholder="id particion" value={idParticionInput} onChange={(e) => setIdParticionInput(e.target.value)} />
                    </div>
                    <Link to={`/logueado/${idParticionInput}/${usuario}`}><button type="button" style={{ fontSize: '1.2em' }} onClick={() => { handleSubmit(); mandaridparticion(idParticionInput); }}>LOGIN</button></Link>
                </form>
                <div className="form-section" style={{ fontSize: '1em' }}>
                    <p>Deseas crear otro usuario? <Link to="/">Comandos</Link> </p>
                </div>
            </div>
        </div>
    );
}

export default Pantalla2Login;
