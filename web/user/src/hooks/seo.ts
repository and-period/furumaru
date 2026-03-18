import type { MaybeRefOrGetter } from 'vue'

const SITE_NAME = '産地直送のお取り寄せ通販【ふるマル】'
const DEFAULT_DESCRIPTION = '産地直送のお取り寄せ通販のふるマルです。生産者のこだわりが「伝える」以上に「伝わる」ライブマルシェ'

export interface BreadcrumbItem {
  name: string
  path: string
}

function useSiteUrl(): string {
  const config = useRuntimeConfig()
  return (config.public.SITE_URL as string) || 'https://www.furumaru.and-period.co.jp'
}

interface SeoPageOptions {
  title: MaybeRefOrGetter<string>
  description?: MaybeRefOrGetter<string>
  ogImage?: MaybeRefOrGetter<string>
  ogType?: 'website' | 'article' | 'product'
  path?: MaybeRefOrGetter<string>
}

/**
 * ページ単位の SEO メタタグを設定する composable
 */
export function useSeoHead(options: SeoPageOptions) {
  const route = useRoute()
  const siteUrl = useSiteUrl()
  const defaultOgImage = `${siteUrl}/ogp/ogp.jpg`

  const resolvedTitle = computed(() => toValue(options.title))
  const resolvedDescription = computed(() => toValue(options.description) || DEFAULT_DESCRIPTION)
  const resolvedOgImage = computed(() => toValue(options.ogImage) || defaultOgImage)
  const resolvedPath = computed(() => toValue(options.path) || route.path)
  const canonicalUrl = computed(() => `${siteUrl}${resolvedPath.value}`)

  useSeoMeta({
    title: resolvedTitle,
    description: resolvedDescription,
    ogTitle: resolvedTitle,
    ogDescription: resolvedDescription,
    ogImage: resolvedOgImage,
    ogType: options.ogType || 'website',
    ogUrl: canonicalUrl,
    ogSiteName: SITE_NAME,
    twitterCard: 'summary_large_image',
    twitterTitle: resolvedTitle,
    twitterDescription: resolvedDescription,
    twitterImage: resolvedOgImage,
  })

  useHead({
    link: [
      { rel: 'canonical', href: canonicalUrl },
    ],
  })
}

interface ProductJsonLdOptions {
  name: MaybeRefOrGetter<string>
  description: MaybeRefOrGetter<string>
  images: MaybeRefOrGetter<string[]>
  price: MaybeRefOrGetter<number>
  inventory: MaybeRefOrGetter<number>
  ratingAverage?: MaybeRefOrGetter<number>
  ratingCount?: MaybeRefOrGetter<number>
  producerName?: MaybeRefOrGetter<string>
  url: MaybeRefOrGetter<string>
}

/**
 * Product の JSON-LD 構造化データを出力する composable
 */
export function useProductJsonLd(options: ProductJsonLdOptions) {
  const siteUrl = useSiteUrl()

  const jsonLd = computed(() => {
    const data: Record<string, unknown> = {
      '@context': 'https://schema.org',
      '@type': 'Product',
      'name': toValue(options.name),
      'description': toValue(options.description),
      'image': toValue(options.images),
      'url': `${siteUrl}${toValue(options.url)}`,
      'offers': {
        '@type': 'Offer',
        'price': toValue(options.price),
        'priceCurrency': 'JPY',
        'availability': toValue(options.inventory) > 0
          ? 'https://schema.org/InStock'
          : 'https://schema.org/OutOfStock',
      },
    }

    const ratingAvg = toValue(options.ratingAverage)
    const ratingCnt = toValue(options.ratingCount)
    if (ratingAvg && ratingCnt && ratingCnt > 0) {
      data.aggregateRating = {
        '@type': 'AggregateRating',
        'ratingValue': ratingAvg,
        'reviewCount': ratingCnt,
      }
    }

    const producer = toValue(options.producerName)
    if (producer) {
      data.brand = {
        '@type': 'Organization',
        'name': producer,
      }
    }

    return data
  })

  useHead({
    script: [
      {
        type: 'application/ld+json',
        innerHTML: computed(() => JSON.stringify(jsonLd.value)),
      },
    ],
  })
}

