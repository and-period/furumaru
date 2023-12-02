<script lang="ts" setup>
interface Props {
  pageName: string
  errorMessage: string
  code: string
  buttonText: string
  message: string
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
  <div class="mx-auto block sm:min-w-[560px]">
    <the-marche-logo class="mb-10" />
    <the-card>
      <the-card-title>
        {{ pageName }}
      </the-card-title>
      <the-card-content>
        <the-alert v-show="errorMessage" class="mb-2">
          {{ errorMessage }}
        </the-alert>

        <the-stack>
          <p>{{ message }}</p>
          <the-verify-code-form
            v-model:code="codeValue"
            :button-text="buttonText"
            @submit="handleSubmit"
          />
        </the-stack>
      </the-card-content>
    </the-card>
  </div>
</template>
