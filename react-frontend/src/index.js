import React, { useState, useEffect } from 'react';
import { createRoot } from 'react-dom/client';
import { CarouselProvider, Slider, Slide, ButtonBack, ButtonNext } from 'pure-react-carousel';
import 'pure-react-carousel/dist/react-carousel.es.css';
import axios from 'axios';
import './App.css';

function App() {
  const [imageUrls, setImageUrls] = useState([]);
  
  // Determina la URL base según el entorno
  const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8000';

  useEffect(() => {
    // Usa la URL base configurada
    axios.get(`${API_URL}/images`, {
      headers: {
        'Content-Type': 'application/json',
        // Agrega headers adicionales si son necesarios
      },
      withCredentials: true // Si estás usando cookies
    })
    .then(response => {
      setImageUrls(response.data);
    })
    .catch(error => {
      console.error('Error fetching image data:', error);
    });
  }, []);

  return (
   hola ql
    <div className="carousel-container">
      <CarouselProvider
        naturalSlideWidth={100}
        naturalSlideHeight={125}
        totalSlides={imageUrls.length}
        isPlaying={true}
        interval={3000}
        infinite={true}
      >
        <Slider>
          {imageUrls.map((imgUrl, index) => (
            <Slide index={index} key={imgUrl.id}>
              <img 
                className="centered-slide" 
                src={imgUrl.link} 
                alt={`Image ${imgUrl.id}`}
                onError={(e) => {
                  console.error(`Error loading image: ${imgUrl.link}`);
                  e.target.src = 'placeholder.jpg'; // Imagen de respaldo
                }}
              />
            </Slide>
          ))}
        </Slider>
        <ButtonBack>Back</ButtonBack>
        <ButtonNext>Next</ButtonNext>
      </CarouselProvider>
    </div>
  );
}

const container = document.querySelector('#root');
if (container) {
  createRoot(container).render(<App />);
}

export default App;
