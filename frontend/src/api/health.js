import axios from 'axios'
import { loadConfig } from '../config'

const healthClient = axios.create()

healthClient.interceptors.request.use(async (requestConfig) => {
  const { apiUrl } = await loadConfig()
  // /healthz and /readyz are served on the same origin as the API (see
  // backend/internal/server/router.go), so there's no separate health
  // address to configure — just the API's origin without the /api path.
  requestConfig.baseURL = new URL(apiUrl).origin
  return requestConfig
})

// Polls /readyz on the backend. A response (whether 200 or 503) means the
// backend is reachable; the status code tells us whether the backend can in
// turn reach the database. No response at all means the frontend can't reach
// the backend.
export async function checkStatus() {
  try {
    const res = await healthClient.get('/readyz')
    return { backendConnected: true, databaseConnected: res.status === 200 }
  } catch (err) {
    if (err.response) {
      return { backendConnected: true, databaseConnected: false }
    }
    return { backendConnected: false, databaseConnected: false }
  }
}

// Fetches the backend's build version from /version. Returns null if it
// can't be reached — the version this rarely changes, so it's fetched once
// rather than polled like checkStatus above.
export async function fetchVersion() {
  try {
    const res = await healthClient.get('/version')
    return res.data.version
  } catch {
    return null
  }
}
