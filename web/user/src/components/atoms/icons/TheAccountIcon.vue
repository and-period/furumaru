<script setup lang="ts">
interface Props {
  color?: 'main' | 'white'
}

const props = withDefaults(defineProps<Props>(), {
  color: 'main',
})

const bgFillColor = computed(() => {
  return props.color === 'main' ? '#604C3F' : '#D9D9D9'
})

const pathFillColor = computed(() => {
  return props.color === 'main' ? '#F9F6EA' : '#FFFFFF'
})
</script>

<template>
  <svg
    width="24"
    height="24"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M16 32C24.8366 32 32 24.8366 32 16C32 7.16344 24.8366 0 16 0C7.16344 0 0 7.16344 0 16C0 24.8366 7.16344 32 16 32Z"
      :fill="bgFillColor"
    />
    <path
      d="M16.0018 16.817C18.6652 16.817 20.8244 14.6579 20.8244 11.9945C20.8244 9.33102 18.6652 7.17188 16.0018 7.17188C13.3383 7.17188 11.1792 9.33102 11.1792 11.9945C11.1792 14.6579 13.3383 16.817 16.0018 16.817Z"
      :fill="pathFillColor"
    />
    <path
      d="M16.0018 18.1714C13.1936 18.1714 9.49639 19.4026 8.32514 22.0485C8.0923 22.5777 8.49448 23.1704 9.06952 23.1704H22.9305C23.509 23.1704 23.9112 22.5777 23.6748 22.0485C22.5036 19.4026 18.8064 18.1714 15.9982 18.1714H16.0018Z"
      :fill="pathFillColor"
    />
  </svg>
</template>
