<script setup lang="ts">
import type { Experience } from '~/types/api'
import { priceFormatter } from '~/lib/price'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const dt = (str: keyof I18n['experiences']['purchase']) => {
  return i18n.t(`experiences.purchase.${str}`)
}

interface Props {
  experience: Experience
  adultCount: number
  juniorHighSchoolCount: number
  elementarySchoolCount: number
  preschoolCount: number
  seniorCount: number
}

const props = defineProps<Props>()

// 購入対象の種別情報
const targetTypes = computed(() => {
  const targets = []
  if (props.adultCount > 0) {
    targets.push({
      type: 'adult',
      count: props.adultCount,
      price: props.experience.priceAdult,
    })
  }
  if (props.juniorHighSchoolCount > 0) {
    targets.push({
      type: 'juniorHighSchool',
      count: props.juniorHighSchoolCount,
      price: props.experience.priceJuniorHighSchool,
    })
  }
  if (props.elementarySchoolCount > 0) {
    targets.push({
      type: 'elementarySchool',
      count: props.elementarySchoolCount,
      price: props.experience.priceElementarySchool,
    })
  }
  if (props.preschoolCount > 0) {
    targets.push({
      type: 'preschool',
      count: props.preschoolCount,
      price: props.experience.pricePreschool,
    })
  }
  if (props.seniorCount > 0) {
    targets.push({
      type: 'senior',
      count: props.seniorCount,
      price: props.experience.priceSenior,
    })
  }

  return targets
})

const totalPrice = computed(() => {
  return (
    props.experience.priceAdult * props.adultCount
    + props.experience.priceJuniorHighSchool * props.juniorHighSchoolCount
    + props.experience.priceElementarySchool * props.elementarySchoolCount
    + props.experience.pricePreschool * props.preschoolCount
    + props.experience.priceSenior * props.seniorCount
  )
})
</script>

<template>
  <div class="bg-base p-10 text-main flex flex-col gap-4">
    <div class="text-[14px] font-bold tracking-[1.6px] md:text-[16px]">
      購入内容
    </div>
    <!-- 体験情報 -->
    <div class="flex gap-4">
      <nuxt-img
        width="64px"
        height="64px"
        provider="cloudFront"
        class="aspect-square border object-contain w-[64px] h-[64px]"
        :src="experience.thumbnailUrl"
        :alt="experience.title"
      />
      <div class="w-full flex flex-col gap-4">
        <div class=" font-semibold">
          {{ experience.title }}
        </div>
        <div
          class="text-[12px] tracking-[1.2px]"
        >
          {{ experience.description }}
        </div>
      </div>
    </div>

    <!-- 購入枚数 -->
    <div class=" divide-y border-y">
      <div
        v-for="target in targetTypes"
        :key="target.type"
        class="grid py-4 text-[12px] tracking-[1.2px] grid-cols-5"
      >
        <div class=" col-span-3">
          {{ dt(target.type) }}
        </div>
        <div class="md:text-right">
          {{ `${dt('quantityLabel')} ${target.count}` }}
        </div>
        <div class="text-right">
          {{ priceFormatter(target.price) }}
        </div>
      </div>
    </div>

    <div
      class="mt-6 grid grid-cols-2 text-[14px] font-bold tracking-[1.4px]"
    >
      <div>{{ dt("totalPriceLabel") }}</div>
      <div class="text-right">
        {{ priceFormatter(totalPrice) }}
      </div>
    </div>
  </div>
</template>
