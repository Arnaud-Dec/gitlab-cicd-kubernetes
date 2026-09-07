import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { notifications, notify, clearNotifications } from '../notifications'

describe('notify', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    clearNotifications()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('adds a notification with the given fields', () => {
    notify({ method: 'GET', url: '/api/todos', status: 200, ok: true })

    expect(notifications).toHaveLength(1)
    expect(notifications[0]).toMatchObject({ method: 'GET', url: '/api/todos', status: 200, ok: true })
  })

  it('carries the hostname when given one', () => {
    notify({ method: 'GET', url: '/api/todos', status: 200, ok: true, hostname: 'todo-backend-7d9f8c6b8-xk2p9' })

    expect(notifications[0].hostname).toBe('todo-backend-7d9f8c6b8-xk2p9')
  })

  it('leaves hostname undefined when none is given', () => {
    notify({ method: 'GET', url: '/api/todos', status: 200, ok: true })

    expect(notifications[0].hostname).toBeUndefined()
  })

  it('assigns each notification a unique, increasing id', () => {
    notify({ method: 'GET', url: '/a', status: 200, ok: true })
    notify({ method: 'GET', url: '/b', status: 200, ok: true })

    expect(notifications[1].id).toBeGreaterThan(notifications[0].id)
  })

  it('auto-removes the notification after the display duration', () => {
    notify({ method: 'DELETE', url: '/api/todos/1', status: 204, ok: true })
    expect(notifications).toHaveLength(1)

    vi.advanceTimersByTime(3000)

    expect(notifications).toHaveLength(0)
  })

  it('removes notifications independently of each other', () => {
    notify({ method: 'GET', url: '/a', status: 200, ok: true })
    vi.advanceTimersByTime(1500)
    notify({ method: 'GET', url: '/b', status: 200, ok: true })

    vi.advanceTimersByTime(1500)
    expect(notifications).toHaveLength(1)
    expect(notifications[0].url).toBe('/b')

    vi.advanceTimersByTime(1500)
    expect(notifications).toHaveLength(0)
  })
})
