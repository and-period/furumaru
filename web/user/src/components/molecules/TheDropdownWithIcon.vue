<script lang="ts" setup>
interface Props {
  triggerAriaLabel?: string
}

withDefaults(defineProps<Props>(), {
  triggerAriaLabel: 'メニューを開閉',
})

interface Expose {
  open: () => void
  close: () => void
}

const isShow = ref<boolean>(false)
const dropdownArea = ref<HTMLElement | null>(null)
const triggerButtonRef = ref<HTMLElement | null>(null)

const handleIconClick = () => {
  isShow.value = !isShow.value
}

const handleCloseIconClick = () => {
  isShow.value = false
}

// isShowがtrueになってから0.3秒経過しているかのフラグ
let isShowFlag = false

// フォーカストラップ: ドロップダウン内のフォーカス可能要素を循環
const trapFocus = (e: KeyboardEvent) => {
  if (!isShow.value || !dropdownArea.value) return
  if (e.key !== 'Tab') return

  const focusable = dropdownArea.value.querySelectorAll<HTMLElement>(
    'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
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

// Escape キーでドロップダウンを閉じる
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isShow.value) {
    isShow.value = false
  }
}

watch(isShow, (newValue) => {
  if (newValue === true) {
    setTimeout(() => {
      isShowFlag = true
    }, 300)
    // 開いたらドロップダウン内の最初のフォーカス可能要素にフォーカス
    nextTick(() => {
      const firstFocusable = dropdownArea.value?.querySelector<HTMLElement>(
        'button, a[href], input, select, textarea, [tabindex]:not([tabindex="-1"])',
      )
      firstFocusable?.focus()
    })
  }
  else {
    isShowFlag = false
    // 閉じたらトリガーボタンにフォーカスを戻す
    nextTick(() => {
      triggerButtonRef.value?.focus()
    })
  }
})

const clickOutside = (e: MouseEvent) => {
  if (e.target instanceof Node && !dropdownArea.value?.contains(e.target)) {
    if (isShow.value === true) {
      // isShowがtrueになってから0.3秒経過していない場合は、クリックイベントを無視する
      if (!isShowFlag) {
        return
      }
      isShow.value = false
    }
  }
}

const handleOpen = () => {
  isShow.value = true
}

onMounted(() => {
  addEventListener('click', clickOutside)
  addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
  removeEventListener('click', clickOutside)
  removeEventListener('keydown', handleKeydown)
})

defineExpose<Expose>({
  open: handleOpen,
  close: handleCloseIconClick,
})
</script>

<template>
  <div
    ref="dropdownArea"
    class="relative"
  >
    <the-icon-button
      ref="triggerButtonRef"
      :aria-label="triggerAriaLabel"
      :aria-expanded="isShow"
      aria-haspopup="true"
      @click="handleIconClick"
    >
      <slot name="icon" />
    </the-icon-button>
    <the-dropdown-area
      v-show="isShow"
      class="absolute right-0 min-w-[240px] before:absolute before:-right-7 before:-top-8 before:-translate-x-1/2 before:border-[24px] before:border-transparent before:border-b-white before:content-[''] md:-right-8 md:before:right-1"
      @keydown="trapFocus"
    >
      <div class="px-4 text-right">
        <the-icon-button
          aria-label="閉じる"
          @click="handleCloseIconClick"
        >
          <the-close-icon />
        </the-icon-button>
      </div>
      <slot name="content" />
    </the-dropdown-area>
  </div>
</template>
