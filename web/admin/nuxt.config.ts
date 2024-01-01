export default defineNuxtConfig({
  ssr: false,
  srcDir: 'src',
  app: {
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
    }
  },
  css: ['~/assets/main.scss', '~/assets/variables.scss'],
  plugins: [
    '~/plugins/axios',
    '~/plugins/firebase',
    '~/plugins/google-analytics',
    '~/plugins/vuetify',
    '~/plugins/api-client'
  ],
  modules: [
    '@nuxt/devtools',
    '@nuxtjs/google-fonts',
    '@nuxtjs/stylelint-module',
    ['@pinia/nuxt', { autoImports: ['defineStore'] }]
  ],
  devtools: {
    enabled: true,
    vscode: {}
  },
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
      ENVIRONMENT: process.env.ENVIRONMENT || ''
    }
  }
})
