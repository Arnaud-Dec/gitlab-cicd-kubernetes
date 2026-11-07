import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

// Build-time version stamp: current date and time, to the minute, computed
// once when this config is evaluated — i.e. at `vite build` time for a
// production build, or `vite dev`/vitest startup time for local dev. Mirrors
// the backend's -ldflags-injected version.Version (see backend/internal/version).
function buildVersion() {
  const now = new Date()
  const pad = (n) => String(n).padStart(2, '0')
  return `${now.getUTCFullYear()}-${pad(now.getUTCMonth() + 1)}-${pad(now.getUTCDate())}T${pad(now.getUTCHours())}:${pad(now.getUTCMinutes())}`
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  define: {
    __APP_VERSION__: JSON.stringify(buildVersion()),
  },
  test: {
    environment: 'jsdom',
  },
})
