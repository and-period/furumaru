<script setup lang="ts">
import { computed, watch, ref, nextTick, onMounted, onBeforeUnmount, useId } from 'vue'

interface Props {
  open: boolean
  title?: string
  ariaLabel?: string
  confirmText?: string
  cancelText?: string
  variant?: 'default' | 'danger'
}

interface Emits {
  (e: 'update:open', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  ariaLabel: '',
  confirmText: '確認',
  cancelText: 'キャンセル',
  variant: 'default',
})

const emits = defineEmits<Emits>()

const dialogRef = ref<HTMLElement | null>(null)
const titleId = `fm-dialog-title-${useId()}`

// title か ariaLabel のどちらかが必須
if (!props.title && !props.ariaLabel) {
  console.warn('[FmDialog] title または ariaLabel のどちらかを指定してください。ダイアログにアクセシブルネームがありません。')
}

const confirmVariantClass = computed(() =>
  props.variant === 'danger'
    ? 'bg-error text-white hover:bg-error/90'
    : 'bg-orange text-white hover:bg-orange/90',
)

const handleConfirm = () => {
  emits('confirm')
  emits('update:open', false)
}

const handleCancel = () => {
  emits('cancel')
  emits('update:open', false)
}

const handleOverlayClick = (e: MouseEvent) => {
  if (e.target === e.currentTarget) {
    handleCancel()
  }
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    handleCancel()
    return
  }

  // フォーカストラップ
  if (e.key !== 'Tab' || !dialogRef.value) return

  const focusable = dialogRef.value.querySelectorAll<HTMLElement>(
    'button:not([disabled]), a[href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
  )
  if (focusable.length === 0) return

  const first = focusable[0]
  const last = focusable[focusable.length - 1]

  if (e.shiftKey && document.activeElement === first) {
    e.preventDefault()
    last.focus()
  }
  else if (!e.shiftKey && document.activeElement === last) {
    e.preventDefault()
    first.focus()
  }
}

// 開閉時のフォーカス管理とスクロールロック
let previousFocus: HTMLElement | null = null
let previousOverflow: string = ''

const applyOpenState = (isOpen: boolean) => {
  if (typeof document === 'undefined') return

  if (isOpen) {
    previousFocus = document.activeElement as HTMLElement
    previousOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    nextTick(() => {
      const firstFocusable = dialogRef.value?.querySelector<HTMLElement>(
        'button:not([disabled]), a[href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
      )
      ;(firstFocusable ?? dialogRef.value)?.focus()
    })
  }
  else {
    document.body.style.overflow = previousOverflow
    nextTick(() => {
      previousFocus?.focus()
    })
  }
}

watch(() => props.open, applyOpenState)

onMounted(() => {
  if (props.open) {
    applyOpenState(true)
  }
})

onBeforeUnmount(() => {
  if (typeof document !== 'undefined') {
    document.body.style.overflow = previousOverflow
  }
})
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50"
      @click="handleOverlayClick"
      @keydown="handleKeydown"
    >
      <div
        ref="dialogRef"
        role="dialog"
        aria-modal="true"
        tabindex="-1"
        :aria-labelledby="title ? titleId : undefined"
        :aria-label="!title && ariaLabel ? ariaLabel : undefined"
        class="mx-4 w-full max-w-md rounded-lg bg-white p-6 shadow-xl"
      >
        <h2
          v-if="title"
          :id="titleId"
          class="mb-4 text-lg font-semibold text-main"
        >
          {{ title }}
        </h2>

        <div class="mb-6 text-sm text-typography">
          <slot />
        </div>

        <div class="flex justify-end gap-3">
          <slot name="actions">
            <button
              type="button"
              class="px-4 py-2 text-sm font-medium text-typography hover:bg-gray-100 transition-colors"
              @click="handleCancel"
            >
              {{ cancelText }}
            </button>
            <button
              type="button"
              :class="[
                'px-4 py-2 text-sm font-medium transition-colors',
                confirmVariantClass,
              ]"
              @click="handleConfirm"
            >
              {{ confirmText }}
            </button>
          </slot>
        </div>
      </div>
    </div>
  </Teleport>
</template>
