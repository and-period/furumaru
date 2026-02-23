export interface AddSnackbarPayload {
  message: string
  color: 'primary' | 'success' | 'info' | 'error'
}

export interface Snackbar extends AddSnackbarPayload {
  isOpen: boolean
  timeout: number | string
}

export const useCommonStore = defineStore('common', () => {
  const snackbars = ref<Snackbar[]>([])

  /**
   * スナックバーを表示する関数
   * @param snackbar スナックバーの情報
   */
  function addSnackbar(snackbar: AddSnackbarPayload): void {
    snackbars.value.push({
      isOpen: true,
      ...snackbar,
      timeout: snackbar.color === 'error' ? -1 : 5000,
    })
  }

  /**
   * スナックバーを削除する関数
   * @param index 削除するスナックバーのindex
   */
  function hideSnackbar(index: number): void {
    snackbars.value.splice(index, 1)
  }

  function $reset(): void {
    snackbars.value = []
  }

  return { snackbars, addSnackbar, hideSnackbar, $reset }
})
