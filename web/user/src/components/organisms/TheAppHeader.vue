<template>
  <v-app-bar flat fixed app color="base">
    <v-toolbar-title class="accent--text title">Online Marche</v-toolbar-title>
    <v-spacer />
    <slot />
    <v-tooltip bottom>
      <template #activator="{ on, attrs }">
        <v-btn
          icon
          class="mr-4"
          v-bind="attrs"
          v-on="on"
          @click="handleCartClick"
        >
          <v-badge overlap :content="cartContent" :value="!cartIsEmpty">
            <v-icon>mdi-cart-outline</v-icon>
          </v-badge>
        </v-btn>
      </template>
      <span>{{ cartIsEmpty ? cartEmptyMessage : cartNotEmptyMessage }}</span>
    </v-tooltip>
    <v-menu offset-y>
      <template #activator="{ on }">
        <the-header-menu-button :img-src="profileImgUrl" v-on="on" />
      </template>
      <v-list>
        <v-list-item v-for="n in 3" :key="n">
          <v-list-item-title>{{ n }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </v-app-bar>
</template>

<script lang="ts">
import {
  computed,
  ComputedRef,
  defineComponent,
  SetupContext,
} from '@nuxtjs/composition-api'

export default defineComponent({
  props: {
    profileImgUrl: {
      type: String,
      required: false,
      default: null,
    },
    cartItemCount: {
      type: Number,
      required: true,
      default: 0,
    },
    cartEmptyMessage: {
      type: String,
      required: true,
    },
    cartNotEmptyMessage: {
      type: String,
      required: true,
    },
  },

  setup(props, { emit }: SetupContext) {
    const { cartItemCount } = props

    const cartIsEmpty: ComputedRef<boolean> = computed(() => {
      return cartItemCount === 0
    })

    const cartContent: ComputedRef<number | string> = computed(() => {
      return cartItemCount > 99 ? '99+' : cartItemCount
    })

    const handleCartClick = (): void => {
      emit('click:cart')
    }

    return {
      cartIsEmpty,
      cartContent,
      handleCartClick,
    }
  },
})
</script>

<style lang="scss" scoped>
.title {
  font-family: 'STIX Two Text' !important;
  font-style: normal;
  font-weight: 400;
}
</style>
