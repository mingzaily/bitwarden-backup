import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: './dist',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        // Vite 8's Rolldown expects a function for manual chunk routing.
        // Keep Vue isolated so the application chunks stay cache-friendly.
        manualChunks(id) {
          if (id.includes('node_modules/vue') || id.includes('node_modules/vue-router')) {
            return 'vue-vendor'
          }
          return undefined
        }
      }
    }
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
