<template>
  <main class="couple-bind-page" :style="pageVars">
    <header class="page-header" aria-label="页面顶部">
      <button class="back-btn" type="button" aria-label="返回" @click="router.back()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <path d="m15 18-6-6 6-6" />
        </svg>
      </button>
      <h1 class="page-title">情侣绑定</h1>
    </header>

    <nav class="segment" aria-label="绑定方式切换">
      <button class="tab-btn" :class="{ active: mode === 'share' }" type="button" @click="switchMode('share')">分享邀请码</button>
      <button class="tab-btn" :class="{ active: mode === 'join' }" type="button" @click="switchMode('join')">输入邀请码</button>
    </nav>

    <section class="bind-card" aria-label="情侣绑定内容">
      <div v-show="mode === 'share'" class="panel">
        <div class="heart-art" aria-hidden="true">
          <i class="sparkle s1"></i>
          <i class="sparkle s2"></i>
          <i class="sparkle s3"></i>
          <i class="mini-heart one"></i>
          <i class="mini-heart two"></i>
          <i class="mini-heart three"></i>
          <i class="heart-core"></i>
          <i class="heart-ring"></i>
        </div>

        <p class="desc">将邀请码发给你的另一半，TA 输入后即可绑定</p>

        <div class="code-box">
          <span class="code-label">你的邀请码</span>
          <strong class="code-value">{{ inviteCode || 'A8K9P2' }}</strong>
          <button class="copy-btn" type="button" :disabled="!inviteCode" @click="copyCode">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <rect x="9" y="9" width="11" height="11" rx="2" />
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
            </svg>
            <span>{{ copied ? '已复制' : '复制' }}</span>
          </button>
        </div>

        <button class="primary-btn" type="button" :disabled="loading" @click="generateCode">
          {{ loading ? '生成中...' : inviteCode ? '重新生成邀请码' : '生成邀请码' }}
        </button>

        <button class="wechat-btn" type="button" :disabled="!inviteCode" @click="shareToWechat">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M21 11.5a8.4 8.4 0 0 1-8.5 8.3 8.9 8.9 0 0 1-3.5-.7L3 21l1.8-4.5A8 8 0 0 1 3 11.5a8.4 8.4 0 0 1 8.5-8.3 8.4 8.4 0 0 1 8.5 8.3Z" />
            <path d="M8.7 10.2h.01M14.9 10.2h.01" />
          </svg>
          <span>分享给微信好友</span>
        </button>

        <p v-if="waitingBind" class="status-text">等待对方绑定中...</p>

        <p class="hint">
          <ShieldCheckIcon />
          <span>绑定后可一起规划菜单、共享购物清单</span>
        </p>
      </div>

      <div v-show="mode === 'join'" class="panel">
        <div class="heart-art" aria-hidden="true">
          <i class="sparkle s1"></i>
          <i class="sparkle s2"></i>
          <i class="sparkle s3"></i>
          <i class="mini-heart one"></i>
          <i class="mini-heart two"></i>
          <i class="mini-heart three"></i>
          <i class="heart-core"></i>
          <i class="heart-ring"></i>
        </div>

        <p class="desc">输入另一半发来的邀请码，确认后即可绑定</p>

        <input
          v-model="inputCode"
          class="invite-input"
          type="text"
          inputmode="text"
          maxlength="6"
          placeholder="请输入 6 位邀请码"
          @input="normalizeInput"
          @keyup.enter="handleBind"
        />
        <p class="error-text" :class="{ show: !!error }">{{ error || '请输入完整邀请码' }}</p>

        <button class="primary-btn" type="button" :disabled="loading" @click="handleBind">
          {{ loading ? '绑定中...' : '确认绑定' }}
        </button>

        <button class="wechat-btn" type="button" @click="askInviteCode">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M21 11.5a8.4 8.4 0 0 1-8.5 8.3 8.9 8.9 0 0 1-3.5-.7L3 21l1.8-4.5A8 8 0 0 1 3 11.5a8.4 8.4 0 0 1 8.5-8.3 8.4 8.4 0 0 1 8.5 8.3Z" />
            <path d="M8.7 10.2h.01M14.9 10.2h.01" />
          </svg>
          <span>向 TA 索要邀请码</span>
        </button>

        <p class="hint">
          <ShieldCheckIcon />
          <span>绑定后可一起规划菜单、共享购物清单</span>
        </p>
      </div>
    </section>

    <div class="toast" :class="{ show: !!toastText }">{{ toastText }}</div>

    <div v-if="showSuccess" class="dialog-overlay">
      <div class="dialog-box">
        <div class="heart-art compact" aria-hidden="true">
          <i class="sparkle s1"></i>
          <i class="sparkle s2"></i>
          <i class="sparkle s3"></i>
          <i class="mini-heart one"></i>
          <i class="mini-heart two"></i>
          <i class="mini-heart three"></i>
          <i class="heart-core"></i>
          <i class="heart-ring"></i>
        </div>
        <h3 class="success-title">绑定成功！</h3>
        <p class="success-desc">{{ successMessage }}</p>
        <button class="success-btn" type="button" @click="goToCouple">进入情侣主页</button>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getInviteCode, bindCouple, getCoupleInfo } from '@/api/couple'
