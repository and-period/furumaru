<script lang="ts" setup>
import type { VDataTable } from 'vuetify/lib/components/index.mjs'
import { storeToRefs } from 'pinia'

import { useAlert, usePagination } from '~/lib/hooks'
import { useContactStore } from '~/store'

const router = useRouter()
const contactStore = useContactStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { contacts, total } = storeToRefs(contactStore)

const loading = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])

watch(pagination.itemsPerPage, (): void => {
  fetchState.refresh()
})
watch(sortBy, (): void => {
  fetchState.refresh()
})

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchContacts()
})

const fetchContacts = async (): Promise<void> => {
  try {
    const orders: string[] = sortBy.value.map((item) => {
      switch (item.order) {
        case 'asc':
          return item.key
        case 'desc':
          return `-${item.key}`
        default:
          return item.order ? item.key : `-${item.key}`
      }
    }) || []

    await contactStore.fetchContacts(pagination.itemsPerPage.value, pagination.offset.value, orders)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleUpdatePage = async (page: number): Promise<void> => {
  pagination.updateCurrentPage(page)
  await fetchContacts()
}

const handleClickRow = (contactId: string): void => {
  router.push(`/contacts/${contactId}`)
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-contact-list
    v-model:sort-by="sortBy"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :contacts="contacts"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    @click:row="handleClickRow"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
  />
</template>
