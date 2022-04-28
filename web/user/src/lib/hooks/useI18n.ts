import { useContext } from '@nuxtjs/composition-api'

export function useI18n() {
  const { i18n } = useContext()

  return { i18n }
}
