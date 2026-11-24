import * as Sentry from '@sentry/nuxt'

export default defineNuxtPlugin((nuxtApp) => {
  const router = useRouter()
  const runtimeConfig = useRuntimeConfig()

  if (!runtimeConfig.public.SENTRY_DSN) {
    return
  }

  Sentry.init({
    dsn: runtimeConfig.public.SENTRY_DSN,
    environment: runtimeConfig.public.ENVIRONMENT,
    integrations: [
      Sentry.browserProfilingIntegration(),
      Sentry.browserTracingIntegration({ router }),
      Sentry.replayIntegration({ maskAllText: false, blockAllMedia: false }),
    ],
    tracesSampleRate: runtimeConfig.public.SENTRY_TRACES_SAMPLE_RATE,
    // @sentry/browser v10 以降、profilesSampleRate は profileSessionSampleRate + profileLifecycle に置き換え
    profileSessionSampleRate: runtimeConfig.public.SENTRY_PROFILES_SAMPLE_RATE,
    profileLifecycle: 'trace',
    replaysSessionSampleRate: runtimeConfig.public.SENTRY_REPLAYS_SESSION_SAMPLE_RATE,
    replaysOnErrorSampleRate: runtimeConfig.public.SENTRY_REPLAYS_ON_ERROR_SAMPLE_RATE,
  })

  // @sentry/nuxt v10 以降は init の app オプションが廃止され、Vue のエラーハンドラは
  // @sentry/nuxt/module 経由でのみ設定されるため、Nuxt のフックで明示的に送信する
  nuxtApp.hook('vue:error', (error) => {
    Sentry.captureException(error)
  })
  nuxtApp.hook('app:error', (error) => {
    Sentry.captureException(error)
  })
})
