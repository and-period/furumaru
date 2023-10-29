<script setup lang="ts">
import { MOCK_ALL_PRODUCT_ITEMS, ProductItemMock } from '~/constants/mock'

const router = useRouter()
const route = useRoute()

const id = computed<string>(() => {
  const ids = route.params.id
  if (Array.isArray(ids)) {
    return ids[0]
  } else {
    return ids
  }
})

const selectItem = computed<ProductItemMock | undefined>(() => {
  return MOCK_ALL_PRODUCT_ITEMS.find((item) => {
    return item.id === id.value
  })
})

const priceString = computed<string>(() => {
  if (selectItem.value) {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(selectItem.value.price)
  } else {
    return ''
  }
})

const handleClickAddCartButton = () => {
  router.push('/purchase')
}
</script>

<template>
  <div class="grid grid-cols-2 bg-white px-[112px] pb-6 pt-[40px] text-main">
    <div class="w-full">
      <div class="h=[500px] mx-auto aspect-square w-[500px]">
        <img :src="selectItem?.imgSrc" />
      </div>
    </div>

    <div class="flex w-full flex-col gap-4">
      <div class="break-words text-[24px] tracking-[2.4px]">
        {{ selectItem?.name }}
      </div>

      <div class="mt-4 flex flex-col leading-[32px]">
        <div class="text-[16px] tracking-[1.6px]">
          生産者:
          <a href="#" class="font-bold underline">{{ selectItem?.cnName }}</a>
        </div>
        <div class="text-[14px] tracking-[1.4px]">
          {{ selectItem?.address }}
        </div>
      </div>

      <div
        class="mt-8 w-full rounded-2xl bg-base px-[20px] py-[28px] text-main"
      >
        <p class="mb-[12px] text-[14px] tracking-[1.4px]">おすすめポイント</p>
        <ol
          class="recommend-list flex flex-col divide-y divide-dashed divide-main px-[4px] pl-[24px]"
        >
          <li class="py-3">キリっとした酸味と華やかな香りが特徴です</li>
          <li class="py-3">島ならではの資材を使った土づくりをしています</li>
          <li class="py-3">防腐剤・ワックスは一切使用しておりません</li>
        </ol>
      </div>

      <div>
        <div
          class="mt-[60px] text-[32px] after:ml-2 after:text-[16px] after:content-['(税込)']"
        >
          {{ priceString }}
        </div>

        <div class="mt-8 inline-flex items-center">
          <label class="mr-2 block text-[16px]">数量</label>
          <select class="h-full border-[1px] border-main px-2" @click.stop>
            <option value="0">0</option>
          </select>
        </div>
      </div>

      <button
        class="mt-8 w-full bg-main py-4 text-center text-white"
        @click="handleClickAddCartButton"
      >
        買い物カゴに入れる
      </button>
    </div>
  </div>
</template>

<style scoped>
.recommend-list {
  list-style: none;
  counter-reset: li;
  position: relative;
}

.recommend-list li {
  padding-left: 16px;
}

.recommend-list li::before {
  content: counter(li);
  counter-increment: li;
  position: absolute;
  left: 0;
  background-color: #604c3f;
  color: #f9f6ea;
  border-radius: 100%;
  width: 24px;
  height: 24px;
  text-align: center;
}
</style>
