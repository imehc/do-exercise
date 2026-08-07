import path from 'path'
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { lingui, linguiTransformerBabelPreset } from '@lingui/vite-plugin'
import babel from '@rolldown/plugin-babel'
import tailwindcss from '@tailwindcss/vite'
import { tanstackRouter } from '@tanstack/router-plugin/vite'
import { collectStaticIconNames } from './scripts/collect-static-icon-names.ts'

const staticIconNames = collectStaticIconNames(
  path.resolve(import.meta.dirname, './src')
)

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tanstackRouter({
      target: 'react',
      autoCodeSplitting: true,
    }),
    react(),
    babel({
      presets: [linguiTransformerBabelPreset()],
    }),
    lingui(),
    tailwindcss(),
  ],
  server: {
    port: 6021,
    host: true,
    proxy: {
      '/apis': {
        target: 'http://127.0.0.1:6020',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/apis/, ''),
      },
      '/oss': {
        target: 'http://127.0.0.1:9000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/oss/, ''),
      },
    },
  },
  build: {
    // sourcemap: true,
    rolldownOptions: {
      preserveEntrySignatures: 'allow-extension',
      output: {
        codeSplitting: {
          groups: [
            {
              // 图标选择器通过 icon.tsx 的 glob 动态加载 6000+ 图标，逐个成包会
              // 产生大量碎片请求，滚动时明显卡顿。这里按首字母合并成 26 个 chunk。
              //
              // 业务代码具名引入的图标除外，原因见 scripts/collect-static-icon-names.ts。
              name: (moduleId) => {
                const matched = moduleId.match(
                  /@tabler[/\\]icons-react[/\\]dist[/\\]esm[/\\]icons[/\\](Icon(.)[A-Za-z0-9]*)\.mjs$/
                )
                if (!matched) return null
                if (staticIconNames.has(matched[1])) return null

                const initial = matched[2].toLowerCase()
                return initial >= 'a' && initial <= 'z'
                  ? `icons-${initial}`
                  : 'icons-misc'
              },
              // 只收图标模块本身。默认会连依赖闭包一起并入，导致
              // createReactComponent -> react 被吸进 icons-a，首屏白拉 38KB。
              includeDependenciesRecursively: false,
            },
          ],
        },
        strictExecutionOrder: true,
        minify: {
          compress: {
            dropConsole: true,
            dropDebugger: true,
          },
          codegen: {
            legalComments: 'none',
          },
        },
      },
    },
  },
  resolve: {
    alias: {
      '~': path.resolve(import.meta.dirname, './src'),
    },
  },
})
