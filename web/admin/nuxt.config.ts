import { defineNuxtConfig } from 'nuxt/config'
import colors from 'vuetify/es5/util/colors'
import ja from 'vuetify/src/locale/ja'

export default defineNuxtConfig({
  dev: false,
  telemetry: false,
  ssr: false,
  srcDir: 'src',
  target: 'static',
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
  css: ['~/assets/main.scss'],
  plugins: [
    '~/plugins/firebase',
    '~/plugins/google-analytics',
    '~/plugins/auth',
    '~/plugins/api-error-handler',
    '~/plugins/api-client',
  ],
  components: [
    { path: '~/components' },
    { path: '~/components/atoms' },
    { path: '~/components/molecules' },
    { path: '~/components/organisms' },
    { path: '~/components/templates' },
  ],
  modules: [
    // '@nuxtjs/axios',
    '@nuxtjs/google-fonts',
    '@nuxtjs/stylelint-module',
    // '@nuxtjs/vuetify',
    '@pinia/nuxt',
  ],
  axios: {
    baseURL: '/',
  },
  env: {
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
  },
  router: {
    middleware: ['auth', 'notification'],
  },
  vuetify: {
    customVariables: ['~/assets/variables.scss'],
    treeShake: true,
    lang: {
      locales: { ja },
      current: 'ja',
    },
    theme: {
      dark: false,
      themes: {
        dark: {
          primary: colors.green.accent1,
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },
        light: {
          primary: colors.lightGreen.darken2,
          primaryLight: colors.lightGreen.lighten2,
          accent: colors.amber.darken1,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          unknown: colors.grey.darken2,
          success: colors.green.accent3,
        },
      },
      options: { customProperties: true },
    },
  },
  build: {},
  googleFonts: {
    download: true,
    inject: true,
    overwriting: true,
    display: 'swap',
    families: {
      'BIZ+UDGothic': true,
    },
  },
})
