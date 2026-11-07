import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import NotificationBubble from '../components/NotificationBubble.vue'
import { notify, clearNotifications } from '../notifications'

describe('NotificationBubble', () => {
  beforeEach(() => {
    clearNotifications()
  })

  it('renders nothing when there are no notifications', () => {
    const wrapper = mount(NotificationBubble)
    expect(wrapper.findAll('.notification')).toHaveLength(0)
  })

  it('renders the method, URI, and status code of a successful request', async () => {
    const wrapper = mount(NotificationBubble)
    notify({ method: 'POST', url: '/api/todos', status: 201, ok: true })
    await wrapper.vm.$nextTick()

    const item = wrapper.find('.notification')
    expect(item.classes()).toContain('ok')
    expect(item.text()).toContain('POST')
    expect(item.text()).toContain('/api/todos')
    expect(item.text()).toContain('201')
  })

  it('styles a failed request differently', async () => {
    const wrapper = mount(NotificationBubble)
    notify({ method: 'GET', url: '/api/todos/999', status: 404, ok: false })
    await wrapper.vm.$nextTick()

    const item = wrapper.find('.notification')
    expect(item.classes()).toContain('error')
    expect(item.text()).toContain('404')
  })

  it('shows the hostname when the backend sent one', async () => {
    const wrapper = mount(NotificationBubble)
    notify({ method: 'GET', url: '/api/todos', status: 200, ok: true, hostname: 'todo-backend-abc123' })
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.hostname').exists()).toBe(true)
    expect(wrapper.find('.hostname').text()).toBe('todo-backend-abc123')
  })

  it('omits the hostname element when none was sent', async () => {
    const wrapper = mount(NotificationBubble)
    notify({ method: 'GET', url: '/api/todos', status: 200, ok: true })
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.hostname').exists()).toBe(false)
  })

  it('shows ERR for requests with no status (network failure)', async () => {
    const wrapper = mount(NotificationBubble)
    notify({ method: 'GET', url: '/api/todos', status: null, ok: false })
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.notification').text()).toContain('ERR')
  })

  it('renders multiple notifications', async () => {
    const wrapper = mount(NotificationBubble)
    notify({ method: 'GET', url: '/api/todos', status: 200, ok: true })
    notify({ method: 'POST', url: '/api/todos', status: 201, ok: true })
    await wrapper.vm.$nextTick()

    expect(wrapper.findAll('.notification')).toHaveLength(2)
  })
})
