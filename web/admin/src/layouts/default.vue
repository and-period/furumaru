<script lang="ts" setup>
import { mdiHome, mdiMenu, mdiOrderBoolAscendingVariant, mdiCart, mdiAntenna, mdiAccountDetails, mdiForum, mdiBellRing, mdiCash100, mdiAccount, mdiAccountStarOutline, mdiCog, mdiBell } from '@mdi/js'
import { useCommonStore, useMessageStore } from '~/store'

interface NavigationDrawerItem {
  to: string
  icon: string
  title: string
}

const drawer = ref<boolean>(true)
const router = useRouter()

const commonStore = useCommonStore()
const messageStore = useMessageStore()

const snackbars = computed(() => {
  return commonStore.snackbars.filter(item => item.isOpen)
})
const hasUnread = computed<boolean>(() => messageStore.hasUnread)

const navigationDrawerHomeItem: NavigationDrawerItem = {
  to: '/',
  icon: mdiHome,
  title: 'ホーム'
}

const navigationDrawerList: NavigationDrawerItem[] = [
  {
    to: '/orders',
    icon: mdiOrderBoolAscendingVariant,
    title: '注文'
  },
  {
    to: '/products',
    icon: mdiCart,
    title: '商品管理'
  },
  {
    to: '/livestreaming',
    icon: mdiAntenna,
    title: 'ライブ配信'
  },
  // {
  //   to: '/analytics',
  //   icon: mdiPoll,
  //   title: '分析',
  // },
  {
    to: '/customers',
    icon: mdiAccountDetails,
    title: '顧客管理'
  },
  {
    to: '/contacts',
    icon: mdiForum,
    title: 'お問い合わせ管理'
  },
  {
    to: '/notifications',
    icon: mdiBellRing,
    title: 'お知らせ管理'
  },
  {
    to: '/promotions',
    icon: mdiCash100,
    title: 'セール情報管理'
  },
  {
    to: '/producers',
    icon: mdiAccount,
    title: '生産者管理'
  },
  {
    to: '/coordinators',
    icon: mdiAccountStarOutline,
    title: 'コーディネータ管理'
  }
]

const navigationDrawerSettingsList: NavigationDrawerItem[] = [
  {
    to: '/settings',
    icon: mdiCog,
    title: 'システム設定'
  }
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
</script>

<template>
  <v-app>
    <v-app-bar color="primary">
      <template #prepend>
        <v-app-bar-nav-icon @click="handleClickNavIcon">
          <v-icon :icon="mdiMenu" color="white" />
        </v-app-bar-nav-icon>
      </template>
      <v-toolbar-title>
        <nuxt-link to="/">
          <atoms-app-title class="pt-2" />
        </nuxt-link>
      </v-toolbar-title>
      <template #append>
        <v-btn icon @click="handleClickMessage">
          <v-badge v-model="hasUnread" color="info" dot floating>
            <v-icon :icon="mdiBell" color="white" />
          </v-badge>
        </v-btn>
      </template>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer">
      <v-list>
        <v-list-item
          :to="navigationDrawerHomeItem.to"
          exact
          :prepend-icon="navigationDrawerHomeItem.icon"
          color="primary"
        >
          <v-list-item-title>
            {{ navigationDrawerHomeItem.title }}
          </v-list-item-title>
        </v-list-item>
      </v-list>

      <v-divider />
      <v-list>
        <v-list-item
          v-for="(item, i) in navigationDrawerList"
          :key="i"
          :to="item.to"
          :prepend-icon="item.icon"
          color="primary"
        >
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
      <v-divider />

      <v-list>
        <v-list-item
          v-for="(item, i) in navigationDrawerSettingsList"
          :key="i"
          :to="item.to"
          exact
          :prepend-icon="item.icon"
          color="primary"
        >
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-snackbar
      v-for="(snackbar, i) in snackbars"
      :key="i"
      v-model="snackbar.isOpen"
      :color="snackbar.color"
      location="top"
      variant="elevated"
      :timeout="snackbar.timeout"
      :style="calcStyle(i)"
    >
      {{ snackbar.message }}
      <template #action="{ props }">
        <v-btn variant="text" v-bind="props" @click="commonStore.hideSnackbar(i)">
          閉じる
        </v-btn>
      </template>
    </v-snackbar>

    <v-main class="bg-color">
      <v-container>
        <slot />
      </v-container>
    </v-main>
  </v-app>
</template>

<style lang="scss" scoped>
.bg-color {
  background-color: #eef5f9;
}
</style>
