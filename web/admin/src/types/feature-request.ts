/**
 * 要望リクエスト機能 - ローカル型定義
 *
 * NOTE: バックエンド API 実装後は types/api/v1/ 配下の自動生成型に移行する。
 */

export const FeatureRequestStatus = {
  /** 受付中 */
  Waiting: 1,
  /** 検討中 */
  Reviewing: 2,
  /** 採用決定 */
  Adopted: 3,
  /** 開発中 */
  InProgress: 4,
  /** 完了 */
  Done: 5,
  /** 却下 */
  Rejected: 6,
} as const
export type FeatureRequestStatus = typeof FeatureRequestStatus[keyof typeof FeatureRequestStatus]

export const FeatureRequestCategory = {
  /** UI/UX改善 */
  UI: 1,
  /** 新機能 */
  Feature: 2,
  /** パフォーマンス改善 */
  Performance: 3,
  /** その他 */
  Other: 4,
} as const
export type FeatureRequestCategory = typeof FeatureRequestCategory[keyof typeof FeatureRequestCategory]

export const FeatureRequestPriority = {
  /** 低 */
  Low: 1,
  /** 中 */
  Medium: 2,
  /** 高 */
  High: 3,
} as const
export type FeatureRequestPriority = typeof FeatureRequestPriority[keyof typeof FeatureRequestPriority]

export interface FeatureRequest {
  id: string
  title: string
  description: string
  category: FeatureRequestCategory
  priority: FeatureRequestPriority
  status: FeatureRequestStatus
  /** 管理者コメント */
  note: string
  submittedBy: string
  submitterName: string
  createdAt: number
  updatedAt: number
}

export interface CreateFeatureRequestInput {
  title: string
  description: string
  category: FeatureRequestCategory
  priority: FeatureRequestPriority
  submittedBy: string
  submitterName: string
}

export interface UpdateFeatureRequestInput {
  status: FeatureRequestStatus
  note: string
}