import { useCoupleStore } from '@/stores/couple'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'
import { copyText } from '@/utils/clipboard'

type BindMode = 'share' | 'join'

const ShieldCheckIcon = defineComponent({
  name: 'ShieldCheckIcon',
  setup() {
    return () => h('svg', {
      viewBox: '0 0 24 24',
      fill: 'none',
      stroke: 'currentColor',
      'stroke-linecap': 'round',
      'stroke-linejoin': 'round',
      'aria-hidden': 'true',
    }, [
      h('path', { d: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z' }),
      h('path', { d: 'm9 12 2 2 4-5' }),
    ])
  },
})

const router = useRouter()
const coupleStore = useCoupleStore()

const mode = ref<BindMode>('share')
const inviteCode = ref('')
const inputCode = ref('')
const loading = ref(false)
const copied = ref(false)
const error = ref('')
const waitingBind = ref(false)
const showSuccess = ref(false)
const successMessage = ref('')
const toastText = ref('')

let pollTimer: ReturnType<typeof setInterval> | null = null
let toastTimer: ReturnType<typeof setTimeout> | null = null
let copiedTimer: ReturnType<typeof setTimeout> | null = null

const pageVars = computed(() => ({
  '--kitchen-bg': 'url(' + kitchenBg + ')',
}))

function switchMode(nextMode: BindMode) {
  mode.value = nextMode
  error.value = ''
}

function showToast(message: string) {
  toastText.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastText.value = ''
  }, 1500)
}

async function generateCode() {
  loading.value = true
  error.value = ''
  try {
    const res: any = await getInviteCode()
    inviteCode.value = res.invite_code
    waitingBind.value = true
    startPolling()
    showToast(inviteCode.value ? '已生成邀请码' : '生成邀请码失败')
  } catch (e: any) {
    error.value = e.message || '生成邀请码失败'
    showToast(error.value)
  } finally {
    loading.value = false
  }
}

async function copyCode() {
  if (!inviteCode.value) {
    showToast('请先生成邀请码')
    return
  }
  try {
    if (await copyText(inviteCode.value)) {
      showToast('已复制邀请码')
    } else {
      showToast('邀请码：' + inviteCode.value)
    }
  } catch {
    showToast('邀请码：' + inviteCode.value)
  }
  copied.value = true
  if (copiedTimer) clearTimeout(copiedTimer)
  copiedTimer = setTimeout(() => {
    copied.value = false
  }, 1500)
}

async function shareToWechat() {
  if (!inviteCode.value) {
    showToast('请先生成邀请码')
    return
  }
  const shareText = `我的情侣绑定邀请码：${inviteCode.value}`
  try {
    if (navigator.share) {
      await navigator.share({
        title: '情侣绑定邀请码',
        text: shareText,
      })
      showToast('分享成功')
      return
    }
    if (await copyText(shareText)) {
      showToast('已复制邀请码，可发送给对方')
      return
    }
    showToast('请复制下方邀请码后发送给对方：' + inviteCode.value)
  } catch {
    showToast('分享已取消或未完成')
  }
}

