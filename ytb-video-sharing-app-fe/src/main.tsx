import { createRoot } from 'react-dom/client';
import App from './App.tsx';
import './index.css';
console.log(`${import.meta.env.VITE_WS_URL}`);
createRoot(document.getElementById('root')!).render(<App />);
