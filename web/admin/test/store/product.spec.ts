import { createPinia, setActivePinia } from 'pinia'

import { setupAuthStore } from '../helpers/auth-helpter'
import { axiosMock, baseURL } from '../helpers/axios-helpter'

import { useProductStore } from '~/store'
import { ProductApi } from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
} from '~/types/exception'

jest.mock('~/plugins/firebase', () => {
  const mock = {
    messaging: jest.fn(),
  }
  return jest.fn(() => mock)
})

describe('Product Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('default product store', () => {
    const productStore = useProductStore()
    expect(productStore.products).toEqual([])
    expect(productStore.apiClient('') instanceof ProductApi)
  })

  describe('fetchProducts', () => {
    const productPath = `${baseURL}/v1/products?limit=20&offset=0`
    const product = {
      id: 'kSByoE6FetnPs5Byk3a9Zx',
      name: '新鮮なじゃがいも',
      description: '新鮮なじゃがいもをお届けします。',
      producerId: 'kSByoE6FetnPs5Byk3a9Zx',
      storeName: '&.農園',
      categoryId: 'kSByoE6FetnPs5Byk3a9Zx',
      categoryName: '野菜',
      productTypeId: 'kSByoE6FetnPs5Byk3a9Zx',
      productTypeName: 'じゃがいも',
      productTypeIconUrl: 'https://and-period.jp/icon.png',
      public: true,
      inventory: 30,
      weight: 2.5,
      itemUnit: '袋',
      itemDescription: '1袋あたり2.5kgのじゃがいも',
      media: [
        {
          url: 'https://and-period.jp/thumbnail01.png',
          isThumbnail: true,
        },
        {
          url: 'https://and-period.jp/thumbnail02.png',
          isThumbnail: false,
        },
      ],
      price: 2500,
      deliveryType: 1,
      box60Rate: 80,
      box80Rate: 50,
      box100Rate: 40,
      originPrefecture: '滋賀県',
      originCity: '彦根市',
      createdBy: 'kSByoE6FetnPs5Byk3a9Zx',
      updatedBy: 'kSByoE6FetnPs5Byk3a9Zx',
      createdAt: 1640962800,
      updatedAt: 1640962800,
    }
    axiosMock.onGet(productPath).reply(200, {
      products: [product],
      total: 1,
    })

    it('success', async () => {
      setupAuthStore(true)

      const productStore = useProductStore()
      await productStore.fetchProducts()
      expect(productStore.products.length).toEqual(1)
      expect(productStore.products[0]).toEqual(product)
    })

    it('failed when not authenticated', async () => {
      setupAuthStore(false)

      const productStore = useProductStore()
      try {
        await productStore.fetchProducts()
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when network error', async () => {
      axiosMock.onGet(productPath).networkError()

      setupAuthStore(true)
      const productStore = useProductStore()
      try {
        await productStore.fetchProducts()
      } catch (error) {
        expect(error instanceof ConnectionError).toBeTruthy()
      }
    })

    it('failed when return status code is 401', async () => {
      axiosMock.onGet(productPath).reply(401)

      setupAuthStore(true)
      const productStore = useProductStore()
      try {
        await productStore.fetchProducts()
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when return status code is 500', async () => {
      axiosMock.onGet(productPath).reply(500)

      setupAuthStore(true)
      const productStore = useProductStore()
      try {
        await productStore.fetchProducts()
      } catch (error) {
        expect(error instanceof InternalServerError).toBeTruthy()
      }
    })
  })
})
