<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useCoordinatorStore } from '~/store/coordinator'

const route = useRoute()

const coordinatorStore = useCoordinatorStore()

const { fetchCoordinator } = coordinatorStore
const { coordinator } = storeToRefs(coordinatorStore)

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
    <div class="mx-auto w-[1216px]">
      <img
        class="h-[320px] w-[1216px] object-cover"
      />
      <div class="flex w-full grid-cols-5 flex-col">
        <div class="col-span-2">
          <img
            src="https://furumaru-stg.s3.amazonaws.com/coordinators/header/3QiqwSrso7QpA6WVDKJFGz.png"
            class="block aspect-square w-[168px] rounded-full border-2 border-white"
          />
        </div>
        <div class="col-span-3">test</div>
      </div>
    </div>
  </div>
</template>
