<template>
  <main class="about-shell" :style="pageVars">
    <header class="page-header" aria-label="页面顶部">
      <button class="back-btn" type="button" aria-label="返回" @click="router.back()">
        <svg viewBox="0 0 24 24" aria-hidden="true">
          <path d="m15 18-6-6 6-6" />
        </svg>
      </button>
      <h1 class="page-title">关于我们</h1>
    </header>

    <section class="brand-card glass-card" aria-labelledby="brand-title">
      <div class="app-icon" aria-hidden="true">
        <div class="pot-mark">
          <svg viewBox="0 0 24 24">
            <path d="M12 20C7 16 4 13.2 4 9.8A4.2 4.2 0 0 1 11.2 7L12 7.9l.8-.9A4.2 4.2 0 0 1 20 9.8C20 13.2 17 16 12 20Z" />
          </svg>
        </div>
      </div>
      <div class="brand-copy">
        <h2 id="brand-title">{{ brandInfo.name }}</h2>
        <p class="subtitle">{{ brandInfo.subtitle }}</p>
        <div class="brand-rule" aria-hidden="true"></div>
        <p class="brand-desc">{{ brandInfo.description }}</p>
      </div>
    </section>

    <section v-if="features.length" class="feature-card glass-card" aria-label="核心功能">
      <article
        v-for="feature in features"
        :key="feature.title"
        class="feature-row"
        :style="{ '--feature-color': feature.color, '--feature-bg': feature.bg }"
      >
        <span class="feature-icon" aria-hidden="true">
          <AboutIcon :name="feature.icon" />
        </span>
        <div class="feature-copy">
          <h3>{{ feature.title }}</h3>
          <p>{{ feature.description }}</p>
        </div>
      </article>
    </section>

    <section v-if="infoRows.length" class="list-card glass-card" aria-label="联系与版本信息">
      <button
        v-for="row in infoRows"
        :key="row.label"
        class="info-row"
        type="button"
        @click="handleInfoAction(row.action)"
      >
        <span class="row-icon" :style="{ '--row-color': row.color }" aria-hidden="true">
          <AboutIcon :name="row.icon" />
        </span>
        <span class="row-label">{{ row.label }}</span>
        <span class="row-value">{{ row.value }}</span>
        <ChevronIcon />
      </button>
    </section>

    <section v-if="legalRows.length" class="list-card glass-card" aria-label="协议入口">
      <button
        v-for="row in legalRows"
        :key="row.label"
        class="info-row"
        type="button"
        @click="handleLegalAction(row.path)"
      >
        <span class="row-icon" :style="{ '--row-color': row.color }" aria-hidden="true">
          <AboutIcon :name="row.icon" />
        </span>
        <span class="row-label">{{ row.label }}</span>
        <span class="row-value"></span>
        <ChevronIcon />
      </button>
    </section>

    <footer class="slogan" aria-label="品牌标语">
      <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="M7 20c0-5 2.5-8 8-9" />
        <path d="M9 20c5.5 0 9-3.5 9-9V4h-7c-5.5 0-9 3.5-9 9 0 2.1.7 3.9 2 5.2" />
      </svg>
      <div class="slogan-line">
        <span>{{ brandInfo.slogan }}</span>
      </div>
    </footer>

    <div class="toast" :class="{ show: !!toastText }">{{ toastText }}</div>
  </main>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getAboutInfo, type AboutInfo } from '@/api/app'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'
import { copyText } from '@/utils/clipboard'

type IconName = 'book' | 'calendar' | 'basket' | 'mail' | 'message' | 'info' | 'file' | 'shield'
type InfoAction = 'contact' | 'wechat' | 'version'

interface FeatureItem {
  title: string
  description: string
  icon: IconName
  color: string
  bg: string
}

interface InfoRow {
  label: string
  value: string
  icon: IconName
  color: string
  action: InfoAction
}

interface LegalRow {
  label: string
  icon: IconName
  color: string
  path: string
}

const router = useRouter()
const toastText = ref('')
const aboutInfo = ref<AboutInfo | null>(null)
let toastTimer: ReturnType<typeof setTimeout> | null = null

