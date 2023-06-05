<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useContactStore } from '~/store'
import {
  ContactPriority,
  ContactResponse,
  ContactStatus,
  UpdateContactRequest
} from '~/types/api'

const route = useRoute()
const router = useRouter()
const contactStore = useContactStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const contactId = route.params.id as string

const formData = reactive<UpdateContactRequest>({
  priority: ContactPriority.UNKNOWN,
  status: ContactStatus.UNKNOWN,
  note: ''
})

const contact = computed<ContactResponse>(() => {
  return contactStore.contact
})

const fetchState = useAsyncData(async () => {
  try {
    const res = await contactStore.getContact(contactId)

    formData.priority = res.priority
    formData.status = res.status
    formData.note = res.note
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const handleSubmit = async (): Promise<void> => {
  try {
    await contactStore.updateContact(formData, contactId)
    router.push('/contacts')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-contact-edit
    v-model:form-data="formData"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :contact="contact"
    @submit="handleSubmit"
  />
</template>
