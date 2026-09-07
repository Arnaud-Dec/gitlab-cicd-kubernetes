<script setup>
import { onMounted } from 'vue'
import { useTodos } from './composables/useTodos'
import TodoForm from './components/TodoForm.vue'
import TodoList from './components/TodoList.vue'
import StatusBar from './components/StatusBar.vue'
import NotificationBubble from './components/NotificationBubble.vue'

const { todos, loading, error, load, add, toggle, remove } = useTodos()

onMounted(load)
</script>

<template>
  <main class="app">
    <div class="card">
      <h1>Todo List</h1>
      <StatusBar />
      <TodoForm @add="add" />
      <p v-if="error" class="error">{{ error }}</p>
      <TodoList :todos="todos" :loading="loading" @toggle="toggle" @delete="remove" />
    </div>
  </main>
  <NotificationBubble />
</template>

<style scoped>
.app {
  min-height: 100vh;
  padding: 3rem 1rem;
  display: flex;
  justify-content: center;
}

.card {
  width: 100%;
  max-width: 34rem;
  height: fit-content;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 0.75rem;
  padding: 1.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

h1 {
  margin: 0 0 1.25rem;
  font-size: 1.5rem;
}

.error {
  margin: 0 0 1rem;
  padding: 0.6rem 0.75rem;
  border-radius: 0.5rem;
  background: var(--color-danger-bg);
  color: var(--color-danger);
  font-size: 0.85rem;
}
</style>