const brandInfo = computed(() => ({
  name: aboutInfo.value?.name || '后端未配置应用名称',
  subtitle: aboutInfo.value?.subtitle || '',
  description: aboutInfo.value?.description || '后端未配置品牌介绍',
  slogan: aboutInfo.value?.slogan || '后端未配置品牌标语',
  version: aboutInfo.value?.version || '',
  email: aboutInfo.value?.email || '',
  wechat: aboutInfo.value?.wechat || '',
}))

const features = computed<FeatureItem[]>(() => {
  const list = aboutInfo.value?.features
  if (!Array.isArray(list)) return []
  return list
    .filter((feature) => feature.title && feature.description)
    .map((feature) => ({
      title: feature.title,
      description: feature.description,
      icon: normalizeIcon(feature.icon, 'book'),
      color: feature.color || '#e95645',
      bg: feature.bg || 'rgba(233, 86, 69, 0.14)',
    }))
})

const infoRows = computed<InfoRow[]>(() => {
  const rows: InfoRow[] = []
  if (brandInfo.value.email) {
    rows.push({ label: '联系我们', value: brandInfo.value.email, icon: 'mail', color: '#e95645', action: 'contact' })
  }
  if (brandInfo.value.wechat) {
    rows.push({ label: '微信公众号', value: brandInfo.value.wechat, icon: 'message', color: '#7ea36a', action: 'wechat' })
  }
  if (brandInfo.value.version) {
    rows.push({ label: '当前版本', value: brandInfo.value.version, icon: 'info', color: '#8b715e', action: 'version' })
  }
  return rows
})

const legalRows = computed<LegalRow[]>(() => {
  const rows: LegalRow[] = []
  if (aboutInfo.value?.terms_url) {
    rows.push({ label: '用户协议', icon: 'file', color: '#8b715e', path: aboutInfo.value.terms_url })
  }
  if (aboutInfo.value?.privacy_url) {
    rows.push({ label: '隐私政策', icon: 'shield', color: '#8b715e', path: aboutInfo.value.privacy_url })
  }
  return rows
})

const pageVars = computed(() => ({
  '--about-bg': `url(${kitchenBg})`,
}))

function showToast(message: string) {
  toastText.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastText.value = ''
  }, 1500)
}

async function copyEmail() {
  if (await copyText(brandInfo.value.email)) {
    showToast('邮箱已复制')
    return
  }
  showToast(`邮箱：${brandInfo.value.email}`)
}

function handleInfoAction(action: InfoAction) {
  if (action === 'contact') {
    copyEmail()
    return
  }
  if (action === 'wechat') {
    showToast(`微信公众号：${brandInfo.value.wechat}`)
    return
  }
  showToast(`当前版本 ${brandInfo.value.version}`)
}

function handleLegalAction(path: string) {
  if (/^https?:\/\//.test(path)) {
    window.location.href = path
    return
  }
  router.push(path)
}

function normalizeIcon(value: unknown, fallback: IconName): IconName {
  const icons: IconName[] = ['book', 'calendar', 'basket', 'mail', 'message', 'info', 'file', 'shield']
  return icons.includes(value as IconName) ? value as IconName : fallback
}

async function loadAboutInfo() {
  try {
    aboutInfo.value = await getAboutInfo()
  } catch {
    aboutInfo.value = null
    showToast('应用信息读取失败')
  }
}

onMounted(() => {
  loadAboutInfo()
})

onUnmounted(() => {
  if (toastTimer) clearTimeout(toastTimer)
})

const ChevronIcon = defineComponent({
  name: 'ChevronIcon',
  setup() {
    return () => h('svg', {
      class: 'chevron',
      viewBox: '0 0 24 24',
      'aria-hidden': 'true',
    }, [
      h('path', { d: 'm9 18 6-6-6-6' }),
    ])
  },
})

