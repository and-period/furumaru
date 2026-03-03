import { Chat } from '@ai-sdk/vue'
import { DefaultChatTransport } from 'ai'
import type { Ref } from 'vue'
import type { ProductFormDataForAi, ProductFormUpdate, FormFieldChange } from '~/types/ai'
import { useAuthStore } from '~/store'

const FIELD_LABELS: Record<string, string> = {
  name: '商品名',
  description: '商品説明',
  price: '販売価格',
  cost: '原価',
  inventory: '在庫数',
  weight: '重量',
  itemUnit: '単位',
  itemDescription: '内容量',
  deliveryType: '配送方法',
  storageMethodType: '保存方法',
  expirationDate: '賞味期限',
  recommendedPoint1: 'おすすめポイント1',
  recommendedPoint2: 'おすすめポイント2',
  recommendedPoint3: 'おすすめポイント3',
  originPrefectureCode: '原産地（都道府県）',
  originCity: '原産地（市区町村）',
  scope: '公開範囲',
  box60Rate: '60サイズ占有率',
  box80Rate: '80サイズ占有率',
  box100Rate: '100サイズ占有率',
}

function extractFormDataForAi(formData: Record<string, unknown>): ProductFormDataForAi {
  return {
    name: (formData.name as string) || '',
    description: (formData.description as string) || '',
    price: (formData.price as number) || 0,
    cost: (formData.cost as number) || 0,
    inventory: (formData.inventory as number) || 0,
    weight: (formData.weight as number) || 0,
    itemUnit: (formData.itemUnit as string) || '',
    itemDescription: (formData.itemDescription as string) || '',
    deliveryType: (formData.deliveryType as number) || 0,
    storageMethodType: (formData.storageMethodType as number) || 0,
    expirationDate: (formData.expirationDate as number) || 0,
    recommendedPoint1: (formData.recommendedPoint1 as string) || '',
    recommendedPoint2: (formData.recommendedPoint2 as string) || '',
    recommendedPoint3: (formData.recommendedPoint3 as string) || '',
    originPrefectureCode: (formData.originPrefectureCode as number) || 0,
    originCity: (formData.originCity as string) || '',
    scope: (formData.scope as number) || 0,
    box60Rate: (formData.box60Rate as number) || 0,
    box80Rate: (formData.box80Rate as number) || 0,
    box100Rate: (formData.box100Rate as number) || 0,
  }
}

export function computeFormChanges(
  currentForm: ProductFormDataForAi,
  update: ProductFormUpdate,
): FormFieldChange[] {
  const changes: FormFieldChange[] = []
  for (const [key, newValue] of Object.entries(update)) {
    if (newValue === undefined) {
      continue
    }
    const oldValue = currentForm[key as keyof ProductFormDataForAi]
    if (oldValue !== newValue) {
      changes.push({
        field: key,
        label: FIELD_LABELS[key] || key,
        oldValue,
        newValue,
      })
    }
  }
  return changes
}

export function useProductAiAssistant(formData: Ref<Record<string, unknown>>) {
  const isPanelOpen = ref(false)
  const input = ref('')
  const pendingUpdate = ref<ProductFormUpdate | null>(null)
  const pendingToolName = ref<string>('')
  const pendingChanges = ref<FormFieldChange[]>([])
  const pendingToolCallId = ref<string>('')

  const runtimeConfig = useRuntimeConfig()
  const authStore = useAuthStore()

  const chat = new Chat({
    transport: new DefaultChatTransport({
      api: `${runtimeConfig.public.API_BASE_URL}/v1/ai/chat`,
      streamProtocol: 'data',
      headers: () => ({
        Authorization: `Bearer ${authStore.accessToken}`,
      }),
      body: computed(() => ({
        formData: extractFormDataForAi(formData.value),
      })),
    }),
    onToolCall({ toolCall }) {
      const { toolName, args, toolCallId } = toolCall as unknown as {
        toolName: string
        args: Record<string, unknown>
        toolCallId: string
      }

      if (toolName === 'updateProductForm') {
        const update = args as ProductFormUpdate
        const currentFormForAi = extractFormDataForAi(formData.value)
        const changes = computeFormChanges(currentFormForAi, update)

        if (changes.length > 0) {
          pendingUpdate.value = update
          pendingToolName.value = toolName
          pendingChanges.value = changes
          pendingToolCallId.value = toolCallId
        }
      }

      if (toolName === 'suggestDescription') {
        const { description } = args as { description: string }
        pendingUpdate.value = { description }
        pendingToolName.value = toolName
        const currentFormForAi = extractFormDataForAi(formData.value)
        pendingChanges.value = computeFormChanges(currentFormForAi, { description })
        pendingToolCallId.value = toolCallId
      }

      if (toolName === 'suggestPoints') {
        const { point1, point2, point3 } = args as {
          point1: string
          point2: string
          point3: string
        }
        const update: ProductFormUpdate = {
          recommendedPoint1: point1,
          recommendedPoint2: point2,
          recommendedPoint3: point3,
        }
        pendingUpdate.value = update
        pendingToolName.value = toolName
        const currentFormForAi = extractFormDataForAi(formData.value)
        pendingChanges.value = computeFormChanges(currentFormForAi, update)
        pendingToolCallId.value = toolCallId
      }
    },
  })

  const messages = computed(() => chat.messages)
  const isChatLoading = computed(() => chat.status === 'streaming' || chat.status === 'submitted')
  const chatError = computed(() => chat.error)
  const hasPendingApproval = computed(() => pendingUpdate.value !== null)

  function applyUpdate() {
    if (!pendingUpdate.value) {
      return
    }

    for (const [key, value] of Object.entries(pendingUpdate.value)) {
      if (value !== undefined && key in formData.value) {
        (formData.value as Record<string, unknown>)[key] = value
      }
    }

    clearPending()
  }

  function rejectUpdate() {
    clearPending()
  }

  function clearPending() {
    pendingUpdate.value = null
    pendingToolName.value = ''
    pendingChanges.value = []
    pendingToolCallId.value = ''
  }

  function togglePanel() {
    isPanelOpen.value = !isPanelOpen.value
  }

  async function sendMessage() {
    const text = input.value.trim()
    if (!text) {
      return
    }
    input.value = ''
    await chat.sendMessage({ text })
  }

  return {
    // Panel state
    isPanelOpen,
    togglePanel,
    // Chat state
    messages,
    input,
    isChatLoading,
    chatError,
    sendMessage,
    // Tool approval
    hasPendingApproval,
    pendingChanges,
    pendingToolName,
    applyUpdate,
    rejectUpdate,
  }
}
