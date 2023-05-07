export default defineNuxtConfig({
  ssr: false,
  srcDir: 'src',
  app: {
    head: {
      titleTemplate: '%s - user',
      title: 'user',
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
  plugins: [],
  modules: [
    '@nuxtjs/i18n',
    '@nuxtjs/tailwindcss',
    ['@pinia/nuxt',
      {
        autoImports: [
          // automatically imports `defineStore`
          'defineStore'
        ]
      }]
  ],
  i18n: {
    vueI18n: './i18n.config.ts' // if you are using custom path, default
  },
  components: [
    {
      path: '~/components/',
      pathPrefix: false
    }
  ],
  runtimeConfig: {
    public: {
      API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:18000'
    }
  },
  build: {}
})
