import { notify } from '../notifications'

const HOSTNAME_HEADER = 'x-server-hostname'

function describeRequest(config) {
  const method = (config?.method || '').toUpperCase()
  const url = config?.url || ''
  const baseURL = config?.baseURL || ''
  return { method, url: url.startsWith('http') ? url : baseURL + url }
}

// Reads the backend's X-Server-Hostname response header, if present — see
// backend/internal/server/hostname.go. Handles both axios' AxiosHeaders
// (which exposes a .get method) and a plain object, just in case.
function extractHostname(headers) {
  if (!headers) return null
  if (typeof headers.get === 'function') {
    return headers.get(HOSTNAME_HEADER) ?? null
  }
  return headers[HOSTNAME_HEADER] ?? null
}

// Attaches a response interceptor that reports every request made through
// this axios instance — method, URI, HTTP status, and (if the backend sent
// one) its hostname — as a notification bubble (see notifications.js /
// components/NotificationBubble.vue).
export function attachRequestNotifications(instance) {
  instance.interceptors.response.use(
    (response) => {
      notify({
        ...describeRequest(response.config),
        status: response.status,
        ok: response.status < 400,
        hostname: extractHostname(response.headers),
      })
      return response
    },
    (error) => {
      notify({
        ...describeRequest(error.config),
        status: error.response?.status ?? null,
        ok: false,
        hostname: extractHostname(error.response?.headers),
      })
      return Promise.reject(error)
    },
  )
  return instance
}
