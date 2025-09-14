<script lang="ts" setup>
import {
  mdiHome,
  mdiMenu,
  mdiOrderBoolAscendingVariant,
  mdiCart,
  mdiAntenna,
  mdiAccountDetails,
  mdiVideo,
  mdiBellRing,
  mdiCash100,
  mdiAccount,
  mdiCog,
  mdiBell,
  mdiFootPrint,
  mdiMagnify,
  mdiChevronDown,
  mdiChevronUp,
} from '@mdi/js'
import { storeToRefs } from 'pinia'
import { getResizedImages } from '~/lib/helpers'
import { useAuthStore, useCommonStore, useMessageStore } from '~/store'
import { AdminType } from '~/types/api/v1'

interface NavigationDrawerItem {
  to: string
  icon: string
  title: string
  adminTypes?: AdminType[]
  category?: string
}

interface NavigationGroup {
  title: string
  items: NavigationDrawerItem[]
  expanded?: boolean
}

const searchQuery = ref<string>('')
const expandedGroups = ref<{ [key: string]: boolean }>({
  general: true,
  settings: true,
})
const router = useRouter()
const authStore = useAuthStore()
const commonStore = useCommonStore()
const messageStore = useMessageStore()

const { user, adminType } = storeToRefs(authStore)

// Responsive breakpoints using window size
const windowWidth = ref(1280) // Default to desktop size

const updateWindowWidth = () => {
  if (import.meta.client && window) {
    windowWidth.value = window.innerWidth
  }
}

onMounted(() => {
  updateWindowWidth()
  if (import.meta.client) {
    window.addEventListener('resize', updateWindowWidth)
  }
})

onUnmounted(() => {
  if (import.meta.client && window) {
    window.removeEventListener('resize', updateWindowWidth)
  }
})

const snackbars = computed(() => {
  return commonStore.snackbars.filter(item => item.isOpen)
})
const hasUnread = computed<boolean>(() => messageStore.hasUnread)

const homeDrawer: NavigationDrawerItem = {
  to: '/',
  icon: mdiHome,
  title: 'ホーム',
}

const navigationGroups: NavigationGroup[] = [
  {
    title: '販売・商品管理',
    items: [
      {
        to: '/producers',
        icon: mdiAccount,
        title: '生産者管理',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'general',
      },
      {
        to: '/products',
        icon: mdiCart,
        title: '商品管理',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'general',
      },
      {
        to: '/experiences',
        icon: mdiFootPrint,
        title: '体験管理',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'general',
      },
    ],
  },
  {
    title: 'ライブ・動画',
    items: [
      {
        to: '/schedules',
        icon: mdiAntenna,
        title: 'ライブ配信',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'general',
      },
      {
        to: '/videos',
        icon: mdiVideo,
        title: '動画管理',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'general',
      },
    ],
  },
  {
    title: '注文・顧客管理',
    items: [
      {
        to: '/orders',
        icon: mdiOrderBoolAscendingVariant,
        title: '注文管理',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'general',
      },
      {
        to: '/customers',
        icon: mdiAccountDetails,
        title: '顧客管理',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'general',
      },
    ],
  },
  {
    title: 'マーケティング',
    items: [
      {
        to: '/notifications',
        icon: mdiBellRing,
        title: 'お知らせ情報',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'general',
      },
      {
        to: '/promotions',
        icon: mdiCash100,
        title: 'セール情報',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'general',
      },
    ],
  },
  {
    title: '設定',
    items: [
      {
        to: '/accounts',
        icon: mdiAccount,
        title: 'マイページ',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'settings',
      },
      {
        to: '/system',
        icon: mdiCog,
        title: 'システム設定',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'settings',
      },
      {
        to: '/version',
        icon: mdiCog,
        title: 'バージョン情報',
        adminTypes: [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator],
        category: 'settings',
      },
    ],
  },
]

const getImages = (): string => {
  if (!user.value?.thumbnailUrl) {
    return ''
  }
  return getResizedImages(user.value.thumbnailUrl)
}

