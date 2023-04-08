<script lang="ts" setup>
import { useAuthStore } from '~/store'
import { UpdateAuthEmailRequest, VerifyAuthEmailRequest } from '~/types/api'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = route.params.email as string
const verifyCode = route.params.verifyCode as string

const convertEmail: UpdateAuthEmailRequest = {
  email,
}
const formData = reactive<VerifyAuthEmailRequest>({
  verifyCode,
})

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
</script>

<template>
  <v-card elevation="0">
    <h1 class="text-center">二要素認証</h1>
    <v-card-text>
      <h4 class="text-center">認証コードが{{ email }}に送信されました</h4>
      <div class="ma-auto" style="max-width: 300px">
        <v-otp-input
          v-model="formData.verifyCode"
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
