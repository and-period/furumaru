<script setup lang="ts">
import { useSeoHead } from '~/hooks/seo'

useSeoHead({
  title: 'すべての商品',
  description: '産地直送のお取り寄せ通販ふるマルの商品一覧ページです。全国の生産者から届くこだわりの特産品をご覧ください。',
  path: '/items',
})
</script>

<template>
  <the-item />
</template>
