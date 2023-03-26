export interface ImageItem {
  url: string;
  size: number;
}

export interface MediaItem {
  url: string;
  isThumbnail: boolean;
  images: ImageItem[];
}

export interface ProductItem {
  id: string;
  name: string;
  description: string;
  producerId: string;
  storeName: string;
  categoryId: string;
  categoryName: string;
  productTypeId: string;
  productTypeName: string;
  productTypeIconUrl: string;
  public: boolean;
  inventory: number;
  weight: number;
  itemUnit: string;
  itemDescription: string;
  media: MediaItem[];
  price: number;
  deliveryType: number;
  box60Rate: number;
  box80Rate: number;
  box100Rate: number;
  originPrefecture: string;
  originCity: string;
  createdAt: number;
  updatedAt: number;
}
