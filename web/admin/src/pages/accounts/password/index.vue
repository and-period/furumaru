<script lang="ts" setup>
import { mdiEye, mdiEyeOff } from '@mdi/js'
import { useVuelidate } from '@vuelidate/core'

import { useAlert } from '~/lib/hooks'
import {
  required,
  minLength,
  maxLength,
  sameAs,
  getErrorMessage
} from '~/lib/validations'
import { useAuthStore } from '~/store'
import { UpdateAuthPasswordRequest } from '~/types/api'

const router = useRouter()
const formData = reactive<UpdateAuthPasswordRequest>({
  oldPassword: '',
  newPassword: '',
  passwordConfirmation: ''
})

const rules = computed(() => ({
  oldPassword: { required },
  newPassword: {
    required,
    minLength: minLength(8),
    maxLength: maxLength(32)
  },
  passwordConfirmation: { required, sameAs: sameAs(formData.newPassword) }
}))

const v$ = useVuelidate(rules, formData)

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
</script>

<template>
  <div>
    <v-card-title>パスワード変更</v-card-title>
    <v-alert .sync="isShow" :type="alertType" v-text="alertText" />
    <form @submit.prevent="handleSubmit">
      <v-card>
        <v-card-text>
          <v-text-field
            v-model="v$.oldPassword.$model"
            :error-messages="getErrorMessage(v$.oldPassword.$errors)"
            label="現在のパスワード"
            :append-icon="oldPasswordShow ? 'mdiEye' : mdiEyeOff"
            :type="oldPasswordShow ? 'text' : 'password'"
            @click:append="oldPasswordShow = !oldPasswordShow"
          />
          <v-text-field
            v-model="v$.newPassword.$model"
            :error-messages="getErrorMessage(v$.newPassword.$errors)"
            label="新しいパスワード"
            :append-icon="newPasswordShow ? 'mdiEye' : mdiEyeOff"
            :type="newPasswordShow ? 'text' : 'password'"
            @click:append="newPasswordShow = !newPasswordShow"
          />
          <v-text-field
            v-model="v$.passwordConfirmation.$model"
            label="新しいパスワード(確認用)"
            :append-icon="passwordConfirmationShow ? 'mdiEye' : mdiEyeOff"
            :type="passwordConfirmationShow ? 'text' : 'password'"
            :error-messages="
              getErrorMessage(v$.passwordConfirmation.$errors) === ''
                ? ''
                : 'パスワードが一致しません。'
            "
            @click:append="passwordConfirmationShow = !passwordConfirmationShow"
          />
        </v-card-text>
        <v-card-actions>
          <v-btn block variant="outlined" color="primary" type="submit">
            変更
          </v-btn>
        </v-card-actions>
      </v-card>
    </form>
  </div>
</template>
