import { useAuthStore } from './auth'
import {
  Address,
  AddressesResponse,
  CreateAddressRequest,
  PostalCodeResponse,
} from '~/types/api'

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
    async searchAddressByPostalCode(
      postalCode: string,
    ): Promise<PostalCodeResponse> {
      const res = await this.addressApiClient().v1SearchPostalCode({
        postalCode,
      })
      return res
    },

    async registerAddress(payload: CreateAddressRequest) {
      await this.addressApiClient().v1CreateAddress({ body: payload })
    },

    async fetchAddresses(limit: number = 20, offset: number = 0) {
      const authStore = useAuthStore()

      this.addressesFetchState.isLoading = true
      const client = this.addressApiClient(authStore.accessToken)
      console.log(client)
      const res: AddressesResponse = await this.addressApiClient(
        authStore.accessToken,
      ).v1ListAddresses({
        limit,
        offset,
      })
      this.addresses = res.addresses
      this.total = res.total
      this.addressesFetchState.isLoading = false
    },
  },
})
