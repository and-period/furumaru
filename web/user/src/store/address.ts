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

    async registerAddress(payload: CreateAddressRequest): Promise<Address> {
      const authStore = useAuthStore()

      const res = await this.addressApiClient(
        authStore.accessToken,
      ).v1CreateAddress({
        body: payload,
      })

      return res.address
    },

    async fetchAddresses(limit: number = 20, offset: number = 0) {
      const authStore = useAuthStore()

      this.addressesFetchState.isLoading = true
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

  getters: {
    defaultAddress(): Address | undefined {
      return this.addresses.find((address) => address.isDefault)
    },
  },
})
