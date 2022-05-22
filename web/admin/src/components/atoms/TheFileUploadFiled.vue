<template>
  <label>
    <div class="d-flex justify-center align-center rounded-lg file_upload_area">
      <p class="mb-0">
        <v-icon left>mdi-plus</v-icon>
        {{ text }}
        <input type="file" class="d-none" />
      </p>
    </div>
  </label>
</template>

<script lang="ts">
import { defineComponent } from '@vue/composition-api'

export default defineComponent({
  props: {
    text: {
      type: String,
      required: true,
    },
  },
})
</script>

<style lang="scss" scoped>
.file_upload_area {
  cursor: pointer;
  border: dashed rgba(0, 0, 0, 0.87);
  height: 100px;
}
</style>
