<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useCoordinatorStore } from '~/store/coordinator'

const route = useRoute()

const coordinatorStore = useCoordinatorStore()

const { fetchCoordinator } = coordinatorStore

const coordinator = storeToRefs(coordinatorStore)

const id = computed<string>(() => {
  const ids = route.params.id
  if (Array.isArray(ids)) {
    return ids[0]
  } else {
    return ids
  }
})

fetchCoordinator(id.value)


</script>

<template>
  <div>
    <div class="mx-auto w-[1216px] text-main">
      <img
        class="h-[320px] w-[1216px] object-cover"
        :src="coordinator.coordnatorInfo.value.headerUrl"
      />
      <div class="grid grid-cols-5">
        <div class="col-span-2">
          <div class="flex justify-center">
            <img
              :src="coordinator.coordnatorInfo.value.thumbnailUrl"
              class="block aspect-square w-[168px] rounded-full border-2 border-white"
            />
          </div>
          <p class="mt-4 text-center text-[20px] font-bold tracking-[2.0px]">{{ coordinator.coordnatorInfo.value.marcheName }}</p>
          <div class="flex justify-center pt-2 text-[14px] tracking-[1.4px]">
            <p>{{ coordinator.coordnatorInfo.value.prefecture }}</p>
            <p class="pl-2">{{ coordinator.coordnatorInfo.value.city }}</p>
          </div>
          <div class="my-4 flex justify-center tracking-[2.4px]">
            <p class="mt-auto text-[14px]">コーディネータ</p>
            <p class="ml-2 text-[24px] font-bold">{{ coordinator.coordnatorInfo.value.username }}</p>
          </div>
          <p class="text-[16px] tracking-[1.6px]">{{ coordinator.coordnatorInfo.value.profile }}</p>
        </div>
        <div class="col-span-3">test</div>
      </div>
    </div>
  </div>
</template>
