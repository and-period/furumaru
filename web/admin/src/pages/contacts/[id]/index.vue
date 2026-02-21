<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useCommonStore, useContactStore } from '~/store'
import { ContactStatus } from '~/types/api/v1'
import type { UpdateContactRequest } from '~/types/api/v1'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const contactStore = useContactStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const contactId = route.params.id as string

const { contact } = storeToRefs(contactStore)

const loading = ref<boolean>(false)
const formData = ref<UpdateContactRequest>({
  status: ContactStatus.ContactStatusUnknown,
  note: '',
  categoryId: '',
  content: '',
  email: '',
  phoneNumber: '',
  responderId: '',
  title: '',
  userId: '',
  username: '',
})

const fetchState = useAsyncData('contact-detail', async (): Promise<void> => {
  try {
    await contactStore.getContact(contactId)
    formData.value = { ...contact.value }
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await contactStore.updateContact(contactId, formData.value)
    commonStore.addSnackbar({
      message: 'お問い合わせ情報が更新されました。',
      color: 'info',
    })
    router.push('/contacts')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-contact-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :contact="contact"
    @submit="handleSubmit"
  />
</template>
