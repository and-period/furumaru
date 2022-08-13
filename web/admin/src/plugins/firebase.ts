import { initializeApp } from 'firebase/app'
import { getMessaging, getToken, onMessage } from 'firebase/messaging'

const config = {
  apiKey: process.env.FIREBASE_API_KEY,
  authDomain: process.env.FIREBASE_AUTH_DOMAIN,
  projectId: process.env.FIREBASE_PROJECT_ID,
  storageBucket: process.env.FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.FIREBASE_APP_ID,
  measurementId: process.env.FIREBASE_MEASUREMENT_ID,
}

const app = initializeApp(config)
const messaging = getMessaging(app)

function requestPermission() {
  Notification.requestPermission().then((permission) => {
    if (permission === 'granted') {
      console.log('Notification permission granted.')
    }
  })
}

getToken(messaging, {
  vapidKey: process.env.FIREBASE_VAPID_KEY,
})
  .then((currentToken) => {
    if (currentToken) {
      console.log('token', currentToken)
    } else {
      requestPermission()
    }
  })
  .catch((err) => {
    console.log(err)
  })

onMessage(messaging, (payload) => {
  console.log('メッセージ', payload)
})

export default { app, messaging }
