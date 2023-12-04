import type { PiniaPluginContext } from 'pinia'

function i18nInjector ({ app, store }: PiniaPluginContext) {
  // storeにi18nインスタンスを注入する
  store.i18n = app.$nuxt.$i18n
}

/**
 * Pinia に i18n をinjection する Nuxt プラグイン
 */
const i18nPlugin = defineNuxtPlugin(({ $pinia }) => {
  $pinia.use(i18nInjector)
})

export default i18nPlugin
