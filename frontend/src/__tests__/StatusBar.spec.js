import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import StatusBar from '../components/StatusBar.vue'
import * as health from '../api/health'
import * as config from '../config'
import { version as frontendVersion } from '../version'

vi.mock('../api/health')
vi.mock('../config')

describe('StatusBar', () => {
  beforeEach(() => {
    config.loadConfig.mockResolvedValue({ apiUrl: 'http://localhost:8080/api' })
    health.fetchVersion.mockResolvedValue('2026-09-06T12:00')
  })

  it('shows both statuses as connected', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: true })

    const wrapper = mount(StatusBar)
    await flushPromises()

    expect(wrapper.text()).toContain('Backend: Connected')
    expect(wrapper.text()).toContain('Database: Connected')
  })

  it('shows the database as disconnected when the backend is up but the DB is not', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: false })

    const wrapper = mount(StatusBar)
    await flushPromises()

    expect(wrapper.text()).toContain('Backend: Connected')
    expect(wrapper.text()).toContain('Database: Disconnected')
  })

  it('shows both statuses as disconnected when the backend is unreachable', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: false, databaseConnected: false })

    const wrapper = mount(StatusBar)
    await flushPromises()

    expect(wrapper.text()).toContain('Backend: Disconnected')
    expect(wrapper.text()).toContain('Database: Disconnected')
  })

  it('shows the configured backend address inside the Backend chip', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: true })
    config.loadConfig.mockResolvedValue({ apiUrl: 'http://10.0.0.5:8080/api' })

    const wrapper = mount(StatusBar)
    await flushPromises()

    const backendChip = wrapper.find('.chip.ok')
    expect(backendChip.text()).toContain('http://10.0.0.5:8080/api')
  })

  it('shows a warning icon and an "Insecure" tooltip when the backend address is plain HTTP', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: true })
    config.loadConfig.mockResolvedValue({ apiUrl: 'http://10.0.0.5:8080/api' })

    const wrapper = mount(StatusBar)
    await flushPromises()

    expect(wrapper.find('.warning-icon').exists()).toBe(true)
    expect(wrapper.find('.help-icon').attributes('title')).toBe('Insecure (HTTP)')
  })

  it('hides the warning icon and shows a "Secure" tooltip when the backend address is HTTPS', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: true })
    config.loadConfig.mockResolvedValue({ apiUrl: 'https://api.example.com/api' })

    const wrapper = mount(StatusBar)
    await flushPromises()

    expect(wrapper.find('.warning-icon').exists()).toBe(false)
    expect(wrapper.find('.help-icon').attributes('title')).toBe('Secure (HTTPS)')
  })

  it('shows the frontend build version and the fetched backend version', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: true, databaseConnected: true })
    health.fetchVersion.mockResolvedValue('2026-09-06T12:00')

    const wrapper = mount(StatusBar)
    await flushPromises()

    expect(wrapper.text()).toContain(`Frontend ${frontendVersion}`)
    expect(wrapper.text()).toContain('Backend 2026-09-06T12:00')
  })

  it('shows a placeholder for the backend version when it cannot be fetched', async () => {
    health.checkStatus.mockResolvedValue({ backendConnected: false, databaseConnected: false })
    health.fetchVersion.mockResolvedValue(null)

    const wrapper = mount(StatusBar)
    await flushPromises()

    expect(wrapper.text()).toContain(`Frontend ${frontendVersion}`)
    expect(wrapper.text()).toContain('Backend —')
  })
})
