import { useContext } from '@nuxtjs/composition-api'

export function useI18n() {
  const { i18n } = useContext()

  const t = i18n.t

  return { t }
}
