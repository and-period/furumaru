<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" app clipped>
      <v-list shaped>
        <v-list-item
          v-for="(item, i) in navigationDrawerList"
          :key="i"
          :to="item.to"
          router
          exact
          color="primary"
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

    <v-app-bar flat app clipped-left color="primary" dark>
      <v-app-bar-nav-icon @click="handleClickNavIcon"></v-app-bar-nav-icon>
      <v-toolbar-title>Online Marche</v-toolbar-title>
      <v-spacer />
      <v-btn icon>
        <v-icon>mdi-bell</v-icon>
      </v-btn>
    </v-app-bar>
    <v-main class="bg-color">
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
import { defineComponent, ref } from '@vue/composition-api'

interface NavigationDrawerItem {
  to: string
  icon: string
  title: string
}

export default defineComponent({
  setup() {
    const drawer = ref<boolean>(false)

    const navigationDrawerList = ref<NavigationDrawerItem[]>([
      {
        to: '/',
        icon: 'mdi-home',
        title: 'ホーム',
      },
      {
        to: '/products',
        icon: 'mdi-tag',
        title: '商品管理',
      },
      {
        to: '/livestreaming',
        icon: 'mdi-antenna',
        title: 'ライブ配信',
      },
    ])

    const handleClickNavIcon = () => {
      drawer.value = !drawer.value
    }

    return {
      drawer,
      navigationDrawerList,
      handleClickNavIcon,
    }
  },
})
</script>

<style lang="scss" scoped>
.bg-color {
  background-color: #eef5f9;
}
</style>
