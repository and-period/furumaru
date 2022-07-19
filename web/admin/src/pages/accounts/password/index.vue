<template>
  <div>
    <v-card-title>パスワード変更</v-card-title>
    <v-card>
      <v-container>
        <form @submit.prevent="handleSubmit">
          <v-text-field
            v-model="formData.oldPassword"
            class="mx-4"
            minlength="8"
            maxlength="32"
            label="現在のパスワード"
            :append-icon="oldPasswordShow ? 'mdi-eye' : 'mdi-eye-off'"
            :type="oldPasswordShow ? 'text' : 'password'"
            @click:append="oldPasswordShow = !oldPasswordShow"
          />
          <v-text-field
            v-model="formData.newPassword"
            class="mx-4"
            minlength="8"
            maxlength="32"
            label="新しいパスワード"
            :append-icon="newPasswordShow ? 'mdi-eye' : 'mdi-eye-off'"
            :type="newPasswordShow ? 'text' : 'password'"
            @click:append="newPasswordShow = !newPasswordShow"
          />
          <v-text-field
            v-model="formData.passwordConfirmation"
            class="mx-4"
            min-length="8"
            maxlength="32"
            label="新しいパスワード(確認用)"
            :append-icon="passwordConfirmationShow ? 'mdi-eye' : 'mdi-eye-off'"
            :type="passwordConfirmationShow ? 'text' : 'password'"
            :error-messages="isMatch ? '' : 'パスワードが一致しません'"
            @click:append="passwordConfirmationShow = !passwordConfirmationShow"
          />
          <div class="d-flex justify-end mr-4">
            <v-btn outlined color="primary" type="submit"> 変更 </v-btn>
          </div>
        </form>
      </v-container>
      <v-alert v-model="isShow" :type="alertType" v-text="alertText" />
    </v-card>
  </div>
</template>

<script lang="ts">
import { computed, useRouter } from '@nuxtjs/composition-api'
import { defineComponent, reactive, ref } from '@vue/composition-api'

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

    const isMatch = computed(() => {
      return formData.newPassword === formData.passwordConfirmation
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
      isMatch,
      handleSubmit,
    }
  },
})
</script>
