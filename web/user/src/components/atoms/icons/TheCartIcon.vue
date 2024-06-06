<script lang="ts" setup>
interface Props {
  id: string
  fill?: string
}

withDefaults(defineProps<Props>(), {
  fill: 'white',
})
</script>

<template>
  <svg
    width="24"
    height="24"
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <g :clip-path="`url(#${id})`">
      <path
        d="M1.79559 8.99951C0.740766 8.99951 -0.0882799 9.83561 0.00697222 10.8093L1.13941 22.4865C1.22408 23.3438 2.00021 23.9999 2.92804 23.9999H21.0718C22.0032 23.9999 22.7793 23.3438 22.8604 22.4865L23.9929 10.8093C24.0881 9.83914 23.2591 8.99951 22.2042 8.99951H1.79559Z"
        :fill="fill"
      />
      <path
        d="M19.7914 7.79303H18.1262C18.1262 4.41335 15.378 1.66515 11.9984 1.66515C8.61867 1.66515 5.87047 4.41335 5.87047 7.79303H4.20532C4.20885 3.4961 7.70495 0 12.0019 0C16.2988 0 19.7949 3.4961 19.7949 7.79303H19.7914Z"
        :fill="fill"
      />
    </g>
    <defs>
      <clipPath :id="id">
        <rect
          width="24"
          height="24"
          :fill="fill"
        />
      </clipPath>
    </defs>
  </svg>
</template>
