interface BreadcrumbItem {
  title: string
  to?: string
  disabled?: boolean
}

const routeLabels: Record<string, string> = {
  '/': 'ホーム',
  '/products': '商品管理',
  '/products/new': '新規作成',
  '/orders': '注文管理',
  '/schedules': 'ライブ配信',
  '/schedules/new': '新規作成',
  '/videos': '動画管理',
  '/videos/new': '新規作成',
  '/producers': '生産者管理',
  '/producers/new': '新規作成',
  '/experiences': '体験管理',
  '/experiences/new': '新規作成',
  '/customers': '顧客管理',
  '/notifications': 'お知らせ情報',
  '/notifications/new': '新規作成',
  '/promotions': 'セール情報',
  '/promotions/new': '新規作成',
  '/accounts': 'マイページ',
  '/accounts/email': 'メール設定',
  '/accounts/password': 'パスワード設定',
  '/accounts/providers': 'プロバイダー設定',
  '/system': 'システム設定',
  '/administrators': '管理者管理',
  '/administrators/new': '新規作成',
  '/coordinators': 'コーディネーター管理',
  '/coordinators/new': '新規作成',
  '/categories': 'カテゴリ管理',
  '/contacts': 'お問い合わせ管理',
  '/experience-types': '体験種別管理',
  '/payment-systems': '決済設定',
  '/product-tags': '商品タグ管理',
  '/shippings': '配送設定',
  '/spot-types': 'スポット種別管理',
  '/messages': 'メッセージ',
}

export function useBreadcrumbs(): ComputedRef<BreadcrumbItem[]> {
  const route = useRoute()

  return computed(() => {
    const path = route.path

    if (path === '/') {
      return [{ title: 'ホーム', disabled: true }]
    }

    const items: BreadcrumbItem[] = [
      { title: 'ホーム', to: '/' },
    ]

    const segments = path.split('/').filter(Boolean)
    let currentPath = ''

    for (let i = 0; i < segments.length; i++) {
      currentPath += `/${segments[i]}`
      const isLast = i === segments.length - 1
      const label = routeLabels[currentPath]

      if (label) {
        items.push({
          title: label,
          to: isLast ? undefined : currentPath,
          disabled: isLast,
        })
      }
      else {
        // Dynamic segment (UUID etc.) - use parent label + "詳細"
        const parentPath = currentPath.split('/').slice(0, -1).join('/')
        const parentLabel = routeLabels[parentPath]
        items.push({
          title: parentLabel ? `${parentLabel}詳細` : '詳細',
          disabled: isLast,
          to: isLast ? undefined : currentPath,
        })
      }
    }

    return items
  })
}
