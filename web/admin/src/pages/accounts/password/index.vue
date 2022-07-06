<template>
  <div>
    <p class="text-h6">パスワード変更</p>
    <v-card>
      <v-container>
        <v-text-field
          v-model="formData.oldPassword"
          class="mx-4"
          maxlength="32"
          label="現在のパスワード"
          :append-icon="oldPasswordShow ? 'mdi-eye' : 'mdi-eye-off'"
          :type="oldPasswordShow ? 'text' : 'password'"
          @click:append="oldPasswordShow = !oldPasswordShow"
        />
        <v-text-field
          v-model="formData.newPassword"
          class="mx-4"
          maxlength="32"
          label="新しいパスワード"
          :append-icon="newPasswordShow ? 'mdi-eye' : 'mdi-eye-off'"
          :type="newPasswordShow ? 'text' : 'password'"
          @click:append="newPasswordShow = !newPasswordShow"
        />
        <v-text-field
          v-model="formData.passwordConfirmation"
          class="mx-4"
          maxlength="32"
          label="新しいパスワード(確認用)"
          :append-icon="passwordConfirmationShow ? 'mdi-eye' : 'mdi-eye-off'"
          :type="passwordConfirmationShow ? 'text' : 'password'"
          @click:append="passwordConfirmationShow = !passwordConfirmationShow"
        />
        <div class="d-flex justify-end mr-4">
          <v-btn outlined color="primary" @click="handleSubmit"> 変更 </v-btn>
        </div>
      </v-container>
      <v-alert v-model="isShow" :type="alertType" v-text="alertText" />
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
import { UpdateAuthPasswordRequest } from '~/types/api'

export default defineComponent({
  /*
  TODO: validation追加
  newPassowrd passwordConfirmationの一致、入力文字数の制限、必須文字のcheck
  */

  setup() {
    const router = useRouter()
    const formData = reactive<UpdateAuthPasswordRequest>({
      oldPassword: '',
      newPassword: '',
      passwordConfirmation: '',
    })

    const oldPasswordShow = ref<Boolean>(false)
    const newPasswordShow = ref<Boolean>(false)
    const passwordConfirmationShow = ref<Boolean>(false)
    const { alertType, isShow, alertText, show } = useAlert('error')
    const authStore = useAuthStore()

    const handleSubmit = async (): Promise<void> => {
      try {
        await authStore.passwordUpdate(formData)
        router.push('/')
      } catch (error) {
        console.log(error)
        show('パスワードの更新に失敗しました。')
      }
    }

    return {
      alertType,
      isShow,
      alertText,
      formData,
      oldPasswordShow,
      newPasswordShow,
      passwordConfirmationShow,
      handleSubmit,
    }
  },
})
</script>
