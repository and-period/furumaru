import { defineStore } from 'pinia'

/**
 * 通知系を管理するグローバルステート
 */
export const useNotificationStore = defineStore('notification', {
  state: () => {
    return {
      notifications: [],
    }
  },
})
