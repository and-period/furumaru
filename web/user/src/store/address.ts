import { CreateAddressRequest } from '~/types/api'

export const useAdressStore = defineStore('address', {
  state: () => {
    return {}
  },

  actions: {
    async registerAddress(payload: CreateAddressRequest) {
      await this.addressApiClient().v1CreateAddress({ body: payload })
    },
  },
})
