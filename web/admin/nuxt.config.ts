import { defineNuxtConfig } from 'nuxt/config'

export default defineNuxtConfig({
  dev: false,
  telemetry: false,
  ssr: false,
  srcDir: 'src',
  target: 'static',
  head: {
    titleTemplate: 'ふるマル - 管理者ツール',
    htmlAttrs: {
      lang: 'ja'
    },
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: '' },
      { name: 'format-detection', content: 'telephone=no' }
    ],
    link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }]
  },
  css: ['~/assets/main.scss', '~/assets/variables.scss'],
  plugins: [
    '~/plugins/firebase',
    '~/plugins/google-analytics',
    '~/plugins/auth',
    '~/plugins/api-error-handler',
    '~/plugins/api-client',
    '~/plugins/vuetify'
  ],
  components: [
    { path: '~/components', pathPrefix: false },
    { path: '~/components/', pathPrefix: false }
  ],
  modules: ['@nuxtjs/google-fonts', '@nuxtjs/stylelint-module', '@pinia/nuxt'],
  axios: {
    baseURL: '/'
  },
  router: {
    middleware: ['auth', 'notification']
  },
  build: {},
  googleFonts: {
    download: true,
    inject: true,
    overwriting: true,
    display: 'swap',
    families: {
      'BIZ+UDGothic': true
    }
  },
  runtimeConfig: {
    public: {
      apiBaseUrl: 'http://localhost:18010',
      firebaseApiKey: '',
      firebaseAuthDomain: '',
      firebaseProjectId: '',
      firebaseStorageBucket: '',
      firebaseMessagingSenderId: '',
      firebaseAppId: '',
      firebaseMeasurementId: '',
      firebaseVapidKey: '',
    }
  }
})
