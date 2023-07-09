import { Admin } from '~/types/api'

export const useAdminStore = defineStore('admin', {
  state: () => ({
    admin: {} as Admin,
    admins: [] as Admin[]
  })
})
