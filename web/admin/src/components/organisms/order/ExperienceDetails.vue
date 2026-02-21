<script lang="ts" setup>
import { mdiAccount, mdiAccountSchool, mdiAccountSupervisor, mdiBaby, mdiSchool } from '@mdi/js'
import type { OrderExperience } from '~/types/api/v1'

interface Props {
  experience: OrderExperience | undefined
  loading?: boolean
  variant?: 'default' | 'compact'
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  variant: 'default',
})

// 年齢カテゴリの定義
interface AgeCategory {
  key: keyof Pick<OrderExperience, 'adultCount' | 'preschoolCount' | 'elementarySchoolCount' | 'juniorHighSchoolCount' | 'seniorCount'>
  priceKey: keyof Pick<OrderExperience, 'adultPrice' | 'preschoolPrice' | 'elementarySchoolPrice' | 'juniorHighSchoolPrice' | 'seniorPrice'>
  label: string
  icon: string
  color: string
}

const ageCategories: AgeCategory[] = [
  {
    key: 'adultCount',
    priceKey: 'adultPrice',
    label: '大人',
    icon: mdiAccount,
    color: 'primary',
  },
  {
    key: 'preschoolCount',
    priceKey: 'preschoolPrice',
    label: '未就学児(3歳〜)',
    icon: mdiBaby,
    color: 'info',
  },
  {
    key: 'elementarySchoolCount',
    priceKey: 'elementarySchoolPrice',
    label: '小学生',
    icon: mdiSchool,
    color: 'success',
  },
  {
    key: 'juniorHighSchoolCount',
    priceKey: 'juniorHighSchoolPrice',
    label: '中学生',
    icon: mdiAccountSchool,
    color: 'warning',
  },
  {
    key: 'seniorCount',
    priceKey: 'seniorPrice',
    label: 'シニア(65歳〜)',
    icon: mdiAccountSupervisor,
    color: 'secondary',
  },
]

// 表示対象のカテゴリのみを抽出
const visibleCategories = computed(() => {
  if (!props.experience) return []

  return ageCategories.filter((category) => {
    const count = props.experience?.[category.key] || 0
    return count > 0
  })
})

// 合計金額を計算
const totalAmount = computed(() => {
  if (!props.experience) return 0

  return ageCategories.reduce((total, category) => {
    const count = props.experience?.[category.key] || 0
    const price = props.experience?.[category.priceKey] || 0
    return total + (count * price)
  }, 0)
})

// カテゴリ別金額を計算
const getCategoryTotal = (category: AgeCategory): number => {
  if (!props.experience) return 0
  const count = props.experience[category.key] || 0
  const price = props.experience[category.priceKey] || 0
  return count * price
}
</script>

<template>
  <div>
    <!-- デフォルト表示（詳細カード形式） -->
    <template v-if="variant === 'default'">
      <v-row v-if="visibleCategories.length > 0">
        <v-col
          v-for="category in visibleCategories"
          :key="category.key"
          cols="12"
          sm="6"
          md="4"
        >
          <v-card
            variant="tonal"
            :color="category.color"
            class="text-center"
            :loading="loading"
          >
            <v-card-text>
              <v-icon
                size="32"
                class="mb-2"
                :icon="category.icon"
              />
              <h4 class="text-h5 mb-1">
                {{ props.experience?.[category.key] || 0 }}人
              </h4>
              <p class="text-subtitle-2 mb-1">
                {{ category.label }}
              </p>
              <p class="font-weight-bold">
                ¥{{ getCategoryTotal(category).toLocaleString() }}
              </p>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>

    <!-- コンパクト表示（テーブル形式） -->
    <template v-else-if="variant === 'compact'">
      <v-table v-if="visibleCategories.length > 0">
        <thead>
          <tr>
            <th>カテゴリ</th>
            <th class="text-center">
              人数
            </th>
            <th class="text-center">
              単価
            </th>
            <th class="text-right">
              小計
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="category in visibleCategories"
            :key="category.key"
          >
            <td>
              <div class="d-flex align-center">
                <v-icon
                  :color="category.color"
                  size="20"
                  class="mr-2"
                >
                  ${{ category.icon }}
                </v-icon>
                {{ category.label }}
              </div>
            </td>
            <td class="text-center">
              <v-chip
                :color="category.color"
                size="small"
                variant="tonal"
              >
                {{ props.experience?.[category.key] || 0 }}人
              </v-chip>
            </td>
            <td class="text-center">
              ¥{{ (props.experience?.[category.priceKey] || 0).toLocaleString() }}
            </td>
            <td class="text-right font-weight-medium">
              ¥{{ getCategoryTotal(category).toLocaleString() }}
            </td>
          </tr>
        </tbody>
      </v-table>
    </template>

    <!-- データがない場合 -->
    <v-alert
      v-if="visibleCategories.length === 0"
      type="info"
      variant="tonal"
    >
      体験予約データがありません
    </v-alert>
  </div>
</template>

<style scoped>
.v-table th {
  font-weight: 600 !important;
}

.v-table tbody tr:hover {
  background-color: rgb(0 0 0 / 4%);
}
</style>
