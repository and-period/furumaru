// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs'

export default withNuxt(
  {
    ignores: [
      'public/**/*',
      'src/assets/**/*',
      'src/types/api/**/*',
    ],
  },
  {
    files: ['src/pages/**', 'src/layouts/**', 'src/app.vue'],
    rules: {
      'vue/no-multiple-template-root': 'off',
      'vue/multi-word-component-names': 'off',
    },
  },
  {
    files: ['src/**/*.{vue,js,ts,tsx}'],
    rules: {
      '@typescript-eslint/ban-ts-comment': 'off',
      '@typescript-eslint/consistent-type-imports': 'error',
      '@typescript-eslint/no-unused-vars': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/unified-signatures': 'off',
      'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
      'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
      'no-unused-vars': [
        'error',
        {
          args: 'none',
          varsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '^_',
          destructuredArrayIgnorePattern: '^_',
        },
      ],
      'vue/no-v-html': 'off',
    },
  },
)
