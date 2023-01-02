import 'vuetify/styles'
import '@fortawesome/fontawesome-free/css/all.css'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, fa } from 'vuetify/iconsets/fa'
import colors from 'vuetify/lib/util/colors'

const opts = {
  components,
  directives,
  icons: {
    defaultSet: 'fa',
    aliases,
    sets: {
      fa,
    },
  },
  theme: {
    defaultTheme: 'light',
    themes: {
      dark: {
        dark: true,
        colors: {
          primary: colors.blue.darken2,
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },
      },
      light: {
        dark: false,
        colors: {
          primary: colors.lightGreen.darken1,
          accent: colors.orange.darken1,
          base: '#FAF2E2',
          facebook: '#1877F2',
          line: '#06C755',
        },
      },
    },
  },
}

export default defineNuxtPlugin((nuxtApp) => {
  const vuetify = createVuetify(opts)
  nuxtApp.vueApp.use(vuetify)
  nuxtApp.provide('vuetify', vuetify)
  return {
    provide: {
      injected: () => vuetify,
    },
  }
})
