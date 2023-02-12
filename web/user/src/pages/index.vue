<script lang="ts" setup>

const banners: string[] = [
  '/_nuxt/img/banner.png',
  '/_nuxt/img/banner.png',
  '/_nuxt/img/banner.png'
]
</script>

<template>
  <div>
    <the-carousel :images="banners" />
  </div>
</template>
