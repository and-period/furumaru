import { defineConfig } from 'vite'
import { resolve } from 'path';
import vue from '@vitejs/plugin-vue'
import dts from 'vite-plugin-dts'

// https://vite.dev/config/
export default defineConfig({
  build: {
    lib: {
      entry: resolve(__dirname, './src/index.ts'),
      name: '@furumaru/shared',
      fileName: "index",
      formats: ['es', 'cjs', 'umd'],
    },
    rollupOptions: {
      external: ['vue'],
      output: {
        globals: {
          vue: 'Vue',
        },
      },
    },
  },
  plugins: [
    vue(),
    dts({ 
      include: ['src/**/**.ts', 'src/**/**.vue'],
      tsconfigPath: './tsconfig.app.json',
    })],
})
