<script lang="ts" setup>
import { CreateProductRequestMediaInner } from '~/types/api'

const props = defineProps({
  media: {
    type: Array<CreateProductRequestMediaInner>,
    default: () => {
      return []
    }
  }
})

const emit = defineEmits<{
  (e: 'update:media', media: Array<{ isThumbnail: boolean; url: string }>): void
  (e: 'update:files', files?: FileList): void
}>()

const selected = ref<number>(-1)

selected.value = props.media.findIndex(item => item.isThumbnail)

watch(selected, () => {
  const newVal = props.media.map((item, i) => {
    return {
      ...item,
      isThumbnail: i === selected.value
    }
  })
  emit('update:media', newVal)
})

const handleImageUpload = (files?: FileList) => {
  emit('update:files', files)
}
</script>

<template>
  <div>
    <div>
      <atom-file-upload-filed
        text="商品画像"
        @update:files="handleImageUpload"
      />
    </div>
    <v-radio-group v-model="selected">
      <div
        v-for="(img, i) in props.media"
        :key="i"
        class="d-flex flex-row align-center my-2"
      >
        <v-radio :value="i" />
        <v-img :src="img.url" max-width="200" class="mx-4" />
        <p class="mb-0 img-url">
          {{ img.url }}
        </p>
      </div>
    </v-radio-group>
    <p>※ check された商品画像がサムネイルになります</p>
  </div>
</template>

<style lang="scss" scoped>
.img-url {
  word-break: break-all;
}
</style>
