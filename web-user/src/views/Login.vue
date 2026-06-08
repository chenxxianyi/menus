<template>
  <div class="login-page">
    <div class="login-panel anim-delay-1">
      <div class="login-media" aria-hidden="true">
        <div class="media-shelf">
          <div class="menu-sheet sheet-back">
            <span></span>
            <span></span>
            <span></span>
          </div>
          <div class="menu-sheet sheet-front">
            <div class="plate">
              <span class="plate-core"></span>
            </div>
            <div class="dish-copy">
              <span class="dish-title"></span>
              <span class="dish-line wide"></span>
              <span class="dish-line"></span>
            </div>
          </div>
        </div>
      </div>
      <div class="login-header">
        <p class="login-eyebrow">Daily kitchen</p>
        <h1 class="login-title">今天吃什么</h1>
        <p class="login-sub">把每日菜谱、购物清单和口味偏好放在一处，晚餐决定更轻松。</p>
      </div>

      <form class="login-form anim-delay-2" @submit.prevent="handleSubmit">
        <div class="form-field">
          <input
            v-model="form.username"
            type="text"
            placeholder="用户名"
            class="form-input"
            required
          />
        </div>
        <div class="form-field">
          <input
            v-model="form.password"
            type="password"
            placeholder="密码"
            class="form-input"
            required
          />
        </div>
        <div v-if="isRegister" class="form-field">
          <input
            v-model="form.confirmPassword"
            type="password"
            placeholder="确认密码"
            class="form-input"
            required
          />
        </div>
        <div v-if="error" class="form-error">{{ error }}</div>
        <button type="submit" class="btn-solid login-btn" :disabled="loading">
          {{ loading ? '请稍候...' : (isRegister ? '注册' : '登录') }}
        </button>
      </form>

      <div class="login-footer anim-delay-3">
        <button class="login-toggle" @click="isRegister = !isRegister">
          {{ isRegister ? '已有账号，去登录' : '没有账号，去注册' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { login as loginApi, register as registerApi } from '@/api/auth'

const router = useRouter()
const userStore = useUserStore()

const isRegister = ref(false)
const loading = ref(false)
const error = ref('')

const form = reactive({
  username: '',
  password: '',
  confirmPassword: '',
})

async function handleSubmit() {
  error.value = ''
  if (!form.username || !form.password) {
    error.value = '请填写用户名和密码'
    return
  }
  if (isRegister.value && form.password !== form.confirmPassword) {
    error.value = '两次密码不一致'
    return
  }

  loading.value = true
  try {
    if (isRegister.value) {
      const res: any = await registerApi({
        username: form.username,
        password: form.password,
        email: '',
        confirm_password: form.password,
      })
      localStorage.setItem('token', res.token)
      await userStore.fetchUserInfo()
    } else {
      await userStore.login({ username: form.username, password: form.password })
    }
    router.replace('/')
  } catch (e: any) {
    error.value = e.message || '操作失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* Hallmark - component: login screen - genre: quiet utility - theme: fresh kitchen
 * states: default - hover - focus - active - disabled - loading - error - success
 * contrast: pass
 */
.login-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: grid;
  place-items: center;
  padding: clamp(var(--sp-4), 5vw, var(--sp-8));
  background:
    linear-gradient(160deg, rgba(7, 193, 96, 0.13), transparent 38%),
    linear-gradient(340deg, rgba(245, 158, 11, 0.11), transparent 34%),
    var(--color-bg);
}

.login-panel {
  width: 100%;
  max-width: 392px;
  padding: var(--sp-5);
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: var(--r-xl);
  background: rgba(255, 255, 255, 0.78);
  box-shadow: 0 18px 60px rgba(31, 41, 55, 0.10);
  backdrop-filter: blur(20px);
}

.login-media {
  position: relative;
  overflow: hidden;
  min-height: 218px;
  border: 1px solid rgba(255, 255, 255, 0.82);
  border-radius: var(--r-lg);
  margin-bottom: var(--sp-6);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.82), rgba(255, 255, 255, 0.42)),
    linear-gradient(135deg, #effaf0 0%, #fff8ec 100%);
  box-shadow: inset 0 -1px 0 rgba(31, 41, 55, 0.05);
}

.login-media::before {
  content: "";
  position: absolute;
  inset: var(--sp-4);
  border: 1px solid rgba(7, 193, 96, 0.14);
  border-radius: calc(var(--r-lg) - 4px);
}

.login-media::after {
  content: "";
  position: absolute;
  right: -34px;
  bottom: -42px;
  width: 144px;
  height: 144px;
  border-radius: var(--r-full);
  background: rgba(245, 158, 11, 0.18);
}

.media-shelf {
  position: relative;
  z-index: 1;
  min-height: 218px;
}

.menu-sheet {
  position: absolute;
  left: 50%;
  width: min(78%, 274px);
  border: 1px solid rgba(31, 41, 55, 0.08);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 14px 36px rgba(31, 41, 55, 0.10);
  transform: translateX(-50%) rotate(-3deg);
}

.sheet-back {
  top: 30px;
  height: 116px;
  padding: var(--sp-4);
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.98), rgba(232, 245, 233, 0.86));
}

