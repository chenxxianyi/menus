<template>
  <div class="page feedback-page">
    <header class="page-header">
      <button class="back-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <h1 class="page-title">帮助与反馈</h1>
      <div style="width:34px"></div>
    </header>

    <div class="feedback-form">
      <div class="form-group">
        <label class="form-label">反馈内容</label>
        <textarea
          v-model="content"
          class="form-textarea"
          placeholder="请描述你遇到的问题或建议..."
          rows="5"
          maxlength="500"
        ></textarea>
        <div class="char-count">{{ content.length }}/500</div>
      </div>

      <div class="form-group">
        <label class="form-label">联系方式（可选）</label>
        <input
          v-model="contact"
          class="form-input"
          placeholder="邮箱或手机号，方便我们联系你"
        />
      </div>

      <button class="submit-btn" @click="handleSubmit" :disabled="!content.trim() || submitting">
        {{ submitting ? '提交中...' : '提交反馈' }}
      </button>
    </div>

    <div v-if="showToast" class="toast">感谢你的反馈！</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index'

const router = useRouter()
const content = ref('')
const contact = ref('')
const submitting = ref(false)
const showToast = ref(false)

async function handleSubmit() {
  if (!content.value.trim()) return
  submitting.value = true
  try {
    await api.post('/feedback', {
      content: content.value.trim(),
      contact: contact.value.trim() || undefined,
    })
    showToast.value = true
    setTimeout(() => {
      showToast.value = false
      router.back()
    }, 1500)
  } catch {
    // ignore
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.feedback-page {
  min-height: 100vh;
  background: var(--color-bg);
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--sp-3) var(--sp-4);
}

.back-btn {
  width: 34px; height: 34px;
  display: flex; align-items: center; justify-content: center;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-2);
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease);
}

.back-btn:hover { background: var(--color-surface-2); }
.back-btn:active { background: var(--color-surface-3); }
.back-btn svg { width: 18px; height: 18px; }

.page-title {
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
}

.feedback-form {
  padding: 0 var(--sp-4);
}

.form-group {
  margin-bottom: var(--sp-5);
}

.form-label {
  display: block;
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--sp-2);
}

.form-textarea {
  width: 100%;
  padding: var(--sp-3) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: var(--text-sm);
  font-family: inherit;
  outline: none;
  resize: vertical;
  transition: border-color var(--dur-fast) var(--ease);
}

.form-textarea:focus {
  border-color: var(--color-text);
}

.form-textarea::placeholder {
  color: var(--color-text-3);
}

.char-count {
  text-align: right;
  font-size: var(--text-2xs);
  color: var(--color-text-3);
  margin-top: var(--sp-1);
}

.form-input {
  width: 100%;
  padding: var(--sp-3) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
  color: var(--color-text);
  font-size: var(--text-sm);
  font-family: inherit;
  outline: none;
  transition: border-color var(--dur-fast) var(--ease);
}

.form-input:focus {
  border-color: var(--color-text);
}

.form-input::placeholder {
  color: var(--color-text-3);
}

.submit-btn {
  width: 100%;
  padding: var(--sp-4);
  border: none;
  border-radius: var(--r-md);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-base);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease),
              transform var(--dur-fast) var(--ease);
}

.submit-btn:hover:not(:disabled) {
  background: var(--color-text-2);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(1px);
}

.submit-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.toast {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  padding: var(--sp-4) var(--sp-6);
  background: var(--color-text);
  color: var(--color-text-inv);
  border-radius: var(--r-md);
  font-size: var(--text-sm);
  font-weight: 600;
  z-index: 100;
  animation: fadeInOut 1.5s ease;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translate(-50%, -50%) scale(0.9); }
  15% { opacity: 1; transform: translate(-50%, -50%) scale(1); }
  85% { opacity: 1; }
  100% { opacity: 0; }
}

@media (min-width: 768px) {
  .feedback-page { max-width: 640px; margin: 0 auto; }
}
</style>
