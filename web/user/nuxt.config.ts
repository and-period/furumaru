import { sentryVitePlugin } from '@sentry/vite-plugin'
import en from './src/locales/en_us.json'
import ja from './src/locales/ja_jp.json'

export default defineNuxtConfig({
  ssr: true,
  srcDir: 'src',
  app: {
    head: {
      titleTemplate: '%s - ふるマル',
      htmlAttrs: {
        lang: 'ja',
        prefix: 'og: http://ogp.me/ns#',
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          hid: 'description',
          name: 'description',
          content: '生産者のこだわりが「伝える」以上に「伝わる」ライブマルシェ',
        },
        { name: 'format-detection', content: 'telephone=no' },
        // Google Search Console
        {
          name: 'google-site-verification',
          content: 'xLstKXV5GxV27-afCCeUr5hg8vElOz_Y6sieUFHw8oU',
        },
        { hid: 'og:site_name', property: 'og:site_name', content: 'ふるマル' },
        { hid: 'og:type', property: 'og:type', content: 'website' },
        {
          hid: 'og:url',
          property: 'og:url',
          content: 'https://www.furumaru.and-period.co.jp/',
        },
        { hid: 'og:title', property: 'og:title', content: 'ふるマル' },
        {
          hid: 'og:description',
          property: 'og:description',
          content: '生産者のこだわりが「伝える」以上に「伝わる」ライブマルシェ',
        },
        {
          hid: 'og:image',
          property: 'og:image',
          content: 'https://www.furumaru.and-period.co.jp/ogp/ogp.jpg',
        },
      ],
      link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
    },
  },
  plugins: ['~/plugins/sentry.client'],
  modules: [
    '@nuxt/eslint',
    '@nuxt/image',
    '@nuxtjs/i18n',
    '@nuxtjs/tailwindcss',
    'nuxt-gtag',
    [
      '@pinia/nuxt',
      {
        autoImports: [
          // automatically imports `defineStore`
          'defineStore',
        ],
      },
    ],
    '@pinia-plugin-persistedstate/nuxt',
  ],
  gtag: {
    id: process.env.NUXT_PUBLIC_GTAG_ID,
  },
  eslint: {
    config: {
      stylistic: {
        indent: 2,
        quotes: 'single',
        semi: false,
      },
    },
  },
  image: {
    providers: {
      cloudFront: {
        name: 'cloudFront',
        provider: '~/provider/cloud-front.ts',
      },
    },
  },
  i18n: {
    defaultLocale: 'ja',
    detectBrowserLanguage: {
      useCookie: false,
      alwaysRedirect: true,
    },
    locales: [
      {
        code: 'ja',
        iso: 'ja',
        file: 'ja_jp.json',
      },
      {
        code: 'en',
        iso: 'en',
        file: 'en_us.json',
      },
    ],
    vueI18n: {
      fallbackLocale: 'ja',
      messages: {
        ja,
        en,
      },
    },
  },
  components: [
    {
      path: '~/components/',
      pathPrefix: false,
    },
  ],
  runtimeConfig: {
    MICRO_CMS_DOMAIN: process.env.MICRO_CMS_DOMAIN,
    MICRO_CMS_API_KEY: process.env.MICRO_CMS_API_KEY,
    public: {
      API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:18000',
      ENVIRONMENT: process.env.ENVIRONMENT || '',
      SENTRY_DSN: process.env.SENTRY_DSN || '',
      SENTRY_TRACES_SAMPLE_RATE: parseFloat(
        process.env.SENTRY_TRACES_SAMPLE_RATE || '0.5',
      ),
      SENTRY_PROFILES_SAMPLE_RATE: parseFloat(
        process.env.SENTRY_PROFILES_SAMPLE_RATE || '0.5',
      ),
      SENTRY_REPLAYS_SESSION_SAMPLE_RATE: parseFloat(
        process.env.SENTRY_REPLAYS_SESSION_SAMPLE_RATE || '0.2',
      ),
      SENTRY_REPLAYS_ON_ERROR_SAMPLE_RATE: parseFloat(
        process.env.SENTRY_REPLAYS_ON_ERROR_SAMPLE_RATE || '1.0',
      ),
    },
  },
  build: {},
  nitro: {
    plugins: ['~/server/plugins/sentry'],
  },
  vite: {
    vue: {
      script: {
        defineModel: true,
        propsDestructure: true,
      },
    },
    build: {
      sourcemap: true,
    },
    plugins: [
      sentryVitePlugin({
        org: process.env.SENTRY_ORGANIZATION,
        project: process.env.SENTRY_PROJECT,
        authToken: process.env.SENTRY_AUTH_TOKEN,
      }),
    ],
  },
})
