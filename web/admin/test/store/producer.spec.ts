import { createPinia, setActivePinia } from 'pinia'

import { setupAuthStore } from '../helpers/auth-helpter'
import { axiosMock, baseURL } from '../helpers/axios-helpter'

import { useProducerStore } from '~/store/producer'
import { AuthError } from '~/types/exception'

axiosMock
  .onGet(`${baseURL}/v1/producers?limit=20&offset=0`)
  .reply(200, { producers: [] })

describe('Producer Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('default producer state', () => {
    const producerStore = useProducerStore()
    expect(producerStore.producers).toEqual([])
  })

  describe('fetchProducers', () => {
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
  })

  describe('createProducer', () => {})

  describe('uploadProducerThumbnail', () => {})

  describe('uploadProducerHeader', () => {})

  describe('getProducer', () => {})
})
