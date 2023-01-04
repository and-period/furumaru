import en from './src/locales/en_us.json'
import ja from './src/locales/ja_jp.json'

export default defineNuxtConfig({
  ssr: false,
  srcDir: 'src',
  app: {
    head: {
      titleTemplate: '%s - user',
      title: 'user',
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
  css: ['vuetify/lib/styles/main.sass', '@fortawesome/fontawesome-free/css/all.css', '~/assets/css/variables.scss'],
  plugins: ['~/plugins/vuetify'],
  modules: ['@nuxtjs/i18n'],
  i18n: {
    locales: ['ja', 'en'],
    defaultLocale: 'ja',
    vueI18n: {
      fallbackLocale: 'ja',
      messages: {
        ja,
        en,
      },
    },
  },
  runtimeConfig: {
    API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:18000',
  },
  build: {
    transpile: ['vuetify'],
  },
})
