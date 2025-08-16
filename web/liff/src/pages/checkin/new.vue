<script setup lang="ts">
import FmTextInput from '../../../shared/src/components/FmTextInput.vue';
</script>

<template>
  <div>
    <p class="mt-6 font-inter text-xl text-center w-full text-main font-semibold">
      チェックインの登録
    </p>
    <div>
      <div class="grid grid-cols-2 gap-2 mt-8 max-w-md mx-auto">
        <div>
          <label
            for="last_name"
            class="inline-block text-xs px-2"
          >名前(姓)</label>
          <FmTextInput
            id="last_name"
            name="last_name"
            placeholder="山田"
            required
            class="w-full px-2"
          />
        </div>
        <div>
          <label
            for="first_name"
            class="inline-block text-xs px-2"
          >名前(名)</label>
          <FmTextInput
            id="first_name"
            name="first_name"
            placeholder="太郎"
            required
            class="w-full px-2"
          />
        </div>
      </div>
      <div class="grid grid-cols-2 gap-2 mt-4 max-w-md mx-auto">
        <div>
          <label
            for="last_name_kana"
            class="inline-block text-xs px-2"
          >フリガナ(姓)</label>
          <FmTextInput
            id="last_name_kana"
            name="last_name_kana"
            placeholder="ヤマダ"
            required
            class="w-full px-2"
          />
        </div>
        <div>
          <label
            for="first_name_kana"
            class="inline-block text-xs px-2"
          >フリガナ(名)</label>
          <FmTextInput
            id="first_name_kana"
            name="first_name_kana"
            placeholder="タロウ"
            required
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
            name="stay_date"
            type="date"
            required
            class="w-full px-2 custom-input"
          >
        </div>
      </div>
      <div class="mt-4 max-w-md mx-auto">
        <div>
          <label
            for="phone_number"
            class="inline-block text-xs px-2"
          >電話番号</label>
          <FmTextInput
            id="phone_number"
            name="phone_number"
            placeholder="09012345678"
            required
            class="w-full px-2"
            type="tel"
            pattern="[0-9]*"
            inputmode="tel"
          />
        </div>
      </div>
      <div class="mt-8 max-w-md mx-auto flex justify-center px-2">
        <button
          type="submit"
          class="bg-[#F48D26] text-white font-semibold rounded-[10px] px-8 w-full py-3 shadow-md hover:bg-opacity-90 transition-all duration-200 text-lg tracking-wide"
        >
          登録する
        </button>
      </div>
    </div>
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
