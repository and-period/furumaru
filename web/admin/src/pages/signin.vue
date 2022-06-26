<template>
  <div>
    <v-alert v-model="isShow" :type="alertType" v-text="alertText" />
    <v-card>
      <form @submit.prevent="handleSubmit">
        <v-card-text>
          <v-text-field
            v-model="formData.username"
            label="ユーザーID（メールアドレス)"
            type="email"
          />
          <v-text-field
            v-model="formData.password"
            label="パスワード"
            type="password"
          />
        </v-card-text>
        <v-card-actions>
          <v-btn elevation="0" block color="primary" type="submit"
            >ログイン</v-btn
          >
        </v-card-actions>
      </form>
    </v-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from '@nuxtjs/composition-api'

import { useAlert } from '~/lib/hooks'
import { SignInRequest } from '~/types/api'

export default defineComponent({
  layout: 'auth',
  setup() {
    const formData = reactive<SignInRequest>({
      username: '',
      password: '',
    })

    const { alertType, isShow, alertText, show } = useAlert('error')

    const handleSubmit = () => {
      console.log('未実装', formData)
      show('ユーザーIDまたはパスワードが違います。')
    }

    return {
      alertType,
      isShow,
      alertText,
      formData,
      handleSubmit,
    }
  },
})
</script>
