import type { Composer, UseI18nOptions } from 'vue-i18n'
import { ExperienceStatus, ProductStatus } from '~/types/api'
import type { I18n } from '~/types/locales/i18n'

/**
 * 体験のステータスに応じて文言を出しわける関数
 * @param status
 * @returns
 */
export function experienceStatusToString(
  status: ExperienceStatus,
  i18n: Composer<NonNullable<UseI18nOptions['messages']>, NonNullable<UseI18nOptions['datetimeFormats']>, NonNullable<UseI18nOptions['numberFormats']>, UseI18nOptions['locale'] extends unknown ? string : UseI18nOptions['locale']>,
): string {
  const statusText = (str: keyof I18n['items']['list']) => {
    return i18n.t(`items.list.${str}`)
  }

  switch (status) {
    case ExperienceStatus.FINISHED:
      return statusText('forSaleText')
    case ProductStatus.OUT_OF_SALES:
      return statusText('outOfSalesText')
    case ExperienceStatus.WAITING:
      return statusText('presalesText')
    default:
      return statusText('unknownItemText')
  }
}
