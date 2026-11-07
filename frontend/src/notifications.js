import { reactive } from 'vue'

const DISPLAY_DURATION_MS = 3000

let nextId = 0

export const notifications = reactive([])

// Pushes a request-result bubble and auto-removes it after a few seconds.
export function notify({ method, url, status, ok, hostname }) {
  const id = ++nextId
  notifications.push({ id, method, url, status, ok, hostname })

  setTimeout(() => {
    const index = notifications.findIndex((n) => n.id === id)
    if (index !== -1) notifications.splice(index, 1)
  }, DISPLAY_DURATION_MS)
}

// Test-only: clears all pending notifications so specs don't leak state into
// each other via this module-level singleton.
export function clearNotifications() {
  notifications.splice(0, notifications.length)
}
