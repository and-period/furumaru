module.exports = {
  customSyntax: 'postcss-html',
  extends: [
    'stylelint-config-standard', // CSS の基本ルール
    'stylelint-config-recommended-vue', // <style lang="css|scss"> を含む Vue コンポーネント対応
  ],
  // add your custom config here
  // https://stylelint.io/user-guide/configuration
  rules: {},
}
