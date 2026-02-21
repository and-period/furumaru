<script lang="ts" setup>
import { marked } from 'marked'

const helpMarkdown = ref('')

onMounted(async () => {
  const response = await fetch('/_content/help.md')
  const markdownText = await response.text()
  helpMarkdown.value = await marked(markdownText)
})
</script>

<template>
  <v-container>
    <v-card>
      <v-card-text
        class="markdown-content"
        v-html="helpMarkdown"
      />
    </v-card>
  </v-container>
</template>

<style scoped>
.markdown-content {
  font-size: 16px;
  line-height: 1.8;
}

.markdown-content :deep(h1) {
  font-size: 1.6em;
  font-weight: bold;
  margin: 1em 0 0.5em;
  padding-bottom: 0.3em;
  border-bottom: 2px solid #e0e0e0;
}

.markdown-content :deep(h2) {
  font-size: 1.3em;
  font-weight: bold;
  margin: 1.5em 0 0.5em;
  padding-bottom: 0.2em;
  border-bottom: 1px solid #e0e0e0;
}

.markdown-content :deep(h3) {
  font-size: 1.1em;
  font-weight: bold;
  margin: 1.2em 0 0.4em;
}

.markdown-content :deep(p) {
  margin: 0.8em 0;
}

.markdown-content :deep(ul) {
  padding-left: 1.5em;
  margin: 0.5em 0;
}

.markdown-content :deep(li) {
  margin: 0.3em 0;
}

.markdown-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 1em 0;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid #e0e0e0;
  padding: 0.6em 1em;
  text-align: left;
}

.markdown-content :deep(th) {
  background-color: #f5f5f5;
  font-weight: bold;
}

.markdown-content :deep(hr) {
  margin: 2em 0;
  border: none;
  border-top: 1px solid #e0e0e0;
}

.markdown-content :deep(code) {
  background-color: #f5f5f5;
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-size: 0.9em;
}
</style>
