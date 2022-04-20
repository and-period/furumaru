import { computed, useContext } from '@nuxtjs/composition-api'

export function useIsMobile() {
  const { $vuetify } = useContext()

  const breakpointList = ['sm', 'xs']

  const isMobile = computed(() =>
    breakpointList.includes($vuetify.breakpoint.name)
  )

  return { isMobile }
}
