import { useAuthStore } from '~/store'

export function setupAuthStore(isAuthenticated: boolean) {
  const authStore = useAuthStore()

  if (isAuthenticated) {
    authStore.$patch((state) => {
      state.user = {
        adminId: 'kSByoE6FetnPs5Byk3a9Zx',
        role: 1,
        accessToken: 'xxxxxxxxxx',
        refreshToken: 'xxxxxxxxxx',
        expiresIn: 3600,
        tokenType: 'Bearer',
      }
    })
  }
}
