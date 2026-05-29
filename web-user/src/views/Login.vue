<template>
  <div class="login-page">
    <div class="login-panel anim-delay-1">
      <div class="login-media">
        <img :src="heroArt" alt="login artwork" />
      </div>
      <div class="login-header">
        <p class="login-eyebrow">Editorial kitchen</p>
        <h1 class="login-title">今天吃什么</h1>
        <p class="login-sub">更安静的菜单浏览，更克制的日常决策。</p>
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
import heroArt from '@/assets/hero.png'

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
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: var(--sp-6);
  background: var(--color-bg);
}

.login-panel {
  width: 100%;
  max-width: 360px;
}

.login-media {
  overflow: hidden;
  border-radius: var(--r-lg);
  margin-bottom: var(--sp-6);
  background: var(--color-surface-2);
}

.login-media img {
  display: block;
  width: 100%;
  aspect-ratio: 16 / 10;
  object-fit: cover;
}

.login-header {
  margin-bottom: var(--sp-6);
}

.login-eyebrow {
  color: var(--color-text-3);
  font-size: var(--text-xs);
  font-weight: 600;
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
  color: var(--color-text-3);
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
  min-height: 48px;
  padding: 0 var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-sm);
  background: var(--color-surface);
  color: var(--color-text);
  font-family: var(--font-body);
  font-size: var(--text-base);
  transition: border-color var(--dur-base) var(--ease), box-shadow var(--dur-base) var(--ease);
}

.form-input::placeholder {
  color: var(--color-text-3);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px var(--color-accent-soft);
}

.form-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.form-error {
  margin-bottom: var(--sp-3);
  color: var(--color-error);
  font-size: var(--text-sm);
}

.login-btn {
  width: 100%;
  margin-top: var(--sp-1);
}

.login-footer {
  margin-top: var(--sp-4);
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
  transition: opacity var(--dur-fast) var(--ease);
}

.login-toggle:hover {
  opacity: 0.8;
}

/* ── Responsive ── */
@media (min-width: 768px) {
  .login-panel {
    max-width: 400px;
  }
}
</style>
