module.exports = {
  root: true,
  env: {
    browser: true,
    node: true,
  },
  extends: [
    '@nuxtjs/eslint-config-typescript',
    'plugin:nuxt/recommended',
    'plugin:import/recommended',
    'plugin:import/typescript',
    'prettier',
  ],
  plugins: ['@typescript-eslint'],
  parserOptions: {
    tsconfigRootDir: '.',
  },
  settings: {
    'import/resolver': {
      node: {
        extensions: ['.ts', '.vue'],
      },
      'babel-module': {
        root: '.',
        alias: {
          '~': './src',
          '@': './src',
          '~~': '.',
          '@@': '.',
        },
      },
    },
  },
  // add your custom rules here
  rules: {
    'dot-notation': 'off',
    'no-unused-expressions': 'off',
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'sort-imports': 0,
    'import/order': [
      'error',
      {
        groups: [
          'builtin',
          'external',
          'parent',
          'sibling',
          'index',
          'object',
          'type',
        ],
        alphabetize: { order: 'asc' },
        'newlines-between': 'always',
      },
    ],
    'import/named': 'off',
    'vue/multi-word-component-names': 'off',
    'vue/no-mutating-props': 'off',
  },
}
