<script setup lang='ts'>
import VersionInfo from '../../version.json'
import dayjs from 'dayjs'

const builtAt = computed(() => dayjs(VersionInfo.buildDate))
</script>

<template>
  <templates-setting-version
    :version="VersionInfo.version"
    :built-at="builtAt"
  />
</template>
