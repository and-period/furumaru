<script lang="ts" setup>
interface Props {
  images: string[]
}

const props = defineProps<Props>()

const activeIdx = ref<number>(0)

const imageItems = computed(() => {
  return props.images.map((item, i) => {
    return {
      src: item,
      active: i === activeIdx.value,
      leftContent: activeIdx.value === 0 ? i === props.images.length - 1 : i === activeIdx.value - 1,
      rightContent: activeIdx.value === props.images.length - 1 ? i === 0 : i === activeIdx.value + 1
    }
  })
})

const handleClickLeftArrowButton = () => {
  if (activeIdx.value === 0) {
    activeIdx.value = props.images.length - 1
  } else {
    activeIdx.value = activeIdx.value - 1
  }
}

const handleClickRightArrowButton = () => {
  if (activeIdx.value === props.images.length - 1) {
    activeIdx.value = 0
  } else {
    activeIdx.value = activeIdx.value + 1
  }
}
</script>

<template>
  <div class="flex justify-center items-center gap-x-4">
    <the-icon-button class="bg-white w-10 h-10 z-50 bg-opacity-70 hover:bg-opacity-100" @click="handleClickLeftArrowButton">
      <the-left-arrow-icon />
    </the-icon-button>

    <div class="relative 2xl:w-[1170px] 2xl:h-[480px] lg:w-[780px] lg:h-[320px] sm:w-[624px] sm:h-[256px] w-[312px] h-[128px]">
      <div
        v-for="imageItem, i in imageItems"
        :key="i"
        :class="{
          'absolute w-full h-full transition-all duration-300': true,
          'z-40': imageItem.active,
          'z-0 brightness-75': !imageItem.active,
          '2xl:-translate-x-[1170px] lg:-translate-x-[780px] sm:-translate-x-[624px] -translate-x-[312px]': imageItem.leftContent,
          '2xl:translate-x-[1170px] lg:translate-x-[780px] sm:translate-x-[624px] translate-x-[312px]': imageItem.rightContent,
        }"
        :style="`background-position: center; background-size: cover; background-image: url(${imageItem.src});`"
      />
    </div>

    <the-icon-button class="bg-white w-10 h-10 z-50 bg-opacity-50 hover:bg-opacity-100" @click="handleClickRightArrowButton">
      <the-right-arrow-icon />
    </the-icon-button>
  </div>
</template>
