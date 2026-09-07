<script setup>
defineProps({
  todo: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['toggle', 'delete'])
</script>

<template>
  <li class="todo-item" :class="{ completed: todo.completed }">
    <input
      type="checkbox"
      :checked="todo.completed"
      aria-label="Toggle completed"
      @change="emit('toggle', todo.id)"
    />
    <div class="todo-text">
      <span class="title">{{ todo.title }}</span>
      <span v-if="todo.description" class="description">{{ todo.description }}</span>
    </div>
    <button type="button" aria-label="Delete todo" @click="emit('delete', todo.id)">
      Delete
    </button>
  </li>
</template>

<style scoped>
.todo-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.65rem 0;
  border-bottom: 1px solid var(--color-border);
}

.todo-item:last-child {
  border-bottom: none;
}

.todo-item input[type='checkbox'] {
  margin-top: 0.2rem;
  width: 1.05rem;
  height: 1.05rem;
  accent-color: var(--color-primary);
  cursor: pointer;
  flex-shrink: 0;
}

.todo-text {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.1rem;
}

.title {
  font-size: 0.95rem;
}

.description {
  font-size: 0.82rem;
  color: var(--color-text-muted);
}

.completed {
  opacity: 0.6;
}

.completed .title {
  text-decoration: line-through;
}

.todo-item button {
  padding: 0.3rem 0.65rem;
  border: 1px solid var(--color-border);
  border-radius: 0.4rem;
  background: var(--color-surface);
  color: var(--color-text-muted);
  font-size: 0.78rem;
  cursor: pointer;
  flex-shrink: 0;
  transition:
    color 0.15s,
    border-color 0.15s;
}

.todo-item button:hover {
  color: var(--color-danger);
  border-color: var(--color-danger);
}
</style>
