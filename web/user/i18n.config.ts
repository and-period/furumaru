import en from './src/locales/en_us.json'
import ja from './src/locales/ja_jp.json'

export default defineI18nConfig(() => ({
  legacy: false,
  locale: 'ja',
  messages: {
    ja,
    en,
  },
}))