const AboutIcon = defineComponent({
  name: 'AboutIcon',
  props: {
    name: { type: String, required: true },
  },
  setup(props) {
    return () => {
      const common = {
        viewBox: '0 0 24 24',
        'aria-hidden': 'true',
      }

      if (props.name === 'book') {
        return h('svg', common, [
          h('path', { d: 'M4 19.5A2.5 2.5 0 0 1 6.5 17H20' }),
          h('path', { d: 'M4 4.5A2.5 2.5 0 0 1 6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15Z' }),
          h('path', { d: 'M9 7h6' }),
          h('path', { d: 'M9 11h7' }),
        ])
      }

      if (props.name === 'calendar') {
        return h('svg', common, [
          h('rect', { x: '3', y: '4', width: '18', height: '18', rx: '3' }),
          h('path', { d: 'M8 2v4' }),
          h('path', { d: 'M16 2v4' }),
          h('path', { d: 'M3 10h18' }),
          h('path', { d: 'm9 16 2 2 4-5' }),
        ])
      }

      if (props.name === 'basket') {
        return h('svg', common, [
          h('path', { d: 'M6 8h12l-1.1 11.2A2 2 0 0 1 14.9 21H9.1a2 2 0 0 1-2-1.8L6 8Z' }),
          h('path', { d: 'M9 8V6a3 3 0 0 1 6 0v2' }),
          h('path', { d: 'M9.5 13.5 11 15l3.5-4' }),
        ])
      }

      if (props.name === 'mail') {
        return h('svg', common, [
          h('rect', { x: '3', y: '5', width: '18', height: '14', rx: '2' }),
          h('path', { d: 'm3 7 9 6 9-6' }),
        ])
      }

      if (props.name === 'message') {
        return h('svg', common, [
          h('path', { d: 'M9.5 15.5c-3.3 0-6-2.1-6-4.8s2.7-4.8 6-4.8 6 2.1 6 4.8-2.7 4.8-6 4.8Z' }),
          h('path', { d: 'm6.5 16.8-2.3 1 .8-2' }),
          h('path', { d: 'M14.5 10.1c3.3.3 5.7 2.4 5.7 5 0 1.3-.6 2.5-1.7 3.4l.5 1.7-2-.8c-.8.3-1.6.5-2.5.5-2.7 0-5-1.4-5.7-3.4' }),
        ])
      }

      if (props.name === 'info') {
        return h('svg', common, [
          h('circle', { cx: '12', cy: '12', r: '9' }),
          h('path', { d: 'M12 11v5' }),
          h('path', { d: 'M12 8h.01' }),
        ])
      }

      if (props.name === 'file') {
        return h('svg', common, [
          h('path', { d: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8Z' }),
          h('path', { d: 'M14 2v6h6' }),
          h('path', { d: 'M8 13h8' }),
          h('path', { d: 'M8 17h5' }),
        ])
      }

      return h('svg', common, [
        h('path', { d: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z' }),
        h('path', { d: 'm9 12 2 2 4-5' }),
      ])
    }
  },
})
</script>

<style scoped>
.about-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --muted: #a19489;
  --cream: rgba(255, 250, 240, 0.82);
  --coral: #e95645;
  --sage: #8fa783;
  --border: rgba(255, 255, 255, 0.62);
  --line: rgba(120, 90, 65, 0.12);
  --shadow: 0 22px 50px rgba(80, 50, 30, 0.16);
  position: relative;
  width: min(100%, 430px);
  min-height: calc(100vh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 150px);
  min-height: calc(100dvh + var(--tab-h, 64px) + var(--safe-bottom, 34px) + 150px);
  margin: 0 auto;
  padding: max(52px, env(safe-area-inset-top)) 24px calc(44px + env(safe-area-inset-bottom));
  overflow-x: clip;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 237, 205, 0.34), rgba(255, 247, 233, 0.18)),
    var(--about-bg) center top / cover fixed;
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif;
}

.about-shell::before {
  content: "";
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 18% 4%, rgba(255, 255, 255, 0.76), transparent 30%),
    radial-gradient(circle at 86% 14%, rgba(238, 143, 66, 0.22), transparent 34%),
    radial-gradient(circle at 11% 92%, rgba(233, 86, 69, 0.15), transparent 31%),
    linear-gradient(90deg, rgba(255, 239, 214, 0.58), rgba(255, 245, 230, 0.18) 55%, rgba(172, 91, 33, 0.18));
  backdrop-filter: blur(4px) saturate(1.12);
  -webkit-backdrop-filter: blur(4px) saturate(1.12);
}

.page-header,
.glass-card,
.slogan {
  position: relative;
  z-index: 1;
}

