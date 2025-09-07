<script setup lang="ts">
import { ref, computed } from 'vue';
import { FmTextInput } from '@furumaru/shared';
import ApiClientFactory from '@/plugins/helper/factory';
import { buildApiErrorMessage } from '@/plugins/helper/error';
import { AuthUserApi } from '@/types/api/facility';
import type { RequestCreateAuthUserRequest } from '@/types/api/facility';
import liff from '@line/liff';

const route = useRoute();
const router = useRouter();

// フォーム状態
const lastName = ref('');
const firstName = ref('');
const lastNameKana = ref('');
const firstNameKana = ref('');
const stayDate = ref(''); // yyyy-MM-dd
const phoneNumber = ref('');

const isSubmitting = ref(false);
const errorMessage = ref<string | null>(null);

// ひらがな検証
const kanaPattern = '[ぁ-んー]+'; // HTMLパターン用
const kanaRegex = /^[ぁ-んー]+$/; // JS検証用
const lastNameKanaValid = computed(() => kanaRegex.test(lastNameKana.value || ''));
const firstNameKanaValid = computed(() => kanaRegex.test(firstNameKana.value || ''));
const lastNameKanaError = computed(() =>
  lastNameKana.value && !lastNameKanaValid.value ? 'ひらがなで入力してください。' : '',
);
const firstNameKanaError = computed(() =>
  firstNameKana.value && !firstNameKanaValid.value ? 'ひらがなで入力してください。' : '',
);

const canSubmit = computed(() => {
  return (
    !!lastName.value
    && !!firstName.value
    && !!lastNameKana.value
    && !!firstNameKana.value
    && lastNameKanaValid.value
    && firstNameKanaValid.value
    && !!stayDate.value
    && !!phoneNumber.value
    && !isSubmitting.value
  );
});

function toUnixSecondsFromDateInput(dateStr: string): number {
  // ローカルタイムの 00:00 として扱い、UNIX秒へ
  const d = new Date(`${dateStr}T00:00:00`);
  return Math.floor(d.getTime() / 1000);
}

// 先頭0始まりの国内番号を国番号+81形式へ正規化（入力は0始まりでOK）
function normalizePhoneNumberToJpIntl(raw: string): string {
  // 全角数字を半角に統一
  const half = raw.replace(/[０-９]/g, ch => String.fromCharCode(ch.charCodeAt(0) - 0xFEE0));
  // 記号や空白・ハイフンを除去（ただし先頭の+は保持）
  const cleaned = half
    .trim()
    .replace(/(?!^)[^0-9]/g, '') // 先頭以外の非数字を除去
    .replace(/^\+?([^0-9]*)(.*)$/, '+$2') // 先頭+がなければ仮に付加（後で整形）
    .replace(/^\+\+/, '+');

  // 先頭+81形式はそのまま
  if (cleaned.startsWith('+81')) return cleaned;
  // 先頭が+でない場合や+0始まりの場合に備えて再取得
  const numeric = cleaned.startsWith('+') ? cleaned.slice(1) : cleaned;
  // 81から始まる（+なし）場合は+81を付与
  if (numeric.startsWith('81')) return `+${numeric}`;
  // 0始まりは+81に変換（先頭の0を除去）
  if (numeric.startsWith('0')) return `+81${numeric.slice(1)}`;
  // その他は+を先頭に付けて返す
  return `+${numeric}`;
}

