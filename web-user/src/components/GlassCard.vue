<template>
  <div class="glass-card" :class="[`glass-card--${variant}`, { 'glass-card--sm': sm }]">
    <slot />
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  variant?: 'default' | 'green' | 'purple' | 'orange' | 'blue' | 'yellow' | 'red'
  sm?: boolean
}>(), {
  variant: 'default',
  sm: false,
})
</script>

<style scoped>
.glass-card {
  border: 1px solid var(--card-border);
  border-radius: var(--card-radius);
  background: var(--card-surface);
  box-shadow: var(--card-shadow);
  backdrop-filter: blur(var(--card-blur));
  -webkit-backdrop-filter: blur(var(--card-blur));
}

.glass-card--green { background: var(--color-green-soft); border-color: var(--color-green-soft); }
.glass-card--purple { background: var(--color-purple-soft); border-color: var(--color-purple-soft); }
.glass-card--orange { background: var(--color-orange-soft); border-color: var(--color-orange-soft); }
.glass-card--blue { background: var(--color-blue-soft); border-color: var(--color-blue-soft); }
.glass-card--yellow { background: var(--color-yellow-soft); border-color: var(--color-yellow-soft); }
.glass-card--red { background: var(--color-red-soft); border-color: var(--color-red-soft); }
.glass-card--sm { padding: var(--sp-4); border-radius: var(--card-radius-inner); }
</style>
