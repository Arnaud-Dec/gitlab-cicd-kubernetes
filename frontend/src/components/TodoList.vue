<script setup>
import TodoItem from './TodoItem.vue'

defineProps({
  todos: {
    type: Array,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['toggle', 'delete'])
</script>

<template>
  <p v-if="loading" class="empty-state">Loading todos...</p>
  <p v-else-if="todos.length === 0" class="empty-state">No todos yet. Add one above.</p>
  <ul v-else class="todo-list">
    <TodoItem
      v-for="todo in todos"
      :key="todo.id"
      :todo="todo"
      @toggle="emit('toggle', $event)"
      @delete="emit('delete', $event)"
    />
  </ul>
</template>

<style scoped>
.todo-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.empty-state {
  margin: 0;
  padding: 1.25rem 0;
  text-align: center;
  color: var(--color-text-muted);
  font-size: 0.9rem;
}
</style>
