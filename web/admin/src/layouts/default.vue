<script lang="ts" setup>
import {
  mdiHome,
  mdiMenu,
  mdiOrderBoolAscendingVariant,
  mdiCart,
  mdiAntenna,
  mdiAccountDetails,
  mdiForum,
  mdiBellRing,
  mdiCash100,
  mdiAccount,
  mdiCog,
  mdiBell
} from '@mdi/js'
import { storeToRefs } from 'pinia'
import { getResizedImages } from '~/lib/helpers'
import { useAuthStore, useCommonStore, useMessageStore } from '~/store'
import { AdminRole } from '~/types/api'

interface NavigationDrawerItem {
  to: string;
  icon: string;
  title: string;
  roles?: AdminRole[];
}

const drawer = ref<boolean>(true)
const router = useRouter()
const authStore = useAuthStore()
const commonStore = useCommonStore()
const messageStore = useMessageStore()

const { user, role } = storeToRefs(authStore)

const snackbars = computed(() => {
  return commonStore.snackbars.filter(item => item.isOpen)
})
const hasUnread = computed<boolean>(() => messageStore.hasUnread)

const homeDrawer: NavigationDrawerItem = {
  to: '/',
  icon: mdiHome,
  title: 'ホーム'
}

const generalDrawers: NavigationDrawerItem[] = [
  {
    to: '/orders',
    icon: mdiOrderBoolAscendingVariant,
    title: '注文管理',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  {
    to: '/products',
    icon: mdiCart,
    title: '商品管理',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  {
    to: '/schedules',
    icon: mdiAntenna,
    title: 'ライブ配信',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  // {
  //   to: '/contacts',
  //   icon: mdiForum,
  //   title: 'お問い合わせ',
  //   roles: [AdminRole.ADMINISTRATOR]
  // },
  {
    to: '/notifications',
    icon: mdiBellRing,
    title: 'お知らせ情報',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  {
    to: '/promotions',
    icon: mdiCash100,
    title: 'セール情報',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  {
    to: '/producers',
    icon: mdiAccount,
    title: '生産者管理',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  {
    to: '/customers',
    icon: mdiAccountDetails,
    title: '顧客管理',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  }
]

const settingDrawers: NavigationDrawerItem[] = [
  {
    to: '/accounts',
    icon: mdiAccount,
    title: 'マイページ',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  {
    to: '/system',
    icon: mdiCog,
    title: 'システム設定',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  },
  {
    to: '/version',
    icon: mdiCog,
    title: 'バージョン情報',
    roles: [AdminRole.ADMINISTRATOR, AdminRole.COORDINATOR]
  }
]

const getImages = (): string => {
  if (!user.value?.thumbnails) {
    return ''
  }
  return getResizedImages(user.value.thumbnails)
}

const getGeneralDrawers = (): NavigationDrawerItem[] => {
  return generalDrawers.filter((drawer: NavigationDrawerItem): boolean => {
    return drawer.roles?.includes(role.value) || false
  })
}

const getSettingDrawers = (): NavigationDrawerItem[] => {
  return settingDrawers.filter((drawer: NavigationDrawerItem): boolean => {
    return drawer.roles?.includes(role.value) || false
  })
}

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
        <v-list-item exact>
          <template #prepend>
            <v-avatar v-if="user?.thumbnailUrl">
              <v-img cover :src="user?.thumbnailUrl" :srcset="getImages()" />
            </v-avatar>
            <v-icon v-else :icon="mdiAccount" />
          </template>

          <div>{{ user?.username || "" }}</div>
          <div class="text-caption text-grey">
            {{ user?.email || "" }}
          </div>
        </v-list-item>
      </v-list>

      <v-divider />

      <v-list>
        <v-list-item
          :to="homeDrawer.to"
          exact
          :prepend-icon="homeDrawer.icon"
          color="primary"
        >
          <v-list-item-title>{{ homeDrawer.title }}</v-list-item-title>
        </v-list-item>
      </v-list>

      <v-divider />

      <v-list>
        <v-list-item
          v-for="(item, i) in getGeneralDrawers()"
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
          v-for="(item, i) in getSettingDrawers()"
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
      <template #actions>
        <v-btn
          variant="text"
          color="white"
          @click="commonStore.hideSnackbar(i)"
        >
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
