import { useApiClient } from '~/composables/useApiClient'
import { ContactApi } from '~/types/api/v1'
import type {
  Contact,
  ContactResponse,
  UpdateContactRequest,
  V1ContactsContactIdGetRequest,
  V1ContactsContactIdPatchRequest,
  V1ContactsGetRequest,
} from '~/types/api/v1'

export const useContactStore = defineStore('contact', () => {
  const { create, errorHandler } = useApiClient()
  const contactApi = () => create(ContactApi)

  const contact = ref<Contact>({} as Contact)
  const contacts = ref<Contact[]>([])
  const total = ref<number>(0)

  async function fetchContacts(limit = 20, offset = 0, orders: string[] = []): Promise<void> {
    try {
      const params: V1ContactsGetRequest = { limit, offset }
      const res = await contactApi().v1ContactsGet(params)
      contacts.value = res.contacts
      total.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function getContact(contactId: string): Promise<ContactResponse> {
    try {
      const params: V1ContactsContactIdGetRequest = { contactId }
      const res = await contactApi().v1ContactsContactIdGet(params)
      contact.value = res.contact
      return res
    }
    catch (err) {
      return errorHandler(err, { 404: '対象のお問い合わせが存在しません' })
    }
  }

  async function updateContact(contactId: string, payload: UpdateContactRequest): Promise<void> {
    try {
      const params: V1ContactsContactIdPatchRequest = {
        contactId,
        updateContactRequest: payload,
      }
      await contactApi().v1ContactsContactIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        404: '対象のお問い合わせが存在しません',
      })
    }
  }

  return {
    contact,
    contacts,
    total,
    fetchContacts,
    getContact,
    updateContact,
  }
})
