<template>
  <div>
    <p class="text-h6">生産者登録</p>
    <v-card elevation="0">
      <v-card-text>
        <v-text-field label="店舗名" />
        <div class="mb-2">
          <the-file-upload-filed text="生産者画像" />
        </div>
        <div class="d-flex">
          <v-text-field class="mr-4" label="生産者名:姓" />
          <v-text-field label="生産者名:名" />
        </div>
        <div class="d-flex">
          <v-text-field class="mr-4" label="生産者名:姓（ふりがな）" />
          <v-text-field label="生産者名:名（ふりがな）" />
        </div>
        <v-text-field label="連絡先（Email）" />
        <v-text-field label="連絡先（電話番号）" />
        <v-text-field label="郵便番号" />
        <v-text-field label="都道府県" />
        <v-text-field label="市区町村" />
        <v-text-field label="番地" />
        <v-text-field label="アパート名・部屋番号（任意）" />
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary">登録</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>
