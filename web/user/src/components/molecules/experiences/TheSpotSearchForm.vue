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
      class="w-full rounded-full bg-white h-10 relative z-10 shadow-md"
      :class="{ 'shadow-none': hasResults }"
      @submit.prevent="handleSubmit"
    >
      <div
        class="w-full flex items-center px-4 h-full"
      >
        <input
          v-model="modelValue"
          type="text"
          placeholder="検索"
          class="w-full py-1 px-2 focus:outline-none text-[16px] z-10"
        >
        <button
          type="submit"
          class="z-10 bg-white"
        >
          <the-search-icon class="h-4" />
        </button>
      </div>
      <div
        v-show="hasResults"
        class="w-full bg-white absolute h-4 bottom-0 z-0"
      />
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
