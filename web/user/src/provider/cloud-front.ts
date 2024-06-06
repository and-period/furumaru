import { joinURL } from 'ufo'
import type { ProviderGetImage, ImageModifiers } from '@nuxt/image'

/**
 * nuxt imageで指定されたパラメータをクエリパラメータに変数する関数
 * @param modifiers
 * @returns
 */
const transformQueryParamsFromModifiers = (
  modifiers: Partial<ImageModifiers>,
): string => {
  const params = Object.keys(modifiers)
    .filter(key => modifiers[key] !== undefined)
    .map((key) => {
      const value = encodeURIComponent(modifiers[key])
      return `${encodeURIComponent(key)}=${value}`
    })

  return params.join('&')
}

export const getImage: ProviderGetImage = (
  src,
  { modifiers = {}, baseURL } = {},
) => {
  if (!baseURL) {
    // also support runtime config
    baseURL = useRuntimeConfig().public.siteUrl
  }

  const queryParams = transformQueryParamsFromModifiers(modifiers)

  return {
    url: joinURL(baseURL, src + (queryParams ? '?' + queryParams : '')),
  }
}
