import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Navbar from "./components/Navbar";
import Pantalla1 from "./pages/Pantalla1";
import Pantalla2 from "./pages/Pantalla2";
import Pantalla3 from "./pages/Pantalla3";
import './css/App.css'; 


function App() {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route path='/' element={<Pantalla1/>}/>
        <Route path='/pantalla2' element={<Pantalla2/>}/>
        <Route path='/pantalla3' element={<Pantalla3/>}/>
      </Routes>
    </Router>
  );
}

export default App;