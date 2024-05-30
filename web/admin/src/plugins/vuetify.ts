import { createVuetify } from 'vuetify'
import { md3 } from 'vuetify/blueprints'
import { ja } from 'vuetify/locale'
import * as components from 'vuetify/components'
import * as labs from 'vuetify/labs/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg'
// @ts-ignore
import colors from 'vuetify/lib/util/colors.mjs'

export default defineNuxtPlugin((nuxtApp) => {
  const vuetify = createVuetify({
    locale: {
      locale: 'ja',
      messages: { ja },
    },
    blueprint: md3,
    components: {
      ...labs,
      ...components,
    },
    directives,
    icons: {
      defaultSet: 'mdi',
      aliases,
      sets: {
        mdi,
      },
    },
    defaults: {
      VAutocomplete: {
        variant: 'underlined',
      },
      VTextarea: {
        variant: 'underlined',
      },
      VTextField: {
        variant: 'underlined',
      },
      VSelect: {
        variant: 'underlined',
      },
      VCombobox: {
        variant: 'underlined',
      },
      VCard: {
        elevation: 0,
      },
      VTab: {
        VBtn: {
          rounded: '0',
        },
      },
      VBtnToggle: {
        rounded: '0',
        density: 'compact',
        variant: 'outlined',
        divided: true,
        VBtn: {
          size: 'small',
        },
      },
    },
    theme: {
      defaultTheme: 'light',
      themes: {
        light: {
          dark: false,
          colors: {
            primary: colors.lightGreen.darken2,
            primaryLight: colors.lightGreen.lighten2,
            accent: colors.amber.darken1,
            secondary: colors.amber.darken3,
            info: colors.teal.lighten1,
            warning: colors.amber.base,
            error: colors.deepOrange.accent4,
            unknown: colors.grey.darken2,
            success: colors.green.accent3,
          },
        },
        dark: {
          dark: true,
          colors: {
            primary: colors.green.accent1,
            accent: colors.grey.darken3,
            secondary: colors.amber.darken3,
            info: colors.teal.lighten1,
            warning: colors.amber.base,
            error: colors.deepOrange.accent4,
            success: colors.green.accent3,
          },
        },
      },
    },
  })

  nuxtApp.vueApp.use(vuetify)
})