function askInviteCode() {
  showToast('请让对方在“情侣绑定”页面生成邀请码后发送给你。')
}

function normalizeInput() {
  inputCode.value = inputCode.value
    .toUpperCase()
    .replace(/[^A-Z0-9]/g, '')
    .slice(0, 6)
  error.value = ''
}

function startPolling() {
  stopPolling()
  pollTimer = setInterval(async () => {
    try {
      const res: any = await getCoupleInfo()
      if (res && res.couple_id) {
        stopPolling()
        waitingBind.value = false
        await coupleStore.fetchCoupleInfo()
        successMessage.value = `你和 ${res.partner?.nickname || 'TA'} 已成功绑定情侣关系`
        showSuccess.value = true
      }
    } catch {
      // Keep polling while the invite is waiting.
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
  normalizeInput()
  if (inputCode.value.length < 6) {
    error.value = '请输入完整邀请码'
    showToast(error.value)
    return
  }

  loading.value = true
  error.value = ''
  try {
    const res: any = await bindCouple(inputCode.value)
    await coupleStore.fetchCoupleInfo()
    successMessage.value = `你和 ${res.partner?.nickname || '对方'} 已成功绑定情侣关系`
    showSuccess.value = true
  } catch (e: any) {
    error.value = e.message || '绑定失败'
    showToast(error.value)
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
  if (toastTimer) clearTimeout(toastTimer)
  if (copiedTimer) clearTimeout(copiedTimer)
})
</script>

<style scoped>
.couple-bind-page {
  --text: #2e241f;
  --sub: #7a6a5f;
  --muted: #a19489;
  --cream: rgba(255, 250, 240, 0.78);
  --coral: #e95645;
  --coral-2: #ef5548;
  --coral-deep: #df4437;
  --sage: #8fa783;
  --border: rgba(255, 255, 255, 0.62);
  position: relative;
  width: min(100%, 430px);
  min-height: calc(100vh + var(--tab-h, 64px) + var(--safe-bottom, 34px));
  min-height: calc(100dvh + var(--tab-h, 64px) + var(--safe-bottom, 34px));
  margin: 0 auto;
  padding: max(52px, env(safe-area-inset-top)) 24px calc(34px + env(safe-area-inset-bottom));
  overflow-x: clip;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 237, 205, 0.34), rgba(255, 247, 233, 0.2)),
    var(--kitchen-bg) center top / cover fixed;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.couple-bind-page::before {
  content: "";
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 22% 8%, rgba(255, 255, 255, 0.68), transparent 30%),
    radial-gradient(circle at 88% 14%, rgba(236, 143, 71, 0.24), transparent 32%),
    radial-gradient(circle at 12% 92%, rgba(232, 154, 69, 0.2), transparent 31%),
    linear-gradient(90deg, rgba(255, 238, 212, 0.56), rgba(255, 244, 228, 0.18) 52%, rgba(203, 118, 55, 0.18));
  backdrop-filter: blur(4px) saturate(1.15);
  -webkit-backdrop-filter: blur(4px) saturate(1.15);
}

.page-header,
.segment,
.bind-card {
  position: relative;
  z-index: 1;
}

button,
input {
  font: inherit;
}

button {
  border: 0;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

button:disabled {
  cursor: not-allowed;
}

svg {
  display: block;
}

.page-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-btn {
  position: absolute;
  left: 0;
  top: 4px;
  width: 52px;
  height: 52px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.65);
  border-radius: 16px;
  color: #4a352a;
  background: rgba(255, 250, 240, 0.86);
  box-shadow:
    0 12px 28px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  transition: transform 180ms ease, box-shadow 180ms ease;
}

.back-btn svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.55;
}

.page-title {
  margin: 0;
  color: #2e241f;
  font-size: 27px;
  font-weight: 900;
  line-height: 1;
  letter-spacing: 0;
}

.segment {
  height: 66px;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  margin-top: 36px;
  padding: 4px;
  border: 1px solid rgba(255, 255, 255, 0.58);
  border-radius: 999px;
  background: rgba(255, 250, 240, 0.68);
  box-shadow:
    0 14px 30px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.76);
  backdrop-filter: blur(18px) saturate(1.08);
  -webkit-backdrop-filter: blur(18px) saturate(1.08);
}

