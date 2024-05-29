<script setup>
definePageMeta({
  name: 'EmptyLayout',
  layout: 'empty',
})

const props = defineProps({
  error: {
    type: Object,
    default: null,
  },
})

const pageNotFound = ref<string>('404 Not Found')
const otherError = ref<string>('An error occurred')

useHead({
  title: props.error?.statusCode === 404 ? pageNotFound : otherError,
})
</script>

<template>
  <v-app dark>
    <h1 v-if="error.statusCode === 404">
      {{ pageNotFound }}
    </h1>
    <h1 v-else>
      {{ otherError }}
    </h1>
    <nuxt-link to="/">
      Home page
    </nuxt-link>
  </v-app>
</template>

<style scoped>
h1 {
  font-size: 20px;
}
</style>
