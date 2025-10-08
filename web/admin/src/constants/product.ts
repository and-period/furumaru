import { DeliveryType, ProductScope, ProductStatus, StorageMethodType } from '~/types/api/v1'

export const productStatuses = [
  { title: '予約販売', value: ProductStatus.ProductStatusPresale },
  { title: '販売中', value: ProductStatus.ProductStatusForSale },
  { title: '販売期間外', value: ProductStatus.ProductStatusOutOfSale },
  { title: '非公開', value: ProductStatus.ProductStatusPrivate },
  { title: '不明', value: ProductStatus.ProductStatusUnknown },
]

export const storageMethodTypes = [
  { title: '常温保存', value: StorageMethodType.StorageMethodTypeNormal },
  { title: '冷暗所保存', value: StorageMethodType.StorageMethodTypeCoolDark },
  { title: '冷蔵保存', value: StorageMethodType.StorageMethodTypeRefrigerated },
  { title: '冷凍保存', value: StorageMethodType.StorageMethodTypeFrozen },
]

export const deliveryTypes = [
  { title: '通常便', value: DeliveryType.DeliveryTypeNormal },
  { title: '冷蔵便', value: DeliveryType.DeliveryTypeRefrigerated },
  { title: '冷凍便', value: DeliveryType.DeliveryTypeFrozen },
]

export const productScopes = [
  { title: '全体公開', value: ProductScope.ProductScopePublic },
  { title: 'LINE限定', value: ProductScope.ProductScopeLimited },
  { title: '下書き', value: ProductScope.ProductScopePrivate },
]

export const productItemUnits = ['個', '瓶']