.tab-btn {
  position: relative;
  z-index: 1;
  border-radius: 999px;
  color: #6f5b4d;
  background: transparent;
  font-size: 19px;
  font-weight: 800;
  letter-spacing: 0;
  transition: color 180ms ease, transform 180ms ease;
}

.tab-btn.active {
  color: var(--coral);
}

.tab-btn.active::before {
  content: "";
  position: absolute;
  inset: 0;
  z-index: -1;
  border-radius: inherit;
  background: rgba(255, 250, 240, 0.88);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.82);
}

.tab-btn.active::after {
  content: "";
  position: absolute;
  left: 50%;
  bottom: 4px;
  width: 42px;
  height: 3px;
  border-radius: 999px;
  background: linear-gradient(90deg, var(--coral), var(--coral-2));
  transform: translateX(-50%);
}

.bind-card {
  min-height: 560px;
  margin-top: 28px;
  padding: 34px 24px 28px;
  border: 1px solid var(--card-border);
  border-radius: var(--card-radius-feature);
  background: var(--card-surface);
  box-shadow: var(--card-shadow-feature);
  backdrop-filter: blur(var(--card-blur));
  -webkit-backdrop-filter: blur(var(--card-blur));
}

.heart-art {
  position: relative;
  width: min(220px, 68vw);
  height: 170px;
  margin: 1px auto 0;
}

.heart-art.compact {
  width: 150px;
  height: 112px;
  margin-top: -6px;
  transform: scale(0.68);
  transform-origin: center;
}

.heart-art::before {
  content: "";
  position: absolute;
  left: 24px;
  right: 16px;
  top: 60px;
  height: 62px;
  border: 2px solid rgba(242, 126, 96, 0.26);
  border-radius: 50%;
  transform: rotate(-13deg);
}

.heart-core,
.heart-ring,
.mini-heart {
  position: absolute;
  transform: rotate(45deg);
}

