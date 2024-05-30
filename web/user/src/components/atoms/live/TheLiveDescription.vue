<script setup lang="ts">
import { datetimeformatterFromUnixtime } from '~/lib/dayjs'

interface Props {
  title: string
  description: string
  isArchive: boolean
  isLiveStreaming: boolean
  startAt: number
  marcheName: string
  coordinatorId: string
  coordinatorName: string
  coordinatorImgSrc: string
  coordinatorAddress: string
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:coordinator', id: string): void
}

const emits = defineEmits<Emits>()

const handleCLickCoordinator = () => {
  emits('click:coordinator', props.coordinatorId)
}

const showDetail = ref<boolean>(false)

const handleClickShowDetailButton = () => {
  showDetail.value = !showDetail.value
}
</script>

<template>
  <div class="mt-2 px-4">
    <div class="flex items-center gap-2">
      <template v-if="isArchive">
        <div
          class="flex max-w-fit items-center justify-center rounded border-2 border-main px-2 font-bold text-main"
        >
          アーカイブ
        </div>
      </template>
      <template v-else-if="isLiveStreaming">
        <div
          class="flex max-w-fit items-center justify-center rounded border-2 border-orange bg-orange px-2 font-bold text-white"
        >
          <div class="mr-2 pt-[2px]">
            <the-live-icon />
          </div>
          <div class="align-middle">
            LIVE
          </div>
        </div>
      </template>

      <div class="text-[14px] tracking-[1.4px] after:content-['〜']">
        {{ datetimeformatterFromUnixtime(startAt) }}
      </div>
    </div>
    <p class="mt-2 line-clamp-1 tracking-[1.6px]">
      {{ title }}
    </p>
    <div class="mt-4 flex items-center gap-2">
      <nuxt-img
        width="40"
        height="40"
        format="webp"
        provider="cloudFront"
        :src="coordinatorImgSrc"
        class="h-10 w-10 rounded-full hover:cursor-pointer"
        :alt="`${coordinatorName}のプロフィール画像`"
        @click="handleCLickCoordinator"
      />
      <div class="text-[12px] tracking-[1.2px]">
        <p class="mb-1">
          {{ marcheName }}/{{ coordinatorAddress }}
        </p>
        <p>
          コーディネーター：
          <span
            class="cursor-pointer hover:underline"
            @click="handleCLickCoordinator"
          >{{ coordinatorName }}</span>
        </p>
      </div>
    </div>

    <div>
      <p
        v-if="showDetail"
        class="mt-6 whitespace-pre-wrap text-[14px] tracking-[1.4px]"
        v-html="description"
      />
      <button
        class="inline-flex w-full items-center justify-center gap-2 text-[12px] tracking-[1.2px]"
        @click="handleClickShowDetailButton"
      >
        <div>
          {{ showDetail ? 'マルシェの詳細を隠す' : 'マルシェの詳細を見る' }}
        </div>
        <div>
          <the-up-arrow-icon v-if="showDetail" />
          <the-down-arrow-icon v-if="!showDetail" />
        </div>
      </button>
    </div>
  </div>
</template>
