import 'pinia'

import { CustomErrorMessage } from '~/plugins/api-error-handler'
import { useAuthStore } from '~/store'

declare module 'vue/types/vue' {
  interface Vue {
    $auth: ReturnType<typeof useAuthStore>
  }
}

declare module '@nuxt/types' {
  interface Context {
    $auth: ReturnType<typeof useAuthStore>
  }

  interface NuxtAppOptions {
    $auth: ReturnType<typeof useAuthStore>
  }
}

declare module 'pinia' {
  export interface PiniaCustomProperties {
    /**
     * apiクライアントのエラーをハンドリングする関数
     * @param error 発生したエラーオブジェクト
     * @param customObject エラーメッセージをカスタマイズするオブジェクト
     * @returns Promise.rejectを返す。呼び出す側で再度returnすることを想定している。
     */
    errorHandler: (
      error: unknown,
      customObject?: CustomErrorMessage
    ) => Promise<never>
  }
}
