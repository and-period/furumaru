require('dotenv').config()
const fs = require('fs')

/**
 * firebase-messaging-sw.js - 生成用
 */
fs.writeFileSync(
  './src/public/firebase-messaging-sw.js',
  `importScripts("https://www.gstatic.com/firebasejs/9.0.0/firebase-app-compat.js");
importScripts("https://www.gstatic.com/firebasejs/9.0.0/firebase-messaging-compat.js");

const config = {
  apiKey: "${process.env.NUXT_FIREBASE_API_KEY}",
  authDomain: "${process.env.NUXT_FIREBASE_AUTH_DOMAIN}",
  projectId: "${process.env.NUXT_FIREBASE_PROJECT_ID}",
  storageBucket: "${process.env.NUXT_FIREBASE_STORAGE_BUCKET}",
  messagingSenderId: "${process.env.NUXT_FIREBASE_MESSAGING_SENDER_ID}",
  appId: "${process.env.NUXT_FIREBASE_APP_ID}",
  measurementId: "${process.env.NUXT_FIREBASE_MEASUREMENT_ID}",
}

firebase.initializeApp(config);

const messaging = firebase.messaging();

messaging.onBackgroundMessage((payload) => {
  console.log('[firebase-messaging-sw.js] Received background message ', payload);
  const notificationTitle = payload.notification.title;
  const notificationOptions = {
    body: payload.notification.body,
  };

  self.registration.showNotification(notificationTitle,
    notificationOptions);
});
`,
)
