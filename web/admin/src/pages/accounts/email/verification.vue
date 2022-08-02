<template>
  <v-card elevation="0">
    <h1 class="text-center">二要素認証</h1>
    <v-card-text>
      <h4 class="text-center">認証コードが{{ email }}に送信されました</h4>
      <div class="ma-auto" style="max-width: 300px">
        <v-otp-input
          v-model="formData.email"
          type="number"
          length="6"
        ></v-otp-input>
      </div>
      <div class="text-center">
        <a class="orange--text text--darken-4" @click="handleClickAddBtn"
          >認証コードを再送する</a
        >
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn block outlined color="primary" @click="verificationBtn"
        >認証</v-btn
      >
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import {
  defineComponent,
  reactive,
  useRoute,
  useRouter,
} from '@nuxtjs/composition-api'

import { useAuthStore } from '~/store/auth'
import { UpdateAuthEmailRequest, VerifyAuthEmailRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const formData = reactive<VerifyAuthEmailRequest>({
      verifyCode: '',
    })
    const router = useRouter()
    const route = useRoute()
    const email = route.value.params.email
    const authStore = useAuthStore()
    const verifyCode = route.value.params.verifyCode
    const convertEmail: UpdateAuthEmailRequest = {
      email,
    }

    const handleClickAddBtn = async (): Promise<void> => {
      try {
        await authStore.emailUpdate(convertEmail)
      } catch (error) {
        console.log(error)
      }
    }

    const verificationBtn = async (): Promise<void> => {
      try {
        await authStore.codeVerify(formData)
        router.push('/')
      } catch (error) {
        console.log(error)
      }
    }

    return {
      verificationBtn,
      handleClickAddBtn,
      verifyCode,
      formData,
    }
  },
})
</script>
