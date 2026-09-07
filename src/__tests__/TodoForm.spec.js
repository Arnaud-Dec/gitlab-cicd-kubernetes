import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TodoForm from '../components/TodoForm.vue'

describe('TodoForm', () => {
  it('emits add with the entered title and description on submit', async () => {
    const wrapper = mount(TodoForm)
    const [titleInput, descriptionInput] = wrapper.findAll('input')

    await titleInput.setValue('Buy milk')
    await descriptionInput.setValue('2%')
    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.emitted('add')).toEqual([[{ title: 'Buy milk', description: '2%' }]])
  })

  it('clears the inputs after submit', async () => {
    const wrapper = mount(TodoForm)
    const [titleInput] = wrapper.findAll('input')

    await titleInput.setValue('Buy milk')
    await wrapper.find('form').trigger('submit.prevent')

    expect(titleInput.element.value).toBe('')
  })

  it('does not emit add when the title is empty', async () => {
    const wrapper = mount(TodoForm)

    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.emitted('add')).toBeUndefined()
  })

  it('disables the submit button when the title is empty', async () => {
    const wrapper = mount(TodoForm)

    expect(wrapper.find('button').attributes('disabled')).toBeDefined()

    const [titleInput] = wrapper.findAll('input')
    await titleInput.setValue('Buy milk')

    expect(wrapper.find('button').attributes('disabled')).toBeUndefined()
  })
})
