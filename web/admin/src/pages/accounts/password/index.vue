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
      </v-container>
    </v-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from '@nuxtjs/composition-api'

import { UpdateAuthPasswordRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const formData = reactive<UpdateAuthPasswordRequest>({
      oldPassword: '',
      newPassword: '',
      passwordConfirmation: '',
    })
    const oldPasswordShow = ref<Boolean>(false)
    const newPasswordShow = ref<Boolean>(false)
    const passwordConfirmationShow = ref<Boolean>(false)

    const handleSubmit = async () => {}

    return {
      formData,
      oldPasswordShow,
      newPasswordShow,
      passwordConfirmationShow,
      handleSubmit,
    }
  },
})
</script>
