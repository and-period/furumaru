import { createPinia, setActivePinia } from 'pinia'

import { useProducerStore } from '~/store/producer'

describe('Producer Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('default producer state', () => {
    const producerStore = useProducerStore()
    expect(producerStore.producers).toEqual([])
  })

  describe('fetchProducers', () => {})

  describe('createProducer', () => {})

  describe('uploadProducerThumbnail', () => {})

  describe('uploadProducerHeader', () => {})

  describe('getProducer', () => {})
})
