<script setup lang="ts">
import { datetimeformatterFromUnixtime } from '~/lib/dayjs'

interface Props {
  title: string
  description: string
  isArchive: boolean
  isLiveStreaming: boolean
  startAt: number
  marcheName: string
  cordinatorId: string
  cordinatorName: string
  cordinatorImgSrc: string
  cordinatorAddress: string
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:cordinator', id: string): void
}

const emits = defineEmits<Emits>()

const handleCLickCorodinator = () => {
  emits('click:cordinator', props.cordinatorId)
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
          <div class="align-middle">LIVE</div>
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
      <img
        :src="cordinatorImgSrc"
        class="h-10 w-10 rounded-full hover:cursor-pointer"
        :alt="`${cordinatorName}のプロフィール画像`"
        @click="handleCLickCorodinator"
      />
      <div class="text-[12px] tracking-[1.2px]">
        <p class="mb-1">{{ marcheName }}/{{ cordinatorAddress }}</p>
        <p>
          コーディネーター：
          <span
            class="cursor-pointer hover:underline"
            @click="handleCLickCorodinator"
            >{{ cordinatorName }}</span
          >
        </p>
      </div>
    </div>

    <div>
      <p
        v-if="showDetail"
        class="mt-6 whitespace-pre-wrap text-[14px] tracking-[1.4px]"
        v-html="description"
      ></p>
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