button {
  border: 0;
  font: inherit;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

svg {
  display: block;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.page-header {
  height: 106px;
  display: grid;
  place-items: start center;
  text-align: center;
}

.back-btn {
  position: absolute;
  top: 0;
  left: 0;
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
  transition: transform 180ms ease, box-shadow 180ms ease, opacity 180ms ease;
}

.back-btn svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.55;
}

.page-title {
  margin: 15px 0 0;
  color: var(--text);
  font-size: 29px;
  font-weight: 950;
  line-height: 1;
  letter-spacing: 0;
}

.glass-card {
  border: 1px solid var(--card-border);
  background: var(--card-surface);
  box-shadow: var(--card-shadow);
  backdrop-filter: blur(var(--card-blur));
  -webkit-backdrop-filter: blur(var(--card-blur));
}

.brand-card {
  min-height: 164px;
  display: flex;
  align-items: center;
  gap: 22px;
  padding: 28px;
  border-radius: var(--card-radius-feature);
  box-shadow: var(--card-shadow-feature);
}

.app-icon {
  position: relative;
  width: 108px;
  height: 108px;
  display: grid;
  flex: 0 0 108px;
  place-items: center;
  overflow: hidden;
  border-radius: 26px;
  color: #fff8ec;
  background:
    radial-gradient(circle at 26% 18%, rgba(255, 255, 255, 0.28), transparent 28%),
    linear-gradient(135deg, #f47a66, #df4638);
  box-shadow:
    0 14px 30px rgba(233, 86, 69, 0.22),
    inset 0 1px 0 rgba(255, 255, 255, 0.34);
}

.app-icon::before,
.app-icon::after {
  content: "";
  position: absolute;
  border-radius: 999px;
  background: rgba(255, 248, 236, 0.86);
  box-shadow: 0 5px 12px rgba(80, 50, 30, 0.08);
}

.app-icon::before {
  left: 19px;
  top: 39px;
  width: 15px;
  height: 15px;
}

.app-icon::after {
  right: 15px;
  top: 24px;
  width: 24px;
  height: 22px;
  clip-path: polygon(50% 94%, 14% 58%, 8% 35%, 20% 15%, 42% 18%, 50% 30%, 58% 18%, 80% 15%, 92% 35%, 86% 58%);
}

.pot-mark {
  position: relative;
  width: 68px;
  height: 54px;
  display: grid;
  place-items: center;
  margin-top: 16px;
  border-radius: 17px 17px 21px 21px;
  color: var(--coral);
  background: #fff4e4;
  box-shadow:
    0 10px 18px rgba(80, 50, 30, 0.18),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
}

.pot-mark::before {
  content: "";
  position: absolute;
  top: -6px;
  width: 52px;
  height: 10px;
  border-radius: 999px;
  background: #fff8ec;
}

.pot-mark svg {
  position: relative;
  z-index: 1;
  width: 28px;
  height: 28px;
  fill: currentColor;
  stroke-width: 0;
}

.brand-copy {
  min-width: 0;
}

.brand-copy h2 {
  margin: 0;
  color: var(--text);
  font-size: 32px;
  font-weight: 950;
  line-height: 1.05;
  letter-spacing: 0;
  white-space: nowrap;
}

.brand-copy .subtitle {
  margin: 8px 0 0;
  color: var(--sub);
  font-size: 20px;
  font-weight: 620;
  line-height: 1;
}

.brand-rule {
  width: 24px;
  height: 4px;
  margin-top: 16px;
  border-radius: 999px;
  background: var(--coral);
}

.brand-desc {
  margin: 16px 0 0;
  color: var(--sub);
  font-size: 16px;
  font-weight: 620;
  line-height: 1.6;
}

.feature-card {
  margin-top: 22px;
  padding: 20px 24px;
  border-radius: var(--card-radius);
}

.feature-row {
  display: flex;
  align-items: flex-start;
  gap: 18px;
  padding: 18px 0;
  border-bottom: 1.5px dashed rgba(120, 90, 65, 0.14);
}

.feature-row:last-child {
  border-bottom: 0;
}

.feature-icon {
  width: 56px;
  height: 56px;
  display: grid;
  flex: 0 0 56px;
  place-items: center;
  border-radius: 50%;
  color: var(--feature-color);
  background: var(--feature-bg);
}

.feature-icon svg {
  width: 29px;
  height: 29px;
  stroke-width: 2.25;
}

.feature-copy h3 {
  margin: 0;
  color: var(--text);
  font-size: 21px;
  font-weight: 900;
  line-height: 1.1;
  letter-spacing: 0;
}

.feature-copy p {
  margin: 10px 0 0;
  color: var(--sub);
  font-size: 16px;
  font-weight: 570;
  line-height: 1.58;
}

.list-card {
  margin-top: 22px;
  overflow: hidden;
  border-radius: var(--card-radius);
  background: var(--card-surface);
  box-shadow: var(--card-shadow);
}

.info-row {
  width: 100%;
  min-height: 78px;
  display: flex;
  align-items: center;
  padding: 0 22px;
  border-bottom: 1px solid var(--line);
  color: var(--text);
  background: transparent;
  text-align: left;
  transition: transform 180ms ease, opacity 180ms ease, background 180ms ease;
}

.info-row:last-child {
  border-bottom: 0;
}

.row-icon {
  width: 28px;
  height: 28px;
  flex: 0 0 28px;
  margin-right: 16px;
  color: var(--row-color);
}

.row-icon svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.2;
}

.row-label {
  flex: 0 0 auto;
  color: var(--text);
  font-size: 19px;
  font-weight: 820;
  line-height: 1;
  white-space: nowrap;
}

.row-value {
  min-width: 0;
  margin-left: auto;
  overflow: hidden;
  color: var(--sub);
  font-size: 16px;
  font-weight: 610;
  line-height: 1;
  text-align: right;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.chevron {
  width: 20px;
  height: 20px;
  flex: 0 0 20px;
  margin-left: 10px;
  color: var(--muted);
  stroke-width: 2.3;
}

.slogan {
  display: grid;
  place-items: center;
  gap: 10px;
  margin: 30px 0 2px;
  color: #9a7957;
  text-align: center;
}

.slogan svg {
  width: 29px;
  height: 29px;
  color: var(--sage);
  stroke-width: 2.25;
}

.slogan-line {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 14px;
  width: 100%;
}

.slogan-line::before,
.slogan-line::after {
  content: "";
  width: 70px;
  height: 1px;
  background: rgba(160, 120, 90, 0.22);
}

.slogan-line span {
  color: #9a7957;
  font-size: 18px;
  font-weight: 720;
  line-height: 1;
  white-space: nowrap;
}

.toast {
  position: fixed;
  left: 50%;
  bottom: calc(28px + env(safe-area-inset-bottom));
  z-index: 5;
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
  white-space: nowrap;
}

.toast.show {
  opacity: 1;
  transform: translate(-50%, 0);
}

.back-btn:active,
.info-row:active {
  transform: scale(0.98);
}

.info-row:active {
  opacity: 0.75;
  background: rgba(255, 255, 255, 0.24);
}

@media (hover: hover) {
  .back-btn:hover {
    transform: translateY(-1px);
  }

  .info-row:hover {
    background: rgba(255, 255, 255, 0.22);
  }
}

@media (max-width: 380px) {
  .about-shell {
    padding-right: 18px;
    padding-left: 18px;
  }

  .brand-card {
    gap: 18px;
    padding: 25px 22px;
  }

  .app-icon {
    width: 100px;
    height: 100px;
    flex-basis: 100px;
  }

  .brand-copy h2 {
    font-size: 29px;
  }

  .brand-copy .subtitle {
    font-size: 19px;
  }

  .feature-card {
    padding: 18px 22px;
  }

  .info-row {
    padding: 0 20px;
  }

  .row-label {
    font-size: 18px;
  }

  .row-value {
    max-width: 150px;
  }
}

@media (max-width: 350px) {
  .page-title {
    font-size: 27px;
  }

  .brand-card {
    align-items: flex-start;
    flex-direction: column;
  }

  .app-icon {
    width: 92px;
    height: 92px;
    flex-basis: 92px;
  }

  .feature-row {
    gap: 14px;
  }

  .row-icon {
    margin-right: 12px;
  }

  .row-value {
    max-width: 116px;
    font-size: 15px;
  }

  .slogan-line::before,
  .slogan-line::after {
    width: 44px;
  }
}

@media (min-width: 431px) {
  .about-shell {
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
