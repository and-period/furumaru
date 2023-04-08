<script lang="ts" setup>
import { useAuthStore } from '~/store'
import { UpdateAuthEmailRequest } from '~/types/api'

const formData = reactive<UpdateAuthEmailRequest>({
  email: ''
})
const router = useRouter()
const authStore = useAuthStore()

const handleClickAddBtn = async (): Promise<void> => {
  try {
    await authStore.emailUpdate(formData)
    router.push({
      name: 'accounts-email-verification',
      params: { email: formData.email }
    })
  } catch (error) {
    console.log(error)
  }
}
</script>

<template>
  <div>
    <p class="text-h6">
      メールアドレス変更
    </p>
    <v-card elevation="0">
      <v-card-text>
        <v-text-field v-model="formData.email" label="新規メールアドレス" />
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" @click="handleClickAddBtn">
          変更
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>
