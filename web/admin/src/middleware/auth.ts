import { Context, Middleware } from '@nuxt/types'
import Cookies from 'universal-cookie'

import { useAuthStore } from '~/store/auth'

const routing: Middleware = async (context: Context) => {
  const publicPages = ['/signin']

  const { isAuthenticated, getAuthByRefreshToken } = useAuthStore()

  if (!publicPages.includes(context.route.path)) {
    if (!isAuthenticated) {
      const cookies = new Cookies()
      const token = cookies.get('refreshToken')
      if (token) {
        await getAuthByRefreshToken(token).catch(() => {
          context.redirect('/signin')
        })
      } else {
        context.redirect('/signin')
      }
    }
  }
}

export default routing
