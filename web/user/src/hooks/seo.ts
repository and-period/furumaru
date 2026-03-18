import type { MaybeRefOrGetter } from 'vue'

const SITE_NAME = '産地直送のお取り寄せ通販【ふるマル】'
const SITE_URL = 'https://www.furumaru.and-period.co.jp'
const DEFAULT_OG_IMAGE = `${SITE_URL}/ogp/ogp.jpg`
const DEFAULT_DESCRIPTION = '産地直送のお取り寄せ通販のふるマルです。生産者のこだわりが「伝える」以上に「伝わる」ライブマルシェ'

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

  const resolvedTitle = computed(() => toValue(options.title))
  const resolvedDescription = computed(() => toValue(options.description) || DEFAULT_DESCRIPTION)
  const resolvedOgImage = computed(() => toValue(options.ogImage) || DEFAULT_OG_IMAGE)
  const resolvedPath = computed(() => toValue(options.path) || route.path)
  const canonicalUrl = computed(() => `${SITE_URL}${resolvedPath.value}`)

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
  const jsonLd = computed(() => {
    const data: Record<string, unknown> = {
      '@context': 'https://schema.org',
      '@type': 'Product',
      'name': toValue(options.name),
      'description': toValue(options.description),
      'image': toValue(options.images),
      'url': `${SITE_URL}${toValue(options.url)}`,
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
  const jsonLd = computed(() => {
    const data: Record<string, unknown> = {
      '@context': 'https://schema.org',
      '@type': 'Event',
      'name': toValue(options.name),
      'description': toValue(options.description),
      'image': toValue(options.images),
      'url': `${SITE_URL}${toValue(options.url)}`,
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

    const lat = toValue(options.latitude)
    const lng = toValue(options.longitude)
    if (lat && lng) {
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

interface BreadcrumbItem {
  name: string
  path: string
}

/**
 * BreadcrumbList の JSON-LD 構造化データを出力する composable
 */
export function useBreadcrumbJsonLd(items: MaybeRefOrGetter<BreadcrumbItem[]>) {
  const jsonLd = computed(() => ({
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    'itemListElement': toValue(items).map((item, index) => ({
      '@type': 'ListItem',
      'position': index + 1,
      'name': item.name,
      'item': `${SITE_URL}${item.path}`,
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
  useHead({
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'WebSite',
          'name': SITE_NAME,
          'url': SITE_URL,
        }),
      },
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'Organization',
          'name': 'ふるマル',
          'url': SITE_URL,
          'logo': `${SITE_URL}/ogp/ogp.jpg`,
          'sameAs': [
            'https://www.instagram.com/and_period',
            'https://twitter.com/and_period',
          ],
        }),
      },
    ],
  })
}
