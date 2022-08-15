import { Context, Plugin } from '@nuxt/types'
import { Analytics, getAnalytics, logEvent } from 'firebase/analytics'

import firebase from './firebase'

const googleAnalyticsPlugin: Plugin = ({ app }: Context) => {
  // 開発環境のイベントはGAに送信しない
  if (process.env.NODE_ENV !== 'production') return

  const analytics: Analytics = getAnalytics(firebase.app)

  app.router?.afterEach((to, _) => {
    // GAにページ遷移情報を保存する
    console.log(analytics)
    logEvent(analytics, 'page_view', {
      page_location: location.hostname,
      page_path: to.fullPath,
      page_title: to.name || 'admin',
    })
  })
}

export default googleAnalyticsPlugin
