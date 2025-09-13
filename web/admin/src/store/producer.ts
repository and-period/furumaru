import { fileUpload } from './helper'
import { useCoordinatorStore } from './coordinator'
import { useShopStore } from './shop'
import type {
  CreateProducerRequest,
  ProducerResponse,
  Producer,
  UpdateProducerRequest,
  V1ProducersGetRequest,
  V1ProducersProducerIdGetRequest,
  V1ProducersPostRequest,
  V1UploadProducersThumbnailPostRequest,
  V1UploadProducersHeaderPostRequest,
  V1UploadProducersPromotionVideoPostRequest,
  V1UploadProducersBonusVideoPostRequest,
  V1ProducersProducerIdPatchRequest,
  V1ProducersProducerIdDeleteRequest,
} from '~/types/api/v1'

export const useProducerStore = defineStore('producer', {
  state: () => ({
    producer: {} as Producer,
    producers: [] as Producer[],
    totalItems: 0,
  }),

  actions: {
    /**
     * 登録済みの生産者一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchProducers(limit = 20, offset = 0, options = ''): Promise<void> {
      try {
        const params: V1ProducersGetRequest = {
          limit,
          offset,
        }
        const res = await this.producerApi().v1ProducersGet(params)

        const coordinatorStore = useCoordinatorStore()
        const shopStore = useShopStore()
        this.producers = res.producers
        this.totalItems = res.total
        coordinatorStore.coordinators = res.coordinators
        shopStore.shops = res.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 生産者を検索する非同期関数
     * @param name 生産者名(あいまい検索)
     * @param producerIds stateの更新時に残しておく必要がある生産者情報
     */
    async searchProducers(
      name = '',
      producerIds: string[] = [],
    ): Promise<void> {
      try {
        const params: V1ProducersGetRequest = {
          username: name,
        }
        const res = await this.producerApi().v1ProducersGet(params)
        const producers: Producer[] = []
        this.producers.forEach((producer: Producer): void => {
          if (!producerIds.includes(producer.id)) {
            return
          }
          producers.push(producer)
        })
        res.producers.forEach((producer: Producer): void => {
          if (producers.find((v): boolean => v.id === producer.id)) {
            return
          }
          producers.push(producer)
        })
        const shopStore = useShopStore()
        this.producers = producers
        this.totalItems = res.total
        shopStore.shops = res.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 生産者IDから生産者の情報を取得する非同期関数
     * @param producerId 生産者ID
     * @returns 生産者の情報
     */
    async getProducer(producerId: string): Promise<ProducerResponse> {
      try {
        const params: V1ProducersProducerIdGetRequest = {
          producerId,
        }
        const res = await this.producerApi().v1ProducersProducerIdGet(params)

        const coordinatorStore = useCoordinatorStore()
        const shopStore = useShopStore()
        this.producer = res.producer
        coordinatorStore.coordinators = res.coordinators
        shopStore.shops = res.shops
        return res
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '生産者の情報は閲覧権限がありません。',
          404: 'この生産者は存在しません。',
        })
      }
    },

    /**
     * 生産者を新規登録する非同期関数
     * @param payload
     */
    async createProducer(payload: CreateProducerRequest): Promise<void> {
      try {
        const params: V1ProducersPostRequest = {
          createProducerRequest: payload,
        }
        await this.producerApi().v1ProducersPost(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          409: 'このメールアドレスはすでに登録されているため、登録できません。',
        })
      }
    },

    /**
     * 生産者のサムネイル画像をアップロードする関数
     * @param payload サムネイル画像のファイルオブジェクト
     * @returns アップロード後のサムネイル画像のパスを含んだオブジェクト
     */
    async uploadProducerThumbnail(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const params: V1UploadProducersThumbnailPostRequest = {
          getUploadURLRequest: {
            fileType: contentType,
          },
        }
        const res = await this.uploadApi().v1UploadProducersThumbnailPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    /**
     * 生産者のヘッダー画像をアップロードする関数
     * @param payload ヘッダー画像のファイルオブジェクト
     * @returns アップロード後のヘッダー画像のパスを含んだオブジェクト
     */
    async uploadProducerHeader(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const params: V1UploadProducersHeaderPostRequest = {
          getUploadURLRequest: {
            fileType: contentType,
          },
        }
        const res = await this.uploadApi().v1UploadProducersHeaderPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    /**
     * 生産者の紹介画像をアップロードする非同期関数
     * @param payload 紹介画像
     * @returns アップロードされた動画のURI
     */
    async uploadProducerPromotionVideo(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const params: V1UploadProducersPromotionVideoPostRequest = {
          getUploadURLRequest: {
            fileType: contentType,
          },
        }
        const res = await this.uploadApi().v1UploadProducersPromotionVideoPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    /**
     * 生産者のサンキュー画像をアップロードする非同期関数
     * @param payload サンキュー画像
     * @returns アップロードされた動画のURI
     */
    async uploadProducerBonusVideo(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const params: V1UploadProducersBonusVideoPostRequest = {
          getUploadURLRequest: {
            fileType: contentType,
          },
        }
        const res = await this.uploadApi().v1UploadProducersBonusVideoPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    /**
     * 生産者を更新する非同期関数
     * @param producerId 更新対象の生産者ID
     * @param payload
     * @returns
     */
    async updateProducer(producerId: string, payload: UpdateProducerRequest) {
      try {
        const params: V1ProducersProducerIdPatchRequest = {
          producerId,
          updateProducerRequest: payload,
        }
        await this.producerApi().v1ProducersProducerIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '生産者の情報を更新する権限がありません。',
          404: 'この生産者は存在しません。',
        })
      }
    },

    /**
     * 生産者を削除する非同期関数
     * @param producerId 削除する生産者のID
     * @returns
     */
    async deleteProducer(producerId: string) {
      try {
        const params: V1ProducersProducerIdDeleteRequest = {
          producerId,
        }
        await this.producerApi().v1ProducersProducerIdDelete(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          403: '生産者を削除する権限がありません。',
          404: 'この生産者は存在しません。',
        })
      }
    },
  },
})
