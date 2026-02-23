import { useApiClient } from '~/composables/useApiClient'
import { PostalCodeApi } from '~/types/api/v1'
import type { Address, V1PostalCodesPostalCodeGetRequest } from '~/types/api/v1'
import type { Prefecture } from '~/types'
import type { SearchAddress } from '~/types/props'

export const useAddressStore = defineStore('address', () => {
  const { create, errorHandler } = useApiClient()
  const postalCodeApi = () => create(PostalCodeApi)

  const address = ref<Address>({} as Address)

  async function searchAddressByPostalCode(postalCode: string): Promise<SearchAddress> {
    try {
      const params: V1PostalCodesPostalCodeGetRequest = { postalCode }
      const res = await postalCodeApi().v1PostalCodesPostalCodeGet(params)
      return {
        prefecture: res.prefectureCode as Prefecture,
        city: res.city,
        town: res.town,
      }
    }
    catch (error) {
      return errorHandler(error)
    }
  }

  return {
    address,
    searchAddressByPostalCode,
  }
})