.heart-core {
  left: 55px;
  top: 46px;
  width: 88px;
  height: 88px;
  border-radius: 18px 18px 20px 18px;
  background:
    radial-gradient(circle at 32% 22%, rgba(255, 255, 255, 0.72), transparent 23%),
    linear-gradient(135deg, #ff9b80 0%, #ff705e 45%, #e64b39 100%);
  box-shadow:
    0 18px 30px rgba(214, 68, 48, 0.28),
    inset -12px -12px 22px rgba(166, 45, 33, 0.12),
    inset 12px 12px 18px rgba(255, 255, 255, 0.2);
}

.heart-core::before,
.heart-core::after,
.heart-ring::before,
.heart-ring::after,
.mini-heart::before,
.mini-heart::after {
  content: "";
  position: absolute;
  border-radius: 50%;
}

.heart-core::before {
  width: 88px;
  height: 88px;
  left: -44px;
  top: 0;
  background: inherit;
}

.heart-core::after {
  width: 88px;
  height: 88px;
  left: 0;
  top: -44px;
  background: inherit;
}

.heart-ring {
  left: 121px;
  top: 75px;
  width: 57px;
  height: 57px;
  border: 12px solid #fff0da;
  border-radius: 15px;
  background: transparent;
  box-shadow:
    0 13px 24px rgba(148, 88, 50, 0.18),
    inset 0 0 0 1px rgba(169, 92, 49, 0.08);
}

.heart-ring::before,
.heart-ring::after {
  width: 57px;
  height: 57px;
  border: 12px solid #fff0da;
  background: transparent;
}

.heart-ring::before {
  left: -40px;
  top: -12px;
}

.heart-ring::after {
  left: -12px;
  top: -40px;
}

.mini-heart {
  width: 21px;
  height: 21px;
  border-radius: 5px;
  background: linear-gradient(135deg, #ffad95, #ef5a4d);
  box-shadow: 0 7px 15px rgba(233, 86, 69, 0.2);
}

.mini-heart::before,
.mini-heart::after {
  width: 21px;
  height: 21px;
  background: inherit;
}

.mini-heart::before {
  left: -10px;
  top: 0;
}

.mini-heart::after {
  left: 0;
  top: -10px;
}

.mini-heart.one {
  left: 20px;
  top: 66px;
  transform: rotate(35deg) scale(0.72);
}

.mini-heart.two {
  right: 16px;
  top: 28px;
  transform: rotate(31deg) scale(0.82);
}

.mini-heart.three {
  left: 57px;
  bottom: 22px;
  transform: rotate(36deg) scale(0.68);
}

.sparkle {
  position: absolute;
  width: 22px;
  height: 22px;
  color: rgba(239, 166, 124, 0.64);
}

.sparkle::before,
.sparkle::after {
  content: "";
  position: absolute;
  left: 50%;
  top: 50%;
  border-radius: 999px;
  background: currentColor;
  transform: translate(-50%, -50%);
}

.sparkle::before {
  width: 3px;
  height: 22px;
}

.sparkle::after {
  width: 22px;
  height: 3px;
}

.sparkle.s1 {
  left: 11px;
  top: 42px;
  transform: scale(0.65);
}

.sparkle.s2 {
  right: 7px;
  top: 98px;
  transform: scale(0.88);
}

.sparkle.s3 {
  left: 20px;
  top: 104px;
  transform: rotate(15deg) scale(0.78);
}

.desc {
  margin: 24px 0 0;
  color: var(--sub);
  font-size: 14.8px;
  font-weight: 540;
  line-height: 1.6;
  text-align: center;
  letter-spacing: 0;
  white-space: nowrap;
}

.code-box {
  position: relative;
  height: 128px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 28px;
  border: 1.5px dashed rgba(120, 90, 65, 0.22);
  border-radius: 22px;
  background: rgba(255, 248, 236, 0.5);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.56);
}

.code-label {
  position: absolute;
  top: 23px;
  left: 0;
  right: 0;
  color: #7a6a5f;
  font-size: 15px;
  font-weight: 740;
  text-align: center;
}

.code-value {
  margin-top: 32px;
  color: #6b5142;
  font-size: clamp(30px, 9vw, 36px);
  font-weight: 950;
  line-height: 1;
  letter-spacing: 7px;
  transform: translate(-23px, 4px);
}

.copy-btn {
  position: absolute;
  right: 16px;
  top: 50%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 5px 2px;
  color: #7a6a5f;
  background: transparent;
  font-size: 13px;
  font-weight: 720;
  transform: translateY(-50%);
  transition: color 180ms ease, transform 180ms ease, opacity 180ms ease;
}

.copy-btn:disabled {
  opacity: 0.48;
}

.copy-btn svg {
  width: 27px;
  height: 27px;
  stroke-width: 2.1;
}

.primary-btn,
.wechat-btn {
  width: 100%;
  height: 60px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 11px;
  border-radius: 18px;
  transition: transform 180ms ease, box-shadow 180ms ease, opacity 180ms ease;
}

.primary-btn {
  margin-top: 32px;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 16px 30px rgba(233, 86, 69, 0.28);
  font-size: 21px;
  font-weight: 900;
  letter-spacing: 0;
}

.wechat-btn {
  margin-top: 18px;
  border: 1px solid rgba(255, 255, 255, 0.58);
  color: #6b5142;
  background: rgba(255, 255, 255, 0.46);
  box-shadow:
    0 12px 24px rgba(80, 50, 30, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.68);
  font-size: 18px;
  font-weight: 820;
}

.primary-btn:disabled,
.wechat-btn:disabled {
  opacity: 0.58;
}

.wechat-btn svg {
  width: 29px;
  height: 29px;
  stroke-width: 2.1;
}

.status-text {
  margin: 14px 0 -6px;
  color: var(--coral-deep);
  font-size: 13px;
  font-weight: 760;
  text-align: center;
  animation: pulse 1.5s ease-in-out infinite;
}

.hint {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin: 26px 0 0;
  color: var(--sub);
  font-size: 15px;
  font-weight: 650;
  line-height: 1.45;
  text-align: center;
}

.hint svg {
  width: 23px;
  height: 23px;
  flex: 0 0 23px;
  color: var(--sage);
  stroke-width: 2.2;
}

.invite-input {
  width: 100%;
  height: 64px;
  margin-top: 28px;
  border: 1px solid rgba(120, 90, 65, 0.14);
  border-radius: 18px;
  outline: 0;
  color: var(--text);
  background: rgba(255, 255, 255, 0.55);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.58);
  font-size: 22px;
  font-weight: 850;
  letter-spacing: 4px;
  text-align: center;
  text-transform: uppercase;
}

