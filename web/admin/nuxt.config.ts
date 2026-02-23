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
        { name: 'description', content: '' },
        { name: 'format-detection', content: 'telephone=no' },
      ],
      link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
    },
  },

  css: ['~/assets/main.scss'],

  plugins: [
    '~/plugins/firebase',
    '~/plugins/google-analytics',
    '~/plugins/sentry.client',
    '~/plugins/api-client',
  ],

  modules: [
    '@nuxt/devtools',
    '@nuxt/eslint',
    '@nuxtjs/google-fonts',
    '@nuxtjs/stylelint-module',
    ['@pinia/nuxt', { autoImports: ['defineStore', 'acceptHMRUpdate'] }],
    '@vueuse/nuxt',
    'vuetify-nuxt-module',
  ],

  vuetify: {
    moduleOptions: {
      styles: {
        configFile: 'assets/variables.scss',
      },
    },
    vuetifyOptions: './vuetify.config.ts',
  },

  imports: {
    autoImport: true,
    dirs: ['stores', 'composables', 'utils'],
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

  googleFonts: {
    download: false,
    inject: true,
    overwriting: true,
    display: 'swap',
    preload: true,
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
      COGNITO_AUTH_DOMAIN: process.env.COGNITO_AUTH_DOMAIN || '',
      COGNITO_CLIENT_ID: process.env.COGNITO_CLIENT_ID || '',
      GOOGLE_SIGNIN_REDIRECT_URI: process.env.GOOGLE_SIGNIN_REDIRECT_URI || '',
      GOOGLE_CONNECT_REDIRECT_URI: process.env.GOOGLE_CONNECT_REDIRECT_URI || '',
      LINE_SIGNIN_REDIRECT_URI: process.env.LINE_SIGNIN_REDIRECT_URI || '',
      LINE_CONNECT_REDIRECT_URI: process.env.LINE_CONNECT_REDIRECT_URI || '',
    },
  },

  devtools: {
    timeline: {
      enabled: true,
    },
  },

  vite: {
    build: {
      sourcemap: true,
    },
    plugins: [
      process.env.SENTRY_AUTH_TOKEN
      && sentryVitePlugin({
        org: process.env.SENTRY_ORGANIZATION,
        project: process.env.SENTRY_PROJECT,
        authToken: process.env.SENTRY_AUTH_TOKEN,
      }),
    ],
  },

  compatibilityDate: '2024-10-27',

  future: {
    compatibilityVersion: 4,
  },
})
