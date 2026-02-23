import type { Admin } from '~/types/api/v1'

export const useAdminStore = defineStore('admin', () => {
  const admin = ref<Admin>({} as Admin)
  const admins = ref<Admin[]>([])

  return {
    admin,
    admins,
  }
})
