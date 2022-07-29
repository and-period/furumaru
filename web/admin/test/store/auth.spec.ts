import { setActivePinia, createPinia } from 'pinia'

import { useAuthStore } from '~/store/auth'
import { SignInRequest } from '~/types/api/api'

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
    beforeAll(() => {
      jest.mock('universal-cookie', () => {
        const mock = {
          set: jest.fn(),
        }
        return jest.fn(() => mock)
      })
      jest.mock('~/types/api/api', () => {
        return {
          // ...jest.requireActual('~/types/api/api'),
          AuthApi: {
            v1SignIn: (_: SignInRequest) => jest.fn(),
          },
        }
      })
    })

    it('signIn success', async () => {
      const authStore = useAuthStore()
      const redirectPath = await authStore.signIn({
        username: 'admin@example.com',
        password: '122345678',
      })
      expect(redirectPath).toBe('/')
    })
  })
})
