<script lang="ts" setup>
const streaming = ref({
  status: 'idle',
  ingestEndpoint: '',
  streamKey: ''
})
</script>

<template>
  <v-row>
    <v-col sm="12" md="12" lg="8">
      <v-card>
        <v-card-text>
          <v-container>
            <video id="video" controls />
          </v-container>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col sm="12" md="12" lg="4">
      <v-card>
        <v-card-text>
          <v-select
            v-model="streaming.status"
            label="配信状況"
            :items="statuses"
            item-title="title"
            item-value="value"
            variant="plain"
            readonly
          />
          <v-text-field
            v-model="streaming.ingestEndpoint"
            label="配信エンドポイント"
            readonly
          />
          <v-text-field
            v-model="streaming.streamKey"
            label="ストリームキー"
            readonly
          />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<style scoped>
#video {
  margin-bottom: 1.5em;
  width: 100%;
}
</style>
