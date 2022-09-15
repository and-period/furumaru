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
        <v-select :items="status" label="ステータス" />
        <div class="d-flex align-center">
          <v-select
            v-model="discountTypeString"
            :items="discountMethod"
            label="割引方法"
          />
          <v-text-field
            v-if="discountTypeString != '配送料無料'"
            v-model="formData.discountRate"
            class="ml-4"
            label="割引値"
            :error-messages="getErrorMessage()"
          />
        </div>

        <p class="text-h6">投稿期間</p>
        <div class="d-flex align-center justify-center">
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
                class="mr-2"
                label="投稿開始日"
                readonly
                outlined
                v-bind="attrs"
                v-on="on"
              />
            </template>
            <v-date-picker
              v-model="postDate"
              scrollable
              @input="postMenu = false"
            >
              <v-spacer></v-spacer>
              <v-btn text color="primary" @click="postMenu = false">
                閉じる
              </v-btn>
            </v-date-picker>
          </v-menu>
          <v-text-field class="postTime" type="time" required outlined />
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
                v-model="useStartDate"
                label="使用開始日"
                readonly
                outlined
                v-bind="attrs"
                class="mr-2"
                v-on="on"
              />
            </template>
            <v-date-picker
              v-model="useStartDate"
              scrollable
              @input="useStartMenu = false"
            >
              <v-spacer></v-spacer>
              <v-btn text color="primary" @click="useStartMenu = false">
                閉じる
              </v-btn>
            </v-date-picker>
          </v-menu>
          <v-text-field type="time" required outlined />
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
                v-model="useEndDate"
                label="使用終了日"
                readonly
                outlined
                v-bind="attrs"
                class="mr-2"
                v-on="on"
              />
            </template>
            <v-date-picker
              v-model="useEndDate"
              scrollable
              @input="useEndMenu = false"
            >
              <v-spacer></v-spacer>
              <v-btn text color="primary" @click="endEndMenu = false">
                閉じる
              </v-btn>
            </v-date-picker>
          </v-menu>
          <v-text-field type="time" required outlined />
        </div>
      </v-card-text>
    </v-card>
    <v-btn block outlined color="primary" type="submit" class="mt-4">
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
    const formDataValue = computed({
      get: (): CreatePromotionRequest => props.formData,
      set: (val: CreatePromotionRequest) => emit('update:formData', val),
    })

    const selectedDiscountMethod = ref<string>('')
    const postMenu = ref<boolean>(false)
    const useStartMenu = ref<boolean>(false)
    const discountTypeString = ref<string>('')
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
      if (discountTypeString.value === '%') {
        if (
          props.formData.discountRate >= 0 &&
          props.formData.discountRate <= 100
        ) {
          return ''
        }
        return '0~100の値を指定してください'
      } else if (discountTypeString.value === '円') {
        if (props.formData.discountRate >= 0) {
          return ''
        }
        return '0以上の値を指定してください'
      } else {
        return ''
      }
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
      discountTypeString,
      formDataValue,
      getErrorMessage,
      handleGenerate,
      handleSubmit,
    }
  },
})
</script>
