<template>
  <div>
    <div>
      <the-file-upload-filed
        text="商品画像"
        @update:files="handleImageUpload"
      />
    </div>
    <v-radio-group v-model="selected">
      <div
        v-for="(img, i) in media"
        :key="i"
        class="d-flex flex-row align-center my-2"
      >
        <v-radio :value="i" />
        <v-img :src="img.url" max-width="200" class="mx-4" />
        <p class="mb-0 img-url">{{ img.url }}</p>
      </div>
    </v-radio-group>
    <p>※ check された商品画像がサムネイルになります</p>
  </div>
</template>

<script lang="ts">
import { defineComponent, PropType, ref, watch } from '@vue/composition-api'

import { CreateProductRequestMediaInner } from '~/types/api'

export default defineComponent({
  props: {
    media: {
      type: Array as PropType<CreateProductRequestMediaInner[]>,
      default: () => {
        return []
      },
    },
  },

  setup(props, { emit }) {
    const selected = ref<number>(-1)

    selected.value = props.media.findIndex((item) => item.isThumbnail)

    watch(selected, () => {
      const newVal = props.media.map((item, i) => {
        return {
          ...item,
          isThumbnail: i === selected.value,
        }
      })
      emit('update:media', newVal)
    })

    const handleImageUpload = (files: FileList) => {
      emit('update:files', files)
    }

    return { selected, handleImageUpload }
  },
})
</script>

<style lang="scss" scoped>
.img-url {
  word-break: break-all;
}
</style>
