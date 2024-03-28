<script setup lang="ts">
import { useAuthStore } from '~/store/auth'

const {} = useAuthStore()
</script>

<template>
  <div class="container mx-auto">
    <the-account-edit-card title="プロフィール写真の変更">
    </the-account-edit-card>
  </div>
</template>
