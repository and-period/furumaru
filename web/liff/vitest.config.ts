import { defineVitestConfig } from '@nuxt/test-utils/config';

export default defineVitestConfig({
  test: {
    environment: 'nuxt',
    coverage: {
      provider: 'v8',
      exclude: ['nuxt.config.ts', 'vitest.config.ts'],
      all: true,
      include: [
        '**/composables/**/*.ts',
      ],
      reporter: ['text', 'json', 'html'],
    },
  },
});
