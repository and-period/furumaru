import { computed } from 'vue'

export function useIsMobile() {
  const { $vuetify } = useNuxtApp()

  const breakpointList = ['sm', 'xs']

  // FIXME: $vuetifyの取得がうまくできていない
  const isMobile = computed(() => breakpointList.includes($vuetify.breakpoint?.name))

  return { isMobile }
}
