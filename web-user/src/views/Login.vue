<template>
  <div class="login-page">
    <main class="phone-frame" :style="{ '--login-bg': loginBackgroundImage }">
      <section class="hero-scene" aria-label="今日菜单视觉预览"></section>

      <section class="login-panel" aria-label="登录">
        <div class="brand-icon" aria-hidden="true">
          <svg viewBox="0 0 64 64">
            <path class="pot" d="M17 28h30v16a8 8 0 0 1-8 8H25a8 8 0 0 1-8-8V28Z" />
            <path class="pot-line" d="M13 31h38M22 25h20M28 20c1.2-3 3.5-4.5 7-4.5 2.8 0 5 1.1 6.5 3.2" />
            <path class="leaf" d="M43 49c8-1 12-6 12-14-8 .5-13.2 5-12 14Z" />
            <path class="heart" d="M45.5 14.5c-2.6-3.2-8.4-1.4-8.4 3.3 0 4.6 8.4 9.4 8.4 9.4s8.4-4.8 8.4-9.4c0-4.7-5.8-6.5-8.4-3.3Z" />
          </svg>
        </div>

        <header class="login-header">
          <h1>今天吃什么</h1>
          <p class="english-title">Daily Kitchen</p>
          <p class="login-sub">把菜谱、菜单和采购清单放在一处</p>
        </header>

        <form class="login-form" @submit.prevent="handleSubmit">
          <label class="form-field">
            <span class="field-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24">
                <path d="M20 21a8 8 0 0 0-16 0M12 13a5 5 0 1 0 0-10 5 5 0 0 0 0 10Z" />
              </svg>
            </span>
            <input
              v-model="form.username"
              type="text"
              placeholder="用户名"
              autocomplete="username"
              required
            />
          </label>

          <label class="form-field">
            <span class="field-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24">
                <path d="M7 11V8a5 5 0 0 1 10 0v3M6 11h12v10H6V11Z" />
              </svg>
            </span>
            <input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="密码"
              autocomplete="current-password"
              required
            />
            <button class="icon-button" type="button" :aria-label="showPassword ? '隐藏密码' : '显示密码'" @click="showPassword = !showPassword">
              <svg v-if="showPassword" viewBox="0 0 24 24">
                <path d="M2 12s3.5-6 10-6 10 6 10 6-3.5 6-10 6S2 12 2 12Z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
              <svg v-else viewBox="0 0 24 24">
                <path d="m3 3 18 18M10.6 10.6A3 3 0 0 0 12 15a3 3 0 0 0 2.8-4M9.9 5.3A10.7 10.7 0 0 1 12 5c6.5 0 10 7 10 7a18.8 18.8 0 0 1-3.1 4.1M6.6 6.7C3.6 8.7 2 12 2 12s3.5 7 10 7c1.6 0 3-.4 4.2-1" />
              </svg>
            </button>
          </label>

          <label v-if="isRegister" class="form-field">
            <span class="field-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24">
                <path d="M7 11V8a5 5 0 0 1 10 0v3M6 11h12v10H6V11Z" />
              </svg>
            </span>
            <input
              v-model="form.confirmPassword"
              :type="showPassword ? 'text' : 'password'"
              placeholder="确认密码"
              autocomplete="new-password"
              required
            />
          </label>

          <p v-if="error" class="form-error">{{ error }}</p>

          <button type="submit" class="login-btn" :disabled="loading">
            {{ loading ? '请稍候...' : (isRegister ? '注册' : '登录') }}
          </button>
        </form>

        <p class="login-footer">
          <span>{{ isRegister ? '已有账号，' : '没有账号，' }}</span>
          <button type="button" @click="isRegister = !isRegister">
            {{ isRegister ? '去登录' : '去注册' }}
          </button>
        </p>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { login as loginApi, register as registerApi } from '@/api/auth'
import { getPreferenceStatus } from '@/api/user'
import loginKitchenReference from '@/assets/login-kitchen-reference.png'

const router = useRouter()
const userStore = useUserStore()

const isRegister = ref(false)
const loading = ref(false)
const error = ref('')
const showPassword = ref(false)
const loginBackgroundImage = computed(() => `url(${loginKitchenReference})`)

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
    await goAfterAuth()
  } catch (e: any) {
    error.value = e.message || '操作失败'
  } finally {
    loading.value = false
  }
}

async function goAfterAuth() {
  try {
    const status = await getPreferenceStatus()
    if (!status.completed) {
      router.replace({ path: '/user/preferences', query: { guide: '1' } })
      return
    }
  } catch {
    // 偏好状态读取失败时不阻塞登录。
  }
  router.replace('/')
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: grid;
  place-items: center;
  background:
    radial-gradient(circle at 50% 18%, rgba(255, 248, 232, 0.74), transparent 32%),
    linear-gradient(145deg, #dec49b, #fff6e7 42%, #c6b08d);
  color: #2f2a27;
  font-family: "Noto Sans SC", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

.phone-frame {
  position: relative;
  width: min(100vw, 430px);
  min-height: 100vh;
  min-height: 100dvh;
  overflow: hidden;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.03), rgba(226, 197, 157, 0.2)),
    var(--login-bg) center -44px / cover no-repeat;
  isolation: isolate;
}

.phone-frame::after {
  content: "";
  position: absolute;
  inset: 42% 0 0;
  z-index: 0;
  background: linear-gradient(180deg, rgba(255, 248, 232, 0), rgba(221, 192, 149, 0.32) 52%, rgba(208, 178, 134, 0.58));
  pointer-events: none;
}

