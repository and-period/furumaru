<script lang="ts" setup>
const deliveries = [
  {
    title: '受注待ち',
    value: 'xx',
  },
  {
    title: '発送対応中',
    value: 'yy',
  },
  {
    title: '発送連絡待ち',
    value: 'xx',
  },
]
</script>

<template>
  <v-container>
    <v-row>
      <v-col cols="4">
        <v-card class="py-2">
          <v-card-title class="mb-2">
            注文状況
          </v-card-title>
          <v-card-text
            v-for="summary in deliveries"
            :key="summary.title"
            class="d-flex py-2"
          >
            <div class="me-auto">
              {{ summary.title }}
            </div>
            <div>{{ summary.value }} 件</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
