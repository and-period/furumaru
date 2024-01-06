<script lang="ts" setup>
import type { FooterMenuItem } from '~/types/props'

interface Props {
  menuItems: FooterMenuItem[]
}

defineProps<Props>()
</script>

<template>
  <footer class="relative bg-base text-center">
    <div class="absolute -top-12 w-full text-center">
      <the-footer-cat-icon class="inline bg-transparent" />
    </div>
    <div class="bg-green py-5 text-white">
      <the-stack class="gap-y-4">
        <p class="text-lg font-bold tracking-widest">Follow us !</p>

        <div class="flex justify-center gap-x-3">
          <a href="https://www.instagram.com/furusato_marche" class="block">
            <the-instagram-icon
              id="footer-instagram-icon"
              class="inline h-5 w-5"
            />
          </a>
          <a href="#" class="block">
            <the-facebook-icon
              id="footer-facebook-icon"
              class="inline h-5 w-5"
            />
          </a>
        </div>

        <ul
          class="mt-4 flex list-none flex-col justify-center gap-4 md:flex-row md:gap-10"
        >
          <li v-for="(item, i) in menuItems" :key="i">
            <nuxt-link :to="item.to">
              {{ item.text }}
            </nuxt-link>
          </li>
        </ul>

        <p class="text-sm tracking-widest">
          &copy; FURUSATO MARCHE All rights Reserved.
        </p>
      </the-stack>
    </div>
  </footer>
</template>
