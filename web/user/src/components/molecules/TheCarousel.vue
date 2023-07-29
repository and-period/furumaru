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
      leftContent:
        activeIdx.value === 0
          ? i === props.images.length - 1
          : i === activeIdx.value - 1,
      rightContent:
        activeIdx.value === props.images.length - 1
          ? i === 0
          : i === activeIdx.value + 1,
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
  <div class="flex items-center justify-center gap-x-4">
    <the-icon-button
      class="z-50 h-10 w-10 bg-white bg-opacity-70 hover:bg-opacity-100"
      @click="handleClickLeftArrowButton"
    >
      <the-left-arrow-icon />
    </the-icon-button>

    <div
      class="relative h-[128px] w-[312px] sm:h-[256px] sm:w-[624px] lg:h-[320px] lg:w-[780px]"
    >
      <div
        v-for="(imageItem, i) in imageItems"
        :key="i"
        :class="{
          'absolute h-full w-full transition-all duration-300': true,
          'z-40': imageItem.active,
          'z-0 brightness-75': !imageItem.active,
          '-translate-x-[312px] sm:-translate-x-[624px] lg:-translate-x-[780px]':
            imageItem.leftContent,
          'translate-x-[312px] sm:translate-x-[624px] lg:translate-x-[780px]':
            imageItem.rightContent,
        }"
        :style="`background-position: center; background-size: cover; background-image: url(${imageItem.src});`"
      />
    </div>

    <the-icon-button
      class="z-50 h-10 w-10 bg-white bg-opacity-50 hover:bg-opacity-100"
      @click="handleClickRightArrowButton"
    >
      <the-right-arrow-icon />
    </the-icon-button>
  </div>
</template>
