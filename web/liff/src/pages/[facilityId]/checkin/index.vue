<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { storeToRefs } from 'pinia';
import { FmTextInput } from '@furumaru/shared';
import { useUserStore } from '@/stores/user';
import type { UpdateAuthUserRequest } from '@/types/api/facility';

const userStore = useUserStore();
const { isLoading, error, profile } = storeToRefs(userStore);

const route = useRoute();
const router = useRouter();

const facilityId = computed<string>(() => String(route.params.facilityId || ''));

// フォームデータ
const lastName = ref('');
const firstName = ref('');
const lastNameKana = ref('');
const firstNameKana = ref('');
const phoneNumber = ref('');

// エラーメッセージ
const formError = ref('');

// ユーザー情報を取得してフォームに設定
onMounted(async () => {
  try {
    await userStore.fetchMe(facilityId.value);
    if (profile.value) {
      lastName.value = profile.value.lastname || '';
      firstName.value = profile.value.firstname || '';
      lastNameKana.value = profile.value.lastnameKana || '';
      firstNameKana.value = profile.value.firstnameKana || '';
      phoneNumber.value = profile.value.phoneNumber?.replace('+81', '0') || '';
    }
  }
  catch {
    // エラーは userStore.error に格納される
  }
});

// フォームのバリデーション
const validateForm = (): boolean => {
  if (!lastName.value.trim()) {
    formError.value = '姓を入力してください';
    return false;
  }
  if (!firstName.value.trim()) {
    formError.value = '名を入力してください';
    return false;
  }
  if (!lastNameKana.value.trim()) {
    formError.value = '姓（かな）を入力してください';
    return false;
  }
  if (!firstNameKana.value.trim()) {
    formError.value = '名（かな）を入力してください';
    return false;
  }
  if (!phoneNumber.value.trim()) {
    formError.value = '電話番号を入力してください';
    return false;
  }

  // 電話番号の簡単な形式チェック
  const phoneRegex = /^[0-9-]+$/;
  if (!phoneRegex.test(phoneNumber.value)) {
    formError.value = '電話番号の形式が正しくありません';
    return false;
  }

  return true;
};

// 更新処理
const handleUpdate = async () => {
  formError.value = '';

  if (!validateForm()) {
    return;
  }

  try {
    // 電話番号を+81形式に変換（先頭の0を+81に置換）
    const formattedPhoneNumber = phoneNumber.value.startsWith('0')
      ? phoneNumber.value.replace(/^0/, '+81')
      : phoneNumber.value;

    const updateData: UpdateAuthUserRequest = {
      lastname: lastName.value.trim(),
      firstname: firstName.value.trim(),
      lastnameKana: lastNameKana.value.trim(),
      firstnameKana: firstNameKana.value.trim(),
      phoneNumber: formattedPhoneNumber,
      // lastCheckInAtは現在の値を保持
      lastCheckInAt: profile.value?.lastCheckInAt || Math.floor(Date.now() / 1000),
    };

    await userStore.updateMe(facilityId.value, updateData);

    // マイページに戻る
    await router.push(`/${facilityId.value}/mypage`);
  }
  catch (e) {
    console.error('Update failed:', e);
    formError.value = 'ユーザー情報の更新に失敗しました';
  }
};

// キャンセル処理
const handleCancel = () => {
  router.back();
};
</script>

<template>
  <div>
    <p class="mt-6 font-inter text-xl text-center w-full text-main font-semibold">
      ユーザー情報編集
    </p>

    <!-- ローディング表示 -->
    <div
      v-if="isLoading"
      class="mt-8 text-center text-gray-500"
    >
      読み込み中…
    </div>

    <!-- エラー表示 -->
    <div
      v-else-if="error"
      class="mt-6 text-center text-red-600 text-sm"
    >
      {{ error }}
    </div>

    <!-- フォームエラー表示 -->
    <div
      v-if="formError"
      class="mt-6 max-w-md mx-auto text-center text-red-600 text-sm"
    >
      {{ formError }}
    </div>

    <!-- フォーム表示 -->
    <div v-else-if="profile">
      <div class="grid grid-cols-2 gap-2 mt-8 max-w-md mx-auto">
        <div>
          <label class="inline-block text-xs px-2">名前(姓)</label>
          <FmTextInput
            v-model="lastName"
            name="lastName"
            class="w-full px-2"
          />
        </div>
        <div>
          <label class="inline-block text-xs px-2">名前(名)</label>
          <FmTextInput
            v-model="firstName"
            name="firstName"
            class="w-full px-2"
          />
        </div>
      </div>
      <div class="grid grid-cols-2 gap-2 mt-4 max-w-md mx-auto">
        <div>
          <label class="inline-block text-xs px-2">ふりがな(姓)</label>
          <FmTextInput
            v-model="lastNameKana"
            name="lastNameKana"
            class="w-full px-2"
          />
        </div>
        <div>
          <label class="inline-block text-xs px-2">ふりがな(名)</label>
          <FmTextInput
            v-model="firstNameKana"
            name="firstNameKana"
            class="w-full px-2"
          />
        </div>
      </div>
      <div class="mt-4 max-w-md mx-auto">
        <div>
          <label class="inline-block text-xs px-2">電話番号</label>
          <FmTextInput
            v-model="phoneNumber"
            name="phoneNumber"
            class="w-full px-2"
          />
        </div>
      </div>

      <!-- ボタン群（常に編集可能） -->
      <div class="mt-8 max-w-md mx-auto px-2 space-y-3">
        <button
          type="button"
          class="bg-[#F48D26] text-white font-semibold rounded-[10px] px-8 w-full py-3 shadow-md hover:bg-opacity-90 transition-all duration-200 text-lg tracking-wide"
          :disabled="isLoading"
          @click="handleUpdate"
        >
          <span v-if="isLoading">更新中...</span>
          <span v-else>保存する</span>
        </button>
        <button
          type="button"
          class="bg-gray-500 text-white font-semibold rounded-[10px] px-8 w-full py-3 shadow-md hover:bg-opacity-90 transition-all duration-200 text-lg tracking-wide"
          :disabled="isLoading"
          @click="handleCancel"
        >
          キャンセル
        </button>
      </div>
    </div>
  </div>
</template>
