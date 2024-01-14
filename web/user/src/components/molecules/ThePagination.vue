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
      <button @click="handleClickPreviosPageButton">
        <the-left-arrow-icon class="h-3" />
      </button>
      <button
        v-for="page in pageArray"
        :key="page"
        :class="{
          'h-8 w-8 rounded-full p-1': true,
          'bg-main text-white': page === currentPage,
        }"
        @click="handleClickPage(page)"
      >
        {{ page }}
      </button>
      <button @click="handleClickNextPageButton">
        <the-right-arrow-icon class="h-3" />
      </button>
    </div>
  </div>
</template>
