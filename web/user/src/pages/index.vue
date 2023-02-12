<script lang="ts" setup>

const banners: string[] = [
  '/img/banner.png',
  '/img/banner.png',
  '/img/banner.png'
]
</script>

<template>
  <div>
    <the-carousel :images="banners" />
  </div>
</template>
