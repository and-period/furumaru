<script lang="ts" setup>
const props = defineProps({
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
})

const emit = defineEmits<{
  (e: 'update:file', files?: FileList): void
}>()

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
</script>

<template>
  <div>
    <p>ヘッダー画像の設定</p>
    <v-card
      class="d-flex flex-column text-center align-center"
      role="button"
      min-width="180"
      flat
      :img="props.imgUrl"
      @click="handleClick"
    >
      <v-card-text>
        <v-avatar size="96">
          <v-icon x-large>mdi-plus</v-icon>
        </v-avatar>
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
    <p v-show="props.error" class="red--text ma-0">{{ props.message }}</p>
  </div>
</template>