.hero-scene {
  position: relative;
  z-index: 1;
  height: min(47vh, 405px);
  min-height: 348px;
}

.login-panel {
  position: relative;
  z-index: 2;
  width: calc(100% - 56px);
  max-width: 344px;
  margin: -2px auto 38px;
  padding: 19px 23px 17px;
  border: 1px solid rgba(255, 255, 255, 0.86);
  border-radius: 25px;
  background:
    radial-gradient(circle at 72% 12%, rgba(219, 236, 196, 0.48), transparent 26%),
    radial-gradient(circle at 12% 8%, rgba(255, 246, 229, 0.74), transparent 34%),
    rgba(255, 252, 247, 0.88);
  box-shadow:
    0 22px 52px rgba(75, 49, 27, 0.26),
    inset 0 1px 0 rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(24px) saturate(1.18);
}

.brand-icon {
  width: 64px;
  height: 64px;
  display: grid;
  place-items: center;
  margin: -2px auto 12px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.72);
  box-shadow:
    0 12px 24px rgba(151, 82, 50, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.96);
}

.brand-icon svg {
  width: 42px;
  height: 42px;
}

.brand-icon path {
  fill: none;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.pot {
  fill: #e5483e;
  stroke: #e5483e;
}

.pot-line {
  stroke: #e5483e;
  stroke-width: 4;
}

.leaf {
  fill: #6d8d5b;
  stroke: #6d8d5b;
}

.heart {
  fill: #e5483e;
  stroke: #e5483e;
}

.login-header {
  text-align: center;
}

.login-header h1 {
  margin: 0;
  color: #2f2b28;
  font-size: clamp(32px, 8.4vw, 40px);
  font-weight: 800;
  line-height: 1.06;
  letter-spacing: 0;
}

.english-title {
  margin: 6px 0 0;
  color: #748762;
  font-family: Georgia, "Times New Roman", serif;
  font-size: 19px;
  font-weight: 600;
  line-height: 1.1;
}

.login-sub {
  margin: 9px 0 18px;
  color: #8b8179;
  font-size: 15px;
  font-weight: 500;
  line-height: 1.3;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 11px;
}

.form-field {
  min-height: 52px;
  display: grid;
  grid-template-columns: 36px minmax(0, 1fr) auto;
  align-items: center;
  gap: 8px;
  padding: 0 14px;
  border: 1px solid rgba(109, 95, 82, 0.16);
  border-radius: 15px;
  background: rgba(255, 255, 255, 0.67);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.76),
    0 7px 18px rgba(94, 63, 38, 0.06);
  transition:
    border-color 180ms ease,
    box-shadow 180ms ease,
    background 180ms ease;
}

.form-field:focus-within {
  border-color: rgba(229, 72, 62, 0.42);
  background: rgba(255, 255, 255, 0.84);
  box-shadow:
    0 0 0 4px rgba(229, 72, 62, 0.1),
    0 10px 20px rgba(94, 63, 38, 0.08);
}

.field-icon,
.icon-button {
  width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  color: #aaa29b;
}

.field-icon svg,
.icon-button svg {
  width: 21px;
  height: 21px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.form-field input {
  width: 100%;
  min-width: 0;
  border: 0;
  outline: 0;
  background: transparent;
  color: #3e3731;
  font: inherit;
  font-size: 17px;
  font-weight: 600;
  letter-spacing: 0;
}

.form-field input::placeholder {
  color: #aaa29b;
  font-weight: 600;
}

.icon-button {
  padding: 0;
  border: 0;
  border-radius: 999px;
  background: transparent;
  cursor: pointer;
  transition:
    color 160ms ease,
    background 160ms ease;
}

.icon-button:hover {
  color: #7f766f;
  background: rgba(127, 118, 111, 0.08);
}

.form-error {
  margin: 0;
  padding: 10px 12px;
  border-radius: 12px;
  background: rgba(229, 72, 62, 0.1);
  color: #d83e36;
  font-size: 13px;
  font-weight: 600;
}

.login-btn {
  min-height: 52px;
  margin-top: 5px;
  border: 0;
  border-radius: 14px;
  background: linear-gradient(135deg, #dd332d 0%, #f25549 100%);
  color: #fff;
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 0;
  box-shadow:
    0 14px 26px rgba(222, 61, 52, 0.25),
    inset 0 1px 0 rgba(255, 255, 255, 0.22);
  cursor: pointer;
  transition:
    transform 160ms ease,
    box-shadow 160ms ease,
    opacity 160ms ease;
}

.login-btn:hover:not(:disabled) {
  box-shadow:
    0 16px 30px rgba(222, 61, 52, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.25);
}

.login-btn:active:not(:disabled) {
  transform: translateY(1px);
}

.login-btn:disabled {
  opacity: 0.68;
  cursor: not-allowed;
}

.login-footer {
  margin: 12px 0 0;
  color: #817870;
  text-align: center;
  font-size: 16px;
  font-weight: 500;
}

.login-footer button {
  padding: 0;
  border: 0;
  background: transparent;
  color: #df4038;
  font: inherit;
  font-weight: 800;
  cursor: pointer;
}

@media (min-width: 431px) {
  .phone-frame {
    min-height: min(100vh, 932px);
    box-shadow: 0 24px 70px rgba(77, 53, 31, 0.36);
  }
}

@media (max-width: 374px) {
  .hero-scene {
    min-height: 360px;
  }

  .login-panel {
    width: calc(100% - 38px);
    padding: 20px 18px 18px;
    margin-bottom: 42px;
  }

  .login-header h1 {
    font-size: 32px;
  }
}
</style>
