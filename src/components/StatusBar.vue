<script setup>
import { ref, computed, onMounted } from 'vue'
import { useHealthStatus } from '../composables/useHealthStatus'
import { loadConfig } from '../config'
import { fetchVersion } from '../api/health'
import { version as frontendVersion } from '../version'

const { backendConnected, databaseConnected } = useHealthStatus()
const backendAddress = ref('')
const backendVersion = ref(null)

onMounted(async () => {
  const config = await loadConfig()
  backendAddress.value = config.apiUrl
  backendVersion.value = await fetchVersion()
})

const isSecure = computed(() => {
  if (!backendAddress.value) return null
  try {
    return new URL(backendAddress.value).protocol === 'https:'
  } catch {
    return null
  }
})

const securityLabel = computed(() => (isSecure.value ? 'Secure (HTTPS)' : 'Insecure (HTTP)'))
</script>

<template>
  <div class="status-bar">
    <span class="chip" :class="backendConnected ? 'ok' : 'down'">
      <span class="dot" />
      Backend: {{ backendConnected ? 'Connected' : 'Disconnected' }}
      <span v-if="backendAddress" class="address">{{ backendAddress }}</span>
      <span v-if="isSecure === false" class="warning-icon" aria-hidden="true">⚠</span>
      <span v-if="isSecure !== null" class="help-icon" :title="securityLabel">?</span>
    </span>
    <span class="chip" :class="databaseConnected ? 'ok' : 'down'">
      <span class="dot" />
      Database: {{ databaseConnected ? 'Connected' : 'Disconnected' }}
    </span>
  </div>
  <p class="versions">
    Frontend <span class="version">{{ frontendVersion }}</span> · Backend
    <span class="version">{{ backendVersion ?? '—' }}</span>
  </p>
</template>

<style scoped>
.status-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 0.4rem;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.25rem 0.65rem;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 600;
  white-space: nowrap;
}

.dot {
  width: 0.45rem;
  height: 0.45rem;
  border-radius: 50%;
  background: currentColor;
  flex-shrink: 0;
}

.chip.ok {
  color: var(--color-success);
  background: var(--color-success-bg);
}

.chip.down {
  color: var(--color-danger);
  background: var(--color-danger-bg);
}

.address {
  font-weight: 400;
  opacity: 0.75;
  font-family: ui-monospace, Consolas, monospace;
  font-size: 0.72rem;
}

.warning-icon {
  color: var(--color-warning);
}

.help-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1rem;
  height: 1rem;
  border-radius: 50%;
  border: 1px solid currentColor;
  font-size: 0.65rem;
  font-weight: 700;
  opacity: 0.7;
  cursor: help;
}

.versions {
  margin: 0 0 1.25rem;
  font-size: 0.72rem;
  color: var(--color-text-muted);
}

.version {
  font-family: ui-monospace, Consolas, monospace;
}
</style>