.invite-input::placeholder {
  color: #aa9d92;
  font-size: 17px;
  font-weight: 700;
  letter-spacing: 0;
}

.error-text {
  min-height: 22px;
  margin: 10px 0 -8px;
  color: var(--coral-deep);
  font-size: 14px;
  font-weight: 740;
  text-align: center;
  opacity: 0;
  transition: opacity 160ms ease;
}

.error-text.show {
  opacity: 1;
}

.toast {
  position: fixed;
  left: 50%;
  bottom: calc(32px + env(safe-area-inset-bottom));
  z-index: 9;
  padding: 10px 16px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 999px;
  color: #fff;
  background: rgba(46, 36, 31, 0.78);
  box-shadow: 0 12px 24px rgba(46, 36, 31, 0.18);
  font-size: 14px;
  font-weight: 740;
  opacity: 0;
  pointer-events: none;
  transform: translate(-50%, 12px);
  transition: opacity 180ms ease, transform 180ms ease;
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
}

.toast.show {
  opacity: 1;
  transform: translate(-50%, 0);
}

.dialog-overlay {
  position: fixed;
  inset: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: rgba(46, 36, 31, 0.38);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.dialog-box {
  width: min(100%, 320px);
  padding: 22px 22px 24px;
  border: 1px solid rgba(255, 255, 255, 0.66);
  border-radius: var(--card-radius-feature);
  background: var(--card-surface-strong);
  box-shadow: var(--card-shadow-dialog);
  text-align: center;
}

.success-title {
  margin: -12px 0 8px;
  color: var(--text);
  font-size: 22px;
  font-weight: 900;
}

.success-desc {
  margin: 0 0 20px;
  color: var(--sub);
  font-size: 14px;
  line-height: 1.55;
}

.success-btn {
  width: 100%;
  height: 52px;
  border-radius: 16px;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  font-size: 17px;
  font-weight: 860;
}

.back-btn:active,
.tab-btn:active,
.primary-btn:active,
.wechat-btn:active,
.success-btn:active {
  transform: scale(0.98);
}

.copy-btn:active {
  transform: translateY(-50%) scale(0.96);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.52; }
}

@media (hover: hover) {
  .back-btn:hover,
  .tab-btn:hover,
  .primary-btn:hover,
  .wechat-btn:hover {
    transform: translateY(-1px);
  }

  .primary-btn:hover {
    box-shadow: 0 18px 34px rgba(233, 86, 69, 0.32);
  }

  .copy-btn:hover {
    color: var(--coral);
  }
}

@media (max-width: 380px) {
  .couple-bind-page {
    padding-left: 18px;
    padding-right: 18px;
  }

  .segment {
    margin-top: 30px;
  }

  .bind-card {
    padding: 30px 18px 26px;
    border-radius: var(--card-radius-feature);
  }

  .tab-btn {
    font-size: 18px;
  }

  .desc {
    font-size: 13.8px;
  }

  .code-value {
    letter-spacing: 6px;
  }

  .copy-btn {
    right: 12px;
  }

  .hint {
    font-size: 14px;
  }
}

@media (max-width: 350px) {
  .page-title {
    font-size: 25px;
  }

  .back-btn {
    width: 48px;
    height: 48px;
  }

  .segment {
    height: 62px;
  }

  .heart-art {
    transform: scale(0.9);
    transform-origin: center top;
    margin-bottom: -12px;
  }

  .desc {
    font-size: 13px;
  }
}

@media (min-width: 431px) {
  .couple-bind-page {
    box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.18);
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
  }
}
</style>