async function onSubmit() {
  if (!canSubmit.value) return;
  isSubmitting.value = true;
  errorMessage.value = null;
  try {
    const facilityId = String(route.params.facilityId ?? '');
    if (!facilityId) {
      throw new Error('facilityId is missing');
    }

    const liffToken = liff.getIDToken();
    if (!liffToken) {
      throw new Error('LIFFの認証トークン取得に失敗しました。');
    }

    const payload: RequestCreateAuthUserRequest = {
      authToken: liffToken,
      firstname: firstName.value,
      firstnameKana: firstNameKana.value,
      lastname: lastName.value,
      lastnameKana: lastNameKana.value,
      phoneNumber: normalizePhoneNumberToJpIntl(phoneNumber.value),
      lastCheckInAt: toUnixSecondsFromDateInput(stayDate.value),
    };

    const factory = new ApiClientFactory();
    const api = factory.createFacility<AuthUserApi>(AuthUserApi);
    await api.facilitiesFacilityIdUsersPost({
      facilityId,
      requestCreateAuthUserRequest: payload,
    });

    await router.push({ path: `/${facilityId}/checkin` });
  }
  catch (e) {
    console.error('check-in submit failed:', e);
    errorMessage.value = await buildApiErrorMessage(e);
  }
  finally {
    isSubmitting.value = false;
  }
}

// エラーメッセージは共通関数で整形
</script>

<template>
  <div>
    <p class="mt-6 font-inter text-xl text-center w-full text-main font-semibold">
      チェックインの登録
    </p>
    <form @submit.prevent="onSubmit">
      <div class="grid grid-cols-2 gap-2 mt-8 max-w-md mx-auto">
        <div>
          <FmTextInput
            id="last_name"
            v-model="lastName"
            label="名前(姓)"
            name="last_name"
            placeholder="山田"
            required
            class="w-full px-2"
            :disabled="isSubmitting"
          />
        </div>
        <div>
          <FmTextInput
            id="first_name"
            v-model="firstName"
            label="名前(名)"
            name="first_name"
            placeholder="太郎"
            required
            class="w-full px-2"
            :disabled="isSubmitting"
          />
        </div>
      </div>
      <div class="grid grid-cols-2 gap-2 mt-4 max-w-md mx-auto">
        <div>
          <FmTextInput
            id="last_name_kana"
            v-model="lastNameKana"
            label="フリガナ(姓)"
            name="last_name_kana"
            placeholder="やまだ"
            required
            class="w-full px-2"
            :pattern="kanaPattern"
            :error="!!lastNameKanaError"
            :error-message="lastNameKanaError"
            :disabled="isSubmitting"
          />
        </div>
        <div>
          <FmTextInput
            id="first_name_kana"
            v-model="firstNameKana"
            label="フリガナ(名)"
            name="first_name_kana"
            placeholder="たろう"
            required
            class="w-full px-2"
            :pattern="kanaPattern"
            :error="!!firstNameKanaError"
            :error-message="firstNameKanaError"
            :disabled="isSubmitting"
          />
        </div>
      </div>
      <div class="mt-4 max-w-md mx-auto">
        <div>
          <FmTextInput
            id="stay_date"
            v-model="stayDate"
            label="宿泊日"
            name="stay_date"
            required
            class="w-full px-2"
            type="date"
            :disabled="isSubmitting"
          />
        </div>
      </div>
      <div class="mt-4 max-w-md mx-auto">
        <div>
          <FmTextInput
            id="phone_number"
            v-model="phoneNumber"
            label="電話番号"
            name="phone_number"
            placeholder="09012345678"
            required
            class="w-full px-2"
            type="tel"
            pattern="[0-9]*"
            inputmode="tel"
            :disabled="isSubmitting"
          />
        </div>
      </div>
      <div class="mt-8 max-w-md mx-auto flex justify-center px-2">
        <button
          type="submit"
          class="bg-[#F48D26] text-white font-semibold rounded-[10px] px-8 w-full py-3 shadow-md hover:bg-opacity-90 transition-all duration-200 text-lg tracking-wide"
          :disabled="!canSubmit"
        >
          {{ isSubmitting ? '送信中…' : '登録する' }}
        </button>
      </div>
      <div
        v-if="errorMessage"
        class="max-w-md mx-auto px-2 mt-3"
      >
        <p
          role="alert"
          class="text-sm text-red-600"
        >
          {{ errorMessage }}
        </p>
      </div>
    </form>
  </div>
</template>
