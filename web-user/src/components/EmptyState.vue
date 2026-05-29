<template>
  <div class="empty-state">
    <div class="empty-state-icon">
      <svg v-if="type === 'no-data'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
        <polyline points="7 10 12 15 17 10"/>
        <line x1="12" y1="15" x2="12" y2="3"/>
      </svg>
      <svg v-else-if="type === 'no-result'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <circle cx="11" cy="11" r="8"/>
        <line x1="21" y1="21" x2="16.65" y2="16.65"/>
      </svg>
      <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
      </svg>
    </div>
    <p class="empty-state-text">{{ displayText }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  type?: 'no-data' | 'no-result' | 'no-favorite'
  text?: string
}>(), {
  type: 'no-data',
})

const defaultTexts: Record<string, string> = {
  'no-data': '暂无数据',
  'no-result': '没有找到相关内容',
  'no-favorite': '还没有收藏哦',
}

const displayText = computed(() => props.text || defaultTexts[props.type])
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--color-ink-3);
}
.empty-state-icon {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
  opacity: 0.4;
}
.empty-state-icon svg {
  width: 100%;
  height: 100%;
}
.empty-state-text {
  font-size: var(--text-sm);
  font-weight: 500;
}
</style>
