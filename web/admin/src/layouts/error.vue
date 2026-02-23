<script setup>
import { mdiFileQuestion, mdiAlertCircle } from '@mdi/js'

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

const is404 = computed(() => props.error?.statusCode === 404)

useHead({
  title: is404.value ? 'ページが見つかりません' : 'エラーが発生しました',
})
</script>

<template>
  <v-app>
    <v-main
      class="d-flex align-center justify-center"
      style="min-height: 100vh;"
    >
      <v-container class="d-flex justify-center">
        <v-card
          max-width="480"
          class="text-center pa-8"
          elevation="2"
          rounded="lg"
        >
          <div class="mb-4">
            <v-icon
              :icon="is404 ? mdiFileQuestion : mdiAlertCircle"
              :color="is404 ? 'warning' : 'error'"
              size="80"
            />
          </div>
          <div class="text-h3 font-weight-bold mb-4 text-medium-emphasis">
            {{ error?.statusCode || 500 }}
          </div>
          <v-card-title class="text-h6 pb-2">
            {{ is404 ? 'ページが見つかりません' : 'エラーが発生しました' }}
          </v-card-title>
          <v-card-text class="text-body-1 text-medium-emphasis pb-6">
            {{
              is404
                ? 'お探しのページは移動または削除された可能性があります。'
                : 'サーバーでエラーが発生しました。しばらくしてから再度お試しください。'
            }}
          </v-card-text>
          <v-card-actions class="justify-center">
            <v-btn
              to="/"
              color="primary"
              variant="elevated"
              size="large"
            >
              ホームに戻る
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-container>
    </v-main>
  </v-app>
</template>
