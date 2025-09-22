import { initializeApp } from 'firebase/app'
import type { FirebaseApp, FirebaseOptions } from 'firebase/app'
import { getMessaging, isSupported } from 'firebase/messaging'
import type { Messaging } from 'firebase/messaging'

/* eslint-disable import/no-mutable-exports */
let app: FirebaseApp
let messaging: Messaging | null = null

export default defineNuxtPlugin(async () => {
  const runtimeConfig = useRuntimeConfig()

  const config: FirebaseOptions = {
    apiKey: runtimeConfig.public.FIREBASE_API_KEY,
    authDomain: runtimeConfig.public.FIREBASE_AUTH_DOMAIN,
    projectId: runtimeConfig.public.FIREBASE_PROJECT_ID,
    storageBucket: runtimeConfig.public.FIREBASE_STORAGE_BUCKET,
    messagingSenderId: runtimeConfig.public.FIREBASE_MESSAGING_SENDER_ID,
    appId: runtimeConfig.public.FIREBASE_APP_ID,
    measurementId: runtimeConfig.public.FIREBASE_MEASUREMENT_ID,
  }

  app = initializeApp(config)

  // Firebase Messaging のサポートチェック
  try {
    const supported = await isSupported()
    if (supported) {
      messaging = getMessaging(app)
    }
    else {
      console.log('Firebase Messaging is not supported in this browser')
    }
  }
  catch (error) {
    console.log('Firebase Messaging initialization failed:', error)
  }

  return {
    provide: {
      firebaseApp: app,
      firebaseMessaging: messaging,
    },
  }
})

export { app, messaging }
