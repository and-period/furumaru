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
            label="割引コード(8文字)"
            required
            maxlength="8"
            :error-messages="
              formData.code.length === 8 ? '' : '割引コードは8文字です'
            "
          />
          <v-btn outlined small color="primary" @click="handleGenerate">
            自動生成
          </v-btn>
          <v-spacer />
        </div>
        <v-select
          v-model="formData.public"
          :items="statusList"
          label="ステータス"
          item-text="status"
          item-value="value"
        />
        <div class="d-flex align-center">
          <v-select
            v-model="formData.discountType"
            :items="discountMethodList"
            item-text="method"
            item-value="value"
            label="割引方法"
          />
          <v-text-field
            v-if="formData.discountType != 3"
            v-model="formData.discountRate"
            class="ml-4"
            type="number"
            label="割引値"
            :error-messages="getErrorMessage()"
          />
        </div>

        <p class="text-h6">投稿開始</p>
        <div class="d-flex align-center justify-center">
          <v-menu
            v-model="publishMenu"
            :close-on-content-click="false"
            :nudge-right="40"
            transition="scale-transition"
            offset-y
            min-width="auto"
          >
            <template #activator="{ on, attrs }">
              <v-text-field
                v-model="timeData.publishedDate"
                class="mr-2"
                label="投稿開始日"
                readonly
                outlined
                v-bind="attrs"
                v-on="on"
              />
            </template>
            <v-date-picker
              v-model="timeData.publishedDate"
              scrollable
              @input="publishMenu = false"
            >
              <v-spacer></v-spacer>
              <v-btn text color="primary" @click="pubulishMenu = false">
                閉じる
              </v-btn>
            </v-date-picker>
          </v-menu>
          <v-text-field
            v-model="timeData.publishedTime"
            type="time"
            required
            outlined
          />
          <p class="text-h6 mb-6 ml-4">〜</p>
          <v-spacer />
        </div>

        <p class="text-h6">使用期間</p>
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
                v-model="timeData.startDate"
                label="使用開始日"
                readonly
                outlined
                v-bind="attrs"
                class="mr-2"
                v-on="on"
              />
            </template>
            <v-date-picker
              v-model="timeData.startDate"
              scrollable
              @input="useStartMenu = false"
            >
              <v-spacer></v-spacer>
              <v-btn text color="primary" @click="useStartMenu = false">
                閉じる
              </v-btn>
            </v-date-picker>
          </v-menu>
          <v-text-field
            v-model="timeData.startTime"
            type="time"
            required
            outlined
          />
          <p class="text-h6 mx-4 mb-6">〜</p>
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
                v-model="timeData.endDate"
                label="使用終了日"
                readonly
                outlined
                v-bind="attrs"
                class="mr-2"
                v-on="on"
              />
            </template>
            <v-date-picker
              v-model="timeData.endDate"
              scrollable
              @input="useEndMenu = false"
            >
              <v-spacer></v-spacer>
              <v-btn text color="primary" @click="endEndMenu = false">
                閉じる
              </v-btn>
            </v-date-picker>
          </v-menu>
          <v-text-field
            v-model="timeData.endTime"
            type="time"
            required
            outlined
          />
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" type="submit" class="mt-4">
          {{ btnText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </form>
</template>

<script lang="ts">
import { computed, ref } from '@nuxtjs/composition-api'
import { defineComponent, PropType } from '@vue/composition-api'
import dayjs from 'dayjs'

import { CreatePromotionRequest } from '~/types/api'
import { PromotionTime } from '~/types/props'

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
    timeData: {
      type: Object as PropType<PromotionTime>,
      default: () => {
        return {
          publishedDate: '',
          publishedTime: '',
          startDate: '',
          startTime: '',
          endDate: '',
          endTime: '',
        }
      },
    },
  },

  setup(props, { emit }) {
    const formDataValue = computed({
      get: (): CreatePromotionRequest => props.formData,
      set: (val: CreatePromotionRequest) => emit('update:formData', val),
    })

    const timeDataValue = computed({
      get: (): PromotionTime => props.timeData,
      set: (val: PromotionTime) => emit('update:timeData', val),
    })

    const selectedDiscountMethod = ref<string>('')
    const publishMenu = ref<boolean>(false)
    const useStartMenu = ref<boolean>(false)
    const useEndMenu = ref<boolean>(false)

    const btnText = computed(() => {
      return props.formType === 'create' ? '登録' : '更新'
    })

    const handleSubmit = () => {
      formDataValue.value.publishedAt = dayjs(
        timeDataValue.value.publishedDate +
          ' ' +
          timeDataValue.value.publishedTime
      ).unix()
      formDataValue.value.startAt = dayjs(
        timeDataValue.value.startDate + ' ' + timeDataValue.value.startTime
      ).unix()
      formDataValue.value.endAt = dayjs(
        timeDataValue.value.endDate + ' ' + timeDataValue.value.endTime
      ).unix()

      emit('submit')
    }

    const handleGenerate = () => {
      const code = generateRandomString()
      props.formData.code = code
    }

    const generateRandomString = (): string => {
      const characters =
        'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
      let result = ''
      const charactersLength = characters.length
      for (let i = 0; i < 8; i++) {
        result += characters.charAt(
          Math.floor(Math.random() * charactersLength)
        )
      }

      return result
    }

    const getErrorMessage = () => {
      switch (formDataValue.value.discountType) {
        case 1:
          if (formDataValue.value.discountRate >= 0) {
            return ''
          } else {
            return '0以上の値を指定してください'
          }
        case 2:
          if (
            formDataValue.value.discountRate >= 0 &&
            formDataValue.value.discountRate <= 100
          ) {
            return ''
          } else {
            return '0~100の値を指定してください'
          }
        default:
          return ''
      }
    }

    const statusList = [
      { status: '有効', value: true },
      { status: '無効', value: false },
    ]

    const discountMethodList = [
      { method: '円', value: 1 },
      { method: '%', value: 2 },
      { method: '送料無料', value: 3 },
    ]

    return {
      selectedDiscountMethod,
      publishMenu,
      useStartMenu,
      useEndMenu,
      discountMethodList,
      btnText,
      statusList,
      formDataValue,
      timeDataValue,
      getErrorMessage,
      handleGenerate,
      handleSubmit,
    }
  },
})
</script>
