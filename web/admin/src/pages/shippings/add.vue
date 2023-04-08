<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useCommonStore, useShippingStore } from '~/store'
import { CreateShippingRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const router = useRouter()

const formData = reactive<CreateShippingRequest>({
  name: '',
  box60Rates: [
    {
      name: '',
      price: 0,
      prefectures: []
    }
  ],
  box60Refrigerated: 0,
  box60Frozen: 0,
  box80Rates: [
    {
      name: '',
      price: 0,
      prefectures: []
    }
  ],
  box80Refrigerated: 0,
  box80Frozen: 0,
  box100Rates: [
    {
      name: '',
      price: 0,
      prefectures: []
    }
  ],
  box100Refrigerated: 0,
  box100Frozen: 0,
  hasFreeShipping: false,
  freeShippingRates: 0
})

const addBox60RateItem = () => {
  formData.box60Rates.push({
    name: '',
    price: 0,
    prefectures: []
  })
}

const addBox80RateItem = () => {
  formData.box80Rates.push({
    name: '',
    price: 0,
    prefectures: []
  })
}

const addBox100RateItem = () => {
  formData.box100Rates.push({
    name: '',
    price: 0,
    prefectures: []
  })
}

const handleClickRemoveItemButton = () => {
  // FIXME: template側で指定されていたため、定義のみ行う
}

const { createShipping } = useShippingStore()

const { alertType, isShow, alertText, show } = useAlert('error')

const { addSnackbar } = useCommonStore()

const handleSubmit = async (): Promise<void> => {
  try {
    await createShipping(formData)
    addSnackbar({
      color: 'info',
      message: `${formData.name}を登録しました。`
    })
    router.push('/shippings')
  } catch (error) {
    if (error instanceof ApiBaseError) {
      show(error.message)
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    }
  }
}

const handleClickCloseButton = (rate: '60' | '80' | '100', index: number) => {
  switch (rate) {
    case '60':
      formData.box60Rates.splice(index, 1)
      break
    case '80':
      formData.box80Rates.splice(index, 1)
      break
    case '100':
      formData.box100Rates.splice(index, 1)
      break
  }
}
</script>

<template>
  <div>
    <v-card-title>配送情報登録</v-card-title>

    <v-alert v-model="isShow" :type="alertType" v-text="alertText" />

    <the-shipping-form
      v-model="formData"
      @click:addBox60RateItem="addBox60RateItem"
      @click:addBox80RateItem="addBox80RateItem"
      @click:addBox100RateItem="addBox100RateItem"
      @click:removeItemButton="handleClickRemoveItemButton"
      @submit="handleSubmit"
    />
  </div>
</template>
