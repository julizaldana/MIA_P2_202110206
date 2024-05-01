import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Navbar from "./components/Navbar";
import Pantalla1 from "./pages/Pantalla1";
import Pantalla2 from "./pages/Pantalla2";
import Pantalla3 from "./pages/Pantalla3";
import Pantalla2p from "./pages/Pantalla2p";
import Pantalla2Login from "./pages/Pantalla2Login";
import Pantalla2Fs from "./pages/Pantalla2Fs";
import Footer from "./components/footer"
import './css/App.css'; 


function App() {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route path='/' element={<Pantalla1/>}/>
        <Route path='/pantalla2' element={<Pantalla2/>}/>
        <Route path='/pantalla3' element={<Pantalla3/>}/>
        <Route path="/particiones/:nombreDisco" element={<Pantalla2p/>}/>
        <Route path="/login/:idParticion" element={<Pantalla2Login/>}/>
        <Route path="/logueado/:idParticion/:user" element={<Pantalla2Fs/>}/>
      </Routes>
      <Footer/>
    </Router>
  );
}

export default App;