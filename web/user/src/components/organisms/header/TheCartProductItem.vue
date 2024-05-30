<script setup lang="ts">
interface Props {
  id: string
  name: string
  price: number
  imgSrc: string
  quantity: number
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:removeButton', id: string): void
}

const emits = defineEmits<Emits>()

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(props.price)
})

const handleRemoveButton = () => {
  emits('click:removeButton', props.id)
}
</script>

<template>
  <div>
    <p>{{ name }}</p>

    <p
      class="my-4 font-bold after:ml-2 after:text-[16px] after:content-['(税込)']"
    >
      {{ priceString }}
    </p>

    <div class="flex gap-x-3 text-sm">
      <nuxt-img
        provider="cloudFront"
        :src="imgSrc"
        :alt="name"
        width="72px"
        height="72px"
        class="aspect-square h-[72px] w-[72px]"
      />
      <div class="flex grow flex-col">
        <div class="flex grow items-start">
          <div class="inline-flex">
            <div class="mr-2 block whitespace-nowrap">
              数量
              {{ quantity }}
            </div>
          </div>
        </div>
        <div class="text-right">
          <button
            type="button"
            class="underline"
            @click="handleRemoveButton"
          >
            削除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
