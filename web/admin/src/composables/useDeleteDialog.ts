export function useDeleteDialog<T>() {
  const dialogVisible = ref(false)
  const selectedItem = ref<T | null>(null) as Ref<T | null>

  const open = (item: T) => {
    selectedItem.value = item
    dialogVisible.value = true
  }

  const close = () => {
    dialogVisible.value = false
    selectedItem.value = null
  }

  return { dialogVisible, selectedItem, open, close }
}
