import { ref, onMounted, onUnmounted } from 'vue'
import { checkStatus } from '../api/health'

const POLL_INTERVAL_MS = 2000

export function useHealthStatus() {
  const backendConnected = ref(false)
  const databaseConnected = ref(false)

  let intervalId = null

  async function poll() {
    const status = await checkStatus()
    backendConnected.value = status.backendConnected
    databaseConnected.value = status.databaseConnected
  }

  onMounted(() => {
    poll()
    intervalId = setInterval(poll, POLL_INTERVAL_MS)
  })

  onUnmounted(() => {
    if (intervalId) clearInterval(intervalId)
  })

  return { backendConnected, databaseConnected }
}
