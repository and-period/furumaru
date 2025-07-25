// this import can be removed if you don't need to display the version on the page
import pkg from './package.json';

export default defineNuxtConfig({
  modules: ['@nuxtjs/tailwindcss'],
  plugins: [],
  app: {
    head: {
      title: 'ふるまる - LINE連携アプリ',
      htmlAttrs: {
        lang: 'ja',
      },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { hid: 'description', name: 'description', content: '' },
        { name: 'format-detection', content: 'telephone=no' },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      ],
    },
  },
  runtimeConfig: {
    public: {
      LIFF_ID: process.env.LIFF_ID,
      VERSION: pkg.version || '0.1.0',
    },
  },
  srcDir: 'src',
  tailwindcss: {
    exposeConfig: true,
    viewer: true,
  },
});
