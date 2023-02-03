<script lang="ts" setup>import { FooterMenuItem } from '~/types/props'

interface Props {
  menuItems: FooterMenuItem[]
}

defineProps<Props>()
</script>

<template>
  <footer class="text-center bg-base relative">
    <div class="text-center absolute w-full -top-12">
      <the-footer-cat-icon class="inline bg-transparent" />
    </div>
    <div class="text-white bg-green py-5">
      <the-stack class="gap-y-4">
        <p class="text-lg font-bold tracking-widest">
          Follow us !
        </p>

        <div class="flex justify-center gap-x-3">
          <a href="#" class="block">
            <the-instagram-icon class="inline w-5 h-5" />
          </a>
          <a href="#" class="block">
            <the-facebook-icon class="inline w-5 h-5" />
          </a>
        </div>

        <ul class="list-none flex justify-center gap-x-10 mt-4">
          <li v-for="item,i in menuItems" :key="i">
            <a href="#" @click="item.onClick">
              {{ item.text }}
            </a>
          </li>
        </ul>

        <p class="text-sm tracking-widest">
          &copy; FURUSATO MARCHE All rights Reserved.
        </p>
      </the-stack>
    </div>
  </footer>
</template>
