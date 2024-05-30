<script lang="ts" setup>
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
</script>

<template>
  <v-card>
    <v-card-title>システム設定</v-card-title>

    <v-card-text>
      <v-list>
        <v-list-item
          v-for="(item, i) in props.menus"
          :key="i"
          :border="true"
          :color="item.color"
          :title="item.text"
          @click="onClick(item.action)"
        />
      </v-list>
    </v-card-text>
  </v-card>
</template>
