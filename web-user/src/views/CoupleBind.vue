<template>
  <div class="page couple-bind">
    <header class="page-header">
      <button class="back-btn" @click="$router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"/></svg>
      </button>
      <h1 class="page-title">情侣绑定</h1>
      <div style="width:34px"></div>
    </header>

    <!-- Tab switch -->
    <div class="tab-bar">
      <button class="tab-item" :class="{ active: mode === 'share' }" @click="mode = 'share'">分享邀请码</button>
      <button class="tab-item" :class="{ active: mode === 'join' }" @click="mode = 'join'">输入邀请码</button>
    </div>

    <!-- Share mode -->
    <div v-if="mode === 'share'" class="bind-card">
      <div class="bind-icon">&#10084;</div>
      <p class="bind-desc">将邀请码发给你的另一半，TA 输入后即可绑定</p>
      <div class="code-box" v-if="inviteCode">
        <span class="code-text">{{ inviteCode }}</span>
        <button class="copy-btn" @click="copyCode">复制</button>
      </div>
      <button v-else class="gen-btn" @click="generateCode" :disabled="loading">
        {{ loading ? '生成中...' : '生成邀请码' }}
      </button>
      <p v-if="copied" class="copied-hint">已复制到剪贴板</p>
      <p v-if="waitingBind" class="waiting-hint">等待对方绑定中...</p>
    </div>

    <!-- Join mode -->
    <div v-else class="bind-card">
      <div class="bind-icon">&#128150;</div>
      <p class="bind-desc">输入对方分享的 6 位邀请码</p>
      <div class="input-row">
        <input
          v-model="inputCode"
          class="code-input"
          type="text"
          maxlength="6"
          placeholder="请输入邀请码"
          @keyup.enter="handleBind"
        />
      </div>
      <button class="gen-btn" @click="handleBind" :disabled="loading || inputCode.length < 6">
        {{ loading ? '绑定中...' : '绑定' }}
      </button>
      <p v-if="error" class="error-hint">{{ error }}</p>
    </div>

    <!-- Success Dialog -->
    <div v-if="showSuccess" class="dialog-overlay">
      <div class="dialog-box">
        <div class="success-icon">&#10084;</div>
        <h3 class="success-title">绑定成功！</h3>
        <p class="success-desc">{{ successMessage }}</p>
        <button class="success-btn" @click="goToCouple">进入情侣主页</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { getInviteCode, bindCouple, getCoupleInfo } from '@/api/couple'
import { useCoupleStore } from '@/stores/couple'

const router = useRouter()
const coupleStore = useCoupleStore()

const mode = ref<'share' | 'join'>('share')
const inviteCode = ref('')
const inputCode = ref('')
const loading = ref(false)
const copied = ref(false)
const error = ref('')
const waitingBind = ref(false)
const showSuccess = ref(false)
const successMessage = ref('')

let pollTimer: ReturnType<typeof setInterval> | null = null

async function generateCode() {
  loading.value = true
  try {
    const res: any = await getInviteCode()
    inviteCode.value = res.invite_code
    waitingBind.value = true
    startPolling()
  } catch {
    error.value = '生成邀请码失败'
  } finally {
    loading.value = false
  }
}

