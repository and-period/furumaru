import { defineStore } from 'pinia'

import type { Admin } from '~/types/api/v1'

export const useAdminStore = defineStore('admin', {
  state: () => ({
    admin: {} as Admin,
    admins: [] as Admin[],
  }),
})
