import VueI18n from 'vue-i18n'

export interface HeaderMenuItem {
  name: VueI18n.TranslateResult | string
  onClick: Function
}
