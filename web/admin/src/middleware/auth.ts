import { Context, Middleware } from '@nuxt/types'
import Cookies from 'universal-cookie'

import { useAuthStore } from '~/store/auth'

const routing: Middleware = async (context: Context) => {
  const publicPages = ['/signin']
  const path = context.route.path

  const { isAuthenticated, getAuthByRefreshToken, setRedirectPath } =
    useAuthStore()

  if (!publicPages.includes(context.route.path)) {
    if (!isAuthenticated) {
      const cookies = new Cookies()
      const token = cookies.get('refreshToken')
      setRedirectPath(path)
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
