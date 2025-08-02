import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'
import typescriptEslint from 'typescript-eslint';

export default typescriptEslint.config(
  { ignores: ['*.d.ts', '**/coverage', '**/dist'] },
  {
    extends: [
      ...pluginVue.configs['flat/recommended'],
    ],
    rules: {},
    languageOptions: {
      sourceType: 'module',
      globals: {
        ...globals.browser
      },
      parserOptions: {
        parser: typescriptEslint.parser,
      },
    }
  }
)
