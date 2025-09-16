<script setup lang="ts">
import { onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useUserStore } from '@/stores/user';

const userStore = useUserStore();
const { isLoading, error, profile, phoneNumber, lastCheckInAt } = storeToRefs(userStore);

const route = useRoute();
const facilityId = computed<string>(() => String(route.params.facilityId || ''));

onMounted(async () => {
  try {
    await userStore.fetchMe(facilityId.value);
  }
  catch {
    // エラーは store.error に格納される
  }
});
</script>

<template>
  <div class="px-4 pb-20 max-w-md mx-auto w-full">
    <h2 class="mt-6 font-semibold font-inter text-center w-full text-xl">
      マイページ
    </h2>

    <div
      v-if="isLoading"
      class="mt-8 text-center text-gray-500"
    >
      読み込み中…
    </div>
    <div
      v-else-if="error"
      class="mt-6 text-center text-red-600 text-sm"
    >
      {{ error }}
    </div>

    <div
      v-else-if="profile"
      class="mt-6 space-y-4"
    >
      <div>
        <div class="mt-0.5 text-xs text-gray-500">
          {{ profile.lastnameKana }} {{ profile.firstnameKana }}
        </div>
        <div class="mt-1 font-medium">
          {{ profile.lastname }}{{ profile.firstname }}
        </div>
      </div>

      <div>
        <div class="text-sm text-gray-500">
          メールアドレス
        </div>
        <div class="mt-1 font-medium break-all">
          {{ profile.email }}
        </div>
      </div>

      <div>
        <div class="text-sm text-gray-500">
          電話番号
        </div>
        <div class="mt-1 font-medium">
          {{ phoneNumber }}
        </div>
      </div>

      <div>
        <div class="text-sm text-gray-500">
          最新のチェックイン
        </div>
        <div class="mt-1 font-medium">
          {{ lastCheckInAt }}
        </div>
      </div>
    </div>
  </div>
</template>
