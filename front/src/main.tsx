import { initializeApp } from 'firebase/app'
import { getMessaging, onMessage } from 'firebase/messaging'
import React from 'react'
import ReactDOM from 'react-dom/client'
import { registerSW } from 'virtual:pwa-register'
import App from './App.tsx'
import './index.css'

if ('serviceWorker' in navigator) {
  registerSW()
}

const config = {
  apiKey: 'AIzaSyAIBsve46yZ7Nskxfp2deN16VVV86K7XG8',
  authDomain: 'kiroku-messaging-app.firebaseapp.com',
  projectId: 'kiroku-messaging-app',
  storageBucket: 'kiroku-messaging-app.appspot.com',
  messagingSenderId: '1085771729373',
  appId: '1:1085771729373:web:ebab984f412cd38ac6ce34',
}

// const app = initializeApp(config)
initializeApp(config)
const messaging = getMessaging()
onMessage(messaging, (payload) => {
  console.log(payload)
})

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
