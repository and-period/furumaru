<script lang="ts" setup>
import { CalendarEvent } from '~/types/props'

const props = defineProps({
  events: {
    type: Array<CalendarEvent>,
    default: () => []
  }
})

const calendarValue = ref<string>('')
const type = ref<string>('month')
const calendarRef = ref<{
  title: string
  prev:() => {}
  next: () => {}
    } | null>(null)
const typeItems = [
  {
    text: '日',
    value: 'day'
  },
  {
    text: '週',
    value: 'week'
  },
  {
    text: '月',
    value: 'month'
  }
]

const calendarTitle = computed(() => {
  if (calendarRef && calendarRef.value) {
    return calendarRef.value?.title
  } else {
    return ''
  }
})

const handleClickToDayButton = () => {
  calendarValue.value = ''
}

const handleClickPrevButton = () => {
  if (calendarRef && calendarRef.value) {
    calendarRef.value?.prev()
  }
}

const handleClickNextButton = () => {
  if (calendarRef && calendarRef.value) {
    calendarRef.value?.next()
  }
}
</script>

<template>
  <div>
    <v-toolbar flat>
      <v-toolbar-title>
        {{ calendarTitle }}
      </v-toolbar-title>
      <v-select
        v-model="type"
        outlined
        :items="typeItems"
        dense
        hide-details
        label="表示形式"
        class="ml-4"
      />
      <v-btn icon class="ma-2" @click="handleClickPrevButton">
        <v-icon>mdi-chevron-left</v-icon>
      </v-btn>
      <v-btn icon class="ma-2" @click="handleClickNextButton">
        <v-icon>mdi-chevron-right</v-icon>
      </v-btn>

      <v-spacer />
      <v-btn
        outlined
        color="primary"
        class="mx-2"
        @click="handleClickToDayButton"
      >
        本日
      </v-btn>
    </v-toolbar>
    <v-sheet height="80vh">
      <v-calendar
        ref="calendarRef"
        v-model="calendarValue"
        :type="type"
        :events="props.events"
      />
    </v-sheet>
  </div>
</template>
