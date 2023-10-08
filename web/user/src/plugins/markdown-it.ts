import MarkdownIt from 'markdown-it'

const md = new MarkdownIt({
  html: true,
})

const mdPlugin = defineNuxtPlugin(() => {
  return {
    provide: {
      md,
    },
  }
})

export default mdPlugin
