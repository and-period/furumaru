<script setup lang="ts">
import { ref } from 'vue';
import FmTextInput from '../../../shared/src/components/FmTextInput.vue';

// 仮データ: 実際はpropsやstore等で受け取る想定
const lastName = ref('山田');
const firstName = ref('太郎');
const lastNameKana = ref('ヤマダ');
const firstNameKana = ref('タロウ');
const phoneNumber = ref('09012345678');
const stayDate = ref('2025-08-18');
</script>

<template>
  <div>
    <p class="mt-6 font-inter text-xl text-center w-full text-main font-semibold">
      チェックイン内容の確認
    </p>
    <div class="grid grid-cols-2 gap-2 mt-8 max-w-md mx-auto">
      <div>
        <label class="inline-block text-xs px-2">名前(姓)</label>
        <FmTextInput
          :model-value="lastName"
          :disabled="true"
          class="w-full px-2"
        />
      </div>
      <div>
        <label class="inline-block text-xs px-2">名前(名)</label>
        <FmTextInput
          :model-value="firstName"
          :disabled="true"
          class="w-full px-2"
        />
      </div>
    </div>
    <div class="grid grid-cols-2 gap-2 mt-4 max-w-md mx-auto">
      <div>
        <label class="inline-block text-xs px-2">フリガナ(姓)</label>
        <FmTextInput
          :model-value="lastNameKana"
          :disabled="true"
          class="w-full px-2"
        />
      </div>
      <div>
        <label class="inline-block text-xs px-2">フリガナ(名)</label>
        <FmTextInput
          :model-value="firstNameKana"
          :disabled="true"
          class="w-full px-2"
        />
      </div>
    </div>
    <div class="mt-4 max-w-md mx-auto">
      <div class="px-2">
        <label
          for="stay_date"
          class="inline-block text-xs px-2"
        >宿泊日</label>
        <input
          id="stay_date"
          v-model="stayDate"
          :disabled="false"
          name="stay_date"
          type="date"
          required
          class="w-full px-2 custom-input"
        >
      </div>
    </div>
    <div class="mt-4 max-w-md mx-auto">
      <div>
        <label class="inline-block text-xs px-2">電話番号</label>
        <FmTextInput
          :model-value="phoneNumber"
          :disabled="true"
          class="w-full px-2"
        />
      </div>
    </div>
  </div>
  <div class="mt-8 max-w-md mx-auto flex justify-center px-2">
    <button
      type="submit"
      class="bg-[#F48D26] text-white font-semibold rounded-[10px] px-8 w-full py-3 shadow-md hover:bg-opacity-90 transition-all duration-200 text-lg tracking-wide"
    >
      編集する
    </button>
  </div>
</template>

<style scoped>
.custom-input {
  border-bottom: 1.5px solid #604c3f;
}
.custom-input:focus {
  outline: none;
  box-shadow: none;
}
</style>
