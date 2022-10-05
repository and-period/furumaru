<template>
  <div>
    <v-stepper v-model="stepValue" alt-labels flat>
      <v-stepper-header>
        <v-stepper-step step="1" :complete="stepValue > 1">
          基本設定
        </v-stepper-step>
        <v-divider />
        <v-stepper-step step="2" :complete="stepValue > 2">
          タイムテーブル設定
        </v-stepper-step>
        <v-divider />
        <v-stepper-step step="3"> 開催地情報 </v-stepper-step>
      </v-stepper-header>

      <v-stepper-items>
        <v-stepper-content step="1">
          <v-card>
            <v-card-text>
              <v-text-field label="タイトル" />
              <v-textarea label="詳細情報" />
              <p>サムネイル画像</p>
              <the-file-upload-filed />
              <v-switch
                label="テスト"
                hint="有効にするとテスト用のマルシェを作成します。テスト用のマルシェは一般ユーザーには公開されません。"
                persistent-hint
              />
            </v-card-text>
            <v-card-actions>
              <v-spacer />
              <v-btn outlined color="primary" @click="stepValue = 2">
                次へ
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-stepper-content>
        <v-stepper-content step="2">
          <v-card>
            <v-card-text>
              <v-calendar type="day" hide-header />
            </v-card-text>
            <v-card-actions>
              <v-btn outlined color="primary" @click="stepValue = 1">
                戻る
              </v-btn>
              <v-spacer />
              <v-btn outlined color="primary" @click="stepValue = 3">
                次へ
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-stepper-content>
        <v-stepper-content step="3">
          <v-card>
            <v-card-text>
              <v-text-field label="タイトル" readonly />
              <v-textarea label="詳細情報" readonly />
              <v-divider />
              <v-calendar type="day" hide-header />
            </v-card-text>
            <v-card-actions>
              <v-btn outlined color="primary" @click="stepValue = 2">
                戻る
              </v-btn>
              <v-spacer />
              <v-btn outlined color="primary"> 完了 </v-btn>
            </v-card-actions>
          </v-card>
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from '@nuxtjs/composition-api'

export default defineComponent({
  setup() {
    const stepValue = ref<number>(1)

    return {
      stepValue,
    }
  },
})
</script>

<style lang="scss" scoped>
.v-stepper__header {
  box-shadow: none;
}
</style>
