<script setup lang="ts">
import type { GoogleMapSearchResult } from '~/types/store'

interface Props {
  results: GoogleMapSearchResult[]
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:result', result: GoogleMapSearchResult): void
  (e: 'clear'): void
  (e: 'submit'): void
}

const modelValue = defineModel<string>({ required: true })

const hasResults = computed(() => props.results.length > 0)

const emits = defineEmits<Emits>()

const handleClickResult = (result: GoogleMapSearchResult) => {
  emits('click:result', result)
}

const handleClear = () => {
  emits('clear')
}

const handleSubmit = () => {
  emits('submit')
}
</script>

<template>
  <div>
    <form
      class="w-full rounded-full bg-white px-4 h-10 flex items-center shadow-md"
      :class="{ 'rounded-t-full rounded-b-none': hasResults }"
      @submit.prevent="handleSubmit"
    >
      <input
        v-model="modelValue"
        type="text"
        placeholder="検索"
        class="w-full py-1 px-2 focus:outline-none text-[16px]"
      >
      <button type="submit">
        <the-search-icon class="h-4" />
      </button>
    </form>
    <div
      v-if="hasResults"
    >
      <div
        v-for="result, i in results"
        :key="i"
        class="px-4 p-2 bg-white hover:bg-base"
        @click="handleClickResult(result)"
      >
        <div class=" font-semibold text-[12px]">
          {{ result.formattedAddress }}
        </div>
      </div>
      <button
        class="w-full p-2 rounded-b-full bg-white text-xs"
        @click="handleClear"
      >
        閉じる
      </button>
    </div>
  </div>
</template>
