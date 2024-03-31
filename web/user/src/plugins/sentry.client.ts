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
    tracesSampleRate: Number(runtimeConfig.public.SENTRY_TRACES_SAMPLE_RATE),
    profilesSampleRate: Number(runtimeConfig.public.SENTRY_PROFILES_SAMPLE_RATE),
    replaysSessionSampleRate: Number(runtimeConfig.public.SENTRY_REPLAYS_SESSION_SAMPLE_RATE),
    replaysOnErrorSampleRate: Number(runtimeConfig.public.SENTRY_REPLAYS_ON_ERROR_SAMPLE_RATE)
  })
})
