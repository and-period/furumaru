<script setup lang="ts">
import dayjs from 'dayjs'

interface Props {
  id: string | undefined
  title: string | undefined
  imgSrc: string | undefined
  width: number | undefined
  startAt: number
  endAt: number
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click'): void
}

const emits = defineEmits<Emits>()

const time = computed(() => {
  const startAt = dayjs.unix(props.startAt)
  const current = dayjs()
  const diff = current.diff(startAt)

  const hours = Math.floor(diff / 1000 / 60 / 60)
  const days = Math.floor(diff / 1000 / 60 / 60 / 24)
  const months = Math.floor(diff / 1000 / 60 / 60 / 24 / 30)
  const years = Math.floor(diff / 1000 / 60 / 60 / 24 / 30 / 12)

  if (years > 0) {
    // 1年以上の場合は年を表示する
    return `${years}年前`
  }
  else if (months > 0) {
    // 1ヶ月以上、1年未満の場合は月を表示する
    return `${months}ヶ月前`
  }
  else if (days > 0) {
    // 24時間以上、1ヶ月未満の場合は日付を表示する
    return `${days}日前`
  }
  else if (hours > 0) {
    // 1時間以上、24時間未満の場合は時間を表示する
    return `${hours}時間前`
  }
  else {
    // 1時間未満の場合は「たった今」と表示する
    return 'たった今'
  }
})

const handleClick = () => {
  emits('click')
}
</script>

<template>
  <div
    class="w-full text-main"
    @click="handleClick"
  >
    <div class="w-full">
      <nuxt-img
        provider="cloudFront"
        class="h-[208px] w-full object-cover"
        :src="imgSrc"
        :alt="`archive-${title}-thumbnail`"
        fit="contain"
        height="208px"
      />
    </div>
    <div class="mt-2 flex w-full flex-col gap-2">
      <div class="flex items-center text-sm gap-2">
        <span
          class="rounded border border-main px-1 text-[10px] font-bold tracking-[10%] text-main"
        >
          アーカイブ配信
        </span>
        <span class="text-sm">{{ time }}</span>
      </div>
      <p class="line-clamp-3 break-words text-[14px] tracking-[10%]">
        {{ title }}
      </p>
    </div>
  </div>
</template>
