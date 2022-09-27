<template>
  <div>
    <div class="d-flex mb-4">
      <v-spacer />
      <v-btn color="primary" outlined @click="handleClickAddButton">
        <v-icon>mdi-plus</v-icon>
        ライブマルシェ登録</v-btn
      >
    </div>
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
        >本日</v-btn
      >
    </v-toolbar>
    <v-sheet height="80vh">
      <v-calendar ref="calendarRef" v-model="calendarValue" :type="type" />
    </v-sheet>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  ref,
  useRouter,
  SetupContext,
} from '@nuxtjs/composition-api'

export default defineComponent({
  setup(_, ctx) {
    const router = useRouter()

    const calendarValue = ref<string>('')
    const type = ref<string>('month')
    const calendarRef = ref(null)
    const typeItems = [
      {
        text: '日',
        value: 'day',
      },
      {
        text: '週',
        value: 'week',
      },
      {
        text: '月',
        value: 'month',
      },
    ]

    const calendarTitle = computed(() => {
      if (calendarRef && calendarRef.value) {
        return calendarRef.value?.title
      } else {
        return ''
      }
    })

    const handleClickAddButton = () => {
      router.push('/livestreaming/add')
    }

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

    return {
      calendarValue,
      type,
      typeItems,
      calendarRef,
      calendarTitle,
      handleClickAddButton,
      handleClickToDayButton,
      handleClickPrevButton,
      handleClickNextButton,
    }
  },
})
</script>

<style scoped>
#preview {
  margin-bottom: 1.5rem;
  background: gray;
  width: 100%;
}
</style>
