<script setup lang='ts'>
import dayjs from 'dayjs'
import VersionInfo from '../../version.json'

const builtAt = computed(() => dayjs(VersionInfo.buildDate))
</script>

<template>
  <templates-setting-version
    :version="VersionInfo.version"
    :built-at="builtAt"
  />
</template>
