<script setup lang="ts">
import { mdiCheck } from '@mdi/js'

const router = useRouter()

const handleClickBackToTopButton = () => {
  router.push('/')
}
</script>

<template>
  <v-card>
    <v-card-text>
      <div class="text-center">
        <v-icon
          :icon="mdiCheck"
          color="primary"
          size="64"
        />
        <p class="text-h5">
          YouTube 連携が完了しました
        </p>
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn
        variant="outlined"
        block
        @click="handleClickBackToTopButton"
      >
        トップに戻る
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
