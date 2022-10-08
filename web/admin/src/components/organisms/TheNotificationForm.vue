<template>
  <form @submit.prevent="handleSubmit">
    <v-card elevation="0">
      <v-card-text>
        <v-select
          v-model="formData.public"
          :items="statusList"
          label="ステータス"
          item-text="public"
          item-value="value"
        ></v-select>
        <v-text-field
          v-model="formData.title"
          label="タイトル"
          required
          maxlength="128"
        />
        <v-textarea v-model="formData.body" label="本文" maxlength="2000" />
      </v-card-text>
      <v-container class="ml-2">
        <p class="text-h6">公開範囲</p>
        <v-checkbox
          v-model="formData.targets"
          label="ユーザー"
          :value="Number(1)"
        ></v-checkbox>
        <v-checkbox
          v-model="formData.targets"
          label="生産者"
          :value="Number(2)"
        ></v-checkbox>
        <v-checkbox
          v-model="formData.targets"
          label="コーディネーター"
          :value="Number(3)"
        ></v-checkbox>
        <p class="text-h6">投稿予約時間</p>
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
          <v-text-field v-model="postTime" type="time" required outlined />
          <p class="text-h6 mb-6 ml-4">〜</p>
          <v-spacer />
        </div>
      </v-container>
      <v-card-actions>
        <v-btn block outlined color="primary" type="submit" class="mt-4">
          {{ btnText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </form>
</template>

<script lang="ts">
import { computed, PropType, ref } from '@nuxtjs/composition-api'
import { defineComponent } from '@vue/composition-api'
import dayjs from 'dayjs'

import { CreateNotificationRequest } from '~/types/api'

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
      type: Object as PropType<CreateNotificationRequest>,
      default: () => {
        return {
          title: '',
          body: '',
          targets: [],
          public: false,
          publishedAt: dayjs().unix(),
        }
      },
    },
  },

  setup(props, { emit }) {
    const formDataValue = computed({
      get: (): CreateNotificationRequest => props.formData,
      set: (val: CreateNotificationRequest) => emit('update:formData', val),
    })

    const btnText = computed(() => {
      return props.formType === 'create' ? '登録' : '更新'
    })
    const postMenu = ref<boolean>(false)
    const postDate = ref<string>('')
    const postTime = ref<string>('')

    const statusList = [
      { public: '公開', value: true },
      { public: '非公開', value: false },
    ]

    const handleSubmit = () => {
      const unixTime = dayjs(postDate.value + ' ' + postTime.value).unix()
      formDataValue.value.publishedAt = unixTime
      emit('submit')
    }

    return {
      formDataValue,
      btnText,
      statusList,
      postDate,
      postTime,
      postMenu,
      handleSubmit,
    }
  },
})
</script>
