import { useAuthStore } from './auth'
import type {
  Address,
  AddressesResponse,
  CreateAddressRequest,
  GuestCheckoutAddress,
  PostalCodeResponse,
} from '~/types/api'

export const useAddressStore = defineStore('address', {
  state: () => {
    return {
      total: 0,
      address: undefined as Address | undefined,
      addressFetchState: {
        isLoading: false,
      },
      addresses: [] as Address[],
      guestAddress: undefined as GuestCheckoutAddress | undefined,
      email: null as string | null,
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

    async fetchAddress(id: string) {
      const authStore = useAuthStore()

      this.addressFetchState.isLoading = true

      const res = await this.addressApiClient(
        authStore.accessToken,
      ).v1GetAddress({
        addressId: id,
      })
      this.address = res.address
      this.addressFetchState.isLoading = false
    },
  },

  getters: {
    defaultAddress(): Address | undefined {
      return this.addresses.find((address) => address.isDefault)
    },
  },
})
