<template>
  <div>
    <v-card-title>パスワード変更</v-card-title>
    <v-card>
      <v-container>
        <form @submit.prevent="handleSubmit">
          <v-text-field
            v-model="v$.oldPassword.$model"
            :error-messages="getErrorMessage(v$.oldPassword.$errors)"
            class="mx-4"
            label="現在のパスワード"
            :append-icon="oldPasswordShow ? 'mdi-eye' : 'mdi-eye-off'"
            :type="oldPasswordShow ? 'text' : 'password'"
            @click:append="oldPasswordShow = !oldPasswordShow"
          />
          <v-text-field
            v-model="v$.newPassword.$model"
            :error-messages="getErrorMessage(v$.newPassword.$errors)"
            class="mx-4"
            label="新しいパスワード"
            :append-icon="newPasswordShow ? 'mdi-eye' : 'mdi-eye-off'"
            :type="newPasswordShow ? 'text' : 'password'"
            @click:append="newPasswordShow = !newPasswordShow"
          />
          <v-text-field
            v-model="v$.passwordConfirmation.$model"
            class="mx-4"
            label="新しいパスワード(確認用)"
            :append-icon="passwordConfirmationShow ? 'mdi-eye' : 'mdi-eye-off'"
            :type="passwordConfirmationShow ? 'text' : 'password'"
            :error-messages="
              getErrorMessage(v$.passwordConfirmation.$errors) === ''
                ? ''
                : 'パスワードが一致しません。'
            "
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
import { useVuelidate } from '@vuelidate/core'

import { useAlert } from '~/lib/hooks'
import {
  required,
  minLength,
  maxLength,
  sameAs,
  getErrorMessage,
} from '~/lib/validations'
import { useAuthStore } from '~/store/auth'
import { UpdateAuthPasswordRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const router = useRouter()
    const formData = reactive<UpdateAuthPasswordRequest>({
      oldPassword: '',
      newPassword: '',
      passwordConfirmation: '',
    })

    const rules = computed(() => ({
      oldPassword: { required },
      newPassword: {
        required,
        minLength: minLength(8),
        maxLength: maxLength(32),
      },
      passwordConfirmation: { required, sameAs: sameAs(formData.newPassword) },
    }))

    const v$ = useVuelidate(rules, formData)

    const isMatch = computed(() => {
      return formData.newPassword === formData.passwordConfirmation
    })

    const oldPasswordShow = ref<boolean>(false)
    const newPasswordShow = ref<boolean>(false)
    const passwordConfirmationShow = ref<boolean>(false)

    const { alertType, isShow, alertText, show } = useAlert('error')
    const authStore = useAuthStore()

    const handleSubmit = async (): Promise<void> => {
      const result = await v$.value.$validate()
      if (!result) {
        return
      }
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
      v$,
      formData,
      oldPasswordShow,
      newPasswordShow,
      passwordConfirmationShow,
      isMatch,
      handleSubmit,
      getErrorMessage,
    }
  },
})
</script>
