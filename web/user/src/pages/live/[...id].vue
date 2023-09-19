<script setup lang="ts">
import { MOCK_LIVE_VIDEO } from '~/constants/mock'

const live = MOCK_LIVE_VIDEO
</script>

<template>
  <div
    class="mx-auto grid max-w-[1440px] grid-flow-col auto-rows-max grid-cols-3 gap-8 text-main xl:px-14"
  >
    <div class="col-span-2">
      <the-live-video-player
        :video-src="live.videoSrc"
        :title="live.title"
        :start-at="live.startAt"
        :is-archive="live.published"
        :marche-name="live.marcheName"
        :description="live.description"
        :address="live.address"
        :cn-img-src="live.cnImgSrc"
        :cn-name="live.cnName"
      />
    </div>

    <div class="col-span-1">
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
  </div>
</template>
