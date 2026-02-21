<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { getErrorMessage } from '~/lib/validations'
import type { DateTimeInput } from '~/types/props'
import { TimeDataValidationRules } from '~/types/validations'

const props = defineProps({
  startTime: {
    type: Object as PropType<DateTimeInput>,
    required: true,
  },
  endTime: {
    type: Object as PropType<DateTimeInput>,
    required: true,
  },
})

const emit = defineEmits<{
  (e: 'update:start-time', item: DateTimeInput): void
  (e: 'update:end-time', item: DateTimeInput): void
  (e: 'change:start-at'): void
  (e: 'change:end-at'): void
}>()

const startTimeDataValue = computed({
  get: (): DateTimeInput => props.startTime,
  set: (item: DateTimeInput) => emit('update:start-time', item),
})
const endTimeDataValue = computed({
  get: (): DateTimeInput => props.endTime,
  set: (item: DateTimeInput) => emit('update:end-time', item),
})

const startTimeDataValidate = useVuelidate(TimeDataValidationRules, startTimeDataValue)
const endTimeDataValidate = useVuelidate(TimeDataValidationRules, endTimeDataValue)

const onChangeStartAt = async (): Promise<void> => {
  await startTimeDataValidate.value.$validate()
  emit('change:start-at')
}

const onChangeEndAt = async (): Promise<void> => {
  await endTimeDataValidate.value.$validate()
  emit('change:end-at')
}
</script>

<template>
  <div class="d-flex flex-column ga-2">
    <v-label class="text-body-2 font-weight-medium">
      使用期間 *
    </v-label>

    <div class="d-flex flex-column flex-md-row align-center ga-2">
      <!-- 開始日時 -->
      <div class="d-flex align-center ga-2 flex-grow-1">
        <v-text-field
          v-model="startTimeDataValidate.date.$model"
          :error-messages="getErrorMessage(startTimeDataValidate.date.$errors)"
          type="date"
          variant="outlined"
          density="compact"
          hide-details="auto"
          @update:model-value="onChangeStartAt"
        />
        <v-text-field
          v-model="startTimeDataValidate.time.$model"
          :error-messages="getErrorMessage(startTimeDataValidate.time.$errors)"
          type="time"
          variant="outlined"
          density="compact"
          hide-details="auto"
          @update:model-value="onChangeStartAt"
        />
      </div>

      <!-- 区切り -->
      <div class="text-body-2 text-grey px-2">
        〜
      </div>

      <!-- 終了日時 -->
      <div class="d-flex align-center ga-2 flex-grow-1">
        <v-text-field
          v-model="endTimeDataValidate.date.$model"
          :error-messages="getErrorMessage(endTimeDataValidate.date.$errors)"
          type="date"
          variant="outlined"
          density="compact"
          hide-details="auto"
          @update:model-value="onChangeEndAt"
        />
        <v-text-field
          v-model="endTimeDataValidate.time.$model"
          :error-messages="getErrorMessage(endTimeDataValidate.time.$errors)"
          type="time"
          variant="outlined"
          density="compact"
          hide-details="auto"
          @update:model-value="onChangeEndAt"
        />
      </div>
    </div>
  </div>
</template>
