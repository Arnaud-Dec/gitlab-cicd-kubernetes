import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useTodos } from '../composables/useTodos'
import * as todosApi from '../api/todos'

vi.mock('../api/todos')

describe('useTodos', () => {
  beforeEach(() => {
    vi.resetAllMocks()
  })

  it('load() populates todos from the API', async () => {
    todosApi.fetchTodos.mockResolvedValue([{ id: 1, title: 'A', completed: false }])

    const { todos, loading, load } = useTodos()
    const promise = load()
    expect(loading.value).toBe(true)
    await promise

    expect(todos.value).toHaveLength(1)
    expect(loading.value).toBe(false)
  })

  it('load() sets error on failure', async () => {
    todosApi.fetchTodos.mockRejectedValue(new Error('network error'))

    const { error, load } = useTodos()
    await load()

    expect(error.value).toBe('Failed to load todos.')
  })

  it('add() appends the created todo', async () => {
    todosApi.createTodo.mockResolvedValue({ id: 2, title: 'New', completed: false })

    const { todos, add } = useTodos()
    await add({ title: 'New', description: '' })

    expect(todos.value).toEqual([{ id: 2, title: 'New', completed: false }])
  })

  it('toggle() replaces the matching todo', async () => {
    todosApi.fetchTodos.mockResolvedValue([{ id: 1, title: 'A', completed: false }])
    todosApi.toggleTodo.mockResolvedValue({ id: 1, title: 'A', completed: true })

    const { todos, load, toggle } = useTodos()
    await load()
    await toggle(1)

    expect(todos.value[0].completed).toBe(true)
  })

  it('remove() deletes the matching todo', async () => {
    todosApi.fetchTodos.mockResolvedValue([{ id: 1, title: 'A', completed: false }])
    todosApi.deleteTodo.mockResolvedValue(undefined)

    const { todos, load, remove } = useTodos()
    await load()
    await remove(1)

    expect(todos.value).toHaveLength(0)
  })
})
