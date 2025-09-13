import type {
  Contact,
  ContactResponse,
  UpdateContactRequest,
  V1ContactsContactIdGetRequest,
  V1ContactsContactIdPatchRequest,
  V1ContactsGetRequest,
} from '~/types/api/v1'

export const useContactStore = defineStore('contact', {
  state: () => ({
    contact: {} as Contact,
    contacts: [] as Contact[],
    total: 0,
  }),

  actions: {
    /**
     * お問い合わせの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchContacts(limit = 20, offset = 0, orders: string[] = []): Promise<void> {
      try {
        const params: V1ContactsGetRequest = {
          limit,
          offset,
        }
        const res = await this.contactApi().v1ContactsGet(params)
        this.contacts = res.contacts
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * お問い合わせを取得する非同期関数
     * @param contactId お問い合わせID
     */
    async getContact(contactId: string): Promise<ContactResponse> {
      try {
        const params: V1ContactsContactIdGetRequest = {
          contactId,
        }
        const res = await this.contactApi().v1ContactsContactIdGet(params)
        this.contact = res.contact
        return res
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のお問い合わせが存在しません' })
      }
    },

    async updateContact(contactId: string, payload: UpdateContactRequest): Promise<void> {
      try {
        const params: V1ContactsContactIdPatchRequest = {
          contactId,
          updateContactRequest: payload,
        }
        await this.contactApi().v1ContactsContactIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: '対象のお問い合わせが存在しません',
        })
      }
    },
  },
})
