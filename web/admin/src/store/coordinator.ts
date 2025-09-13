import { fileUpload } from './helper'
import { useProductTypeStore } from './product-type'
import { useShopStore } from './shop'
import type {
  Coordinator,
  CreateCoordinatorRequest,
  Producer,
  UpdateCoordinatorRequest,
  V1CoordinatorsCoordinatorIdDeleteRequest,
  V1CoordinatorsCoordinatorIdGetRequest,
  V1CoordinatorsCoordinatorIdPatchRequest,
  V1CoordinatorsGetRequest,
  V1CoordinatorsPostRequest,
  V1UploadCoordinatorsBonusVideoPostRequest,
  V1UploadCoordinatorsHeaderPostRequest,
  V1UploadCoordinatorsPromotionVideoPostRequest,
  V1UploadCoordinatorsThumbnailPostRequest,
} from '~/types/api/v1'

export const useCoordinatorStore = defineStore('coordinator', {
  state: () => ({
    coordinator: {} as Coordinator,
    coordinators: [] as Coordinator[],
    producers: [] as Producer[],
    totalItems: 0,
    producerTotalItems: 0,
  }),

  actions: {
    /**
     * コーディネーターの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     */
    async fetchCoordinators(limit = 20, offset = 0): Promise<void> {
      try {
        const params: V1CoordinatorsGetRequest = {
          limit,
          offset,
        }
        const res = await this.coordinatorApi().v1CoordinatorsGet(params)

        const productTypeStore = useProductTypeStore()
        const shopStore = useShopStore()
        this.coordinators = res.coordinators
        this.totalItems = res.total
        productTypeStore.productTypes = res.productTypes
        shopStore.shops = res.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターを検索する非同期関数
     * @param name コーディネーター名(あいまい検索)
     * @param coordinatorIds stateの更新時に残しておく必要があるコーディネーター情報
     */
    async searchCoordinators(
      name = '',
      coordinatorIds: string[] = [],
    ): Promise<void> {
      try {
        const params: V1CoordinatorsGetRequest = {
          username: name,
        }
        const res = await this.coordinatorApi().v1CoordinatorsGet(params)
        const coordinators: Coordinator[] = []
        this.coordinators.forEach((coordinator: Coordinator): void => {
          if (!coordinatorIds.includes(coordinator.id)) {
            return
          }
          coordinators.push(coordinator)
        })
        res.coordinators.forEach((coordinator: Coordinator): void => {
          if (coordinators.find((v): boolean => v.id === coordinator.id)) {
            return
          }
          coordinators.push(coordinator)
        })
        const shopStore = useShopStore()
        this.coordinators = coordinators
        this.totalItems = res.total
        shopStore.shops = res.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * コーディネーターの詳細情報を取得する非同期関数
     * @param coordinatorId 対象のコーディネーターのID
     */
    async getCoordinator(coordinatorId: string): Promise<void> {
      try {
        const params: V1CoordinatorsCoordinatorIdGetRequest = {
          coordinatorId,
        }
        const res = await this.coordinatorApi().v1CoordinatorsCoordinatorIdGet(params)

        const productTypeStore = useProductTypeStore()
        const shopStore = useShopStore()
        this.coordinator = res.coordinator
        productTypeStore.productTypes = res.productTypes
        shopStore.shop = res.shop
      }
      catch (err) {
        return this.errorHandler(err, {
          404: 'コーディネーター情報が見つかりません。',
        })
      }
    },

    /**
     * コーディネーターを登録する非同期関数
     * @param payload
     */
    async createCoordinator(payload: CreateCoordinatorRequest) {
      try {
        const params: V1CoordinatorsPostRequest = {
          createCoordinatorRequest: payload,
        }
        const res = await this.coordinatorApi().v1CoordinatorsPost(params)
        return res
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          409: 'このメールアドレスはすでに登録されているため、登録できません。',
        })
      }
    },

    /**
     * コーディネーターの情報を更新する非同期関数
     * @param payload
     * @param coordinatorId 更新するコーディネーターのID
     */
    async updateCoordinator(
      coordinatorId: string,
      payload: UpdateCoordinatorRequest,
    ): Promise<void> {
      try {
        const params: V1CoordinatorsCoordinatorIdPatchRequest = {
          coordinatorId,
          updateCoordinatorRequest: payload,
        }
        await this.coordinatorApi().v1CoordinatorsCoordinatorIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、入力内容に誤りがあります。',
          404: '対象のコーディネーターが存在しません',
        })
      }
    },

    /**
     * コーディネーターのサムネイル画像をアップロードするためのURLを取得する非同期関数
     * @param payload サムネイル画像
     * @returns アップロードされた画像のURI
     */
    async uploadCoordinatorThumbnail(payload: File): Promise<string> {
      try {
        const params: V1UploadCoordinatorsThumbnailPostRequest = {
          getUploadURLRequest: {
            fileType: payload.type,
          },
        }
        const res = await this.uploadApi().v1UploadCoordinatorsThumbnailPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'ファイルのアップロードに失敗しました',
        })
      }
    },

    /**
     * コーディネーターのヘッダー画像をアップロードする非同期関数
     * @param payload ヘッダー画像
     * @returns アップロードされた画像のURI
     */
    async uploadCoordinatorHeader(payload: File): Promise<string> {
      try {
        const params: V1UploadCoordinatorsHeaderPostRequest = {
          getUploadURLRequest: {
            fileType: payload.type,
          },
        }
        const res = await this.uploadApi().v1UploadCoordinatorsHeaderPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'ファイルのアップロードに失敗しました',
        })
      }
    },

    /**
     * コーディネーターの紹介画像をアップロードする非同期関数
     * @param payload 紹介画像
     * @returns アップロードされた動画のURI
     */
    async uploadCoordinatorPromotionVideo(payload: File): Promise<string> {
      try {
        const params: V1UploadCoordinatorsPromotionVideoPostRequest = {
          getUploadURLRequest: {
            fileType: payload.type,
          },
        }
        const res = await this.uploadApi().v1UploadCoordinatorsPromotionVideoPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'ファイルのアップロードに失敗しました',
        })
      }
    },

    /**
     * コーディネーターのサンキュー画像をアップロードする非同期関数
     * @param payload サンキュー画像
     * @returns アップロードされた動画のURI
     */
    async uploadCoordinatorBonusVideo(payload: File): Promise<string> {
      try {
        const params: V1UploadCoordinatorsBonusVideoPostRequest = {
          getUploadURLRequest: {
            fileType: payload.type,
          },
        }
        const res = await this.uploadApi().v1UploadCoordinatorsBonusVideoPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'ファイルのアップロードに失敗しました',
        })
      }
    },

    /**
     * コーディーネータを削除する非同期関数
     * @param id 削除するコーディネーターのID
     * @returns
     */
    async deleteCoordinator(id: string) {
      try {
        const params: V1CoordinatorsCoordinatorIdDeleteRequest = {
          coordinatorId: id,
        }
        await this.coordinatorApi().v1CoordinatorsCoordinatorIdDelete(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: '対象のコーディネーターが存在しません',
        })
      }
      this.fetchCoordinators()
    },
  },
})
