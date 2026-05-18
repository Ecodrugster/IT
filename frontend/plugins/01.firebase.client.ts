import { initializeApp } from 'firebase/app'
import { getAuth } from 'firebase/auth'
import { getFirestore } from 'firebase/firestore'
import type { FirebaseApp } from 'firebase/app'
import type { Auth } from 'firebase/auth'
import type { Firestore } from 'firebase/firestore'

type FirebaseProvides = {
  firebaseApp: FirebaseApp | null
  auth: Auth | null
  firestore: Firestore | null
}

export default defineNuxtPlugin<FirebaseProvides>(() => {
  const config = useRuntimeConfig()

  const firebaseConfig = {
    apiKey: config.public.firebaseApiKey,
    authDomain: config.public.firebaseAuthDomain,
    projectId: config.public.firebaseProjectId,
    storageBucket: config.public.firebaseStorageBucket,
    messagingSenderId: config.public.firebaseMessagingSenderId,
    appId: config.public.firebaseAppId,
  }

  const requiredFields = ['apiKey', 'authDomain', 'projectId', 'appId'] as const
  const missingFields = requiredFields.filter((key) => !firebaseConfig[key])

  if (missingFields.length > 0) {
    console.warn(
      `[firebase] Missing config fields: ${missingFields.join(', ')}. Auth and Firestore are disabled until env variables are set.`
    )

    return {
      provide: {
        firebaseApp: null,
        auth: null,
        firestore: null
      }
    }
  }

  const app = initializeApp(firebaseConfig)
  const auth = getAuth(app)
  const db = getFirestore(app)

  return {
    provide: {
      firebaseApp: app,
      auth,
      firestore: db
    }
  }
})
