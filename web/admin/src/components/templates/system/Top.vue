<script lang="ts" setup>
import {
  mdiAccountSupervisor,
  mdiAccountGroup,
  mdiTagMultiple,
  mdiTag,
  mdiCalendarHeart,
  mdiMapMarkerMultiple,
  mdiTruck,
  mdiTruckOutline,
  mdiCreditCardSettings,
  mdiCog,
  mdiChevronRight,
} from '@mdi/js'
import type { SettingMenu } from '~/types/props'

const props = defineProps({
  menus: {
    type: Array<SettingMenu>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'click', action: () => void): void
}>()

const onClick = (action: () => void): void => {
  emit('click', action)
}

const getMenuIcon = (text: string): string => {
  if (text.includes('管理者管理')) return mdiAccountSupervisor
  if (text.includes('コーディネーター管理')) return mdiAccountGroup
  if (text.includes('カテゴリー・品目管理')) return mdiTagMultiple
  if (text.includes('商品タグ管理')) return mdiTag
  if (text.includes('体験種別管理')) return mdiCalendarHeart
  if (text.includes('スポット種別管理')) return mdiMapMarkerMultiple
  if (text.includes('デフォルト配送設定管理')) return mdiTruckOutline
  if (text.includes('配送設定管理')) return mdiTruck
  if (text.includes('決済システム管理')) return mdiCreditCardSettings
  if (text.includes('バージョン情報')) return mdiCog
  return mdiCog
}

const getMenuDescription = (text: string): string => {
  if (text.includes('管理者管理')) return 'システム管理者の追加・編集・削除'
  if (text.includes('コーディネーター管理')) return 'コーディネーター情報の管理'
  if (text.includes('カテゴリー・品目管理')) return '商品カテゴリと品目の設定'
  if (text.includes('商品タグ管理')) return '商品に付与するタグの管理'
  if (text.includes('体験種別管理')) return '体験商品の種別設定'
  if (text.includes('スポット種別管理')) return 'スポット情報の種別設定'
  if (text.includes('デフォルト配送設定管理')) return 'システム全体のデフォルト配送設定'
  if (text.includes('配送設定管理')) return 'コーディネーター別配送設定'
  if (text.includes('決済システム管理')) return '決済サービスの設定・管理'
  if (text.includes('バージョン情報')) return 'システムのバージョン情報を確認'
  return '各種システム設定'
}

const getMenuCategory = (text: string): string => {
  if (text.includes('管理者管理') || text.includes('コーディネーター管理')) {
    return 'ユーザー管理'
  }
  if (text.includes('カテゴリー') || text.includes('タグ') || text.includes('体験種別') || text.includes('スポット種別')) {
    return '商品・コンテンツ管理'
  }
  if (text.includes('配送') || text.includes('決済') || text.includes('バージョン情報')) {
    return 'システム設定'
  }
  return 'その他'
}

const getPermissionBadge = (text: string): { color: string, label: string } => {
  if (text.includes('管理者管理') || text.includes('デフォルト配送設定管理') || text.includes('決済システム管理')) {
    return { color: 'warning', label: '管理者専用' }
  }
  if (text.includes('配送設定管理')) {
    return { color: 'info', label: 'コーディネーター専用' }
  }
  return { color: 'success', label: '共通' }
}

const groupedMenus = computed(() => {
  const groups: Record<string, SettingMenu[]> = {
    ユーザー管理: [],
    商品・コンテンツ管理: [],
    システム設定: [],
    その他: [],
  }

  props.menus.forEach((menu) => {
    const category = getMenuCategory(menu.text)
    groups[category].push(menu)
  })

  // 空のカテゴリを除外
  return Object.entries(groups).filter(([_, items]) => items.length > 0)
})
</script>

<template>
  <v-container class="pa-6">
    <div class="mb-8">
      <h1 class="text-h4 font-weight-bold mb-2">
        <v-icon
          :icon="mdiCog"
          size="32"
          class="mr-3 text-primary"
        />
        システム設定
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        システム全体の設定と各種管理機能にアクセスできます
      </p>
    </div>

    <div
      v-for="[category, items] in groupedMenus"
      :key="category"
      class="mb-8"
    >
      <div class="category-header mb-4">
        <h2 class="text-h6 font-weight-bold text-grey-darken-2">
          {{ category }}
        </h2>
        <v-divider class="mt-2" />
      </div>

      <v-row>
        <v-col
          v-for="(item, i) in items"
          :key="i"
          cols="12"
          sm="6"
          md="4"
        >
          <v-card
            class="system-card h-100"
            elevation="2"
            hover
            @click="onClick(item.action)"
          >
            <v-card-text class="text-center pa-6">
              <div class="mb-4">
                <v-avatar
                  color="primary"
                  variant="flat"
                  size="56"
                >
                  <v-icon
                    :icon="getMenuIcon(item.text)"
                    size="28"
                    color="white"
                  />
                </v-avatar>
              </div>

              <div class="mb-3">
                <v-chip
                  :color="getPermissionBadge(item.text).color"
                  size="small"
                  variant="outlined"
                  class="mb-2"
                >
                  {{ getPermissionBadge(item.text).label }}
                </v-chip>
              </div>

              <h3 class="text-subtitle-1 font-weight-medium mb-2">
                {{ item.text }}
              </h3>
              <p class="text-body-2 text-grey-darken-1">
                {{ getMenuDescription(item.text) }}
              </p>
            </v-card-text>
            <v-card-actions class="pa-4 pt-0">
              <v-spacer />
              <v-btn
                color="primary"
                variant="text"
                size="small"
              >
                管理する
                <v-icon
                  :icon="mdiChevronRight"
                  end
                />
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </div>
  </v-container>
</template>

<style scoped>
.system-card {
  transition: all 0.3s ease;
  cursor: pointer;
  border-radius: 12px;
}

.system-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgb(0 0 0 / 15%) !important;
}

.category-header {
  position: sticky;
  top: 0;
  background: white;
  z-index: 1;
  padding: 8px 0;
}

@media (width <= 600px) {
  .system-card {
    margin-bottom: 16px;
  }

  .category-header {
    position: static;
  }
}
</style>
