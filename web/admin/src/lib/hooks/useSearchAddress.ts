import { useAddressStore } from '~/store'
import { ApiBaseError } from '~/types'
import type { SearchAddress } from '~/types/props'

export interface UseSearchAddress {
  loading: Ref<boolean>
  errorMessage: Ref<string>
  searchAddressByPostalCode: (postalCode: string) => Promise<SearchAddress>
}

/**
 * 郵便場号から住所を取得するカスタムフック
 */
export function useSearchAddress(): UseSearchAddress {
  const loading = ref<boolean>(false)
  const errorMessage = ref<string>('')

  const searchAddressByPostalCode = async (postalCode: string): Promise<SearchAddress> => {
    loading.value = true
    errorMessage.value = ''
    try {
      const store = useAddressStore()

      const res = await store.searchAddressByPostalCode(postalCode)
      return res
    }
    catch (err) {
      if (!(err instanceof ApiBaseError)) {
        errorMessage.value = '不明なエラーが発生しました。お手数ですがご自身で入力してください。'
        throw err
      }

      let msg: string
      switch (err.status) {
        case 400:
          msg = '入力内容に誤りがあります'
          break
        case 404:
          msg = '対応する住所が見つかりませんでした。'
          break
        default:
          msg = '不明なエラーが発生しました。お手数ですがご自身で入力してください'
      }
      errorMessage.value = msg
      throw err
    }
    finally {
      loading.value = false
    }
  }

  return {
    loading,
    errorMessage,
    searchAddressByPostalCode,
  }
}
