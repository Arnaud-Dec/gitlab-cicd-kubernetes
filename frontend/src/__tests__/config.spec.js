import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

const DEFAULT_CONFIG = {
  apiUrl: 'http://localhost:8080/api',
}

describe('loadConfig', () => {
  beforeEach(() => {
    vi.resetModules()
    vi.stubGlobal('fetch', vi.fn())
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('fetches and returns /config.json', async () => {
    fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ apiUrl: 'http://10.0.0.5:8080/api' }),
    })

    const { loadConfig } = await import('../config')
    const config = await loadConfig()

    expect(fetch).toHaveBeenCalledWith('/config.json')
    expect(config).toEqual({ apiUrl: 'http://10.0.0.5:8080/api' })
  })

  it('memoizes the result across calls, fetching only once', async () => {
    fetch.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(DEFAULT_CONFIG),
    })

    const { loadConfig } = await import('../config')
    await loadConfig()
    await loadConfig()

    expect(fetch).toHaveBeenCalledTimes(1)
  })

  it('falls back to localhost defaults when the fetch fails', async () => {
    fetch.mockRejectedValue(new Error('network error'))

    const { loadConfig } = await import('../config')
    const config = await loadConfig()

    expect(config).toEqual(DEFAULT_CONFIG)
  })

  it('falls back to localhost defaults when the response is not ok', async () => {
    fetch.mockResolvedValue({ ok: false, status: 404 })

    const { loadConfig } = await import('../config')
    const config = await loadConfig()

    expect(config).toEqual(DEFAULT_CONFIG)
  })
})
