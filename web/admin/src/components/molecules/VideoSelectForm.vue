<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'

const props = defineProps({
  label: {
    type: String,
    default: ''
  },
  videoUrl: {
    type: String,
    default: ''
  },
  accept: {
    type: Array<String>,
    default: (): string[] => ['video/*']
  },
  error: {
    type: Boolean,
    default: false
  },
  message: {
    type: String,
    default: ''
  }
})

const emit = defineEmits<{
  (e: 'update:file', files?: FileList): void
}>()

const inputRef = ref<HTMLInputElement | null>(null)

const onClick = () => {
  if (inputRef.value === null) {
    return
  }
  inputRef.value.click()
}

const acceptedFiles = (): string => {
  return props.accept.join(',')
}

const onChangeFile = () => {
  if (inputRef.value && inputRef.value.files) {
    emit('update:file', inputRef.value.files)
  }
}
</script>

<template>
  <div>
    <p>{{ props.label }}の設定</p>
    <v-card
      class="d-flex flex-column text-center align-center"
      role="button"
      min-width="180"
      flat
      :img="props.videoUrl"
      @click="onClick"
    >
      <v-card-text>
        <v-avatar size="96">
          <v-icon x-large :icon="mdiPlus" />
        </v-avatar>
        <input
          ref="inputRef"
          type="file"
          class="d-none"
          :accept="acceptedFiles()"
          @change="onChangeFile"
        >
        <p class="ma-0">{{ props.label }}を選択</p>
      </v-card-text>
    </v-card>
    <p v-show="props.error" class="red--text ma-0">{{ props.message }}</p>
  </div>
</template>
