import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import Encoder from './VideoClient/Encoder.tsx';
import Manifest from './VideoClient/Manifest.tsx';
const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <Encoder />
    <Manifest />
  </React.StrictMode>
);

