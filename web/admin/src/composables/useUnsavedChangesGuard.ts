import type { RouteLocationNormalized } from 'vue-router'
import type { ComputedRef, Ref } from 'vue'

interface UseUnsavedChangesGuardReturn {
  isDirty: ComputedRef<boolean>
  captureSnapshot: () => void
  markAsSaved: () => void
  showLeaveDialog: Ref<boolean>
  confirmLeave: () => void
  cancelLeave: () => void
}

export function useUnsavedChangesGuard<T>(
  formData: Ref<T>,
): UseUnsavedChangesGuardReturn {
  const router = useRouter()
  const snapshot = ref<string>('')
  const showLeaveDialog = ref<boolean>(false)
  const bypassGuard = ref<boolean>(false)
  let pendingTo: RouteLocationNormalized | null = null

  const isDirty = computed<boolean>(() => {
    return snapshot.value !== '' && JSON.stringify(formData.value) !== snapshot.value
  })

  const captureSnapshot = (): void => {
    snapshot.value = JSON.stringify(formData.value)
    bypassGuard.value = false
  }

  const markAsSaved = (): void => {
    bypassGuard.value = true
  }

  const confirmLeave = (): void => {
    showLeaveDialog.value = false
    bypassGuard.value = true
    if (pendingTo) {
      router.push(pendingTo)
      pendingTo = null
    }
  }

  const cancelLeave = (): void => {
    showLeaveDialog.value = false
    pendingTo = null
  }

  // Browser tab close / refresh guard
  const beforeUnloadHandler = (e: BeforeUnloadEvent): void => {
    if (isDirty.value && !bypassGuard.value) {
      e.preventDefault()
    }
  }

  onMounted(() => {
    window.addEventListener('beforeunload', beforeUnloadHandler)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('beforeunload', beforeUnloadHandler)
  })

  // Vue Router navigation guard
  onBeforeRouteLeave((to, _from, next) => {
    if (isDirty.value && !bypassGuard.value) {
      pendingTo = to
      showLeaveDialog.value = true
      next(false)
    }
    else {
      next()
    }
  })

  return {
    isDirty,
    captureSnapshot,
    markAsSaved,
    showLeaveDialog,
    confirmLeave,
    cancelLeave,
  }
}
