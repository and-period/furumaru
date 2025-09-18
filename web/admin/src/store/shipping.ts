import type { CreateShippingRequest, Shipping, ShippingsResponse, UpdateDefaultShippingRequest, UpdateShippingRequest, UpsertShippingRequest, V1CoordinatorsCoordinatorIdShippingsActivationGetRequest, V1CoordinatorsCoordinatorIdShippingsGetRequest, V1CoordinatorsCoordinatorIdShippingsPatchRequest, V1CoordinatorsCoordinatorIdShippingsPostRequest, V1CoordinatorsCoordinatorIdShippingsShippingIdDeleteRequest, V1CoordinatorsCoordinatorIdShippingsShippingIdGetRequest, V1CoordinatorsCoordinatorIdShippingsShippingIdPatchRequest, V1ShippingsDefaultPatchRequest } from '~/types/api/v1'

export const useShippingStore = defineStore('shipping', {
  state: () => ({
    shipping: {} as Shipping,
    shippings: [] as Array<Shipping>,
    total: 0,
  }),

  actions: {
    /**
     * コーディネーターが登録している配送先一覧を取得する非同期関数
     */
    async fetchShippings(coordinatorId: string, limit?: number, offset?: number): Promise<ShippingsResponse> {
      try {
        const params: V1CoordinatorsCoordinatorIdShippingsGetRequest = {
          coordinatorId,
          limit,
          offset,
        }
        const res = await this.shippingApi().v1CoordinatorsCoordinatorIdShippingsGet(params)
        this.shippings = res.shippings
        this.total = res.total
        return res
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のコーディネーターが見つかりません。' })
      }
    },

    /**
     * 新規の配送設定を作成する非同期関数
     * @param coordinatorId コーディネーターID
     * @param payload 配送先情報
     * @returns
     */
    async createShipping(coordinatorId: string, payload: CreateShippingRequest): Promise<void> {
      try {
        const params: V1CoordinatorsCoordinatorIdShippingsPostRequest = {
          coordinatorId,
          createShippingRequest: payload,
        }
        await this.shippingApi().v1CoordinatorsCoordinatorIdShippingsPost(params)
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 指定した配送設定を更新する非同期関数
     * @param coordinatorId コーディネーターID
     * @param shippingId 配送設定ID
     * @param payload 配送先情報
     * @returns
     */
    async updateShipping(coordinatorId: string, shippingId: string, payload: UpdateShippingRequest): Promise<void> {
      try {
        const params: V1CoordinatorsCoordinatorIdShippingsShippingIdPatchRequest = {
          coordinatorId,
          shippingId,
          updateShippingRequest: payload,
        }
        await this.shippingApi().v1CoordinatorsCoordinatorIdShippingsShippingIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 指定した配送設定を削除する非同期関数
     * @param coordinatorId コーディネーターID
     * @param shippingId 配送設定ID
     * @returns
     */
    async deleteShipping(coordinatorId: string, shippingId: string): Promise<void> {
      try {
        const params: V1CoordinatorsCoordinatorIdShippingsShippingIdDeleteRequest = {
          coordinatorId,
          shippingId,
        }
        await this.shippingApi().v1CoordinatorsCoordinatorIdShippingsShippingIdDelete(params)
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象の配送設定が見つかりません。' })
      }
    },

    /**
     * 指定した配送設定を取得する非同期関数
     * @param coordinatorId コーディネーターID
     * @param shippingId 配送設定ID
     * @returns
     */
    async fetchShipping(coordinatorId: string, shippingId: string): Promise<Shipping> {
      try {
        const params: V1CoordinatorsCoordinatorIdShippingsShippingIdGetRequest = {
          coordinatorId,
          shippingId,
        }
        const res = await this.shippingApi().v1CoordinatorsCoordinatorIdShippingsShippingIdGet(params)
        return res.shipping
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '配送設定が見つかりません。',
        })
      }
    },

    /**
     * デフォルト配送設定を取得する非同期関数
     * @returns
     */
    async fetchDefaultShipping(): Promise<void> {
      try {
        const res = await this.shippingApi().v1ShippingsDefaultGet()
        this.shipping = res.shipping
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * デフォルトの配送設定を変更する非同期関数
     * @param payload
     * @returns
     */
    async updateDefaultShipping(payload: UpdateDefaultShippingRequest): Promise<void> {
      try {
        const params: V1ShippingsDefaultPatchRequest = {
          updateDefaultShippingRequest: payload,
        }
        await this.shippingApi().v1ShippingsDefaultPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, { 400: '必須項目が不足しているか、入力内容に誤りがあります。' })
      }
    },

    /**
     * 指定したコーディネーターの配送設定を取得する非同期関数
     * @param coordinatorId
     * @returns
     */
    async fetchActiveShipping(coordinatorId: string): Promise<void> {
      try {
        const params: V1CoordinatorsCoordinatorIdShippingsActivationGetRequest = {
          coordinatorId,
        }
        const res = await this.shippingApi().v1CoordinatorsCoordinatorIdShippingsActivationGet(params)
        this.shipping = res.shipping
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のコーディネーターが見つかりません。' })
      }
    },

    /**
     * 指定した配送設定を更新する非同期関数
     * @param coordinatorId
     * @param shippingId
     * @param payload
     * @returns
     */
    async updateShipping(coordinatorId: string, shippingId: string, payload: UpdateShippingRequest): Promise<void> {
      try {
        const params: V1CoordinatorsCoordinatorIdShippingsShippingIdPatchRequest = {
          coordinatorId,
          shippingId,
          updateShippingRequest: payload,
        }
        await this.shippingApi().v1CoordinatorsCoordinatorIdShippingsShippingIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、入力内容に誤りがあります。',
          404: '対象のコーディネーターが見つかりません。',
        })
      }
    },

    /**
     * 指定したコーディネーターの配送設定を変更する非同期関数
     * @param coordinatorId コーディネーターID
     * @param payload
     * @returns
     */
    async upsertShipping(coordinatorId: string, payload: UpsertShippingRequest): Promise<void> {
      try {
        const params: V1CoordinatorsCoordinatorIdShippingsPatchRequest = {
          coordinatorId,
          upsertShippingRequest: payload,
        }
        await this.shippingApi().v1CoordinatorsCoordinatorIdShippingsPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、入力内容に誤りがあります。',
          404: '対象のコーディネーターが見つかりません。',
        })
      }
    },
  },
})
