<template>
  <v-card elevation="0">
    <h1 class="text-center">二要素認証</h1>
    <v-card-text>
      <h4 class="text-center">認証コードが{{ email }}に送信されました</h4>
      <div class="ma-auto" style="max-width: 300px">
        <v-otp-input v-model="otp" type="number" length="6"></v-otp-input>
      </div>
      <div class="text-center">
        <a class="orange--text text--darken-4" @click="handleClickAddBtn">認証コードを再送する</a>
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn block outlined color="primary" @click="verificationBtn">認証</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import { defineComponent, reactive, useRoute, useRouter } from '@nuxtjs/composition-api'
import { useAuthStore } from '~/store/auth'
import { UpdateAuthEmailRequest, VerifyAuthEmailRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const formData = reactive<UpdateAuthEmailRequest>({
      email: '',
    })
    const router = useRouter()
    const route = useRoute()
    const email = route.value.params.email
    const authStore = useAuthStore()
    const verifyCode = route.value.params.verifyCode


    const handleClickAddBtn = async (): Promise<void> => {
      try {
        await authStore.emailUpdate(<UpdateAuthEmailRequest><unknown>email)
        router.push({
          //name: 'accounts-email-verification',
          params: { email: formData.email}
        })
      } catch (error){
        console.log(error)
      }
    }
    return {
      handleClickAddBtn,
      email,
    }

    const verificationBtn = async ():Promise<void> => {
      try{
        await authStore.codeVerify(<VerifyAuthEmailRequest><unknown>verifyCode)
        router.push({
          params: {verifyCode: verifyCode}
        })

      } catch (error){
        console.log(error)
      }
    }
    return {
      verificationBtn,
      verifyCode,
    }
  },
})
</script>

<!--
38行目は型変換の仕方がわからず、クイックフィクスに頼った結果<unknown>がつきました。
53行目では入力された認証コードを渡すのかなと思ってapi.ts(パシリの中身？)で宣言されてそうなの入れました
認証ボタンにイベントを入れたつもりなんですが、51~が参照されていないのがなぜかわかりません
40行目でもう一度verification画面に遷移させる必要はあるのでしょうか
-->