const filteredNavigationGroups = computed(() => {
  return navigationGroups.map(group => ({
    ...group,
    items: group.items.filter((item: NavigationDrawerItem): boolean => {
      return item.adminTypes?.includes(adminType.value) || false
    }),
  })).filter(group => group.items.length > 0)
})

const searchFilteredGroups = computed(() => {
  if (!searchQuery.value.trim()) {
    return filteredNavigationGroups.value
  }

  const query = searchQuery.value.toLowerCase().trim()
  return filteredNavigationGroups.value.map(group => ({
    ...group,
    items: group.items.filter(item =>
      item.title.toLowerCase().includes(query),
    ),
  })).filter(group => group.items.length > 0)
})

const toggleGroup = (groupTitle: string) => {
  expandedGroups.value[groupTitle] = !expandedGroups.value[groupTitle]
}

const isGroupExpanded = (groupTitle: string) => {
  return expandedGroups.value[groupTitle] !== false
}

// Responsive drawer behavior based on Vuetify breakpoints
// xs: <600px, sm: 600-959px, md: 960-1279px, lg: 1280-1919px, xl: >=1920px
const isDesktop = computed(() => windowWidth.value >= 1280)
const isTablet = computed(() => windowWidth.value >= 960 && windowWidth.value < 1280)
const isMobile = computed(() => windowWidth.value < 960)

// Simple drawer state (like original)
const drawer = ref<boolean>(true)

// Simplified drawer configuration
const drawerPermanent = computed(() => isDesktop.value)
const drawerTemporary = computed(() => !isDesktop.value)

const handleClickNavIcon = () => {
  drawer.value = !drawer.value
}

