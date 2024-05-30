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
      replayIntegration({ maskAllText: false, blockAllMedia: false }),
    ],
    tracesSampleRate: runtimeConfig.public.SENTRY_TRACES_SAMPLE_RATE,
    profilesSampleRate: runtimeConfig.public.SENTRY_PROFILES_SAMPLE_RATE,
    replaysSessionSampleRate: runtimeConfig.public.SENTRY_REPLAYS_SESSION_SAMPLE_RATE,
    replaysOnErrorSampleRate: runtimeConfig.public.SENTRY_REPLAYS_ON_ERROR_SAMPLE_RATE,
  })
})
