<script lang="ts" setup>
import { mdiFormatBold, mdiFormatItalic, mdiFormatStrikethrough, mdiCodeTags, mdiCodeNotEqualVariant, mdiUndo, mdiRedo, mdiFormatListBulleted, mdiFormatListNumbered, mdiArrowSplitHorizontal, mdiKeyboardReturn } from '@mdi/js'

import { StarterKit } from '@tiptap/starter-kit'
import { useEditor, EditorContent } from '@tiptap/vue-3'

const props = defineProps({
  label: {
    type: String,
    default: '',
  },
  modelValue: {
    type: String,
    default: '',
  },
  errorMessage: {
    type: String,
    default: '',
  },
})

const emits = defineEmits<{
  (e: 'update:modelValue', html: string): void
}>()

const editor = useEditor({
  content: props.modelValue,
  extensions: [StarterKit],
  onUpdate: () => {
    if (editor.value) {
      emits('update:modelValue', editor.value.getHTML())
    }
  },
})

const fontSize = [
  { text: '標準', value: 0 },
  { text: '見出し「1」', value: 1 },
  { text: '見出し「2」', value: 2 },
  { text: '見出し「3」', value: 3 },
  { text: '見出し「4」', value: 4 },
  { text: '見出し「5」', value: 5 },
  { text: '見出し「6」', value: 6 },
]

const menus = [
  {
    label: 'bold',
    icon: mdiFormatBold,
    onClick: () => editor.value?.chain().focus().toggleBold().run(),
  },
  {
    label: 'italic',
    icon: mdiFormatItalic,
    onClick: () => editor.value?.chain().focus().toggleItalic().run(),
  },
  {
    label: 'strike',
    icon: mdiFormatStrikethrough,
    onClick: () => editor.value?.chain().focus().toggleStrike().run(),
  },
  {
    label: 'code',
    icon: mdiCodeTags,
    onClick: () => editor.value?.chain().focus().toggleCode().run(),
  },
  {
    label: 'codeBlock',
    icon: mdiCodeNotEqualVariant,
    onClick: () => editor.value?.chain().focus().toggleCodeBlock().run(),
  },
]

const activeMenus = computed<number[]>(() => {
  return menus
    .map((menu, i) => (editor.value?.isActive(menu.label) ? i : -1))
    .filter(item => item !== -1)
})

const activeStyle = computed(() => {
  if (!editor.value) {
    return 0
  }
  if (editor.value.isActive('heading', { level: 1 })) {
    return 1
  }
  else if (editor.value.isActive('heading', { level: 2 })) {
    return 2
  }
  else if (editor.value.isActive('heading', { level: 3 })) {
    return 3
  }
  else if (editor.value.isActive('heading', { level: 4 })) {
    return 4
  }
  else if (editor.value.isActive('heading', { level: 5 })) {
    return 5
  }
  else if (editor.value.isActive('heading', { level: 6 })) {
    return 6
  }
  return 0
})

const handleChangeTextStyle = (level: number): void => {
  if (!editor.value) {
    return
  }

  if (level === 0) {
    editor.value.chain().focus().setParagraph().run()
    return
  }
  // @ts-ignore
  editor.value.chain().focus().toggleHeading({ level }).run()
}
</script>

<template>
  <div v-if="editor">
    <p :class="errorMessage && 'text-error'">
      {{ label }}
    </p>
    <div
      class="d-flex flex-column gap pa-2 editor-menu"
      :class="errorMessage && 'error'"
    >
      <div class="d-flex align-center gap">
        <v-btn-toggle>
          <v-btn aria-label="元に戻す" @click="editor?.chain().focus().undo().run()">
            <v-icon :icon="mdiUndo" />
          </v-btn>
          <v-btn aria-label="やり直す" @click="editor?.chain().focus().redo().run()">
            <v-icon :icon="mdiRedo" />
          </v-btn>
        </v-btn-toggle>

        <v-select
          :model-value="activeStyle"
          :items="fontSize"
          item-title="text"
          item-value="value"
          density="compact"
          hide-details
          single-line
          @update:model-value="handleChangeTextStyle"
        />

        <v-btn-toggle
          v-model="activeMenus"
          multiple
        >
          <v-btn
            v-for="(menu, i) in menus"
            :key="i"
            :aria-label="menu.label"
            @click="menu.onClick"
          >
            <v-icon>{{ menu.icon }}</v-icon>
          </v-btn>
        </v-btn-toggle>
      </div>
      <div class="d-flex gap">
        <v-btn-toggle>
          <v-btn
            :class="{ 'is-active': editor.isActive('bulletList') }"
            aria-label="箇条書きリスト"
            @click="editor?.chain().focus().toggleBulletList().run()"
          >
            <v-icon :icon="mdiFormatListBulleted" />
          </v-btn>
          <v-btn
            :class="{ 'is-active': editor.isActive('orderedList') }"
            aria-label="番号付きリスト"
            @click="editor?.chain().focus().toggleOrderedList().run()"
          >
            <v-icon :icon="mdiFormatListNumbered" />
          </v-btn>
        </v-btn-toggle>

        <v-btn-toggle>
          <v-btn aria-label="水平線" @click="editor?.chain().focus().setHorizontalRule().run()">
            <v-icon :icon="mdiArrowSplitHorizontal" />
          </v-btn>
          <v-btn aria-label="改行" @click="editor?.chain().focus().setHardBreak().run()">
            <v-icon :icon="mdiKeyboardReturn" />
          </v-btn>
        </v-btn-toggle>
      </div>
    </div>
    <editor-content
      class="editor"
      :class="errorMessage && 'error'"
      :editor="editor"
    />
    <div
      v-show="errorMessage"
      role="alert"
      class="text-error text-caption mt-1"
    >
      {{ errorMessage }}
    </div>
  </div>
</template>

<style lang="scss" scoped>
$border-color: rgb(224,224,224);
$border-color-error: rgb(var(--v-theme-error));
$background-color-editor-menu: #F5F5F5;
$background-color-editor-menu-item: rgb(255,251,254);

.editor-menu {
  background-color: $background-color-editor-menu;
  border-width: 1px 1px 0 1px;
  border-style: solid;
  border-color: $border-color;
  border-radius: 4px 4px 0 0;

  &.error {
    border-color: $border-color-error;
  }

  :deep(.v-btn) {
    background-color: $background-color-editor-menu-item;
  }

  :deep(.v-field) {
    background-color: $background-color-editor-menu-item;
    padding-left: 4px;
  }

  :deep(.v-field__input) {
    padding-top: 4px;
  }

  :deep(.v-field__append-inner) {
    padding-top: 4px;
  }
}
.editor {
  border: 1px;
  border-style: solid;
  border-color: $border-color;
  border-radius: 0 0 4px 4px;

  &.error {
    border-color: $border-color-error;
  }

  :deep(.ProseMirror) {
    padding: 8px;
    min-height: 200px;
    max-height: 500px;
    overflow: scroll;

    li {
      margin-left: 16px;
      padding-left: 0;
    }

  }

  :deep(.ProseMirror-focused) {
    outline: none;
  }
}

.gap {
  gap: 4px;
}
</style>
