import type {
  CreateFeatureRequestInput,
  FeatureRequest,
  UpdateFeatureRequestInput,
} from '~/types/feature-request'
import { FeatureRequestStatus } from '~/types/feature-request'

/**
 * 要望リクエスト Store
 *
 * NOTE: バックエンド API 実装前の暫定実装として localStorage を使用している。
 *       各 action 内の "TODO: Replace with API call" コメント箇所を実 API 呼び出しに差し替えること。
 */
export const useFeatureRequestStore = defineStore('feature-request', () => {
  const STORAGE_KEY = 'furumaru:feature-requests'

  const featureRequest = ref<FeatureRequest>({} as FeatureRequest)
  const featureRequests = ref<FeatureRequest[]>([])
  const total = ref<number>(0)

  function loadFromStorage(): FeatureRequest[] {
    if (!import.meta.client)
      return []
    try {
      const data = localStorage.getItem(STORAGE_KEY)
      return data ? JSON.parse(data) : []
    }
    catch {
      return []
    }
  }

  function saveToStorage(requests: FeatureRequest[]): void {
    if (!import.meta.client)
      return
    localStorage.setItem(STORAGE_KEY, JSON.stringify(requests))
  }

  /**
   * 要望一覧取得
   * TODO: Replace with API call → GET /v1/feature-requests
   */
  async function fetchFeatureRequests(
    limit = 20,
    offset = 0,
    submittedBy?: string,
  ): Promise<void> {
    const all = loadFromStorage()
    const filtered = submittedBy
      ? all.filter(r => r.submittedBy === submittedBy)
      : all
    total.value = filtered.length
    featureRequests.value = filtered.slice(offset, offset + limit)
  }

  /**
   * 要望詳細取得
   * TODO: Replace with API call → GET /v1/feature-requests/:id
   */
  async function getFeatureRequest(id: string): Promise<void> {
    const all = loadFromStorage()
    const found = all.find(r => r.id === id)
    if (!found)
      throw new Error('対象の要望リクエストが存在しません')
    featureRequest.value = found
  }

  /**
   * 要望新規作成
   * TODO: Replace with API call → POST /v1/feature-requests
   */
  async function createFeatureRequest(
    input: CreateFeatureRequestInput,
  ): Promise<void> {
    const all = loadFromStorage()
    const now = Math.floor(Date.now() / 1000)
    const newRequest: FeatureRequest = {
      id: crypto.randomUUID(),
      title: input.title,
      description: input.description,
      category: input.category,
      priority: input.priority,
      status: FeatureRequestStatus.Waiting,
      note: '',
      submittedBy: input.submittedBy,
      submitterName: input.submitterName,
      createdAt: now,
      updatedAt: now,
    }
    all.unshift(newRequest)
    saveToStorage(all)
  }

  /**
   * 要望更新（ステータス変更・コメント）
   * TODO: Replace with API call → PATCH /v1/feature-requests/:id
   */
  async function updateFeatureRequest(
    id: string,
    input: UpdateFeatureRequestInput,
  ): Promise<void> {
    const all = loadFromStorage()
    const index = all.findIndex(r => r.id === id)
    if (index === -1)
      throw new Error('対象の要望リクエストが存在しません')
    all[index] = {
      ...all[index],
      status: input.status,
      note: input.note,
      updatedAt: Math.floor(Date.now() / 1000),
    }
    saveToStorage(all)
    featureRequest.value = all[index]
  }

  /**
   * 要望削除
   * TODO: Replace with API call → DELETE /v1/feature-requests/:id
   */
  async function deleteFeatureRequest(id: string): Promise<void> {
    const all = loadFromStorage()
    const filtered = all.filter(r => r.id !== id)
    if (filtered.length === all.length)
      throw new Error('対象の要望リクエストが存在しません')
    saveToStorage(filtered)
  }

  return {
    featureRequest,
    featureRequests,
    total,
    fetchFeatureRequests,
    getFeatureRequest,
    createFeatureRequest,
    updateFeatureRequest,
    deleteFeatureRequest,
  }
})
