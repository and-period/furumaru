<script setup lang="ts">
import { computed, ref } from 'vue';

interface Props {
	name: string;
	price: number
	addToCartButtonText?: string;
	selectLabelText?: string
}

const props = withDefaults(defineProps<Props>(), {
	addToCartButtonText: 'カゴに入れる',
	selectLabelText: '数量'
});

const quantity = ref<number>(1)

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(props.price)
})

const stokeValues = computed<number>(() => {
	return 10
})
</script>

<template>
  <div class="flex flex-col text-main w-full font-semibold gap-2">
    <p>
      {{ name }}
    </p>

    <p>
      {{ priceString }}
    </p>

    <div class="inline-flex gap-2 text-xs items-center">
      <div class="inline-flex gap-1 items-center">
        <label for="select">{{ selectLabelText }}</label>
        <select
          id="select"
          v-model="quantity"
          class="border border-main py-1 pl-1"
        >
          <template
            v-for="n in stokeValues"
            :key="n"
          >
            <option :value="n">
              {{ n }}
            </option>
          </template>
        </select>
      </div>
      <button class="bg-orange text-white py-1 px-2">
        {{ addToCartButtonText }}
      </button>
    </div>
  </div>
</template>
