export type AlertType = 'success' | 'info' | 'warning' | 'error'

export function useAlert(alertType: AlertType) {
  const isShow = ref<boolean>(false)
  const alertText = ref<string>('')

  const show = (text?: string) => {
    if (text) {
      alertText.value = text
    }
    isShow.value = true
  }

  const hide = () => {
    isShow.value = false
  }

  return {
    alertType,
    isShow,
    alertText,
    show,
    hide,
  }
}
