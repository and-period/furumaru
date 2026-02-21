import { md3 } from 'vuetify/blueprints'
import { ja } from 'vuetify/locale'
import type { VOptions } from 'vuetify-nuxt-module'

export default {
  locale: {
    locale: 'ja',
    messages: { ja },
  },
  blueprint: md3,
  icons: {
    defaultSet: 'mdi-svg',
  },
  defaults: {
    VAutocomplete: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VAvatar: {
      variant: 'outlined',
      color: 'grey-lighten-1',
    },
    VTextarea: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VCombobox: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VCard: {
      elevation: 0,
      rounded: 'lg',
    },
    VBtn: {
      rounded: 'lg',
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
          primary: '#7cb342',
          primaryLight: '#dcedc8',
          accent: '#ffb300',
          secondary: '#795548',
          info: '#26c6da',
          warning: '#ffc107',
          error: '#ef5350',
          unknown: '#616161',
          success: '#66bb6a',
          white: '#ffffff',
          surface: '#ffffff',
          background: '#f1f8e9',
        },
      },
      dark: {
        dark: true,
        colors: {
          primary: '#9ccc65',
          primaryLight: '#558b2f',
          accent: '#424242',
          secondary: '#a1887f',
          info: '#00acc1',
          warning: '#ffc107',
          error: '#ef9a9a',
          unknown: '#9e9e9e',
          success: '#81c784',
          surface: '#1e1e1e',
          background: '#121212',
        },
      },
    },
  },
} satisfies VOptions
