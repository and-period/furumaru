<script setup lang="ts">
export interface CheckoutStep {
  label: string
}

interface Props {
  currentStep: number
  steps?: CheckoutStep[]
}

const props = withDefaults(defineProps<Props>(), {
  steps: () => [
    { label: 'カート確認' },
    { label: '配送先' },
    { label: 'お支払い' },
    { label: '完了' },
  ],
})

const isCompleted = (index: number): boolean => {
  return index + 1 < props.currentStep
}

const isCurrent = (index: number): boolean => {
  return index + 1 === props.currentStep
}
</script>

<template>
  <nav
    aria-label="購入手続きの進捗"
    class="mx-auto my-6 max-w-2xl px-4"
  >
    <ol class="flex items-center justify-between">
      <li
        v-for="(step, index) in props.steps"
        :key="index"
        class="flex flex-1 items-center"
        :aria-current="isCurrent(index) ? 'step' : undefined"
      >
        <!-- Step circle + label -->
        <div class="flex flex-col items-center">
          <!-- Circle -->
          <div
            class="flex h-8 w-8 items-center justify-center rounded-full text-[12px] font-bold md:h-10 md:w-10 md:text-[14px]"
            :class="{
              'bg-green-500 text-white': isCompleted(index),
              'bg-orange text-white': isCurrent(index),
              'bg-gray-200 text-gray-500': !isCompleted(index) && !isCurrent(index),
            }"
          >
            <template v-if="isCompleted(index)">
              <!-- Checkmark icon -->
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="2.5"
                stroke="currentColor"
                class="h-4 w-4 md:h-5 md:w-5"
                aria-hidden="true"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M4.5 12.75l6 6 9-13.5"
                />
              </svg>
            </template>
            <template v-else>
              {{ index + 1 }}
            </template>
          </div>
          <!-- Label -->
          <span
            class="mt-1 text-center text-[10px] tracking-[0.5px] md:text-[12px] md:tracking-[1px]"
            :class="{
              'font-bold text-green-600': isCompleted(index),
              'font-bold text-orange': isCurrent(index),
              'text-gray-400': !isCompleted(index) && !isCurrent(index),
            }"
          >
            {{ step.label }}
          </span>
        </div>

        <!-- Connector line -->
        <div
          v-if="index < props.steps.length - 1"
          class="mx-1 mb-5 h-[2px] flex-1 md:mx-2"
          :class="{
            'bg-green-500': isCompleted(index),
            'bg-gray-200': !isCompleted(index),
          }"
        />
      </li>
    </ol>
  </nav>
</template>
