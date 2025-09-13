import type { Prefecture } from '../enum'

export interface SearchAddress {
  prefecture: Prefecture
  city: string
  town: string
}
