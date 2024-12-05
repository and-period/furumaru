import { sentryVitePlugin } from '@sentry/vite-plugin'

/* eslint-disable nuxt/nuxt-config-keys-order */
export default defineNuxtConfig({
  ssr: false,
  srcDir: 'src',

  app: {
    head: {
      titleTemplate: 'ふるマル - 管理者ツール',
      htmlAttrs: {
        lang: 'ja',
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { hid: 'description', name: 'description', content: '' },
        { name: 'format-detection', content: 'telephone=no' },
      ],
      link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
    },
  },

  css: ['~/assets/main.scss', '~/assets/variables.scss'],

  plugins: [
    '~/plugins/axios',
    '~/plugins/chartjs.client',
    '~/plugins/firebase',
    '~/plugins/google-analytics',
    '~/plugins/sentry.client',
    '~/plugins/vuetify',
    '~/plugins/api-client',
  ],

  modules: [
    '@nuxt/devtools',
    '@nuxt/eslint',
    '@nuxtjs/google-fonts',
    '@nuxtjs/stylelint-module',
    '@pinia/nuxt',
    '@vueuse/nuxt',
  ],

  eslint: {
    config: {
      stylistic: {
        indent: 2,
        quotes: 'single',
        semi: false,
      },
    },
  },

  googleFonts: {
    download: true,
    inject: true,
    overwriting: true,
    display: 'swap',
    families: {
      'BIZ+UDGothic': true,
    },
  },

  runtimeConfig: {
    public: {
      API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:18010',
      FIREBASE_API_KEY: process.env.FIREBASE_API_KEY || '',
      FIREBASE_AUTH_DOMAIN: process.env.FIREBASE_AUTH_DOMAIN || '',
      FIREBASE_PROJECT_ID: process.env.FIREBASE_PROJECT_ID || '',
      FIREBASE_STORAGE_BUCKET: process.env.FIREBASE_STORAGE_BUCKET || '',
      FIREBASE_MESSAGING_SENDER_ID:
        process.env.FIREBASE_MESSAGING_SENDER_ID || '',
      FIREBASE_APP_ID: process.env.FIREBASE_APP_ID || '',
      FIREBASE_MEASUREMENT_ID: process.env.FIREBASE_MEASUREMENT_ID || '',
      FIREBASE_VAPID_KEY: process.env.FIREBASE_VAPID_KEY || '',
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

  devtools: {
    timeline: {
      enabled: true,
    },
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

  compatibilityDate: '2024-10-27',
})