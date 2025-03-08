<script setup lang="ts">
const { locale } = useI18n()
const currentLocale = computed(() => locale.value)

const { $md } = useNuxtApp()
const content = ref<string>('')

const renderedContent = computed<string>(() => {
  return $md.render(content.value)
})

const fetchContent = async () => {
  const response = await $fetch(`/_content/termsOfUse_${currentLocale.value}.md`, {
    method: 'GET',
    headers: {
      'content-type': 'text/markdown',
    },
    redirect: 'follow',
  })
  content.value = response as string
}

onMounted(() => {
  fetchContent()
})

useSeoMeta({
  title: '利用規約',
})
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <div
      class="prose max-w-none p-4"
      v-html="renderedContent"
    />
  </div>
</template>
