const DEFAULT_CONFIG = {
  apiUrl: 'http://localhost:8080/api',
}

let configPromise = null

// Fetches /config.json once at runtime and memoizes the result for the
// lifetime of the page, so the backend address can be changed per-deployment
// (e.g. by swapping the static file after the frontend is built) without a
// rebuild. Falls back to localhost defaults if the file is missing or
// unparsable.
export function loadConfig() {
  if (!configPromise) {
    configPromise = (async () => {
      try {
        const res = await fetch('/config.json')
        if (!res.ok) throw new Error(`config.json responded with ${res.status}`)
        return await res.json()
      } catch {
        return DEFAULT_CONFIG
      }
    })()
  }
  return configPromise
}