interface ExperienceJsonLdOptions {
  name: MaybeRefOrGetter<string>
  description: MaybeRefOrGetter<string>
  images: MaybeRefOrGetter<string[]>
  price: MaybeRefOrGetter<number>
  startAt?: MaybeRefOrGetter<number>
  endAt?: MaybeRefOrGetter<number>
  locationName?: MaybeRefOrGetter<string>
  address: MaybeRefOrGetter<string>
  postalCode?: MaybeRefOrGetter<string>
  latitude?: MaybeRefOrGetter<number>
  longitude?: MaybeRefOrGetter<number>
  url: MaybeRefOrGetter<string>
}

/**
 * Experience (Event) の JSON-LD 構造化データを出力する composable
 */
export function useExperienceJsonLd(options: ExperienceJsonLdOptions) {
  const siteUrl = useSiteUrl()

  const jsonLd = computed(() => {
    const data: Record<string, unknown> = {
      '@context': 'https://schema.org',
      '@type': 'Event',
      'name': toValue(options.name),
      'description': toValue(options.description),
      'image': toValue(options.images),
      'url': `${siteUrl}${toValue(options.url)}`,
      'eventAttendanceMode': 'https://schema.org/OfflineEventAttendanceMode',
      'offers': {
        '@type': 'Offer',
        'price': toValue(options.price),
        'priceCurrency': 'JPY',
        'availability': 'https://schema.org/InStock',
      },
      'location': {
        '@type': 'Place',
        'name': toValue(options.locationName) || '',
        'address': {
          '@type': 'PostalAddress',
          'streetAddress': toValue(options.address),
          'postalCode': toValue(options.postalCode) || '',
          'addressCountry': 'JP',
        },
      },
    }

    const startAt = toValue(options.startAt)
    if (startAt) {
      data.startDate = new Date(startAt * 1000).toISOString()
    }

    const endAt = toValue(options.endAt)
    if (endAt) {
      data.endDate = new Date(endAt * 1000).toISOString()
    }

    const lat = toValue(options.latitude)
    const lng = toValue(options.longitude)
    if (lat != null && lng != null) {
      (data.location as Record<string, unknown>).geo = {
        '@type': 'GeoCoordinates',
        'latitude': lat,
        'longitude': lng,
      }
    }

    return data
  })

  useHead({
    script: [
      {
        type: 'application/ld+json',
        innerHTML: computed(() => JSON.stringify(jsonLd.value)),
      },
    ],
  })
}

/**
 * BreadcrumbList の JSON-LD 構造化データを出力する composable
 */
export function useBreadcrumbJsonLd(items: MaybeRefOrGetter<BreadcrumbItem[]>) {
  const siteUrl = useSiteUrl()

  const jsonLd = computed(() => ({
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    'itemListElement': toValue(items).map((item, index) => ({
      '@type': 'ListItem',
      'position': index + 1,
      'name': item.name,
      'item': `${siteUrl}${item.path}`,
    })),
  }))

  useHead({
    script: [
      {
        type: 'application/ld+json',
        innerHTML: computed(() => JSON.stringify(jsonLd.value)),
      },
    ],
  })
}

/**
 * WebSite + Organization の JSON-LD を出力する composable（トップページ用）
 */
export function useWebSiteJsonLd() {
  const siteUrl = useSiteUrl()

  useHead({
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'WebSite',
          'name': SITE_NAME,
          'url': siteUrl,
        }),
      },
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'Organization',
          'name': 'ふるマル',
          'url': siteUrl,
          'logo': `${siteUrl}/ogp/ogp.jpg`,
          'sameAs': [
            'https://www.instagram.com/and_period',
            'https://twitter.com/and_period',
          ],
        }),
      },
    ],
  })
}
