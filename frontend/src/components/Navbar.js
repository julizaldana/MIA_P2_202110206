import React from "react";
import { Link, useLocation } from 'react-router-dom';
import "../css/Navbar.css";

const Navbar = () => {
  // Obtenemos la ubicación actual usando useLocation de react-router-dom
  const location = useLocation();

  return (
    <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
      <Link to=''>
        <img src='./mialogo.png' width='70' alt='logo'/>
      </Link>
      <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
        <span className="navbar-toggler-icon"></span>
      </button>
      <div className="collapse navbar-collapse" id="navbarNav">
        <ul className="navbar-nav mx-auto">
          <li className={"nav-item" + (location.pathname === '/' ? ' active' : '')}>
            <Link className="nav-link" to='/'> Inicio </Link>
          </li>
          <li className={"nav-item" + (location.pathname === '/pantalla2' ? ' active' : '')}>
            <Link className="nav-link" to='/pantalla2'> Gestión de Discos, Particiones y Archivos </Link>
          </li>
          <li className={"nav-item" + (location.pathname === '/pantalla3' ? ' active' : '')}>
            <Link className="nav-link" to='/pantalla3'> Reportes </Link>
          </li>
        </ul>
      </div>
    </nav>
  );
}

export default Navbar;
