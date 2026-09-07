import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TodoList from '../components/TodoList.vue'
import TodoItem from '../components/TodoItem.vue'

const todos = [
  { id: 1, title: 'Buy milk', completed: false },
  { id: 2, title: 'Walk dog', completed: true },
]

describe('TodoList', () => {
  it('renders one TodoItem per todo', () => {
    const wrapper = mount(TodoList, { props: { todos, loading: false } })

    expect(wrapper.findAllComponents(TodoItem)).toHaveLength(2)
  })

  it('shows an empty-state message when there are no todos', () => {
    const wrapper = mount(TodoList, { props: { todos: [], loading: false } })

    expect(wrapper.text()).toContain('No todos yet')
    expect(wrapper.findAllComponents(TodoItem)).toHaveLength(0)
  })

  it('shows a loading message while loading', () => {
    const wrapper = mount(TodoList, { props: { todos: [], loading: true } })

    expect(wrapper.text()).toContain('Loading')
  })
})
