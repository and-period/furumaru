// @ts-check
import { createConfigForNuxt } from '@nuxt/eslint-config/flat';

export default createConfigForNuxt({
  features: {
    stylistic: {
      semi: true,
      quotes: 'single',
    },
  },
}).override(
  'nuxt/vue/rules', {
    rules: {
      'vue/no-multiple-template-root': 'off',
      'vue/multi-word-component-names': 'off',
    },
    files: ['src/pages/**', 'src/layouts/**', 'src/app.vue'],
  },
);
