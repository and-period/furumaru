<script setup lang='ts'>
interface Props {
  name: string
  price: number
  imgSrc: string
}

const props = defineProps<Props>()

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(props.price)
})

</script>

<template>
  <div>
    <p>{{ name }}</p>

    <p class="font-bold after:ml-2 my-4 after:content-['(税込)'] after:text-[16px]">
      {{ priceString }}
    </p>

    <div class="flex gap-x-3 text-sm">
      <img :src="imgSrc" :alt="name" class="aspect-square w-[72px] h-[72px]">
      <div class="flex flex-col">
        <div class="flex items-start grow">
          <div class="inline-flex">
            <label class="block mr-2 whitespace-nowrap">数量</label>
            <select class="border-main border-[1px] px-2 h-full ">
              <option value="0">
                0
              </option>
            </select>
          </div>
        </div>
        <div class="text-right">
          <a href="#" class="underline">削除</a>
        </div>
      </div>
    </div>
  </div>
</template>
