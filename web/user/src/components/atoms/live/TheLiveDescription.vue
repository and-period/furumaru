<script setup lang="ts">
import { datetimeformatterFromUnixtime } from '~/lib/dayjs'
import type { I18n } from '~/types/locales'
import type { Snackbar } from '~/types/props'

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

const i18n = useI18n()

const snackbarItems = ref<Snackbar[]>([])

const dt = (str: keyof I18n['lives']['details']) => {
  return i18n.t(`lives.details.${str}`)
}

const coordinatorThumbnailAlt = computed<string>(() => {
  return i18n.t('lives.details.coordinatorThumbnailAlt', {
    coordinatorName: props.coordinatorName,
  })
})

const props = defineProps<Props>()

interface Emits {
  (e: 'click:coordinator', id: string): void
}

const emits = defineEmits<Emits>()

const handleCLickCoordinator = () => {
  emits('click:coordinator', props.coordinatorId)
}

const handleClickCopyButton = async () => {
  snackbarItems.value.push({
    text: i18n.t('lives.details.linkCopied'),
    isShow: true,
  })
  await navigator.clipboard.writeText(window.location.href)
}

const handleClickXButton = () => {
  const shareXUrl = `https://twitter.com/intent/tweet?text=${props.title}&url=${window.location.href}`
  window.open(shareXUrl, '_blank')
}

const handleClickFacebookButton = () => {
  const shareFacebookUrl = `https://www.facebook.com/sharer/sharer.php?u=${window.location.href}`
  window.open(shareFacebookUrl, '_blank')
}

const showDetail = ref<boolean>(false)

const handleClickShowDetailButton = () => {
  showDetail.value = !showDetail.value
}
</script>

<template>
  <template
    v-for="(snackbarItem, i) in snackbarItems"
    :key="i"
  >
    <the-snackbar
      v-model:is-show="snackbarItem.isShow"
      :text="snackbarItem.text"
    />
  </template>
  <div class="mt-2 px-4">
    <div class="md:flex block md:justify-between">
      <div class="flex items-center gap-2">
        <template v-if="isArchive">
          <div
            class="flex max-w-fit items-center justify-center rounded border-2 border-main px-2 font-bold text-main"
          >
            {{ dt('archivedStreamText') }}
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
      <the-dropdown-with-icon ref="dropdownRef">
        <template #icon>
          <div class="flex md:mt-0 mt-2 text-[14px] tracking-[1.4px] ">
            <the-sns-share-icon />
            <p class="ml-2">
              SHARE
            </p>
          </div>
        </template>
        <template #content>
          <div class="flex flex-col">
            <button
              class="px-4 py-2 text-left flex hover:bg-gray-200 items-center"
              @click="handleClickCopyButton"
            >
              <the-share-icon class="mr-2" />
              {{ dt('linkCopy') }}
            </button>
            <button
              class="px-4 py-2 text-left flex hover:bg-gray-200 items-center"
              @click="handleClickXButton"
            >
              <the-share-x-icon class="mr-2" />
              X (Twitter)
            </button>
            <button
              class="px-4 py-2 text-left flex hover:bg-gray-200 items-center"
              @click="handleClickFacebookButton"
            >
              <the-share-facebook-icon class="mr-2" />
              Facebook
            </button>
          </div>
        </template>
      </the-dropdown-with-icon>
    </div>
    <p class="mt-2 line-clamp-1 tracking-[1.6px]">
      {{ title }}
    </p>
    <div class="mt-4 flex items-center gap-2">
      <nuxt-img
        width="40"
        height="40"
        provider="cloudFront"
        :src="coordinatorImgSrc"
        class="h-10 w-10 rounded-full hover:cursor-pointer"
        :alt="coordinatorThumbnailAlt"
        @click="handleCLickCoordinator"
      />
      <div class="text-[12px] tracking-[1.2px]">
        <p class="mb-1">
          {{ marcheName }}/{{ coordinatorAddress }}
        </p>
        <p>
          {{ dt('coordinatorLabel') }}:
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
          {{ showDetail ? dt('hideMarcheDetailsText') : dt('showMarcheDetailsText') }}
        </div>
        <div>
          <the-up-arrow-icon v-if="showDetail" />
          <the-down-arrow-icon v-if="!showDetail" />
        </div>
      </button>
    </div>
  </div>
</template>
