import { defineConfig } from 'vite'
import { resolve } from 'path'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { visualizer } from 'rollup-plugin-visualizer'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss(), visualizer({ open: true })],
  server: {
    host: '0.0.0.0',
    port: 6021,
    proxy: {
      '/api': {
        target: 'http://localhost:6020',
        changeOrigin: true,
        rewrite: path => path.replace(/^\/api/, '')
      }
    }
  },

  resolve: {
    alias: {
      // 设置别名后如果是typescript则必须在tsconfig.json中配置baseUrl和paths，否者ts类型会报错
      '~': resolve(__dirname, 'src'),
      '#': resolve(__dirname, 'mock')
    }
  },

  build: {
    manifest: true
  }
})
