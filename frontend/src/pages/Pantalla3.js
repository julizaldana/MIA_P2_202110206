import React, { useState, useEffect } from 'react';

const Pantalla3 = () => {
    const [reportes, setReportes] = useState([]);
    const [imagenSeleccionada, setImagenSeleccionada] = useState('');
    const [modalAbierto, setModalAbierto] = useState(false);

    useEffect(() => {
        fetch('http://localhost:8080/obtenerreportes')
            .then(response => response.json())
            .then(data => setReportes(data))
            .catch(error => console.error('Error al obtener los reportes:', error));
    }, []);

    const abrirModal = (url) => {
        setImagenSeleccionada(url);
        setModalAbierto(true);
    };

    const cerrarModal = () => {
        setModalAbierto(false);
    };

    return (
        <div>
            <h1>REPORTES</h1>
            <div style={{ display: 'flex', flexWrap: 'wrap' }}>
                {reportes.map(reporte => (
                    <div key={reporte.nombre} style={{ marginRight: '20px', marginBottom: '20px' }}>
                        <h2>{reporte.nombre}</h2>
                        <img
                            src={`data:image/jpeg;base64,${reporte.contenido}`}
                            alt={reporte.nombre}
                            style={{ width: '200px', height: 'auto', border: '1px solid #ccc', borderRadius: '5px', cursor: 'pointer' }}
                            onClick={() => abrirModal(`data:image/jpeg;base64,${reporte.contenido}`)}
                        />
                    </div>
                ))}
            </div>
            {modalAbierto && (
                <div style={{ position: 'fixed', top: 0, left: 0, width: '100%', height: '100%', backgroundColor: 'rgba(0, 0, 0, 0.7)', zIndex: 999 }}>
                    <div style={{ position: 'absolute', top: '50%', left: '50%', transform: 'translate(-50%, -50%)' }}>
                        <button type="button" class="btn btn-danger" onClick={cerrarModal}>Cerrar</button>
                        <img src={imagenSeleccionada} alt="Imagen seleccionada" style={{ maxWidth: '90%', maxHeight: '90%' }} />
                    </div>
                </div>
            )}
        </div>
    );
};

export default Pantalla3;
