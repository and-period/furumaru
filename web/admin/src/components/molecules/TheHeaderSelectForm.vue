<template>
  <div>
    <p>ヘッダー画像の設定</p>
    <v-card
      class="d-flex flex-column text-center align-center"
      role="button"
      min-width="180"
      flat
      :img="imgUrl"
      @click="handleClick"
    >
      <v-card-text>
        <v-icon>mdi-plus</v-icon>
        <input
          ref="inputRef"
          type="file"
          class="d-none"
          accept="image/*"
          @change="handleInputFileChange"
        />
        <p class="ma-0">ヘッダー画像を選択</p>
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
