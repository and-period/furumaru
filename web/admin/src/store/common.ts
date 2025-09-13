export interface AddSnackbarPayload {
  message: string
  color: 'primary' | 'success' | 'info' | 'error'
}

export interface Snackbar extends AddSnackbarPayload {
  isOpen: boolean
  timeout: number | string
}

export const useCommonStore = defineStore('common', {
  state: () => ({
    snackbars: [] as Snackbar[],
  }),

  actions: {
    /**
     * スナックバーを表示する関数
     * @param snackbar スナックバーの情報
     */
    addSnackbar(snackbar: AddSnackbarPayload): void {
      this.snackbars.push({
        isOpen: true,
        ...snackbar,
        timeout: snackbar.color === 'error' ? -1 : 5000,
      })
    },

    /**
     * スナックバーを削除する関数
     * @param index 削除するスナックバーのindex
     */
    hideSnackbar(index: number): void {
      this.snackbars.splice(index, 1)
    },
  },
})
