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
        <br />
        <v-icon left>mdi-plus</v-icon>
        {{ text }}
        <input
          ref="inputRef"
          type="file"
          class="d-none"
          multiple
          @change="handleInputFileChange"
        />
      </p>
    </div>
    <v-list v-if="files">
      <div
        v-for="(file, i) in files"
        :key="i"
        class="d-flex flex-row align-center"
      >
        <v-checkbox />
        <img :src="createPreviewURL(file)" width="200" class="mx-4" />
        <p class="mb-0">{{ file.name }}</p>
      </div>
    </v-list>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, PropType } from '@vue/composition-api'

export default defineComponent({
  props: {
    text: {
      type: String,
      required: true,
    },
    value: {
      type: Object as PropType<any>,
      default: null,
    },
  },

  setup(props, { emit }) {
    const formData = computed({
      get: (): any => props.value,
      set: (val: any) => emit('update:value', val),
    })

    const inputRef = ref<HTMLInputElement | null>(null)
    const active = ref<boolean>(false)
    const files = ref<FileList | null>(null)

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

    const createPreviewURL = (file: File) => {
      return URL.createObjectURL(file)
    }

    return {
      files,
      formData,
      inputRef,
      active,
      handleInputFileChange,
      handleClick,
      handleDragenter,
      handleDragleave,
      handleDrop,
      handleDragover,
      createPreviewURL,
    }
  },
})
</script>

<style lang="scss" scoped>
.file_upload_area {
  border: dashed var(--v-secondary-lighten4);
  height: 100px;
}

.active {
  border: dashed var(--v-primary-darken3);
}
</style>
