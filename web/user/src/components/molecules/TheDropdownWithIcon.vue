<script lang="ts" setup>
interface Expose {
  open: () => void
  close: () => void
}

const isShow = ref<boolean>(false)

const handleIconClick = () => {
  isShow.value = !isShow.value
}

const handleCloseIconClick = () => {
  isShow.value = false
}

const dropdownArea = ref<HTMLElement | null>(null)

// isShowがtrueになってから0.3秒経過しているかのフラグ
let isShowFlag = false
watch(isShow, (newValue) => {
  if (newValue === true) {
    setTimeout(() => {
      isShowFlag = true
    }, 300)
  }
  else {
    isShowFlag = false
  }
})

const clickOutside = (e: MouseEvent) => {
  if (e.target instanceof Node && !dropdownArea.value?.contains(e.target)) {
    if (isShow.value === true) {
      // isShowがtrueになってから0.3秒経過していない場合は、クリックイベントを無視する
      if (!isShowFlag) {
        return
      }
      isShow.value = false
    }
  }
}

const handleOpen = () => {
  isShow.value = true
}

onMounted(() => {
  addEventListener('click', clickOutside)
})

onBeforeUnmount(() => {
  removeEventListener('click', clickOutside)
})

defineExpose<Expose>({
  open: handleOpen,
  close: handleCloseIconClick,
})
</script>

<template>
  <div
    ref="dropdownArea"
    class="relative"
  >
    <the-icon-button @click="handleIconClick">
      <slot name="icon" />
    </the-icon-button>
    <the-dropdown-area
      v-show="isShow"
      class="absolute right-0 min-w-[240px] before:absolute before:-right-7 before:-top-8 before:-translate-x-1/2 before:border-[24px] before:border-transparent before:border-b-white before:content-[''] md:-right-8 md:before:right-1"
    >
      <div class="px-4 text-right">
        <the-icon-button @click="handleCloseIconClick">
          <the-close-icon />
        </the-icon-button>
      </div>
      <slot name="content" />
    </the-dropdown-area>
  </div>
</template>
