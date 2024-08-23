// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs'

export default withNuxt(
  {
    ignores: [
      'bin/**/*',
      'src/assets/**/*',
      'src/public/**/*',
      'src/types/api/**/*',
    ],
  },
  {
    files: ['src/**/*.{vue,js,ts,tsx}'],
    rules: {
      '@typescript-eslint/ban-types': 'off',
      '@typescript-eslint/ban-ts-comment': 'off',
      '@typescript-eslint/unified-signatures': 'off',
      '@typescript-eslint/no-empty-object-type': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      '@typescript-eslint/no-var-requires': 'off',
      'func-call-spacing': 'off',
      'import/no-mutable-exports': 'off',
      'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
      'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
      'no-undef': 'off',
      'vue/no-mutating-props': 'warn',
      'vue/no-ref-as-operand': 'off',
      'vue/no-v-html': 'off',
      'vue/no-v-text-v-html-on-component': 'off',
    },
  },
)
