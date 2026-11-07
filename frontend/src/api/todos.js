import client from './client'

export async function fetchTodos() {
  const { data } = await client.get('/todos')
  return data
}

export async function createTodo({ title, description }) {
  const { data } = await client.post('/todos', { title, description })
  return data
}

export async function toggleTodo(id) {
  const { data } = await client.patch(`/todos/${id}/toggle`)
  return data
}

export async function deleteTodo(id) {
  await client.delete(`/todos/${id}`)
}
