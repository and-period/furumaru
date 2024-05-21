<template>
  <v-app>
    <v-main>
      <v-container>
        <v-card>
          <v-card-text class="markdown-content"v-html="legalNoticeMarkdown"  />
        </v-card>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import { marked } from 'marked'

export default {
  name: 'App',
  data () {
    return {
      legalNoticeMarkdown: ''
    }
  },
  mounted () {
    this.loadMarkdown()
  },
  methods: {
    async loadMarkdown () {
      const response = await fetch('/_content/legalPolicy.md')
      const markdownText = await response.text()
      this.legalNoticeMarkdown = marked(markdownText)
    }
  }
}

definePageMeta({
  layout: 'empty'
})
</script>

<style>
.markdown-content {
  font-size: 16px;
  line-height: 1.5;
}

.markdown-content ol {
  padding: inherit;
}

.markdown-content h1 {
  font-size: 1.5em;
  margin: 1em 0;
}

.markdown-content h2 {
  font-size: 1.2em;
  margin: 1em 0;
}

.markdown-content p {
  margin: 1em 0;
}
</style>
