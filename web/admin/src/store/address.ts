import type { Prefecture } from '~/types'
import type { SearchAddress } from '~/types/props'
import type { Address, V1PostalCodesPostalCodeGetRequest } from '~/types/api/v1'

export const useAddressStore = defineStore('address', {
  state: () => ({
    address: {} as Address,
  }),

  actions: {
    async searchAddressByPostalCode(postalCode: string): Promise<SearchAddress> {
      try {
        const params: V1PostalCodesPostalCodeGetRequest = {
          postalCode,
        }
        const res = await this.postalCodeApi().v1PostalCodesPostalCodeGet(params)
        return {
          prefecture: res.prefectureCode as Prefecture,
          city: res.city,
          town: res.town,
        }
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },
  },
})
