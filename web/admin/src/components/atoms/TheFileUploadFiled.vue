<script lang="ts" setup>
const props = defineProps({
  text: {
    type: String,
    required: true
  },
  value: {
    type: Object,
    default: null
  }
})

const emit = defineEmits<{
  (e: 'update:value', value: any): void
  (e: 'update:files', file?: FileList): void
}>()

const formData = computed({
  get: (): any => props.value,
  set: (val: any) => emit('update:value', val)
})

const inputRef = ref<HTMLInputElement | null>(null)
const active = ref<boolean>(false)
const files = ref<FileList | null>(null)

watch(files, () => {
  emit('update:files', files.value)
})

const handleInputFileChange = () => {
  if (inputRef.value && inputRef.value.files) {
    files.value = inputRef.value?.files
  }
}

const handleClick = () => {
  if (inputRef.value) {
    inputRef.value.click()
  }
}

const handleDragenter = () => {
  active.value = true
}

const handleDragover = () => {
  active.value = true
}

const handleDragleave = () => {
  active.value = false
}

const handleDrop = (e: DragEvent) => {
  if (e.dataTransfer && inputRef.value) {
    files.value = e.dataTransfer.files
  }
  active.value = false
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-center align-center rounded-lg file_upload_area"
      role="button"
      :class="{ active: active }"
      @click="handleClick"
      @dragenter="handleDragenter"
      @dragleave="handleDragleave"
      @drop.prevent="handleDrop"
      @dragover.prevent="handleDragover"
    >
      <p class="mb-0 text-center">
        クリックまたはドラッグ&amp;ドロップでファイルを追加
        <br>
        <v-icon start>
          mdi-plus
        </v-icon>
        {{ props.text }}
        <input
          ref="inputRef"
          type="file"
          class="d-none"
          accept="image/*"
          multiple
          @change="handleInputFileChange"
        >
      </p>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.file_upload_area {
  border: dashed var(--v-secondary-lighten4);
  height: 100px;
}

.active {
  border: dashed var(--v-primary-darken3);
}
</style>
