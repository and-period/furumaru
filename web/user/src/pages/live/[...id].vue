<script setup lang="ts">
import { MOCK_LIVE_VIDEO } from '~/constants/mock'

const live = MOCK_LIVE_VIDEO

useSeoMeta({
  title: 'ライブ配信',
})
</script>

<template>
  <div
    class="mx-auto grid max-w-[1440px] grid-flow-col auto-rows-max grid-cols-3 gap-8 text-main xl:px-14"
  >
    <div class="col-span-3 lg:col-span-2">
      <the-live-video-player
        :video-src="'https://furumaru-dev.s3.ap-northeast-1.amazonaws.com/promotion/promotion.mp4'"
        :title="live.title"
        :start-at="live.startAt"
        :is-archive="live.published"
        :marche-name="live.marcheName"
        :description="live.description"
        :address="live.address"
        :cn-img-src="live.cnImgSrc"
        :cn-name="live.cnName"
      />
      <the-live-timeline class="mt-4" />
    </div>

    <!-- PC画面のみ表示する右サイドバー -->
    <!--
    <div class="col-span-1 hidden lg:block">
      <div class="flex h-[450px] flex-col bg-white p-4">
        <div class="grow overflow-scroll"></div>
        <div class="relative text-[12px] text-main">
          <input
            type="text"
            placeholder="コメントを投稿する"
            class="w-full rounded-[20px] border border-main py-2 pl-4 pr-[56px]"
          />
          <button class="absolute right-4 top-2 min-w-max">送信</button>
        </div>
      </div>
    </div>
    -->
  </div>
</template>
