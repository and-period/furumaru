import en from './src/locales/en_us.json'
import ja from './src/locales/ja_jp.json'

export default defineI18nConfig((nuxt) => {
  return {
    legacy: false,
    locales: [
      {
        code: 'ja',
        iso: 'ja',
        file: 'ja_jp.json'
      },
      {
        code: 'en',
        iso: 'en',
        file: 'en_us.json'
      }
    ],
    defaultLocale: 'ja',
    messages: {
      ja,
      en
    }
  }
})
