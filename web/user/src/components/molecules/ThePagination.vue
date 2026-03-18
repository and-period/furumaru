<script setup lang="ts">
interface Props {
  currentPage: number
  pageArray: number[]
}

const props = defineProps<Props>()

interface Emits {
  (e: 'changePage', page: number): boolean
}

const emits = defineEmits<Emits>()

const handleClickPreviosPageButton = () => {
  if (props.currentPage === 1) return
  handleClickPage(props.currentPage - 1)
}

const handleClickPage = (page: number) => {
  emits('changePage', page)
}

const handleClickNextPageButton = () => {
  const lastPage = props.pageArray[props.pageArray.length - 1]
  if (props.currentPage === lastPage) return
  handleClickPage(props.currentPage + 1)
}
</script>

<template>
  <div class="text-center text-main">
    <div class="inline-flex gap-4">
      <button
        class="min-h-[44px] min-w-[44px] flex items-center justify-center"
        aria-label="前のページ"
        @click="handleClickPreviosPageButton"
      >
        <the-left-arrow-icon class="h-3" />
      </button>
      <button
        v-for="page in pageArray"
        :key="page"
        :class="{
          'min-h-[44px] min-w-[44px] rounded-full flex items-center justify-center': true,
          'bg-main text-white': page === currentPage,
        }"
        :aria-label="`${page}ページ`"
        :aria-current="page === currentPage ? 'page' : undefined"
        @click="handleClickPage(page)"
      >
        {{ page }}
      </button>
      <button
        class="min-h-[44px] min-w-[44px] flex items-center justify-center"
        aria-label="次のページ"
        @click="handleClickNextPageButton"
      >
        <the-right-arrow-icon class="h-3" />
      </button>
    </div>
  </div>
</template>
