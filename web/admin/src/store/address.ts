import type { Address } from '~/types/api'

export const useAddressStore = defineStore('address', {
  state: () => ({
    address: {} as Address
  })
})
