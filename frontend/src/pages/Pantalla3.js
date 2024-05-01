import React, { useState, useEffect } from 'react';

const Pantalla3 = () => {
    const [reportes, setReportes] = useState([]);
    const [imagenSeleccionada, setImagenSeleccionada] = useState('');
    const [modalAbierto, setModalAbierto] = useState(false);
    const [cargando, setCargando] = useState(true);

    useEffect(() => {
        fetch('http://localhost:8080/obtenerreportes')
            .then(response => response.json())
            .then(data => {
                setReportes(data || []);
                setCargando(false); // Cuando se recibe la respuesta, indicar que ya no se está cargando
            })
            .catch(error => {
                console.error('Error al obtener los reportes:', error);
                setCargando(false); // En caso de error, también se indica que ya no se está cargando
            });
    }, []);

    const abrirModal = (url) => {
        setImagenSeleccionada(url);
        setModalAbierto(true);
    };

    const cerrarModal = () => {
        setModalAbierto(false);
    };

    // Mostrar un mensaje si no hay reportes o la lista está vacía
    if (cargando) {
        return <div>Cargando...</div>;
    }

    if (!reportes || reportes.length === 0) {
        return <div><h1>No hay reportes para mostrar</h1></div>;
    }

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
                        <button type="button" className="btn btn-danger" onClick={cerrarModal}>Cerrar</button>
                        <img src={imagenSeleccionada} alt="Imagen seleccionada" style={{ maxWidth: '90%', maxHeight: '90%' }} />
                    </div>
                </div>
            )}
        </div>
    );
};

export default Pantalla3;
