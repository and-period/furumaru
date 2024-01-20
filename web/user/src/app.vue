<script setup lang="ts">
const config = useRuntimeConfig()

useHead({
  script: [
    {
      key: 'newRelic',
      src: `/new-relic/${config.public.ENVIRONMENT}.js`,
      defer: true,
      type: 'text/javascript',
    },
  ],
})
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>
