import { createPinia, setActivePinia } from 'pinia'

import { setupAuthStore } from '../helpers/auth-helpter'
import { axiosMock, baseURL } from '../helpers/axios-helpter'

import { useProducerStore } from '~/store/producer'
import { CreateProducerRequest } from '~/types/api'
import {
  AuthError,
  ConflictError,
  ConnectionError,
  InternalServerError,
  NotFoundError,
  ValidationError,
} from '~/types/exception'

jest.mock('firebase/messaging', () => {
  const mock = {
    getToken: jest.fn(),
    isSupported: jest.fn(),
  }
  return jest.fn(() => mock)
})
jest.mock('~/plugins/firebase', () => {
  const mock = {
    messaging: jest.fn(),
  }
  return jest.fn(() => mock)
})

describe('Producer Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('default producer state', () => {
    const producerStore = useProducerStore()
    expect(producerStore.producers).toEqual([])
  })

  describe('fetchProducers', () => {
    const producersPath = `${baseURL}/v1/producers?limit=20&offset=0`
    axiosMock.onGet(producersPath).reply(200, { producers: [] })

    it('success', async () => {
      setupAuthStore(true)

      const producerStore = useProducerStore()
      await producerStore.fetchProducers()
      expect(producerStore.producers).toEqual([])
    })

    it('failed when not authenticated', async () => {
      setupAuthStore(false)

      const producerStore = useProducerStore()
      try {
        await producerStore.fetchProducers()
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
        if (error instanceof AuthError) {
          expect(error.cause).toBeUndefined()
        }
      }
    })

    it('failed when network error', async () => {
      axiosMock.onGet(producersPath).networkError()

      setupAuthStore(true)
      const producerStore = useProducerStore()
      try {
        await producerStore.fetchProducers()
      } catch (error) {
        expect(error instanceof ConnectionError).toBeTruthy()
      }
    })

    it('failed when return status code is 401', async () => {
      axiosMock.onGet(producersPath).reply(401)

      setupAuthStore(true)
      const producerStore = useProducerStore()
      try {
        await producerStore.fetchProducers()
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
        if (error instanceof AuthError) {
          expect(error.message).toBe('認証エラー。再度ログインをしてください。')
        }
      }
    })

    it('failed when return status code is 500', async () => {
      axiosMock.onGet(producersPath).reply(500)

      setupAuthStore(true)
      const producerStore = useProducerStore()
      try {
        await producerStore.fetchProducers()
      } catch (error) {
        expect(error instanceof InternalServerError).toBeTruthy()
      }
    })
  })

  describe('createProducer', () => {
    axiosMock.onPost(`${baseURL}/v1/producers`).reply(200, {
      id: 'kSByoE6FetnPs5Byk3a9Zx',
      coordinatorId: 'kSByoE6FetnPs5Byk3a9Zx',
      lastname: '&.',
      firstname: '管理者',
      lastnameKana: 'あんどどっと',
      firstnameKana: 'かんりしゃ',
      storeName: '&.農園',
      thumbnailUrl: 'https://and-period.jp/thumbnail.png',
      headerUrl: 'https://and-period.jp/header.png',
      email: 'test-user@and-period.jp',
      phoneNumber: 819012345678,
      postalCode: '1000014',
      prefecture: '東京都',
      city: '千代田区',
      addressLine1: '永田町1-7-1',
      addressLine2: '',
      createdAt: 1640962800,
      updatedAt: 1640962800,
    })

    it('success', () => {
      setupAuthStore(true)

      const validPayload: CreateProducerRequest = {
        coordinatorId: 'kSByoE6FetnPs5Byk3a9Zx',
        lastname: '&.',
        firstname: '管理者',
        lastnameKana: 'あんどどっと',
        firstnameKana: 'かんりしゃ',
        storeName: '&.農園',
        thumbnailUrl: 'https://and-period.jp/thumbnail.png',
        headerUrl: 'https://and-period.jp/header.png',
        email: 'test-user@and-period.jp',
        phoneNumber: '819012345678',
        postalCode: '1000014',
        prefecture: '東京都',
        city: '千代田区',
        addressLine1: '永田町1-7-1',
        addressLine2: '',
      }

      const producerStore = useProducerStore()
      return expect(producerStore.createProducer(validPayload)).resolves.toBe(
        undefined
      )
    })

    it('failed when not authenticated', async () => {
      setupAuthStore(false)

      const validPayload: CreateProducerRequest = {
        coordinatorId: 'kSByoE6FetnPs5Byk3a9Zx',
        lastname: '&.',
        firstname: '管理者',
        lastnameKana: 'あんどどっと',
        firstnameKana: 'かんりしゃ',
        storeName: '&.農園',
        thumbnailUrl: 'https://and-period.jp/thumbnail.png',
        headerUrl: 'https://and-period.jp/header.png',
        email: 'test-user@and-period.jp',
        phoneNumber: '819012345678',
        postalCode: '1000014',
        prefecture: '東京都',
        city: '千代田区',
        addressLine1: '永田町1-7-1',
        addressLine2: '',
      }

      const producerStore = useProducerStore()
      try {
        await producerStore.createProducer(validPayload)
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when return status code is 400', async () => {
      setupAuthStore(true)
      axiosMock.onPost(`${baseURL}/v1/producers`).reply(400)

      const invalidPayload: CreateProducerRequest = {
        coordinatorId: 'kSByoE6FetnPs5Byk3a9Zx',
        lastname: '&.',
        firstname: '管理者',
        lastnameKana: 'アンドドット',
        firstnameKana: 'カンリシャ',
        storeName: '&.農園',
        thumbnailUrl: 'https://and-period.jp/thumbnail.png',
        headerUrl: 'https://and-period.jp/header.png',
        email: 'test-user@and-period.jp',
        phoneNumber: '819012345678',
        postalCode: '1000014',
        prefecture: '東京都',
        city: '千代田区',
        addressLine1: '永田町1-7-1',
        addressLine2: '',
      }

      const producerStore = useProducerStore()
      try {
        await producerStore.createProducer(invalidPayload)
      } catch (error) {
        expect(error instanceof ValidationError).toBeTruthy()
        if (error instanceof ValidationError) {
          expect(error.message).toBe('入力内容に誤りがあります。')
        }
      }
    })

    it('failed when return status code is 400', async () => {
      setupAuthStore(true)
      axiosMock.onPost(`${baseURL}/v1/producers`).reply(409)

      const notUniquePayload: CreateProducerRequest = {
        coordinatorId: 'kSByoE6FetnPs5Byk3a9Zx',
        lastname: '&.',
        firstname: '管理者',
        lastnameKana: 'あんどどっと',
        firstnameKana: 'かんりしゃ',
        storeName: '&.農園',
        thumbnailUrl: 'https://and-period.jp/thumbnail.png',
        headerUrl: 'https://and-period.jp/header.png',
        email: 'test-user@and-period.jp',
        phoneNumber: '819012345678',
        postalCode: '1000014',
        prefecture: '東京都',
        city: '千代田区',
        addressLine1: '永田町1-7-1',
        addressLine2: '',
      }

      const producerStore = useProducerStore()
      try {
        await producerStore.createProducer(notUniquePayload)
      } catch (error) {
        expect(error instanceof ConflictError).toBeTruthy()
        if (error instanceof ConflictError) {
          expect(error.message).toBe(
            'このメールアドレスはすでに登録されているため、登録できません。'
          )
        }
      }
    })

    it('failed when return status code is 401', async () => {
      setupAuthStore(true)
      axiosMock.onPost(`${baseURL}/v1/producers`).reply(401)

      const validPayload: CreateProducerRequest = {
        coordinatorId: 'kSByoE6FetnPs5Byk3a9Zx',
        lastname: '&.',
        firstname: '管理者',
        lastnameKana: 'あんどどっと',
        firstnameKana: 'かんりしゃ',
        storeName: '&.農園',
        thumbnailUrl: 'https://and-period.jp/thumbnail.png',
        headerUrl: 'https://and-period.jp/header.png',
        email: 'test-user@and-period.jp',
        phoneNumber: '819012345678',
        postalCode: '1000014',
        prefecture: '東京都',
        city: '千代田区',
        addressLine1: '永田町1-7-1',
        addressLine2: '',
      }

      const producerStore = useProducerStore()
      try {
        await producerStore.createProducer(validPayload)
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when return status code is 401', async () => {
      setupAuthStore(true)
      axiosMock.onPost(`${baseURL}/v1/producers`).reply(500)

      const validPayload: CreateProducerRequest = {
        coordinatorId: 'kSByoE6FetnPs5Byk3a9Zx',
        lastname: '&.',
        firstname: '管理者',
        lastnameKana: 'あんどどっと',
        firstnameKana: 'かんりしゃ',
        storeName: '&.農園',
        thumbnailUrl: 'https://and-period.jp/thumbnail.png',
        headerUrl: 'https://and-period.jp/header.png',
        email: 'test-user@and-period.jp',
        phoneNumber: '819012345678',
        postalCode: '1000014',
        prefecture: '東京都',
        city: '千代田区',
        addressLine1: '永田町1-7-1',
        addressLine2: '',
      }

      const producerStore = useProducerStore()
      try {
        await producerStore.createProducer(validPayload)
      } catch (error) {
        expect(error instanceof InternalServerError).toBeTruthy()
      }
    })
  })

  describe('uploadProducerThumbnail', () => {
    const dummyFile: File = new File(['dummy'], 'image.png', {
      type: 'image/png',
    })
    axiosMock
      .onPost(`${baseURL}/v1/upload/producers/thumbnail`)
      .reply(200, { url: 'https://and-period.jp/image.png' })

    it('success', async () => {
      setupAuthStore(true)

      const producerStore = useProducerStore()
      const actual = await producerStore.uploadProducerThumbnail(dummyFile)
      expect(actual.url).toBe('https://and-period.jp/image.png')
    })

    it('failed when not authenticated', async () => {
      setupAuthStore(false)

      const producerStore = useProducerStore()
      try {
        await producerStore.uploadProducerThumbnail(dummyFile)
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when return status code is 400', async () => {
      axiosMock.onPost(`${baseURL}/v1/upload/producers/thumbnail`).reply(400)
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.uploadProducerThumbnail(dummyFile)
      } catch (error) {
        expect(error instanceof ValidationError).toBeTruthy()
      }
    })

    it('failed when return status code is 401', async () => {
      axiosMock.onPost(`${baseURL}/v1/upload/producers/thumbnail`).reply(401)
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.uploadProducerThumbnail(dummyFile)
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when return status code is 500', async () => {
      axiosMock.onPost(`${baseURL}/v1/upload/producers/thumbnail`).reply(500)
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.uploadProducerThumbnail(dummyFile)
      } catch (error) {
        expect(error instanceof InternalServerError).toBeTruthy()
      }
    })
  })

  describe('uploadProducerHeader', () => {
    const dummyFile: File = new File(['dummy'], 'image.png', {
      type: 'image/png',
    })
    axiosMock
      .onPost(`${baseURL}/v1/upload/producers/header`)
      .reply(200, { url: 'https://and-period.jp/image.png' })

    it('success', async () => {
      setupAuthStore(true)

      const producerStore = useProducerStore()
      const actual = await producerStore.uploadProducerHeader(dummyFile)
      expect(actual.url).toBe('https://and-period.jp/image.png')
    })

    it('failed when not authenticated', async () => {
      setupAuthStore(false)

      const producerStore = useProducerStore()
      try {
        await producerStore.uploadProducerHeader(dummyFile)
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when return status code is 400', async () => {
      axiosMock.onPost(`${baseURL}/v1/upload/producers/header`).reply(400)
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.uploadProducerHeader(dummyFile)
      } catch (error) {
        expect(error instanceof ValidationError).toBeTruthy()
      }
    })

    it('failed when return status code is 401', async () => {
      axiosMock.onPost(`${baseURL}/v1/upload/producers/header`).reply(401)
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.uploadProducerHeader(dummyFile)
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when return status code is 500', async () => {
      axiosMock.onPost(`${baseURL}/v1/upload/producers/header`).reply(500)
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.uploadProducerHeader(dummyFile)
      } catch (error) {
        expect(error instanceof InternalServerError).toBeTruthy()
      }
    })
  })

  describe('getProducer', () => {
    const producerId = 'kSByoE6FetnPs5Byk3a9Zx'
    axiosMock.onGet(`${baseURL}/v1/producers/${producerId}`).reply(200, {
      id: 'kSByoE6FetnPs5Byk3a9Zx',
      lastname: '&.',
      firstname: '管理者',
      lastnameKana: 'あんどどっと',
      firstnameKana: 'かんりしゃ',
      storeName: '&.農園',
      thumbnailUrl: 'https://and-period.jp/thumbnail.png',
      headerUrl: 'https://and-period.jp/header.png',
      email: 'test-user@and-period.jp',
      phoneNumber: 819012345678,
      postalCode: '1000014',
      prefecture: '東京都',
      city: '千代田区',
      addressLine1: '永田町1-7-1',
      addressLine2: '',
      createdAt: 1640962800,
      updatedAt: 1640962800,
    })

    it('success', async () => {
      setupAuthStore(true)

      const producerStore = useProducerStore()
      const actual = await producerStore.getProducer(producerId)
      expect(actual.id).toBe(producerId)
    })

    it('failed when not authenticated', async () => {
      setupAuthStore(false)

      const producerStore = useProducerStore()
      try {
        await producerStore.getProducer(producerId)
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when network error', async () => {
      axiosMock.onGet(`${baseURL}/v1/producers/${producerId}`).networkError()
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.getProducer(producerId)
      } catch (error) {
        expect(error instanceof ConnectionError).toBeTruthy()
      }
    })

    it('failed when return status code is 401', async () => {
      axiosMock.onGet(`${baseURL}/v1/producers/${producerId}`).reply(401)
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.getProducer(producerId)
      } catch (error) {
        expect(error instanceof AuthError).toBeTruthy()
      }
    })

    it('failed when return status code is 404', async () => {
      axiosMock.onGet(`${baseURL}/v1/producers/${producerId}`).reply(404)
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.getProducer(producerId)
      } catch (error) {
        expect(error instanceof NotFoundError).toBeTruthy()
      }
    })

    it('failed when return status code is 500', async () => {
      axiosMock.onGet(`${baseURL}/v1/producers/${producerId}`).reply(500)
      setupAuthStore(true)

      const producerStore = useProducerStore()
      try {
        await producerStore.getProducer(producerId)
      } catch (error) {
        expect(error instanceof InternalServerError).toBeTruthy()
      }
    })
  })
})
