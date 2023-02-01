<script lang="ts" setup>
import { HeaderMenuItem } from '~/types/props'

interface Props {
  menuItems: HeaderMenuItem[]
}

defineProps<Props>()

const emits = defineEmits<{(name: 'click:cart'): void}>()

const _cartIsEmpty = computed<boolean>(() => {
  return true
})

const handleCartClick = (): void => {
  emits('click:cart')
}
</script>

<template>
  <div class="flex px-10 py-4 bg-base items-center">
    <the-marche-logo class="m-0 " />
    <div class="flex items-center justify-end w-full text-main">
      <nav class="mr-16">
        <ul class="list-none flex gap-x-10">
          <li v-for="item, i in menuItems" :key="i">
            <a href="#" :class="{'border-b border-main pb-1': item.active}" @click="item.onClick">
              {{ item.text }}
            </a>
          </li>
        </ul>
      </nav>

      <div class="flex gap-x-8">
        <the-icon-button>
          <the-account-icon />
        </the-icon-button>

        <the-icon-button>
          <the-ring-icon />
        </the-icon-button>

        <the-icon-button @click="handleCartClick">
          <the-cart-icon />
        </the-icon-button>
      </div>
    </div>
  </div>
</template>
