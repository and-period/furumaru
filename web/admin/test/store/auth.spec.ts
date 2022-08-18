import { setActivePinia, createPinia } from 'pinia'

import { axiosMock, baseURL } from '../helpers/axios-helpter'

import { useAuthStore } from '~/store/auth'
import { useCommonStore } from '~/store/common'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
  ValidationError,
} from '~/types/exception'

axiosMock.onPost(`${baseURL}/v1/auth`).reply(200, {
  adminId: 'kSByoE6FetnPs5Byk3a9Zx',
  role: 1,
  accessToken: 'xxxxxxxxxx',
  refreshToken: 'xxxxxxxxxx',
  expiresIn: 3600,
  tokenType: 'Bearer',
})
axiosMock.onPatch(`${baseURL}/v1/auth/password`).reply(204, {})
axiosMock.onPost(`${baseURL}/v1/auth/refresh-token`).reply(200, {
  adminId: 'kSByoE6FetnPs5Byk3a9Zx',
  role: 1,
  accessToken: 'xxxxxxxxxx',
  refreshToken: 'xxxxxxxxxx',
  expiresIn: 3600,
  tokenType: 'Bearer',
})

jest.mock('firebase/messaging', () => {
  const mock = {
    getToken: jest.fn(),
    isSupported: jest.fn(),
  }
  return jest.fn(() => mock)
})
jest.mock('universal-cookie', () => {
  const mock = {
    set: jest.fn(),
    remove: jest.fn(),
  }
  return jest.fn(() => mock)
})
jest.mock('~/plugins/firebase', () => {
  const mock = {
    messaging: jest.fn(),
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

  describe('passwordUpdate', () => {
    it('success', async () => {
      const authStore = useAuthStore()
      const commonStore = useCommonStore()
      await authStore.passwordUpdate({
        oldPassword: '12345678',
        newPassword: 'newPass1234',
        passwordConfirmation: 'newPass1234',
      })
      expect(commonStore.snackbars.length).toBe(1)
    })

    it('failed when network error', async () => {
      axiosMock.onPatch(`${baseURL}/v1/auth/password`).networkError()

      const authStore = useAuthStore()
      try {
        await authStore.passwordUpdate({
          oldPassword: '12345678',
          newPassword: 'newPass1234',
          passwordConfirmation: 'newPass1234',
        })
      } catch (error) {
        expect(error instanceof ConnectionError).toBeTruthy()
      }
    })

    it('failed when return status code is 401', async () => {
      axiosMock.onPatch(`${baseURL}/v1/auth/password`).reply(401)

      const authStore = useAuthStore()
      try {
        await authStore.passwordUpdate({
          oldPassword: '12345678',
          newPassword: 'newPass1234',
          passwordConfirmation: 'newPass1234',
        })
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
        if (error instanceof Error) {
          expect(error.message).toBe('認証エラー。再度ログインをしてください。')
        }
      }
    })

    it('failed when return status code is 400', async () => {
      axiosMock.onPatch(`${baseURL}/v1/auth/password`).reply(400)

      const authStore = useAuthStore()
      try {
        await authStore.passwordUpdate({
          oldPassword: '12345678',
          newPassword: 'newPass1234',
          passwordConfirmation: 'newPass',
        })
      } catch (error) {
        expect(error instanceof ValidationError).toBeTruthy()
        if (error instanceof Error) {
          expect(error.message).toBe('入力値に誤りがあります。')
        }
      }
    })

    it('failed when return status code is 500', async () => {
      axiosMock.onPatch(`${baseURL}/v1/auth/password`).reply(500)

      const authStore = useAuthStore()
      try {
        await authStore.passwordUpdate({
          oldPassword: '12345678',
          newPassword: 'newPass1234',
          passwordConfirmation: 'newPass1234',
        })
      } catch (error) {
        expect(error instanceof InternalServerError).toBeTruthy()
      }
    })
  })

  describe('getAuthByRefreshToken', () => {
    it('success', async () => {
      const refreshToken = 'token'
      const authStore = useAuthStore()
      await authStore.getAuthByRefreshToken(refreshToken)
      expect(authStore.isAuthenticated).toBeTruthy()
      expect(authStore.user?.refreshToken).toBe(refreshToken)
    })

    it('failed when return network error', async () => {
      axiosMock.onPost(`${baseURL}/v1/auth/refresh-token`).networkError()

      const refreshToken = 'token'
      const authStore = useAuthStore()
      try {
        await authStore.getAuthByRefreshToken(refreshToken)
      } catch (error) {
        expect(error instanceof ConnectionError).toBeTruthy()
      }
    })

    it('failed when return status code is 401', async () => {
      axiosMock.onPost(`${baseURL}/v1/auth/refresh-token`).reply(401)

      const refreshToken = 'token'
      const authStore = useAuthStore()
      try {
        await authStore.getAuthByRefreshToken(refreshToken)
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when return status code is 500', async () => {
      axiosMock.onPost(`${baseURL}/v1/auth/refresh-token`).reply(500)

      const refreshToken = 'token'
      const authStore = useAuthStore()
      try {
        await authStore.getAuthByRefreshToken(refreshToken)
      } catch (error) {
        expect(error instanceof InternalServerError).toBeTruthy()
      }
    })
  })
})
