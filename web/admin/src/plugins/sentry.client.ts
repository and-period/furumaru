import * as Sentry from '@sentry/vue'
import { browserProfilingIntegration, browserTracingIntegration, replayIntegration } from '@sentry/vue'

export default defineNuxtPlugin((nuxtApp) => {
  const router = useRouter()
  const runtimeConfig = useRuntimeConfig()

  if (!runtimeConfig.public.SENTRY_DSN) {
    return
  }

  Sentry.init({
    app: nuxtApp.vueApp,
    dsn: runtimeConfig.public.SENTRY_DSN,
    environment: runtimeConfig.public.ENVIRONMENT,
    integrations: [
      browserProfilingIntegration(),
      browserTracingIntegration({ router }),
      replayIntegration({ maskAllText: false, blockAllMedia: false })
    ],
    tracesSampleRate: 1.0,
    profilesSampleRate: 1.0,
    replaysSessionSampleRate: 1.0,
    replaysOnErrorSampleRate: 1.0
  })
})
