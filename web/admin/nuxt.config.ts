import { NuxtConfig } from '@nuxt/types'
import colors from 'vuetify/es5/util/colors'
import ja from 'vuetify/src/locale/ja'

const config: NuxtConfig = {
  srcDir: 'src',
  ssr: false,
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
    '~/plugins/api-client',
    '~/plugins/api-error-handler',
    '~/plugins/auth',
    '~/plugins/firebase',
    '~/plugins/google-analytics',
  ],
  components: [
    { path: '~/components' },
    { path: '~/components/atoms' },
    { path: '~/components/molecules' },
    { path: '~/components/organisms' },
    { path: '~/components/templates' },
  ],
  buildModules: [
    '@nuxt/typescript-build',
    '@nuxtjs/composition-api/module',
    '@nuxtjs/google-fonts',
    '@nuxtjs/vuetify',
    '@pinia/nuxt',
  ],
  modules: [
    '@nuxtjs/axios',
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
    middleware: 'auth',
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
          accentDarken: colors.orange.darken4,
        },
      },
      options: { customProperties: true },
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
  build: {},
}

export default config
