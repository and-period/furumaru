<script setup lang="ts">
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const ct = (str: keyof I18n['base']['contact']) => {
  return i18n.t(`base.contact.${str}`)
}

useSeoMeta({
  title: 'お問い合わせ',
})
</script>

<template>
  <div class="my-24 flex flex-col items-center px-4 text-main">
    <div>
      <h1 class="text-4xl">
        {{ ct('contactTitle') }}
      </h1>
    </div>
    <div>
      <p class="mt-8">
        {{ ct('description') }}
      </p>
      <div class="mt-8 flex justify-center underline">
        <a href="https://lin.ee/49SOeUC"><img
          :src="ct('lineAddFriendImageUrl')"
          :alt="ct('lineAddFriendImageAlt')"
          height="36"
        ></a>
      </div>
      <p class="mt-8">
        {{ ct('contactNotice') }}
      </p>
    </div>
  </div>
</template>