function copyCode() {
  navigator.clipboard.writeText(inviteCode.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}

// 轮询检测对方是否已绑定
function startPolling() {
  stopPolling()
  pollTimer = setInterval(async () => {
    try {
      const res: any = await getCoupleInfo()
      if (res && res.couple_id) {
        // 对方已绑定成功
        stopPolling()
        waitingBind.value = false
        await coupleStore.fetchCoupleInfo()
        successMessage.value = `你和 ${res.partner?.nickname || 'TA'} 已成功绑定情侣关系`
        showSuccess.value = true
      }
    } catch {
      // ignore
    }
  }, 3000)
}

function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

async function handleBind() {
  if (inputCode.value.length < 6) return
  loading.value = true
  error.value = ''
  try {
    const res: any = await bindCouple(inputCode.value.toUpperCase())
    await coupleStore.fetchCoupleInfo()
    successMessage.value = `你和 ${res.partner?.nickname || '对方'} 已成功绑定情侣关系`
    showSuccess.value = true
  } catch (e: any) {
    error.value = e.message || '绑定失败'
  } finally {
    loading.value = false
  }
}

function goToCouple() {
  showSuccess.value = false
  router.replace('/couple')
}

onMounted(async () => {
  await coupleStore.fetchCoupleInfo()
  if (coupleStore.isBound) {
    router.replace('/couple')
  }
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.couple-bind {
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
}

.back-btn svg { width: 18px; height: 18px; }

.page-title {
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
}

.tab-bar {
  display: flex;
  gap: var(--sp-2);
  padding: 0 var(--sp-4);
  margin-bottom: var(--sp-6);
}

.tab-item {
  flex: 1;
  padding: var(--sp-3) 0;
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-3);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--dur-fast) var(--ease),
              color var(--dur-fast) var(--ease),
              border-color var(--dur-fast) var(--ease),
              box-shadow var(--dur-fast) var(--ease);
}

.tab-item:hover:not(.active) {
  border-color: var(--color-border-med);
  box-shadow: var(--shadow-sm);
}

.tab-item.active {
  background: var(--color-text);
  border-color: var(--color-text);
  color: var(--color-text-inv);
  box-shadow: var(--shadow-md);
}

.bind-card {
  margin: 0 var(--sp-4);
  padding: var(--sp-8) var(--sp-6);
  background: var(--color-surface);
  border-radius: var(--r-lg);
  border: 1px solid var(--color-border);
  text-align: center;
}

.bind-icon {
  font-size: 48px;
  margin-bottom: var(--sp-4);
}

.bind-desc {
  color: var(--color-text-3);
  font-size: var(--text-sm);
  margin-bottom: var(--sp-6);
  line-height: 1.6;
}

.code-box {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--sp-3);
  padding: var(--sp-5);
  background: var(--color-surface-2);
  border-radius: var(--r-md);
  margin-bottom: var(--sp-4);
}

.code-text {
  font-size: var(--text-2xl);
  font-weight: 700;
  letter-spacing: 0.15em;
  color: var(--color-text);
  font-family: var(--font-display);
}

.copy-btn {
  padding: var(--sp-2) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text-2);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
}

.gen-btn {
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
              transform var(--dur-fast) var(--ease),
              box-shadow var(--dur-fast) var(--ease);
}

.gen-btn:hover:not(:disabled) {
  background: var(--color-text-2);
  box-shadow: var(--shadow-md);
}

.gen-btn:active:not(:disabled) {
  transform: translateY(1px);
  box-shadow: none;
}

.gen-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.input-row {
  margin-bottom: var(--sp-4);
}

.code-input {
  width: 100%;
  padding: var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface-2);
  color: var(--color-text);
  font-size: var(--text-xl);
  font-weight: 700;
  text-align: center;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  outline: none;
  transition: border-color var(--dur-fast) var(--ease);
}

.code-input:focus {
  border-color: var(--color-text);
}

.code-input::placeholder {
  font-size: var(--text-sm);
  letter-spacing: 0;
  font-weight: 500;
  color: var(--color-text-3);
}

.copied-hint {
  margin-top: var(--sp-3);
  color: var(--color-success);
  font-size: var(--text-xs);
  font-weight: 600;
}

.error-hint {
  margin-top: var(--sp-3);
  color: var(--color-error);
  font-size: var(--text-xs);
  font-weight: 600;
}

.waiting-hint {
  margin-top: var(--sp-4);
  color: var(--color-accent);
  font-size: var(--text-xs);
  font-weight: 600;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Dialog */
.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(28, 25, 21, 0.55);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.dialog-box {
  width: 80%;
  max-width: 300px;
  background: var(--color-surface);
  border-radius: var(--r-xl);
  padding: var(--sp-8) var(--sp-6);
  text-align: center;
  animation: scaleIn 0.3s var(--ease-out);
}

@keyframes scaleIn {
  from { transform: scale(0.9); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

.success-icon {
  font-size: 56px;
  margin-bottom: var(--sp-4);
  animation: heartbeat 1s ease-in-out;
}

@keyframes heartbeat {
  0% { transform: scale(0.8); }
  25% { transform: scale(1.1); }
  50% { transform: scale(0.95); }
  75% { transform: scale(1.05); }
  100% { transform: scale(1); }
}

.success-title {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: var(--sp-2);
}

.success-desc {
  color: var(--color-text-3);
  font-size: var(--text-sm);
  margin-bottom: var(--sp-6);
  line-height: 1.5;
}

.success-btn {
  width: 100%;
  padding: var(--sp-4);
  border: none;
  border-radius: var(--r-md);
  background: var(--color-text);
  color: var(--color-text-inv);
  font-size: var(--text-base);
  font-weight: 600;
  cursor: pointer;
  transition: opacity var(--dur-fast) var(--ease);
}

.success-btn:active {
  opacity: 0.8;
}
</style>
