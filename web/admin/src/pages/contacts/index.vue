<script lang="ts" setup>
import { VDataTable } from 'vuetify/lib/labs/components'

import { useAlert, usePagination } from '~/lib/hooks'
import { useContactStore } from '~/store'

const router = useRouter()
const contactStore = useContactStore()
const { itemsPerPage, offset, updateCurrentPage, handleUpdateItemsPerPage } = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const sortBy = reactive<VDataTable['sortBy']>([])

const contacts = computed(() => {
  return contactStore.contacts
})
const contactTotal = computed(() => {
  return contactStore.total
})

watch(itemsPerPage, () => {
  fetchState.refresh()
})
watch(sortBy, () => {
  fetchState.refresh()
})

const fetchState = useAsyncData(async () => {
  await fetchContacts()
})

const fetchContacts = async () => {
  try {
    const orders: string[] = sortBy?.map((item) => {
      switch (item.order) {
        case 'asc':
          return item.key
        case 'desc':
          return `-${item.key}`
        default:
          return item.order ? item.key : `-${item.key}`
      }
    }) || []

    await contactStore.fetchContacts(itemsPerPage.value, offset.value, orders)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUpdatePage = async (page: number) => {
  updateCurrentPage(page)
  await fetchContacts()
}

const handleClickSortBy = (item: VDataTable['sortBy']): void => {
  sortBy.splice(0, sortBy.length, ...item)
}

const handleEditCategory = (contactId: string) => {
  router.push(`/contacts/edit/${contactId}`)
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-contact-list
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :contacts="contacts"
    :sort-by="sortBy"
    :table-items-per-page="itemsPerPage"
    :table-items-total="contactTotal"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="handleUpdateItemsPerPage"
    @click:edit="handleEditCategory"
    @update:sort-by="handleClickSortBy"
  />
</template>
