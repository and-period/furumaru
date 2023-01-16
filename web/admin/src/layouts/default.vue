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
      <v-toolbar-title>
        <nuxt-link to="/">
          <the-app-title class="pt-2" />
        </nuxt-link>
      </v-toolbar-title>
      <v-spacer />
      <v-badge :value="hasUnread" color="info" dot overlap>
        <v-btn icon @click="handleClickMessage">
          <v-icon>mdi-bell</v-icon>
        </v-btn>
      </v-badge>
    </v-app-bar>

    <v-snackbar
      v-for="(snackbar, i) in snackbars"
      :key="i"
      v-model="snackbar.isOpen"
      :color="snackbar.color"
      top
      app
      elevation="1"
      :timeout="snackbar.timeout"
      :style="calcStyle(i)"
    >
      {{ snackbar.message }}
      <template #action="{ attrs }">
        <v-btn text v-bind="attrs" @click="commonStore.hideSnackbar(i)"
          >閉じる</v-btn
        >
      </template>
    </v-snackbar>

    <v-main class="bg-color">
      <v-container>
        <Nuxt />
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { useRouter } from '@nuxtjs/composition-api'
import { computed, defineComponent, ref } from '@vue/composition-api'

import { useCommonStore } from '~/store/common'
import { useMessageStore } from '~/store/message'

interface NavigationDrawerItem {
  to: string
  icon: string
  title: string
}

export default defineComponent({
  setup() {
    const drawer = ref<boolean>(true)
    const router = useRouter()

    const commonStore = useCommonStore()
    const messageStore = useMessageStore()

    const snackbars = computed(() => {
      return commonStore.snackbars.filter((item) => item.isOpen)
    })
    const hasUnread = computed<boolean>(() => messageStore.hasUnread)

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
      // {
      //   to: '/analytics',
      //   icon: 'mdi-poll',
      //   title: '分析',
      // },
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
        to: '/promotions',
        icon: 'mdi-cash-100',
        title: 'セール情報管理',
      },
      {
        to: '/producers',
        icon: 'mdi-account',
        title: '生産者管理',
      },
      {
        to: '/coordinators',
        icon: 'mdi-account-star-outline',
        title: 'コーディネータ管理',
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
    const handleClickMessage = () => {
      router.push('/messages')
    }

    const calcStyle = (i: number) => {
      if (i > 0) {
        return `top: ${60 * i}px;`
      }
    }

    return {
      drawer,
      navigationDrawerHomeItem,
      navigationDrawerList,
      navigationDrawerSettingsList,
      handleClickNavIcon,
      handleClickMessage,
      hasUnread,
      snackbars,
      commonStore,
      calcStyle,
    }
  },
})
</script>

<style lang="scss" scoped>
.bg-color {
  background-color: #eef5f9;
}
</style>
