<template>
  <div>
    <v-alert v-model="isShow" :type="alertType" v-text="alertText" />
    <div class="pa-8">
      <the-app-logo-with-title :width="450" class="ma-auto" />
    </div>
    <v-card>
      <form @submit.prevent="handleSubmit">
        <v-card-title>ログイン</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="formData.username"
            label="ユーザーID（メールアドレス)"
            type="email"
            required
          />
          <v-text-field
            v-model="formData.password"
            label="パスワード"
            :append-icon="passwordShow ? 'mdi-eye' : 'mdi-eye-off'"
            :type="passwordShow ? 'text' : 'password'"
            required
            @click:append="passwordShow = !passwordShow"
          />
        </v-card-text>
        <v-card-actions>
          <v-btn block color="primary" type="submit" outlined> ログイン </v-btn>
        </v-card-actions>
      </form>
    </v-card>
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  reactive,
  ref,
  useRouter,
} from '@nuxtjs/composition-api'

import { useAlert } from '~/lib/hooks'
import { useAuthStore } from '~/store/auth'
import { SignInRequest } from '~/types/api'

export default defineComponent({
  layout: 'auth',
  setup() {
    const router = useRouter()
    const formData = reactive<SignInRequest>({
      username: '',
      password: '',
    })
    const passwordShow = ref<boolean>(false)
    const { alertType, isShow, alertText, show } = useAlert('error')
    const authStore = useAuthStore()

    const handleSubmit = async () => {
      try {
        const path = await authStore.signIn(formData)
        router.push(path)
      } catch (error) {
        if (error instanceof Error) {
          show(error.message)
        }
      }
    }

    return {
      alertType,
      isShow,
      alertText,
      formData,
      handleSubmit,
      passwordShow,
    }
  },
})
</script>
