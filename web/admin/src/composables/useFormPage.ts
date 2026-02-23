import type { AlertType } from '~/lib/hooks'
import { useAlert } from '~/lib/hooks'

interface UseFormPageOptions {
  submitFn: () => Promise<void>
  successMessage: string
  redirectPath: string
}

export function useFormPage(options: UseFormPageOptions) {
  const router = useRouter()
  const commonStore = useCommonStore()
  const { alertType, isShow, alertText, show } = useAlert('error')
  const loading = ref<boolean>(false)

  const isLoading = (): boolean => {
    return loading.value
  }

  const handleSubmit = async (): Promise<void> => {
    try {
      loading.value = true
      await options.submitFn()
      commonStore.addSnackbar({
        message: options.successMessage,
        color: 'info',
      })
      router.push(options.redirectPath)
    }
    catch (err) {
      if (err instanceof Error) {
        show(err.message)
      }
      window.scrollTo({
        top: 0,
        behavior: 'smooth',
      })
      console.log(err)
    }
    finally {
      loading.value = false
    }
  }

  return {
    alertType: alertType as AlertType,
    isShow,
    alertText,
    show,
    loading,
    isLoading,
    handleSubmit,
  }
}
