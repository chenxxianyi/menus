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
.glass-card--green { background: rgba(212,237,218,0.4); border-color: rgba(212,237,218,0.5); }
.glass-card--purple { background: rgba(232,223,245,0.4); border-color: rgba(232,223,245,0.5); }
.glass-card--orange { background: rgba(255,224,168,0.4); border-color: rgba(255,224,168,0.5); }
.glass-card--blue { background: rgba(168,216,234,0.4); border-color: rgba(168,216,234,0.5); }
.glass-card--yellow { background: rgba(255,243,205,0.4); border-color: rgba(255,243,205,0.5); }
.glass-card--red { background: rgba(252,165,165,0.4); border-color: rgba(252,165,165,0.5); }
.glass-card--sm { padding: 14px; border-radius: 16px; }
</style>
