<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  label: {
    type: String,
    default: '',
  },
  imgUrl: {
    type: String,
    default: '',
  },
  accept: {
    type: Array<string>,
    default: (): string[] => ['image/*'],
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
      class="text-center"
      role="button"
      tabindex="0"
      :aria-label="props.label ? `${props.label}を${props.imgUrl === '' ? '選択' : '変更'}` : 'ファイルをアップロード'"
      flat
      @click="onClick"
      @keydown.enter.prevent="onClick"
      @keydown.space.prevent="onClick"
    >
      <v-card-text>
        <div class="mb-4">
          <v-avatar
            v-if="props.imgUrl === ''"
            size="80"
            :icon="mdiPlus"
          />
          <v-avatar
            v-else
            size="160"
            :image="props.imgUrl"
          />
        </div>
        <input
          ref="inputRef"
          type="file"
          class="d-none"
          :accept="acceptedFiles()"
          @change="onChangeFile"
        >
        <p class="ma-0">
          {{ props.label }}を{{ props.imgUrl === "" ? "選択" : "変更" }}
        </p>
      </v-card-text>
    </v-card>
    <p
      v-show="props.error"
      role="alert"
      class="text-red ma-0"
    >
      {{ props.message }}
    </p>
  </div>
</template>
