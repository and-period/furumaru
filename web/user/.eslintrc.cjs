// https://github.com/nuxt/eslint-config#nuxteslint-config
module.exports = {
  root: true,
  extends: [
    '@nuxtjs/eslint-config-typescript',
    'plugin:tailwindcss/recommended',
    'prettier',
  ],
  rules: {
    '@typescript-eslint/no-unused-vars': 'off',
    'no-unused-vars': [
      'error',
      {
        args: 'none',
        varsIgnorePattern: '^_',
        caughtErrorsIgnorePattern: '^_',
        destructuredArrayIgnorePattern: '^_',
      },
    ],
  },
  overrides: [
    {
      files: ['src/pages/**', 'src/layouts/**'],
      rules: {
        'vue/no-multiple-template-root': 'off',
        'vue/multi-word-component-names': 'off',
      },
    },
  ],
}
