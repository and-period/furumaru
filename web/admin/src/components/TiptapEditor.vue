<template>
  <div v-if="editor">
    <p>{{ label }}</p>

    <div class="mb-2">
      <v-btn-toggle :value="[]" dense>
        <v-btn @click="editor.chain().focus().undo().run()"
          ><v-icon>mdi-undo</v-icon></v-btn
        >
        <v-btn @click="editor.chain().focus().redo().run()"
          ><v-icon>mdi-redo</v-icon></v-btn
        >
      </v-btn-toggle>

      <v-btn-toggle multiple :value="activeMenus" dense>
        <v-btn v-for="(menu, i) in menus" :key="i" @click="menu.onClick">
          <v-icon>{{ menu.icon }}</v-icon>
        </v-btn>
      </v-btn-toggle>

      <v-btn-toggle dense>
        <v-btn
          :class="{ 'is-active': editor.isActive('bulletList') }"
          @click="editor.chain().focus().toggleBulletList().run()"
        >
          <v-icon>mdi-format-list-bulleted</v-icon>
        </v-btn>
        <v-btn
          :class="{ 'is-active': editor.isActive('orderedList') }"
          @click="editor.chain().focus().toggleOrderedList().run()"
        >
          <v-icon>mdi-format-list-numbered</v-icon>
        </v-btn>
      </v-btn-toggle>

      <v-btn-toggle :value="[]" dense>
        <v-btn @click="editor.chain().focus().setHorizontalRule().run()">
          <v-icon>mdi-arrow-split-horizontal</v-icon>
        </v-btn>
        <v-btn @click="editor.chain().focus().setHardBreak().run()">
          <v-icon>mdi-keyboard-return</v-icon>
        </v-btn>
      </v-btn-toggle>

      <div class="d-flex mt-1">
        <v-overflow-btn
          :value="activeStyle"
          :items="fontSize"
          hide-details
          dense
          class="pa-0 mt-0"
          label="スタイル"
          @change="handleChangeTextStyle"
        />
        <v-spacer />
      </div>
    </div>
    <editor-content class="editor" :editor="editor" />
  </div>
</template>

<script lang="ts">
import StarterKit from '@tiptap/starter-kit'
import { EditorContent, Editor } from '@tiptap/vue-2'
import {
  defineComponent,
  reactive,
  onBeforeUnmount,
  computed,
} from '@vue/composition-api'

export default defineComponent({
  components: {
    EditorContent,
  },
  props: {
    label: {
      type: String,
      default: '',
    },
    value: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const editor = reactive<Editor>(
      new Editor({
        content: props.value,
        extensions: [StarterKit],
        onUpdate: (p) => {
          emit('update:value', p.editor.getHTML())
        },
      })
    )

    const fontSize = [
      { text: 'Paragraph', value: 0 },
      { text: 'H1', value: 1 },
      { text: 'H2', value: 2 },
      { text: 'H3', value: 3 },
      { text: 'H4', value: 4 },
      { text: 'H5', value: 5 },
      { text: 'H6', value: 6 },
    ]

    const menus = [
      {
        label: 'bold',
        icon: 'mdi-format-bold',
        onClick: () => editor.chain().focus().toggleBold().run(),
      },
      {
        label: 'italic',
        icon: 'mdi-format-italic',
        onClick: () => editor.chain().focus().toggleItalic().run(),
      },
      {
        label: 'strike',
        icon: 'mdi-format-strikethrough',
        onClick: () => editor.chain().focus().toggleStrike().run(),
      },
      {
        label: 'code',
        icon: 'mdi-code-tags',
        onClick: () => editor.chain().focus().toggleCode().run(),
      },
      {
        label: 'codeBlock',
        icon: 'mdi-code-not-equal-variant',
        onClick: () => editor.chain().focus().toggleCodeBlock().run(),
      },
    ]

    const activeMenus = computed<number[]>(() => {
      return menus
        .map((menu, i) => (editor.isActive(menu.label) ? i : -1))
        .filter((item) => item !== -1)
    })

    const activeStyle = computed(() => {
      if (!editor) {
        return 0
      }
      if (editor.isActive('heading', { level: 1 })) {
        return 1
      } else if (editor.isActive('heading', { level: 2 })) {
        return 2
      } else if (editor.isActive('heading', { level: 3 })) {
        return 3
      } else if (editor.isActive('heading', { level: 4 })) {
        return 4
      } else if (editor.isActive('heading', { level: 5 })) {
        return 5
      } else if (editor.isActive('heading', { level: 6 })) {
        return 6
      }
      return 0
    })

    onBeforeUnmount(() => {
      editor.destroy()
    })

    const handleChangeTextStyle = (level: number): void => {
      if (!editor) {
        return
      }

      if (level === 0) {
        editor.chain().focus().setParagraph().run()
        return
      }
      // @ts-ignore
      editor.chain().focus().toggleHeading({ level }).run()
    }

    return {
      editor,
      activeMenus,
      activeStyle,
      menus,
      fontSize,
      handleChangeTextStyle,
    }
  },
})
</script>

<style lang="scss" scoped>
.editor {
  ::v-deep .ProseMirror {
    border: solid var(--v-secondary-lighten5);
    border-radius: 4px;
    padding: 4px;
    min-height: 200px;
    max-height: 500px;
    overflow: scroll;
  }

  ::v-deep .ProseMirror-focused {
    outline: none;
  }
}
</style>
