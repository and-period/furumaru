<script setup lang="ts">
import { useSortable } from '@vueuse/integrations/useSortable'
import { mdiDrag } from '@mdi/js'
import { getResizedImages } from '~/lib/helpers'
import type { Experience, ExperienceMedia } from '~/types/api/v1'

const model = defineModel<any>()

interface Props {
  experiences: Experience[]
}

defineProps<Props>()

const sortableRef = ref<HTMLElement | null>(null)
useSortable(sortableRef, model)

const getExperienceThumbnailUrl = (experience: Experience): string => {
  const thumbnail = experience.media?.find((media: ExperienceMedia) => {
    return media.isThumbnail
  })
  return thumbnail ? thumbnail.url : ''
}

const getExperienceThumbnails = (experience: Experience): string => {
  const thumbnail = experience.media?.find((media: ExperienceMedia) => {
    return media.isThumbnail
  })
  return thumbnail ? getResizedImages(thumbnail.url) : ''
}
</script>

<template>
  <v-table>
    <thead>
      <tr>
        <th width="5" />
        <th />
        <th>タイトル</th>
        <th>大人料金</th>
      </tr>
    </thead>

    <tbody ref="sortableRef">
      <tr
        v-for="experience in experiences"
        :key="experience.id"
      >
        <td>
          <v-icon :icon="mdiDrag" />
        </td>
        <td>
          <v-img
            aspect-ratio="1/1"
            :max-height="56"
            :max-width="80"
            :src="getExperienceThumbnailUrl(experience)"
            :srcset="getExperienceThumbnails(experience)"
            :alt="experience.title || '体験画像'"
          />
        </td>
        <td>{{ experience.title }}</td>
        <td>{{ experience.priceAdult.toLocaleString() }}</td>
      </tr>
    </tbody>
  </v-table>
</template>