// Auto-close drawer on route change for mobile/tablet
watch(() => router.currentRoute.value.path, () => {
  if (!isDesktop.value) {
    drawer.value = false
  }
})

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
    <v-app-bar
      color="primary"
      elevation="2"
      :density="$vuetify.display.mobile ? 'compact' : 'default'"
    >
      <template #prepend>
        <v-app-bar-nav-icon @click="handleClickNavIcon">
          <v-icon
            :icon="mdiMenu"
            color="white"
          />
        </v-app-bar-nav-icon>
      </template>
      <v-toolbar-title class="d-flex align-center">
        <nuxt-link
          to="/"
          class="text-decoration-none"
        >
          <atoms-app-title class="pt-2" />
        </nuxt-link>
      </v-toolbar-title>
      <v-spacer />
      <div class="d-flex align-center ga-2">
        <div
          v-if="isDesktop && user"
          class="text-white d-flex align-center ga-2 mr-4"
        >
          <v-avatar size="32">
            <v-img
              v-if="user?.thumbnailUrl"
              :src="user.thumbnailUrl"
              :srcset="getImages()"
              cover
            />
            <v-icon
              v-else
              :icon="mdiAccount"
            />
          </v-avatar>
          <div class="d-flex flex-column">
            <span class="text-body-2 font-weight-medium">{{ user.username }}</span>
            <span class="text-caption opacity-75">{{ user.email }}</span>
          </div>
        </div>
        <v-btn
          icon
          variant="text"
          @click="handleClickMessage"
        >
          <v-badge
            v-model="hasUnread"
            color="error"
            dot
            floating
          >
            <v-icon
              :icon="mdiBell"
              color="white"
            />
          </v-badge>
        </v-btn>
      </div>
    </v-app-bar>

    <v-navigation-drawer
      v-model="drawer"
      :permanent="drawerPermanent"
      :temporary="drawerTemporary"
      width="280"
      class="custom-drawer"
    >
      <!-- User Profile Section (Mobile/Tablet only) -->
      <v-list v-if="!isDesktop && user">
        <v-list-item class="px-2 py-3">
          <template #prepend>
            <v-avatar size="40">
              <v-img
                v-if="user?.thumbnailUrl"
                :src="user?.thumbnailUrl"
                :srcset="getImages()"
                cover
              />
              <v-icon
                v-else
                :icon="mdiAccount"
              />
            </v-avatar>
          </template>
          <div class="ml-3">
            <div class="text-subtitle-2 font-weight-medium">
              {{ user?.username || "" }}
            </div>
            <div class="text-caption text-grey-darken-1">
              {{ user?.email || "" }}
            </div>
          </div>
        </v-list-item>
        <v-divider />
      </v-list>

      <!-- Search Section -->
      <div class="pa-3">
        <v-text-field
          v-model="searchQuery"
          :prepend-inner-icon="mdiMagnify"
          placeholder="メニューを検索..."
          variant="outlined"
          :density="isMobile ? 'comfortable' : 'compact'"
          hide-details
          clearable
          class="mb-2 search-field"
        />
      </div>

      <v-divider />

      <!-- Home Item -->
      <v-list>
        <v-list-item
          :to="homeDrawer.to"
          exact
          :prepend-icon="homeDrawer.icon"
          color="primary"
          class="rounded-lg mx-2 my-1"
        >
          <v-list-item-title class="font-weight-medium">
            {{ homeDrawer.title }}
          </v-list-item-title>
        </v-list-item>
      </v-list>

      <v-divider class="my-2" />

      <!-- Navigation Groups -->
      <div
        v-for="group in searchFilteredGroups"
        :key="group.title"
        class="mb-2"
      >
        <v-list-item
          class="px-4 py-2 cursor-pointer"
          :class="{ 'bg-grey-lighten-4': !isGroupExpanded(group.title) }"
          @click="toggleGroup(group.title)"
        >
          <template #prepend>
            <v-icon
              size="18"
              class="mr-2"
              :icon="isGroupExpanded(group.title) ? mdiChevronUp : mdiChevronDown"
            />
          </template>
          <v-list-item-title class="text-body-2 font-weight-bold text-grey-darken-2">
            {{ group.title }}
          </v-list-item-title>
        </v-list-item>

        <v-expand-transition>
          <v-list
            v-show="isGroupExpanded(group.title)"
            class="pt-0"
          >
            <v-list-item
              v-for="item in group.items"
              :key="item.to"
              :to="item.to"
              :prepend-icon="item.icon"
              color="primary"
              class="rounded-lg mx-2 mb-1 nav-item"
            >
              <v-list-item-title class="text-body-2">
                {{ item.title }}
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-expand-transition>
      </div>

      <!-- Empty State for Search -->
      <div
        v-if="searchQuery.trim() && searchFilteredGroups.length === 0"
        class="pa-4 text-center"
      >
        <v-icon
          :icon="mdiMagnify"
          size="48"
          class="text-grey-lighten-1 mb-2"
        />
        <div class="text-body-2 text-grey-darken-1">
          「{{ searchQuery }}」に一致するメニューが見つかりません
        </div>
      </div>
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

.cursor-pointer {
  cursor: pointer;
}

.nav-item {
  transition: all 0.2s ease-in-out;

  &:hover {
    transform: translateX(4px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
}

.v-navigation-drawer {
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.08) !important;

  .v-list-item {
    &.v-list-item--active {
      background: linear-gradient(135deg, rgba(var(--v-theme-primary), 0.08), rgba(var(--v-theme-primary), 0.12));
      border-left: 3px solid rgb(var(--v-theme-primary));

      .v-list-item-title {
        font-weight: 600;
        color: rgb(var(--v-theme-primary));
      }
    }

    &:not(.v-list-item--active):hover {
      background: rgba(var(--v-theme-primary), 0.04);
    }
  }
}

.v-app-bar {
  backdrop-filter: blur(10px);
  background: linear-gradient(135deg, rgb(var(--v-theme-primary)), rgba(var(--v-theme-primary), 0.9)) !important;
}

.search-field {
  .v-field--focused {
    box-shadow: 0 0 0 2px rgba(var(--v-theme-primary), 0.2);
  }
}

.custom-drawer {
  .v-list-item {
    margin: 0 8px 4px 8px;
    border-radius: 8px;

    &:not(.nav-item):hover {
      background: rgba(var(--v-theme-primary), 0.04);
    }
  }

  // Mobile and tablet responsive adjustments
  @media (max-width: 1279px) {
    .v-list-item {
      min-height: 48px;
    }
  }
}
</style>
