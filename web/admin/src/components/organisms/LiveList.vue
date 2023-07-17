<script lang="ts" setup>
import dayjs from 'dayjs';
import { CreateLiveRequest, Live, Producer, Product, UpdateLiveRequest } from '~/types/api';

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  createFormData: {
    type: Object as PropType<CreateLiveRequest>,
    default: (): CreateLiveRequest => ({
      producerId: '',
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  updateFormData: {
    type: Object as PropType<UpdateLiveRequest>,
    default: (): UpdateLiveRequest => ({
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  lives: {
    type: Array<Live>,
    default: () => []
  },
  producers: {
    type: Array<Producer>,
    default: () => []
  },
  products: {
    type: Array<Product>,
    default: () => []
  }
})
</script>

<template>
  <v-card v-for="(live, i) in props.lives" :key="i" variant="outlined" color="primary">
    aaa
  </v-card>
</template>
