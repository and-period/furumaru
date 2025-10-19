<script setup lang="ts">
const route = useRoute();
const { $liffReady } = useNuxtApp();

onMounted(async () => {
  await $liffReady;
});

const loading = computed(() => 'liff.state' in route.query);

definePageMeta({
  layout: 'init',
  public: true,
});
</script>

<template>
  <div>
    <!-- Loading screen -->
    <div
      v-if="loading"
      class="flex h-screen items-center justify-center"
    >
      <div class="flex flex-col items-center gap-3">
        <div class="h-10 w-10 animate-spin rounded-full border-4 border-gray-200 border-t-blue-500" />
        <p class="text-gray-600">
          読み込み中...
        </p>
      </div>
    </div>

    <!-- Content when not loading -->
    <template v-else>
      <div
        class="flex flex-col h-screen"
      >
        <h2 class="font-semibold font-inter text-center w-full mt-12 py-2 px-4">
          ご利用可能な施設一覧
        </h2>
        <div class="grow p-4">
          <div>
            ご利用希望施設を選択してください。
          </div>

          <div class="mt-4">
            <nuxt-link
              to="/tiRfMs1kx57LHF2Mbqt46t"
              class="block shadow rounded-md hover:shadow-md bg-[url(/img/1.jpg)]"
            >
              <div class="p-4 font-semibold text-white tracking-wide">
                検証用キャンプ場
              </div>
              <div class="text-gray-200 text-xs p-4">
                2025/10/01 ~
              </div>
            </nuxt-link>
          </div>
        </div>
        <div
          class="px-4 text-white bg-green py-4 text-xs"
        >
          <div class="text-center flex flex-col gap-2">
            <div>
              株式会社アンドピリオド
            </div>
            <div>
              © andperiod Inc. All rights reserved.
            </div>
            <div>
              <nuxt-link
                to="/legal-notice"
                class="underline hover:text-white/80"
              >
                特定商取引法に基づく表記
              </nuxt-link>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
