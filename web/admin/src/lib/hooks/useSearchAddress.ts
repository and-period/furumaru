import {
  searchAddressByPostalCode as _searchAddressByPostalCode,
  SearchAddressResponse
} from '../externals'

export interface UseSearchAddress {
  loading: Ref<boolean>
  errorMessage: Ref<string>
  searchAddressByPostalCode: (
    postalCode: number
  ) => Promise<SearchAddressResponse | undefined>
}

/**
 * 郵便場号から住所を取得するカスタムフック
 */
export function useSearchAddress (): UseSearchAddress {
  const loading = ref<boolean>(false)
  const errorMessage = ref<string>('')

  const searchAddressByPostalCode = async (postalCode: number) => {
    loading.value = true
    errorMessage.value = ''
    try {
      return await _searchAddressByPostalCode(Number(postalCode))
    } catch (e) {
      if (e instanceof Error) {
        errorMessage.value = e.message
      } else {
        errorMessage.value =
          '不明なエラーが発生しました。お手数ですがご自身で入力してください。'
      }
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    errorMessage,
    searchAddressByPostalCode
  }
}
