import { FirebaseOptions, initializeApp } from 'firebase/app'
import { getMessaging } from 'firebase/messaging'

const runtimeConfig = useRuntimeConfig()

const config: FirebaseOptions = {
  apiKey: runtimeConfig.public.firebaseApiKey,
  authDomain: runtimeConfig.public.firebaseAuthDomain,
  projectId: runtimeConfig.public.firebaseProjectId,
  storageBucket: runtimeConfig.public.firebaseStorageBucket,
  messagingSenderId: runtimeConfig.public.firebaseMessagingSenderId,
  appId: runtimeConfig.public.firebaseAppId,
  measurementId: runtimeConfig.public.firebaseMeasurementId,
}

const app = initializeApp(config)
const messaging = getMessaging(app)

export default { app, messaging }
