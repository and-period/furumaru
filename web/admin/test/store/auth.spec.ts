import { setActivePinia, createPinia } from 'pinia'

import { useAuthStore } from '~/store/auth'

jest.mock('universal-cookie', () => {
  const mock = {
    set: jest.fn(),
  }
  return jest.fn(() => mock)
})

jest.mock('~/plugins/factory', () => {
  return jest.fn().mockImplementation(() => ({
    create: () => {
      return {
        v1SignIn: () => {
          return {
            data: [],
          }
        },
      }
    },
  }))
})

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('default auth state', () => {
    const authStore = useAuthStore()
    expect(authStore.redirectPath).toBe('/')
    expect(authStore.isAuthenticated).toBeFalsy()
    expect(authStore.user).toBeUndefined()
    expect(authStore.accessToken).toBeUndefined()
  })

  it('update redirectPath', () => {
    const authStore = useAuthStore()
    authStore.setRedirectPath('/signin')
    expect(authStore.redirectPath).toBe('/signin')
  })

  describe('signIn', () => {
    it('signIn success', async () => {
      const authStore = useAuthStore()
      const redirectPath = await authStore.signIn({
        username: 'admin@example.com',
        password: '122345678',
      })
      expect(redirectPath).toBe('/')
      expect(authStore.isAuthenticated).toBeTruthy()
    })

    it('signIn failed', async () => {
      // TODO: 失敗ケースをどう記述するか
      const authStore = useAuthStore()
      const redirectPath = await authStore.signIn({
        username: 'admin@example.com',
        password: '122345678',
      })
      expect(redirectPath).toBe('/')
      expect(authStore.isAuthenticated).toBeTruthy()
    })
  })
})
