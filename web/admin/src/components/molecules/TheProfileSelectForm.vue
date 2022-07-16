<template>
  <div>
    <p>アイコン画像の設定</p>
    <v-card class="text-center" role="button" flat @click="handleClick">
      <v-card-text>
        <v-avatar size="96">
          <v-icon v-if="imgUrl === ''" x-large>mdi-account</v-icon>
          <v-img
            v-else
            :src="imgUrl"
            aspect-ratio="1"
            max-height="150"
            contain
          />
        </v-avatar>
        <input
          ref="inputRef"
          type="file"
          class="d-none"
          accept="image/*"
          @change="handleInputFileChange"
        />
        <p class="ma-0">アイコン画像を選択</p>
      </v-card-text>
    </v-card>
    <p v-show="error" class="red--text ma-0">{{ message }}</p>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from '@vue/composition-api'

export default defineComponent({
  props: {
    imgUrl: {
      type: String,
      default: '',
    },
    error: {
      type: Boolean,
      default: false,
    },
    message: {
      type: String,
      default: '',
    },
  },
  setup(_, { emit }) {
    const inputRef = ref<HTMLInputElement | null>(null)
    const handleClick = () => {
      if (inputRef.value !== null) {
        inputRef.value.click()
      }
    }

    const handleInputFileChange = () => {
      if (inputRef.value && inputRef.value.files) {
        emit('update:file', inputRef.value.files)
      }
    }

    return {
      inputRef,
      handleClick,
      handleInputFileChange,
    }
  },
})
</script>
