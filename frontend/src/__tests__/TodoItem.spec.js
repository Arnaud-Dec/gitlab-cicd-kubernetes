import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TodoItem from '../components/TodoItem.vue'

const todo = { id: 1, title: 'Buy milk', description: '2%', completed: false }

describe('TodoItem', () => {
  it('renders the title and description', () => {
    const wrapper = mount(TodoItem, { props: { todo } })

    expect(wrapper.text()).toContain('Buy milk')
    expect(wrapper.text()).toContain('2%')
  })

  it('emits toggle with the todo id when the checkbox changes', async () => {
    const wrapper = mount(TodoItem, { props: { todo } })

    await wrapper.find('input[type="checkbox"]').setValue(true)

    expect(wrapper.emitted('toggle')).toEqual([[1]])
  })

  it('emits delete with the todo id when the delete button is clicked', async () => {
    const wrapper = mount(TodoItem, { props: { todo } })

    await wrapper.find('button').trigger('click')

    expect(wrapper.emitted('delete')).toEqual([[1]])
  })

  it('applies completed styling when the todo is completed', () => {
    const wrapper = mount(TodoItem, { props: { todo: { ...todo, completed: true } } })

    expect(wrapper.classes()).toContain('completed')
  })
})