.sheet-back span {
  display: block;
  height: 10px;
  border-radius: var(--r-full);
  background: rgba(31, 41, 55, 0.08);
}

.sheet-back span + span {
  margin-top: var(--sp-3);
}

.sheet-back span:nth-child(1) {
  width: 38%;
  background: rgba(7, 193, 96, 0.22);
}

.sheet-back span:nth-child(2) {
  width: 74%;
}

.sheet-back span:nth-child(3) {
  width: 58%;
}

.sheet-front {
  top: 74px;
  min-height: 132px;
  display: grid;
  grid-template-columns: 78px minmax(0, 1fr);
  gap: var(--sp-4);
  align-items: center;
  padding: var(--sp-4);
  transform: translateX(-50%) rotate(2deg);
}

.plate {
  position: relative;
  width: 76px;
  height: 76px;
  display: grid;
  place-items: center;
  border-radius: var(--r-full);
  background:
    radial-gradient(circle, rgba(255, 255, 255, 0.92) 0 42%, transparent 43%),
    conic-gradient(from 20deg, #07c160, #f59e0b, #22c55e, #07c160);
  box-shadow: inset 0 0 0 8px rgba(255, 255, 255, 0.78);
}

.plate-core {
  width: 34px;
  height: 34px;
  border-radius: var(--r-full);
  background:
    radial-gradient(circle at 34% 28%, #ffffff 0 14%, transparent 15%),
    linear-gradient(135deg, #22c55e, #f59e0b);
  box-shadow: 0 8px 16px rgba(7, 193, 96, 0.16);
}

.dish-copy span {
  display: block;
  border-radius: var(--r-full);
}

.dish-title {
  width: 62%;
  height: 14px;
  margin-bottom: var(--sp-3);
  background: rgba(31, 41, 55, 0.82);
}

.dish-line {
  width: 64%;
  height: 9px;
  background: rgba(156, 163, 175, 0.34);
}

.dish-line.wide {
  width: 88%;
  margin-bottom: var(--sp-2);
}

.login-header {
  margin-bottom: var(--sp-6);
}

.login-eyebrow {
  color: var(--color-accent-hover);
  font-size: var(--text-xs);
  font-weight: 600;
  letter-spacing: 0;
  margin-bottom: var(--sp-1);
}

.login-title {
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: 650;
  line-height: 1.08;
}

.login-sub {
  margin-top: var(--sp-2);
  max-width: 23em;
  color: var(--color-text-2);
  font-size: var(--text-sm);
  line-height: 1.65;
}

.login-form {
  display: flex;
  flex-direction: column;
}

.form-field {
  margin-bottom: var(--sp-3);
}

.form-input {
  width: 100%;
  min-height: 52px;
  padding: 0 var(--sp-4);
  border: 1px solid rgba(31, 41, 55, 0.08);
  border-radius: var(--r-sm);
  background: rgba(255, 255, 255, 0.92);
  color: var(--color-text);
  font-family: var(--font-body);
  font-size: var(--text-base);
  box-shadow: 0 1px 0 rgba(31, 41, 55, 0.04);
  transition:
    border-color var(--dur-base) var(--ease),
    box-shadow var(--dur-base) var(--ease),
    background var(--dur-base) var(--ease);
}

.form-input::placeholder {
  color: var(--color-text-3);
}

.form-input:hover {
  border-color: rgba(7, 193, 96, 0.36);
}

.form-input:focus,
.form-input:focus-visible {
  outline: none;
  border-color: var(--color-accent);
  background: var(--color-surface);
  box-shadow: 0 0 0 4px rgba(7, 193, 96, 0.14);
}

.form-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.form-error {
  margin-bottom: var(--sp-3);
  padding: var(--sp-3);
  border-radius: var(--r-sm);
  background: var(--color-error-soft);
  color: var(--color-error);
  font-size: var(--text-sm);
}

.login-btn {
  width: 100%;
  min-height: 52px;
  margin-top: var(--sp-1);
  box-shadow: 0 10px 22px rgba(31, 41, 55, 0.14);
}

.login-footer {
  margin-top: var(--sp-5);
  text-align: center;
}

.login-toggle {
  border: 0;
  background: none;
  color: var(--color-accent);
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition:
    color var(--dur-fast) var(--ease),
    opacity var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease);
}

.login-toggle:hover {
  color: var(--color-accent-hover);
}

.login-toggle:active {
  transform: translateY(1px);
}

/* ── Responsive ── */
@media (min-width: 768px) {
  .login-panel {
    max-width: 412px;
    padding: var(--sp-6);
  }
}

@media (max-width: 374px) {
  .login-page {
    padding: var(--sp-3);
  }

  .login-panel {
    padding: var(--sp-4);
  }

  .login-media,
  .media-shelf {
    min-height: 194px;
  }

  .sheet-front {
    grid-template-columns: 64px minmax(0, 1fr);
  }

  .plate {
    width: 64px;
    height: 64px;
  }
}
</style>
