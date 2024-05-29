<script lang="ts" setup>
interface Props {
  code: string
  buttonText: string
}

interface Emits {
  (e: 'update:code', val: string): void
  (e: 'submit'): void
}

const props = defineProps<Props>()
const emits = defineEmits<Emits>()

const codeValue = computed({
  get: () => props.code,
  set: (val: string) => emits('update:code', val),
})

const handleSubmit = () => {
  emits('submit')
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <the-stack>
      <the-pincode-input
        v-model="codeValue"
        class="mb-4"
      />
      <the-submit-button>{{ buttonText }}</the-submit-button>
    </the-stack>
  </form>
</template>
