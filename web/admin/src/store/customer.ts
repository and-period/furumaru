import { useApiClient } from '~/composables/useApiClient'
import { useAddressStore } from '~/store'
import { UserApi } from '~/types/api/v1'
import type { User, UserOrder, UserToList, V1UsersGetRequest, V1UsersUserIdDeleteRequest, V1UsersUserIdGetRequest, V1UsersUserIdOrdersGetRequest } from '~/types/api/v1'

export const useCustomerStore = defineStore('customer', () => {
  const { create, errorHandler } = useApiClient()
  const userApi = () => create(UserApi)

  const customer = ref<User>({} as User)
  const customers = ref<User[]>([])
  const customersToList = ref<UserToList[]>([])
  const orders = ref<UserOrder[]>([])
  const totalItems = ref<number>(0)
  const totalOrderCount = ref<number>(0)
  const totalPaymentCount = ref<number>(0)
  const totalProductAmount = ref<number>(0)
  const totalPaymentAmount = ref<number>(0)

  async function fetchCustomers(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1UsersGetRequest = { limit, offset }
      const res = await userApi().v1UsersGet(params)
      customersToList.value = res.users
      totalItems.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function getCustomer(customerId: string): Promise<void> {
    try {
      const params: V1UsersUserIdGetRequest = { userId: customerId }
      const res = await userApi().v1UsersUserIdGet(params)
      const addressStore = useAddressStore()
      customer.value = res.user
      addressStore.address = res.address
    }
    catch (err) {
      return errorHandler(err, {
        403: '購入者を閲覧する権限がありません',
        404: '対象の購入者が存在しません',
      })
    }
  }

  async function fetchCustomerOrders(customerId: string, limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1UsersUserIdOrdersGetRequest = { userId: customerId, limit, offset }
      const res = await userApi().v1UsersUserIdOrdersGet(params)
      orders.value = res.orders
      totalOrderCount.value = res.orderTotalCount
      totalPaymentCount.value = res.paymentTotalCount
      totalProductAmount.value = res.productTotalAmount
      totalPaymentAmount.value = res.paymentTotalAmount
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function deleteCustomer(customerId: string): Promise<void> {
    try {
      const params: V1UsersUserIdDeleteRequest = { userId: customerId }
      await userApi().v1UsersUserIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '購入者を削除する権限がありません',
        404: '対象の購入者が存在しません',
      })
    }
  }

  return {
    customer,
    customers,
    customersToList,
    orders,
    totalItems,
    totalOrderCount,
    totalPaymentCount,
    totalProductAmount,
    totalPaymentAmount,
    fetchCustomers,
    getCustomer,
    fetchCustomerOrders,
    deleteCustomer,
  }
})
