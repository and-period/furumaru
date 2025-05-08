import { app as fbApp, app as fbApp } from './firebase'
import { type Analytics, getAnalytics, logEvent } from 'firebase/analytics'

export default defineNuxtPlugin((nuxtApp) => {
  // 開発環境のイベントはGAに送信しない
  if (process.env.NODE_ENV !== 'production') {
    return
  }

  const analytics: Analytics = getAnalytics(fbApp)

  nuxtApp.router?.afterEach((to, _) => {
    // GAにページ遷移情報を保存する
    console.log('analytics', analytics)
    logEvent(analytics, 'page_view', {
      page_location: location.hostname,
      page_path: to.fullPath,
      page_title: to.name || 'admin',
    })
  })
})
