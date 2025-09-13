import type { Address } from '~/types/api/v1'

export const useAddressStore = defineStore('address', {
  state: () => ({
    address: {} as Address,
  }),
})
