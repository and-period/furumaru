import axios from 'axios'

const endPoint = 'https://zipcloud.ibsnet.co.jp/api/search'

export interface ZipCloudApiResponse {
  status: number
  message: string
  results: {
    address1: string
    address2: string
    address3: string
    kana1: string
    kana2: string
    kana3: string
    prefcode: string
    zipcode: string
  }[]
}

export interface SearchAddressResponse {
  prefecture: number
  city: string
  addressLine1: string
}

/**
 * zipCLoudを使って郵便番号から住所を検索する非同期関数
 * 参考 - http://zipcloud.ibsnet.co.jp/doc/api
 * @param code 郵便番号
 * @returns
 */
export async function searchAddressByPostalCode(
  code: number,
): Promise<SearchAddressResponse> {
  try {
    const res = await axios.get<ZipCloudApiResponse>(
      `${endPoint}?zipcode=${code}`,
    )

    if (res.data.status === 400) {
      return Promise.reject(new Error(res.data.message))
    }

    if (res.data.results === null || res.data.results.length === 0) {
      return Promise.reject(new Error('対応する住所が見つかりませんでした。'))
    }

    const result: SearchAddressResponse = {
      prefecture: Number(res.data.results[0].prefcode),
      city: res.data.results[0].address2,
      addressLine1: res.data.results[0].address3,
    }
    return result
  }
  catch (error) {
    throw new Error(
      '住所検索に失敗しました。お手数ですがご自身で入力してください。',
    )
  }
}
