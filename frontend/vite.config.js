import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    // 既定では localhost にのみ bind するため、コンテナのポート公開越しに
    // ホストのブラウザから到達できない。
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      // 本番は Go の単一バイナリが SPA と API を同一オリジンで返す。
      // 開発時のみ Vite から backend コンテナへ中継し、同じ相対パスで
      // fetch できるようにする（compose の service 名で解決）。
      '/api': {
        target: 'http://backend:8080',
        changeOrigin: true,
      },
    },
  },
})
