import React, { useState, useEffect } from 'react';
import { createRoot } from 'react-dom/client';
import { CarouselProvider, Slider, Slide, ButtonBack, ButtonNext } from 'pure-react-carousel';
import 'pure-react-carousel/dist/react-carousel.es.css';
import axios from 'axios';
import './App.css'; // Import the CSS for styling

function App() {
  const [imageUrls, setImageUrls] = useState([]);

  useEffect(() => {
    axios.get('http://localhost:8000/images')
      .then(response => {
        setImageUrls(response.data);
      })
      .catch(error => {
        console.error('Error fetching image data:', error);
      });
  }, []);

  useEffect(() => {
    if(imageUrls.length > 0) {
      console.log(imageUrls);
    }
  }, [imageUrls]);

  return (
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
              <img className="centered-slide" src={imgUrl.link} alt={`Image ${imgUrl.id}`} />
            </Slide>
          ))}
        </Slider>
        <ButtonBack></ButtonBack>
        <ButtonNext></ButtonNext>
      </CarouselProvider>
    </div>
  );
}

// Assuming '#root' is the ID of your root div in the HTML
const container = document.querySelector('#root');
if (container) {
  createRoot(container).render(<App />);
}

export default App;
