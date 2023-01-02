<template>
  <v-btn class="menu-btn" variant="outlined" @click="handleClick">
    <v-icon class="mr-1" color="grey">mdi-menu</v-icon>
    <v-img v-if="props.imgSrc" class="profile-img" alt="profile img" :src="props.imgSrc" width="24" />
    <v-icon v-else color="grey">mdi-account-circle</v-icon>
  </v-btn>
</template>

<script lang="ts" setup>
const props = defineProps({
  imgSrc: {
    type: String,
    required: false,
    default: null,
  },
})

const emits = defineEmits<{
  (name: 'click', e: Event): void
}>()

const handleClick = (e: Event) => {
  emits('click', e)
}
</script>

<style lang="scss" scoped>
.menu-btn {
  border: 1px solid #949494;
  box-sizing: border-box;
  border-radius: 130px;
}

.profile-img {
  border-radius: 50%;
}
</style>
