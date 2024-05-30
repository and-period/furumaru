<script setup lang="ts">
interface Props {
  pageName: string
  buttonText: string
  name: string
  nameLabel: string
  namePlaceholder: string
  nameErrorMessage: string
  username: string
  usernameLabel: string
  usernamePlaceholder: string
  usernameErrorMessage: string
  cautionText: string
}

const porps = defineProps<Props>()

interface Emits {
  (e: 'update:name', val: string): void
  (e: 'update:username', val: string): void
}

const emits = defineEmits<Emits>()

const nameValue = computed({
  get: () => porps.name,
  set: (val: string) => emits('update:name', val),
})

const usernameValue = computed({
  get: () => porps.username,
  set: (val: string) => emits('update:username', val),
})
</script>

<template>
  <div class="mx-auto block sm:w-[560px]">
    <the-marche-logo class="mb-10" />
    <the-card>
      <the-card-title>
        {{ pageName }}
      </the-card-title>
      <the-card-content>
        <form>
          <the-stack class="gap-y-6">
            <the-text-input
              v-model="nameValue"
              :label="nameLabel"
              :placeholder="namePlaceholder"
              :with-label="false"
              :error-message="nameErrorMessage"
              type="text"
              required
            />
            <the-text-input
              v-model="usernameValue"
              :label="usernameLabel"
              :placeholder="usernamePlaceholder"
              :with-label="false"
              :error-message="usernameErrorMessage"
              type="text"
              required
            />

            <p
              class="my-6 whitespace-pre-wrap"
              v-html="cautionText"
            />

            <the-submit-button>
              {{ buttonText }}
            </the-submit-button>
          </the-stack>
        </form>
      </the-card-content>
    </the-card>
  </div>
</template>
