import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import './styles/index.css';

const setRem = () => {
  const w = Math.min(document.documentElement.clientWidth, 750);
  document.documentElement.style.fontSize = `${(w / 375) * 16}px`;
};
setRem();
window.addEventListener('resize', setRem);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
);
