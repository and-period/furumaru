<template>
  <div class="stack" :style="spaceStyle"><slot /></div>
</template>

<script lang="ts">
import { computed, defineComponent } from '@nuxtjs/composition-api'

export default defineComponent({
  props: {
    space: {
      type: Number,
      default: 16,
    },
  },
  setup(props) {
    const spaceStyle = computed(() => {
      return {
        '--space': `${props.space}px`,
      }
    })

    return {
      spaceStyle,
    }
  },
})
</script>

<style lang="scss" scoped>
.stack {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}

.stack > * {
  margin-top: 0;
  margin-bottom: 0;
}

.stack > * + * {
  margin-top: var(--space, 16px);
}
</style>
