import * as Sentry from '@sentry/node'

export default defineNitroPlugin(() => {
  const runtimeConfig = useRuntimeConfig()

  if (!runtimeConfig.public.SENTRY_DSN) {
    return
  }

  Sentry.init({
    dsn: runtimeConfig.public.SENTRY_DSN,
    environment: runtimeConfig.public.ENVIRONMENT,
    tracesSampleRate: runtimeConfig.public.SENTRY_TRACES_SAMPLE_RATE,
    profilesSampleRate: runtimeConfig.public.SENTRY_PROFILES_SAMPLE_RATE,
  })
})
