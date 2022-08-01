import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'
import { setActivePinia, createPinia } from 'pinia'

import { useAuthStore } from '~/store/auth'
import {
  ConnectionError,
  InternalServerError,
  ValidationError,
} from '~/types/exception'

const axiosMock = new MockAdapter(axios)
const baseURL = process.env.API_BASE_URL || 'http://localhost:18010'
axiosMock.onPost(`${baseURL}/v1/auth`).reply(201, {})

jest.mock('universal-cookie', () => {
  const mock = {
    set: jest.fn(),
  }
  return jest.fn(() => mock)
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

    it('signIn failed when network error', async () => {
      axiosMock.onPost(`${baseURL}/v1/auth`).networkError()

      const authStore = useAuthStore()
      try {
        await authStore.signIn({
          username: 'admin@example.com',
          password: '122345678',
        })
      } catch (error) {
        expect(error instanceof ConnectionError).toBeTruthy()
        expect(authStore.isAuthenticated).toBeFalsy()
      }
    })

    it('signIn failed when return status code is 404', async () => {
      axiosMock.onPost(`${baseURL}/v1/auth`).reply(400)

      const authStore = useAuthStore()
      try {
        await authStore.signIn({
          username: '',
          password: '',
        })
      } catch (error) {
        expect(error instanceof ValidationError).toBeTruthy()
        expect(authStore.isAuthenticated).toBeFalsy()
      }
    })

    it('signIn failed when return status code is 401', async () => {
      axiosMock.onPost(`${baseURL}/v1/auth`).reply(401)

      const authStore = useAuthStore()
      try {
        await authStore.signIn({
          username: 'unauthorized@email.com',
          password: '122345678',
        })
      } catch (error) {
        expect(error instanceof ValidationError).toBeTruthy()
        expect(authStore.isAuthenticated).toBeFalsy()
      }
    })

    it('signIn failed when return status code is 500', async () => {
      axiosMock.onPost(`${baseURL}/v1/auth`).reply(500)

      const authStore = useAuthStore()
      try {
        await authStore.signIn({
          username: 'admin@example.com',
          password: '122345678',
        })
      } catch (error) {
        expect(error instanceof InternalServerError).toBeTruthy()
        expect(authStore.isAuthenticated).toBeFalsy()
      }
    })
  })
})
