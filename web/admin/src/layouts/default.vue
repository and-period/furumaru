<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" app clipped>
      <v-list shaped>
        <v-list-item
          :to="navigationDrawerHomeItem.to"
          router
          exact
          color="primary"
        >
          <v-list-item-icon>
            <v-icon>{{ navigationDrawerHomeItem.icon }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{
              navigationDrawerHomeItem.title
            }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>

      <v-divider />
      <v-list shaped>
        <v-list-item
          v-for="(item, i) in navigationDrawerList"
          :key="i"
          :to="item.to"
          router
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
      <v-divider />

      <v-list shaped>
        <v-list-item
          v-for="(item, i) in navigationDrawerSettingsList"
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
  </v-app>
</template>

<script lang="ts">
import { useRouter } from '@nuxtjs/composition-api'
import { defineComponent, ref } from '@vue/composition-api'

import { useAuthStore } from '~/store/auth'

interface NavigationDrawerItem {
  to: string
  icon: string
  title: string
}

export default defineComponent({
  setup() {
    // TODO: 雑に未ログイン時を検証する
    const router = useRouter()
    const { isAuthenticated } = useAuthStore()
    if (!isAuthenticated) {
      router.push('/signin')
    }

    const drawer = ref<boolean>(true)

    const navigationDrawerHomeItem: NavigationDrawerItem = {
      to: '/',
      icon: 'mdi-home',
      title: 'ホーム',
    }

    const navigationDrawerList: NavigationDrawerItem[] = [
      {
        to: '/orders',
        icon: 'mdi-order-bool-ascending-variant',
        title: '注文',
      },
      {
        to: '/products',
        icon: 'mdi-cart',
        title: '商品管理',
      },
      {
        to: '/livestreaming',
        icon: 'mdi-antenna',
        title: 'ライブ配信',
      },
      {
        to: '/analytics',
        icon: 'mdi-poll',
        title: '分析',
      },
      {
        to: '/customers',
        icon: 'mdi-account-details',
        title: '顧客管理',
      },
      {
        to: '/contacts',
        icon: 'mdi-forum',
        title: 'お問い合わせ管理',
      },
      {
        to: '/notifications',
        icon: 'mdi-bell-ring',
        title: 'お知らせ管理',
      },
      {
        to: '/events',
        icon: 'mdi-cash-100',
        title: 'セール情報管理',
      },
      {
        to: '/producers',
        icon: 'mdi-account',
        title: '生産者管理',
      },
    ]

    const navigationDrawerSettingsList: NavigationDrawerItem[] = [
      {
        to: '/settings',
        icon: 'mdi-cog',
        title: 'システム設定',
      },
    ]

    const handleClickNavIcon = () => {
      drawer.value = !drawer.value
    }

    return {
      drawer,
      navigationDrawerHomeItem,
      navigationDrawerList,
      navigationDrawerSettingsList,
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
