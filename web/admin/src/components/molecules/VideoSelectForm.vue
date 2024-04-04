<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
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
    default: (): string[] => ['video/mp4']
  },
  error: {
    type: Boolean,
    default: false
  },
  message: {
    type: String,
    default: ''
  },
  validationErrorMessage: {
    type: String,
    default: ''
  }
})

const emit = defineEmits<{
  (e: 'update:file', files?: FileList): void
}>()

const inputRef = ref<HTMLInputElement | null>(null)

const onClick = (): void => {
  if (!inputRef.value) {
    return
  }
  inputRef.value.click()
}

const acceptedFiles = (): string => {
  return props.accept.join(',')
}

const onChangeFile = (): void => {
  if (inputRef.value && inputRef.value.files) {
    emit('update:file', inputRef.value.files)
  }
}
</script>

<template>
  <div class="d-flex flex-column flex-grow-1 flex-shrink-1 pb-4">
    <p>{{ props.label }}の設定</p>
    <v-card
      :disabled="loading"
      :loading="loading"
      :style="{ border: props.validationErrorMessage !== '' ? '1px solid red' : '' }"
      class="text-center"
      role="button"
      flat
      @click="onClick"
    >
      <v-card-text>
        <div class="mb-4">
          <v-avatar v-if="props.videoUrl === ''" size="80" :icon="mdiPlus" />
          <v-responsive v-else max-height="300" aspect-ratio="16 / 9">
            <video id="wrapper" muted autoplay :src="props.videoUrl" />
          </v-responsive>
        </div>
        <input
          ref="inputRef"
          type="file"
          class="d-none"
          :accept="acceptedFiles()"
          @change="onChangeFile"
        >
        <p class="ma-0">
          {{ props.label }}を{{ props.videoUrl === '' ? '選択' : '変更' }}
        </p>
      </v-card-text>
    </v-card>
    <p v-show="props.validationErrorMessage" class="text-red ma-0">
      {{ props.validationErrorMessage }}
    </p>
    <p v-show="props.error" class="red--text ma-0">
      {{ props.message }}
    </p>
  </div>
</template>

<style scoped>
#wrapper {
  position:relative;
  width: 100%;
  height:100%;
  background-color:#030303;
}
</style>
