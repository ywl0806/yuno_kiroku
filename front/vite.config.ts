import react from '@vitejs/plugin-react-swc'
import path from 'path'
import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    VitePWA({
      registerType: 'autoUpdate',
      strategies: 'injectManifest',
      injectManifest: {
        injectionPoint: undefined,
      },
      devOptions: {
        enabled: true,
        type: 'module',
      },
      srcDir: 'src/service-worker',
      filename: 'sw.ts',
      manifest: {
        name: 'My App',
        short_name: 'My App',
        description: 'My App',
        icons: [
          {
            src: 'app_icon/vite.svg',
            type: 'image/png',
            sizes: '192x192',
          },
          {
            src: 'app_icon/vite.svg',
            sizes: '512x512',
            type: 'image/png',
          },
          {
            src: 'app_icon/vite.svg',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any maskable',
          },
        ],
        start_url: 'index.html',
        display: 'standalone',
        background_color: '#ffffff',
        theme_color: '#000000',
        lang: 'ja',
      },
    }),
  ],
  build: {
    outDir: '../dist',
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
})
