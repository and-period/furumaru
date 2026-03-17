import {
  FeatureRequestCategory,
  FeatureRequestPriority,
  FeatureRequestStatus,
} from '~/types/feature-request'

export function useFeatureRequestLabels() {
  const getStatusLabel = (status: FeatureRequestStatus): string => {
    switch (status) {
      case FeatureRequestStatus.Waiting: return '受付中'
      case FeatureRequestStatus.Reviewing: return '検討中'
      case FeatureRequestStatus.Adopted: return '採用決定'
      case FeatureRequestStatus.InProgress: return '開発中'
      case FeatureRequestStatus.Done: return '完了'
      case FeatureRequestStatus.Rejected: return '却下'
      default: return '不明'
    }
  }

  const getStatusColor = (status: FeatureRequestStatus): string => {
    switch (status) {
      case FeatureRequestStatus.Waiting: return 'warning'
      case FeatureRequestStatus.Reviewing: return 'info'
      case FeatureRequestStatus.Adopted: return 'secondary'
      case FeatureRequestStatus.InProgress: return 'primary'
      case FeatureRequestStatus.Done: return 'success'
      case FeatureRequestStatus.Rejected: return 'error'
      default: return 'default'
    }
  }

  const getCategoryLabel = (category: FeatureRequestCategory): string => {
    switch (category) {
      case FeatureRequestCategory.UI: return 'UI/UX改善'
      case FeatureRequestCategory.Feature: return '新機能'
      case FeatureRequestCategory.Performance: return 'パフォーマンス改善'
      case FeatureRequestCategory.Other: return 'その他'
      default: return '不明'
    }
  }

  const getPriorityLabel = (priority: FeatureRequestPriority): string => {
    switch (priority) {
      case FeatureRequestPriority.Low: return '低'
      case FeatureRequestPriority.Medium: return '中'
      case FeatureRequestPriority.High: return '高'
      default: return '不明'
    }
  }

  const getPriorityColor = (priority: FeatureRequestPriority): string => {
    switch (priority) {
      case FeatureRequestPriority.Low: return 'default'
      case FeatureRequestPriority.Medium: return 'warning'
      case FeatureRequestPriority.High: return 'error'
      default: return 'default'
    }
  }

  return {
    getStatusLabel,
    getStatusColor,
    getCategoryLabel,
    getPriorityLabel,
    getPriorityColor,
  }
}
