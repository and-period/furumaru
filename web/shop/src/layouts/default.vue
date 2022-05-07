<template>
  <v-app dark>
    <v-navigation-drawer permanent app clipped>
      <v-list>
        <v-list-item
          v-for="(item, i) in navigationDrawerList"
          :key="i"
          :to="item.to"
          router
          exact
        >
          <v-list-item-icon>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar flat app clipped-left>
      <v-toolbar-title>Online Marche</v-toolbar-title>
      <v-spacer />
    </v-app-bar>
    <v-main>
      <v-container>
        <Nuxt />
      </v-container>
    </v-main>
    <v-footer absolute app>
      <span>&copy; {{ new Date().getFullYear() }}</span>
    </v-footer>
  </v-app>
</template>

<script lang="ts">
import { defineComponent, ref } from '@nuxtjs/composition-api'

interface NavigationDrawerItem {
  to: string
  icon: string
  title: string
}

export default defineComponent({
  setup() {
    const navigationDrawerList = ref<NavigationDrawerItem[]>([
      {
        to: '/',
        icon: 'mdi-home',
        title: 'Home',
      },
    ])

    return {
      navigationDrawerList,
    }
  },
})
</script>
