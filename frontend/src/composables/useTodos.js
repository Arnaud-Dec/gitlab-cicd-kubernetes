import { ref } from 'vue'
import * as todosApi from '../api/todos'

export function useTodos() {
  const todos = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function load() {
    loading.value = true
    error.value = null
    try {
      todos.value = await todosApi.fetchTodos()
    } catch (err) {
      error.value = 'Failed to load todos.'
    } finally {
      loading.value = false
    }
  }

  async function add({ title, description }) {
    error.value = null
    try {
      const todo = await todosApi.createTodo({ title, description })
      todos.value.unshift(todo)
    } catch (err) {
      error.value = 'Failed to add todo.'
    }
  }

  async function toggle(id) {
    error.value = null
    try {
      const updated = await todosApi.toggleTodo(id)
      const index = todos.value.findIndex((t) => t.id === id)
      if (index !== -1) todos.value[index] = updated
    } catch (err) {
      error.value = 'Failed to update todo.'
    }
  }

  async function remove(id) {
    error.value = null
    try {
      await todosApi.deleteTodo(id)
      todos.value = todos.value.filter((t) => t.id !== id)
    } catch (err) {
      error.value = 'Failed to delete todo.'
    }
  }

  return { todos, loading, error, load, add, toggle, remove }
}
