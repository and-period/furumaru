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
      <v-btn block outlined color="primary" @click="certificationBtn">認証</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import { defineComponent, reactive, useRoute, useRouter } from '@nuxtjs/composition-api'
import { useAuthStore } from '~/store/auth'
import { UpdateAuthEmailRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const formData = reactive<UpdateAuthEmailRequest>({
      email: '',
    })
    const router = useRouter()
    const route = useRoute()
    const email = route.value.params.email
    const authStore = useAuthStore()
//    const useAuthApi = UpdateAuthEmailRequest()

    const handleClickAddBtn = async (): Promise<void> => {
      try {
        await authStore.emailUpdate(formData)
        router.push({
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

    // const certificationBtn = async ():Promise<void> => {
    //   try{
    //     await useAuthApi

    //   } catch (error){
    //     console.log(error)
    //   }
    // }
  },
})
</script>

<!--
37行目でemailを渡すよう言われたと思いますが
await authStore.emailUpdate(email)
としてしまうと↓のエラーが出ます。
Argument of type 'string' is not assignable to parameter of type 'UpdateAuthEmailRequest'.


認証ボタンの方は
api.tsのv1VerifyAuthEmailを使うのかと思ったんですが、どう使えばいいか分かリませんでした

-->
