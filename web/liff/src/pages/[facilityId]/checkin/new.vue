<script setup lang="ts">
import { FmTextInput } from '@furumaru/shared';
</script>

<template>
  <div>
    <p class="mt-6 font-inter text-xl text-center w-full text-main font-semibold">
      チェックインの登録
    </p>
    <div>
      <div class="grid grid-cols-2 gap-2 mt-8 max-w-md mx-auto">
        <div>
          <FmTextInput
            id="last_name"
            label="名前(姓)"
            name="last_name"
            placeholder="山田"
            required
            class="w-full px-2"
            :disabled="false"
          />
        </div>
        <div>
          <FmTextInput
            id="first_name"
            label="名前(名)"
            name="first_name"
            placeholder="太郎"
            required
            class="w-full px-2"
            :disabled="false"
          />
        </div>
      </div>
      <div class="grid grid-cols-2 gap-2 mt-4 max-w-md mx-auto">
        <div>
          <FmTextInput
            id="last_name_kana"
            label="フリガナ(姓)"
            name="last_name_kana"
            placeholder="ヤマダ"
            required
            class="w-full px-2"
            :disabled="false"
          />
        </div>
        <div>
          <FmTextInput
            id="first_name_kana"
            label="フリガナ(名)"
            name="first_name_kana"
            placeholder="タロウ"
            required
            class="w-full px-2"
            :disabled="false"
          />
        </div>
      </div>
      <div class="mt-4 max-w-md mx-auto">
        <div>
          <FmTextInput
            id="stay_date"
            label="宿泊日"
            name="stay_date"
            required
            class="w-full px-2"
            type="date"
            :disabled="false"
          />
        </div>
      </div>
      <div class="mt-4 max-w-md mx-auto">
        <div>
          <FmTextInput
            id="phone_number"
            label="電話番号"
            name="phone_number"
            placeholder="09012345678"
            required
            class="w-full px-2"
            type="tel"
            pattern="[0-9]*"
            inputmode="tel"
            :disabled="false"
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
