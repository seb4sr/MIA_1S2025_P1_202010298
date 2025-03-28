import React, { useState } from 'react';
import "../Stylesheets/navegador.css"

export default function Comandos({newIp="localhost"}){
    const [textValue, setTextValue] = useState('');
    const [textExit, setTextExit] = useState('');

    const handleTextChange = (event) => {
        setTextValue(event.target.value);
    };


    //Limpiar las consolas 
    const handlerLimpiar = () => {
        setTextValue(""); //COnsola entreada
        setTextExit("");  //COnsola salida
    }

    const sendData = async (e) => {
        e.preventDefault();
        const data = {
            text: textValue
        };
        
        try {
            const response = await fetch(`http://localhost:8080/analizar`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
    
            if (!response.ok) {
                throw new Error('Error al enviar datos');
            }
    
            const responseData = await response.json();
            console.log('Respuesta del servidor:', responseData);
            console.log('Respuesta del metodo ',responseData.message)
            setTextExit(responseData.message)
           
        } catch (error) {
            console.error('Error:', error);
        }

    }

    const handleLoadClick = () => {
        const input = document.createElement("input");
        input.type = "file";
        input.addEventListener("change",handleFileChange);
        input.click();
    }

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        const reader = new FileReader();
        reader.onload = (e) => {
            const text = e.target.result;
            setTextValue(text);
        }
        reader.readAsText(file);
    }

    return(
        <div className='contenedorEjecutar'>
            <div id="espacio">&nbsp;&nbsp;&nbsp;</div>
            <table>
                <tbody>

                    <tr><td><p><strong>ENTRADA</strong></p></td></tr>

                    <tr>
                        <td>
                            <textarea
                                className='entrada'
                                value={textValue}
                                onChange={handleTextChange}
                                placeholder='Ingrese comandos'
                                id='inputComands'
                            />
                        </td>
                    </tr>

                    <tr><td><strong><p>SALIDA</p></strong></td></tr>

                    <tr>
                        <td>
                            <textarea
                                className='entrada'
                                value={textExit}
                                id='inputComands'
                                readOnly
                            />
                        </td>
                    </tr>

                    <tr>
                        <td style={{textAlign:'center'}}>
                            <div className="botones">
                                <button type="button" className="btn btn-primary" onClick={(e) => sendData(e)}>Ejecutar</button>
                                <button type="button" className="btn btn-primary" onClick={(e) => handlerLimpiar(e)}>Limpiar consolas</button>
                                <button type="button" className="btn btn-primary" onClick={(e) => handleLoadClick(e)}>Subir archivo</button>
                            </div> 
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    );
}