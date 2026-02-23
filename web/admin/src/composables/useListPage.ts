import type { AlertType } from '~/lib/hooks'
import { useAlert, usePagination } from '~/lib/hooks'

interface UseListPageOptions {
  key: string
  fetchFn: (limit: number, offset: number) => Promise<void>
  deleteFn?: (id: string) => Promise<void>
  deleteSuccessMessage?: string
}

export function useListPage(options: UseListPageOptions) {
  const commonStore = useCommonStore()
  const pagination = usePagination()
  const { alertType, isShow, alertText, show } = useAlert('error')
  const loading = ref<boolean>(false)

  const fetchData = async (): Promise<void> => {
    try {
      await options.fetchFn(pagination.itemsPerPage.value, pagination.offset.value)
    }
    catch (err) {
      if (err instanceof Error) {
        show(err.message)
      }
      console.log(err)
    }
  }

  const fetchState = useAsyncData(options.key, async (): Promise<void> => {
    await fetchData()
  })

  watch(pagination.itemsPerPage, (): void => {
    fetchState.refresh()
  })

  const isLoading = (): boolean => {
    return fetchState?.pending?.value || loading.value
  }

  const handleUpdatePage = async (page: number): Promise<void> => {
    pagination.updateCurrentPage(page)
    await fetchData()
  }

  const handleClickDelete = async (id: string): Promise<void> => {
    if (!options.deleteFn) {
      return
    }
    try {
      loading.value = true
      await options.deleteFn(id)
      commonStore.addSnackbar({
        color: 'info',
        message: options.deleteSuccessMessage || '削除しました。',
      })
      fetchState.refresh()
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

  const execute = async (): Promise<void> => {
    try {
      await fetchState.execute()
    }
    catch (err) {
      console.log('failed to setup', err)
    }
  }

  return {
    pagination,
    alertType: alertType as AlertType,
    isShow,
    alertText,
    show,
    loading,
    isLoading,
    fetchState,
    handleUpdatePage,
    handleClickDelete,
    execute,
  }
}
