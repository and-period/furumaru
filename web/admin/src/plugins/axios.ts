import axios, { type AxiosInstance } from 'axios'
import { useAuthStore } from '~/store'

// eslint-disable-next-line import/no-mutable-exports
let client: AxiosInstance

export const defaultTimeout: number = 20000 // 20sec
export const uploadTimeout: number = 20000 // 200sec

export default defineNuxtPlugin(() => {
  const runtimeConfig = useRuntimeConfig()
  const baseUrl = runtimeConfig.public.API_BASE_URL

  client = axios.create({
    baseURL: baseUrl,
    timeout: defaultTimeout,
    withCredentials: true,
    headers: {}
  })

  client.interceptors.request.use((config) => {
    const store = useAuthStore()

    const token: string | undefined = store.accessToken
    if (token) {
      config.headers.setAuthorization(token)
    }

    return config
  })

  return {
    provide: {
      axios: client
    }
  }
})

export { client }
