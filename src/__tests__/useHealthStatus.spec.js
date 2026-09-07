import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { defineComponent, h } from 'vue'
import { useHealthStatus } from '../composables/useHealthStatus'
import * as health from '../api/health'

vi.mock('../api/health')

function mountHost() {
  const state = {}
  const Host = defineComponent({
    setup() {
      Object.assign(state, useHealthStatus())
      return () => h('div')
    },
  })
  const wrapper = mount(Host)
  return { wrapper, state }
}

describe('useHealthStatus', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.resetAllMocks()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('polls on mount and reflects a fully connected status', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: true })

    const { wrapper, state } = mountHost()
    await flushPromises()

    expect(state.backendConnected.value).toBe(true)
    expect(state.databaseConnected.value).toBe(true)
    expect(health.checkStatus).toHaveBeenCalledTimes(1)

    wrapper.unmount()
  })

  it('reflects a disconnected database when the backend reports not ready', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: false })

    const { wrapper, state } = mountHost()
    await flushPromises()

    expect(state.backendConnected.value).toBe(true)
    expect(state.databaseConnected.value).toBe(false)

    wrapper.unmount()
  })

  it('polls again every 2 seconds', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: false })

    const { wrapper } = mountHost()
    await flushPromises()
    expect(health.checkStatus).toHaveBeenCalledTimes(1)

    await vi.advanceTimersByTimeAsync(2000)
    expect(health.checkStatus).toHaveBeenCalledTimes(2)

    await vi.advanceTimersByTimeAsync(2000)
    expect(health.checkStatus).toHaveBeenCalledTimes(3)

    wrapper.unmount()
  })

  it('stops polling once unmounted', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: true })

    const { wrapper } = mountHost()
    await flushPromises()
    expect(health.checkStatus).toHaveBeenCalledTimes(1)

    wrapper.unmount()
    await vi.advanceTimersByTimeAsync(4000)

    expect(health.checkStatus).toHaveBeenCalledTimes(1)
  })
})
