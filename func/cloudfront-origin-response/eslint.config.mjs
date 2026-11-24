import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  {
    ignores: ['node_modules/**', 'dist/**', 'app.js'],
  },
  tseslint.configs.recommended,
  {
    languageOptions: {
      ecmaVersion: 2022,
      sourceType: 'module',
    },
  },
  // Must be last: enables eslint-plugin-prettier and eslint-config-prettier.
  eslintPluginPrettierRecommended,
);
