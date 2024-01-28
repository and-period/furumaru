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
        prefix: 'og: http://ogp.me/ns#'
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
        { hid: 'og:url', property: 'og:url', content: 'https://www.furumaru.and-period.co.jp/' },
        { hid: 'og:title', property: 'og:title', content: 'ふるマル' },
        { hid: 'og:description', property: 'og:description', content: '生産者のこだわりが「伝える」以上に「伝わる」ライブマルシェ' },
        { hid: 'og:image', property: 'og:image', content: 'https://www.furumaru.and-period.co.jp/ogp/ogp.jpg' },
      ],
      link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
    },
  },
  plugins: [],
  modules: [
    '@nuxtjs/i18n',
    '@nuxtjs/tailwindcss',
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
  i18n: {
    defaultLocale: 'ja',
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
    public: {
      API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:18000',
      ENVIRONMENT: process.env.ENVIRONMENT || '',
    },
  },
  build: {},
})
