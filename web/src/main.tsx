import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import { connect } from './ws'
import './styles.css'

connect() // once, outside React so StrictMode double effects don't open two sockets

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
