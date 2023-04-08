import { createVuetify } from 'vuetify/lib/framework.mjs'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import colors from 'vuetify/lib/util/color'

export default defineNuxtPlugin((nuxtApp) => {
  const vuetify = createVuetify({
    locale: {
      locale: 'ja'
    },
    components,
    directives,
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
            success: colors.green.accent3
          }
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
            success: colors.green.accent3
          }
        }
      }
    }
  })

  nuxtApp.vueApp.use(vuetify)
})
