import { defineStore } from 'pinia'

interface AddSnackbarPayload {
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
        timeout: -1,
      })
    },

    /**
     * スナックバーを削除する関数
     * @param index 削除するスナックバーのindex
     */
    hideSnackbar(index: number): void {
      this.snackbars = this.snackbars.filter((_, i) => i !== index)
    },
  },
})
