<script setup>
import { ref } from 'vue'

const emit = defineEmits(['add'])

const title = ref('')
const description = ref('')

function submit() {
  const trimmedTitle = title.value.trim()
  if (!trimmedTitle) return

  emit('add', { title: trimmedTitle, description: description.value.trim() })
  title.value = ''
  description.value = ''
}
</script>

<template>
  <form class="todo-form" @submit.prevent="submit">
    <input
      v-model="title"
      type="text"
      placeholder="What needs to be done?"
      aria-label="Todo title"
    />
    <input
      v-model="description"
      type="text"
      placeholder="Description (optional)"
      aria-label="Todo description"
    />
    <button type="submit" :disabled="!title.trim()">Add</button>
  </form>
</template>

<style scoped>
.todo-form {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.todo-form input {
  flex: 1 1 8rem;
  min-width: 0;
  padding: 0.55rem 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  background: var(--color-surface);
  color: var(--color-text);
  font-size: 0.9rem;
}

.todo-form input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.15);
}

.todo-form button {
  padding: 0.55rem 1.1rem;
  border: none;
  border-radius: 0.5rem;
  background: var(--color-primary);
  color: white;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.15s;
}

.todo-form button:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.todo-form button:disabled {
  background: var(--color-border);
  color: var(--color-text-muted);
  cursor: not-allowed;
}
</style>
