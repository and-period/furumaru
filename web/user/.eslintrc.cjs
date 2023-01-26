// https://github.com/nuxt/eslint-config#nuxteslint-config
module.exports = {
  root: true,
  extends: [
    '@nuxtjs/eslint-config-typescript'
  ],
  rules: {
    '@typescript-eslint/no-unused-vars': 'off',
    'no-unused-vars': [
      'error', {
        args: 'none',
        varsIgnorePattern: '^_',
        caughtErrorsIgnorePattern: '^_',
        destructuredArrayIgnorePattern: '^_'
      }
    ]
  }
}
