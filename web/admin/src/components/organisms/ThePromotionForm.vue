<template>
  <form @submit.prevent="handleSubmit">
    <v-card elevation="0">
      <v-card-text>
        <v-text-field
          v-model="formData.title"
          label="タイトル"
          required
          maxlength="200"
        />
        <v-textarea
          v-model="formData.description"
          label="説明"
          maxlength="2000"
        />
        <div class="d-flex align-center">
          <v-text-field
            v-model="formData.code"
            class="mr-4"
            label="割引コード"
            required
          />
          <v-btn outlined small color="primary"> 自動生成 </v-btn>
          <v-spacer />
        </div>
      </v-card-text>
      <v-select class="mx-4" :items="status" label="ステータス" />
      <div class="d-flex align-center">
        <v-select class="ml-4 mr-4" :items="discountMethod" label="割引方法" />
        <v-text-field v-model="formData.Type" class="mr-4" label="割引値" />
      </div>
      <p class="text-h6 ml-4">投稿期間</p>
      <div class="d-flex align-center">
        <v-menu
          v-model="postMenu"
          :close-on-content-click="false"
          :nudge-right="40"
          transition="scale-transition"
          offset-y
          min-width="auto"
        >
          <template #activator="{ on, attrs }">
            <v-text-field
              v-model="postDate"
              label="投稿開始日"
              readonly
              outlined
              v-bind="attrs"
              v-on="on"
            ></v-text-field>
          </template>
          <v-date-picker
            v-model="postDate"
            scrollable
            @input="postMenu = false"
          >
            <v-spacer></v-spacer>
            <v-btn text color="primary" @click="postMenu = false">
              Cancel
            </v-btn>
          </v-date-picker>
        </v-menu>
        <input class="postTime" type="time" required />
        <p class="text-h6 ml-4">〜</p>
      </div>
      <p class="text-h6 ml-4">使用期間</p>
      <div class="d-flex align-center">
        <v-menu
          v-model="useStartMenu"
          :close-on-content-click="false"
          :nudge-right="40"
          transition="scale-transition"
          offset-y
          min-width="auto"
        >
          <template #activator="{ on, attrs }">
            <v-text-field
              v-model="useStartDate"
              label="使用開始日"
              readonly
              outlined
              v-bind="attrs"
              v-on="on"
            ></v-text-field>
          </template>
          <v-date-picker
            v-model="useStartDate"
            scrollable
            @input="useStartMenu = false"
          >
            <v-spacer></v-spacer>
            <v-btn text color="primary" @click="useStartMenu = false">
              Cancel
            </v-btn>
          </v-date-picker>
        </v-menu>
        <input class="startUseTime" type="time" required />
        <p class="text-h6 ml-4">〜</p>
        <v-menu
          v-model="useEndMenu"
          :close-on-content-click="false"
          :nudge-right="40"
          transition="scale-transition"
          offset-y
          min-width="auto"
        >
          <template #activator="{ on, attrs }">
            <v-text-field
              v-model="useEndDate"
              label="使用終了日"
              readonly
              outlined
              v-bind="attrs"
              v-on="on"
            ></v-text-field>
          </template>
          <v-date-picker
            v-model="useEndDate"
            scrollable
            @input="useEndMenu = false"
          >
            <v-spacer></v-spacer>
            <v-btn text color="primary" @click="endEndMenu = false">
              Cancel
            </v-btn>
          </v-date-picker>
        </v-menu>
        <input class="endUseTime" type="time" required />
      </div>
    </v-card>
    <v-btn block outlined color="primary" type="submit">
      {{ btnText }}
    </v-btn>
  </form>
</template>

<script lang="ts">
import { computed, ref } from '@nuxtjs/composition-api'
import { defineComponent, PropType } from '@vue/composition-api'
import dayjs from 'dayjs'

import { CreatePromotionRequest } from '~/types/api'

export default defineComponent({
  props: {
    formType: {
      type: String,
      default: 'create',
      validator: (value: string) => {
        return ['create', 'edit'].includes(value)
      },
    },
    formData: {
      type: Object as PropType<CreatePromotionRequest>,
      default: () => {
        return {
          title: '',
          description: '',
          public: false,
          publishedAt: dayjs().unix(),
          discountType: 1,
          discountRate: 0,
          code: '',
          startAt: dayjs().unix(),
          endAt: dayjs().unix(),
        }
      },
    },
  },

  setup(props, { emit }) {
    const selectedDiscountMethod = ref<string>('')
    const postMenu = ref<boolean>(false)
    const useStartMenu = ref<boolean>(false)
    const useEndMenu = ref<boolean>(false)
    const postDate = ref<string>('')
    const useStartDate = ref<string>('')
    const useEndDate = ref<string>('')

    const btnText = computed(() => {
      return props.formType === 'create' ? '登録' : '更新'
    })

    const handleSubmit = () => {
      emit('submit')
    }

    const status = ['有効', '無効']

    const discountMethod = ['%', '円', '配送料無料']

    return {
      selectedDiscountMethod,
      postMenu,
      useStartMenu,
      useEndMenu,
      discountMethod,
      btnText,
      postDate,
      useStartDate,
      useEndDate,
      status,
      handleSubmit,
    }
  },
})
</script>

<style scoped lang="scss">
.postTime {
}

.startUseTime {
}

.endUseTime {
}
</style>
