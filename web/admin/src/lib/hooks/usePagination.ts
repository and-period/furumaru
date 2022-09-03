import { computed, reactive, ref } from '@vue/composition-api'

export function usePagination() {
  const currentPage = ref<number>(1)
  const itemsPerPage = ref<number>(10)
  const options = reactive({
    itemsPerPageOptions: [10, 20, 30, 50],
  })

  const handleUpdateItemsPerPage = (n: number) => {
    itemsPerPage.value = n
  }

  const calcTotalPages = (totalItem: number): number => {
    return Math.ceil(totalItem / itemsPerPage.value)
  }

  const updateCurrentPage = (page: number) => {
    currentPage.value = page
  }

  const offset = computed(() => {
    if (currentPage.value === 1) {
      return 0
    } else {
      return itemsPerPage.value * (currentPage.value - 1) + 1
    }
  })

  return {
    currentPage,
    updateCurrentPage,
    itemsPerPage,
    handleUpdateItemsPerPage,
    calcTotalPages,
    offset,
    options,
  }
}
