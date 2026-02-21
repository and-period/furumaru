import { md3 } from 'vuetify/blueprints'
import { ja } from 'vuetify/locale'
import type { VuetifyOptions } from 'vuetify'

export default {
  locale: {
    locale: 'ja',
    messages: { ja },
  },
  blueprint: md3,
  defaults: {
    VAutocomplete: {
      variant: 'underlined',
    },
    VAvatar: {
      variant: 'outlined',
      color: 'grey-lighten-1',
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
          primary: '#689f38',
          primaryLight: '#c5e1a5',
          accent: '#ffb300',
          secondary: '#ff8f00',
          info: '#26a69a',
          warning: '#ffc107',
          error: '#dd2c00',
          unknown: '#616161',
          success: '#00e676',
          white: '#ffffff',
        },
      },
      dark: {
        dark: true,
        colors: {
          primary: '#b9f6ca',
          accent: '#424242',
          secondary: '#ff8f00',
          info: '#26a69a',
          warning: '#ffc107',
          error: '#dd2c00',
          success: '#00e676',
        },
      },
    },
  },
} satisfies VuetifyOptions
