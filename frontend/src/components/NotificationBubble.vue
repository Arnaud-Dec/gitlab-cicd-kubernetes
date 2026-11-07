<script setup>
import { notifications } from '../notifications'
</script>

<template>
  <div class="notification-stack">
    <div v-for="n in notifications" :key="n.id" class="notification" :class="n.ok ? 'ok' : 'error'">
      <span class="method">{{ n.method }}</span>
      <span class="url">{{ n.url }}</span>
      <span class="status">{{ n.status ?? 'ERR' }}</span>
      <span v-if="n.hostname" class="hostname">{{ n.hostname }}</span>
    </div>
  </div>
</template>

<style scoped>
.notification-stack {
  position: fixed;
  bottom: 1rem;
  right: 1rem;
  display: flex;
  flex-direction: column-reverse;
  gap: 0.4rem;
  z-index: 1000;
  max-width: 24rem;
  pointer-events: none;
}

.notification {
  display: flex;
  gap: 0.6rem;
  align-items: center;
  padding: 0.5rem 0.75rem;
  border-radius: 0.375rem;
  font-size: 0.8rem;
  font-family: ui-monospace, Consolas, monospace;
  color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  animation: notification-in 0.15s ease-out;
}

.notification.ok {
  background: var(--color-success);
}

.notification.error {
  background: var(--color-danger);
}

.method {
  font-weight: 700;
  opacity: 0.85;
}

.url {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status {
  font-weight: 700;
}

.hostname {
  opacity: 0.75;
  max-width: 9rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  border-left: 1px solid rgba(255, 255, 255, 0.35);
  padding-left: 0.6rem;
}

@keyframes notification-in {
  from {
    opacity: 0;
    transform: translateY(6px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
