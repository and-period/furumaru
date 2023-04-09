importScripts("https://www.gstatic.com/firebasejs/9.0.0/firebase-app-compat.js");
importScripts("https://www.gstatic.com/firebasejs/9.0.0/firebase-messaging-compat.js");

const config = {
  apiKey: "AIzaSyAVecTlXgke2_yM00uXRAi_Y_IlF_V6xY4",
  authDomain: "furumaru-stg-admin.firebaseapp.com",
  projectId: "furumaru-stg-admin",
  storageBucket: "furumaru-stg-admin.appspot.com",
  messagingSenderId: "933783288161",
  appId: "1:933783288161:web:8c9ff8cffa7c768f99ad20",
  measurementId: "G-3560RV4HZ1",
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
