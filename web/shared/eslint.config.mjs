import eslint from '@eslint/js';
import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'
import typescriptEslint from 'typescript-eslint';

export default typescriptEslint.config(
  { ignores: ['*.d.ts', '**/coverage', '**/dist'] },
  {
    extends: [
      eslint.configs.recommended,
      ...typescriptEslint.configs.recommended,
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
