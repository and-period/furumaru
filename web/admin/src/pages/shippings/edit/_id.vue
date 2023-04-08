<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useCommonStore } from '~/store/common'
import { useShippingStore } from '~/store/shippings'
import { UpdateShippingRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const route = useRoute()
const router = useRouter()
const id = route.params.id
const { alertType, isShow, alertText, show } = useAlert('error')

const { addSnackbar } = useCommonStore()

const formData = ref<UpdateShippingRequest>({
  name: '',
  box60Rates: [
    {
      name: '',
      price: 0,
      prefectures: [],
    },
  ],
  box60Refrigerated: 0,
  box60Frozen: 0,
  box80Rates: [
    {
      name: '',
      price: 0,
      prefectures: [],
    },
  ],
  box80Refrigerated: 0,
  box80Frozen: 0,
  box100Rates: [
    {
      name: '',
      price: 0,
      prefectures: [],
    },
  ],
  box100Refrigerated: 0,
  box100Frozen: 0,
  hasFreeShipping: false,
  freeShippingRates: 0,
})

const { getShipping, updateShipping } = useShippingStore()

const fetchState = useAsyncData(async () => {
  try {
    const shipping = await getShipping(id)
    formData.value = { ...shipping }
  } catch (error) {
    if (error instanceof ApiBaseError) {
      show(error.message)
    }
  }
})

const addBox60RateItem = () => {
  formData.value.box60Rates.push({
    name: '',
    price: 0,
    prefectures: [],
  })
}

const addBox80RateItem = () => {
  formData.value.box80Rates.push({
    name: '',
    price: 0,
    prefectures: [],
  })
}

const addBox100RateItem = () => {
  formData.value.box100Rates.push({
    name: '',
    price: 0,
    prefectures: [],
  })
}

const handleClickRemoveItemButton = (
  rate: '60' | '80' | '100',
  index: number
) => {
  switch (rate) {
    case '60':
      formData.value.box60Rates.splice(index, 1)
      break
    case '80':
      formData.value.box80Rates.splice(index, 1)
      break
    case '100':
      formData.value.box100Rates.splice(index, 1)
      break
  }
}

const handleSubmit = async () => {
  try {
    await updateShipping(id, formData.value)
    addSnackbar({
      color: 'info',
      message: `${formData.value.name}を更新しました。`,
    })
    router.push('/shippings')
  } catch (error) {
    if (error instanceof ApiBaseError) {
      show(error.message)
      window.scrollTo({
        top: 0,
        behavior: 'smooth',
      })
    }
  }
}
</script>

<template>
  <div>
    <v-card-title>配送設定編集</v-card-title>

    <v-alert v-model="isShow" :type="alertType" v-text="alertText" />

    <the-shipping-form
      v-model="formData"
      :loading="fetchState.pending"
      @click:addBox60RateItem="addBox60RateItem"
      @click:addBox80RateItem="addBox80RateItem"
      @click:addBox100RateItem="addBox100RateItem"
      @click:removeItemButton="handleClickRemoveItemButton"
      @submit="handleSubmit"
    />
  </div>
</template>
