import { Address, AddressesResponse, CreateAddressRequest } from '~/types/api'

export const useAdressStore = defineStore('address', {
  state: () => {
    return {
      total: 0,
      addresses: [] as Address[],
      addressesFetchState: {
        isLoading: false,
      },
    }
  },

  actions: {
    async registerAddress(payload: CreateAddressRequest) {
      await this.addressApiClient().v1CreateAddress({ body: payload })
    },

    async fetchAddresses(limit: number = 20, offset: number = 0) {
      this.addressesFetchState.isLoading = true
      const res: AddressesResponse =
        await this.addressApiClient().v1ListAddresses({
          limit,
          offset,
        })
      this.addresses = res.addresses
      this.total = res.total
      this.addressesFetchState.isLoading = false
    },
  },
})
