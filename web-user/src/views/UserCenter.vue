<template>
  <div class="profile-shell" :style="pageVars">
    <div class="profile-warm-overlay" aria-hidden="true"></div>

    <main class="profile-phone">
      <section class="profile-header" aria-label="个人资料">
        <div class="profile-avatar" aria-label="头像">
          <img v-if="avatarUrl" :src="avatarUrl" alt="用户头像" />
          <span v-else>{{ userInitial }}</span>
        </div>

        <div class="profile-copy">
          <h1>{{ displayName }}</h1>
        </div>

        <button class="edit-btn" type="button" aria-label="编辑个人资料" @click="handleEditProfile">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
            <path d="M18.4 2.6a2.1 2.1 0 0 1 3 3L12 15l-4 1 1-4Z" />
          </svg>
        </button>
      </section>

      <section class="glass-card stats-card" aria-label="个人数据统计">
        <button
          v-for="stat in stats"
          :key="stat.label"
          class="stat-item"
          type="button"
          @click="goPath(stat.path)"
        >
          <span class="stat-value">{{ stat.value }}</span>
          <span class="stat-label">{{ stat.label }}</span>
        </button>
      </section>

      <section class="glass-card menu-group" aria-label="我的内容">
        <button
          v-for="item in contentItems"
          :key="item.label"
          class="menu-item"
          type="button"
          :style="{ '--item-color': item.color, '--item-bg': item.bg }"
          @click="goPath(item.path)"
        >
          <span class="icon-badge" aria-hidden="true">
            <ProfileIcon :name="item.icon" />
          </span>
          <span class="menu-label">{{ item.label }}</span>
          <ChevronIcon />
        </button>
      </section>

      <section class="glass-card menu-group" aria-label="偏好与更多">
        <button
          v-for="item in settingItems"
          :key="item.label"
          class="menu-item"
          type="button"
          :style="{ '--item-color': item.color, '--item-bg': item.bg }"
          @click="handleSettingAction(item)"
        >
          <span class="icon-badge" aria-hidden="true">
            <ProfileIcon :name="item.icon" />
          </span>
          <span class="menu-label">{{ item.label }}</span>
          <ChevronIcon />
        </button>
      </section>

      <button class="logout-btn" type="button" @click="handleLogout">退出登录</button>
    </main>

    <Teleport to="body">
      <div v-if="editVisible" class="edit-mask" @click.self="closeEditDialog">
        <section class="edit-dialog glass-card" aria-label="编辑个人资料">
          <div class="sheet-handle" aria-hidden="true"></div>
          <header class="edit-header">
            <div>
              <h2>编辑资料</h2>
            </div>
            <button class="dialog-close" type="button" aria-label="关闭" @click="closeEditDialog">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round">
                <path d="M18 6 6 18" />
                <path d="m6 6 12 12" />
              </svg>
            </button>
          </header>

          <label class="edit-field">
            <span>昵称</span>
            <input v-model.trim="editForm.nickname" type="text" maxlength="20" placeholder="输入昵称" />
          </label>

          <div class="edit-field">
            <span>头像</span>
            <label class="avatar-uploader" :class="{ uploading: uploadingAvatar }">
              <input type="file" accept="image/png,image/jpeg,image/webp" @change="handleAvatarChange" />
              <span class="avatar-preview">
                <img v-if="editForm.avatar" :src="editForm.avatar" alt="头像预览" />
                <span v-else>{{ userInitial }}</span>
              </span>
              <span class="avatar-upload-copy">
                <strong>{{ uploadingAvatar ? '上传中...' : '选择头像图片' }}</strong>
              </span>
            </label>
          </div>

          <div class="edit-field">
            <span>性别</span>
            <div class="gender-group">
              <button
                v-for="option in genderOptions"
                :key="option.value"
                class="gender-chip"
                :class="{ active: editForm.gender === option.value }"
                type="button"
                @click="editForm.gender = option.value"
              >
                {{ option.label }}
              </button>
            </div>
          </div>

          <div class="dialog-actions">
            <button class="secondary-btn" type="button" @click="closeEditDialog">取消</button>
            <button class="primary-btn" type="button" :disabled="savingProfile" @click="submitProfile">
              {{ savingProfile ? '保存中' : '保存修改' }}
            </button>
          </div>
        </section>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getUserStats, updateProfile, uploadAvatar } from '@/api/user'
import kitchenBg from '@/assets/home/kitchen-bg.jpg'

type IconName = 'heart' | 'file' | 'cart' | 'bowl' | 'clock' | 'info'

interface MenuItem {
  label: string
  icon: IconName
  color: string
  bg: string
  path?: string
  action?: 'history' | 'about'
}

interface ProfileStats {
  favorite_count: number
  menu_count: number
  shopping_list_count: number
}

const router = useRouter()
const userStore = useUserStore()

const profileStats = ref<ProfileStats | null>(null)
const editVisible = ref(false)
const savingProfile = ref(false)
const uploadingAvatar = ref(false)
const editForm = reactive({
  nickname: '',
  avatar: '',
  gender: 0,
})

const genderOptions = [
  { label: '保密', value: 0 },
  { label: '男', value: 1 },
  { label: '女', value: 2 },
]

const pageVars = computed(() => ({
  '--profile-bg': `url(${kitchenBg})`,
}))

const displayName = computed(() => {
  const user = userStore.userInfo
  return user?.email || user?.nickname || user?.username || '用户'
})

const avatarUrl = computed(() => userStore.userInfo?.avatar || '')

const userInitial = computed(() => {
  const name = displayName.value.trim()
  if (!name) return 'U'
  return name.charAt(0).toUpperCase()
})

function statValue(key: keyof ProfileStats) {
  return profileStats.value ? profileStats.value[key] : '--'
}

const stats = computed(() => [
  { label: '我的收藏', value: statValue('favorite_count'), path: '/user/favorites' },
  { label: '定制菜单', value: statValue('menu_count'), path: '/user/menus' },
  { label: '购物清单', value: statValue('shopping_list_count'), path: '/shopping-list' },
])

const contentItems: MenuItem[] = [
  {
    label: '我的收藏',
    icon: 'heart',
    color: '#e95645',
    bg: 'rgba(233, 86, 69, 0.14)',
    path: '/user/favorites',
  },
  {
    label: '我的菜单',
    icon: 'file',
    color: '#5f8df7',
    bg: 'rgba(95, 141, 247, 0.15)',
    path: '/user/menus',
  },
  {
    label: '我的购物清单',
    icon: 'cart',
    color: '#f28a2e',
    bg: 'rgba(242, 138, 46, 0.16)',
    path: '/shopping-list',
  },
]

const settingItems: MenuItem[] = [
  {
    label: '饮食偏好设置',
    icon: 'bowl',
    color: '#6f8b65',
    bg: 'rgba(143, 167, 131, 0.22)',
    path: '/user/preferences',
  },
  {
    label: '浏览历史',
    icon: 'clock',
    color: '#9a7957',
    bg: 'rgba(210, 170, 120, 0.20)',
    path: '/user/history',
  },
  {
    label: '关于我们',
    icon: 'info',
    color: '#8b715e',
    bg: 'rgba(160, 130, 100, 0.18)',
    path: '/about',
  },
]

function handleEditProfile() {
  const user = userStore.userInfo
  editForm.nickname = user?.nickname || ''
  editForm.avatar = user?.avatar || ''
  editForm.gender = Number(user?.gender || 0)
  editVisible.value = true
}

function closeEditDialog() {
  if (savingProfile.value || uploadingAvatar.value) return
  editVisible.value = false
}

async function handleAvatarChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (!file.type.startsWith('image/')) return

  uploadingAvatar.value = true
  try {
    const res = await uploadAvatar(file)
    editForm.avatar = res.url
  } finally {
    uploadingAvatar.value = false
  }
}

async function submitProfile() {
  if (savingProfile.value) return
  savingProfile.value = true
  try {
    await updateProfile({
      nickname: editForm.nickname,
      avatar: editForm.avatar,
      gender: editForm.gender,
    })
    await userStore.fetchUserInfo()
    editVisible.value = false
  } finally {
    savingProfile.value = false
  }
}

function goPath(path?: string) {
  if (path) router.push(path)
}

function handleSettingAction(item: MenuItem) {
  if (item.path) {
    router.push(item.path)
    return
  }
  console.log('profile action:', item.action || item.label)
}

function handleLogout() {
  if (!window.confirm('确认退出登录吗？')) return
  userStore.logout()
  router.replace('/login')
}

onMounted(async () => {
  if (!userStore.userInfo) {
    await userStore.fetchUserInfo()
  }

  try {
    const res: any = await getUserStats()
    profileStats.value = {
      favorite_count: Number(res?.favorite_count ?? 0),
      menu_count: Number(res?.menu_count ?? 0),
      shopping_list_count: Number(res?.shopping_list_count ?? 0),
    }
  } catch {
    profileStats.value = null
  }
})

const ChevronIcon = defineComponent({
  name: 'ChevronIcon',
  setup() {
    return () => h('svg', {
      class: 'chevron',
      viewBox: '0 0 24 24',
      fill: 'none',
      stroke: 'currentColor',
      'stroke-linecap': 'round',
      'stroke-linejoin': 'round',
      'aria-hidden': 'true',
    }, [
      h('path', { d: 'm9 18 6-6-6-6' }),
    ])
  },
})

const ProfileIcon = defineComponent({
  name: 'ProfileIcon',
  props: {
    name: { type: String, required: true },
  },
  setup(props) {
    return () => {
      const common = {
        viewBox: '0 0 24 24',
        fill: 'none',
        stroke: 'currentColor',
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        'aria-hidden': 'true',
      }

      if (props.name === 'heart') {
        return h('svg', common, [
          h('path', { d: 'M19.5 12.6 12 20l-7.5-7.4A5 5 0 1 1 12 6.1a5 5 0 1 1 7.5 6.5Z' }),
        ])
      }

      if (props.name === 'file') {
        return h('svg', common, [
          h('path', { d: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8Z' }),
          h('path', { d: 'M14 2v6h6' }),
          h('path', { d: 'M8 13h8M8 17h6M8 9h2' }),
        ])
      }

      if (props.name === 'cart') {
        return h('svg', common, [
          h('circle', { cx: '8', cy: '21', r: '1' }),
          h('circle', { cx: '19', cy: '21', r: '1' }),
          h('path', { d: 'M2.1 2.1h3l2.2 12.4a2 2 0 0 0 2 1.6h7.9a2 2 0 0 0 1.9-1.5L21 7H6.2' }),
        ])
      }

      if (props.name === 'bowl') {
        return h('svg', common, [
          h('path', { d: 'M4 14h16' }),
          h('path', { d: 'M6 14a6 6 0 0 1 12 0' }),
          h('path', { d: 'M12 5v2' }),
          h('path', { d: 'M8.5 7 10 8.5M15.5 7 14 8.5' }),
          h('path', { d: 'M6 17h12' }),
          h('path', { d: 'M8 20h8' }),
        ])
      }

      if (props.name === 'clock') {
        return h('svg', common, [
          h('circle', { cx: '12', cy: '12', r: '9' }),
          h('path', { d: 'M12 7v5l3 2' }),
        ])
      }

      return h('svg', common, [
        h('circle', { cx: '12', cy: '12', r: '9' }),
        h('path', { d: 'M12 11v5' }),
        h('path', { d: 'M12 8h.01' }),
      ])
    }
  },
})
</script>

<style scoped>
.profile-shell {
  --text: #2e241f;
  --sub: #7a6a5f;
  --cream: rgba(255, 250, 240, 0.78);
  --coral: #e95645;
  --line: rgba(120, 100, 80, 0.14);
  --border: rgba(255, 255, 255, 0.62);
  --shadow: 0 18px 40px rgba(80, 50, 30, 0.15);
  position: relative;
  min-height: 100vh;
  min-height: 100dvh;
  overflow-x: clip;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 238, 211, 0.4), rgba(255, 246, 233, 0.52) 44%, rgba(234, 192, 145, 0.22)),
    var(--profile-bg) center top / cover fixed;
}

.profile-warm-overlay {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 15% 11%, rgba(255, 255, 255, 0.72), transparent 31%),
    radial-gradient(circle at 79% 3%, rgba(255, 235, 204, 0.34), transparent 30%),
    linear-gradient(180deg, rgba(255, 240, 218, 0.56), rgba(255, 250, 242, 0.28) 50%, rgba(139, 83, 39, 0.1));
  backdrop-filter: blur(5px);
  -webkit-backdrop-filter: blur(5px);
}

.profile-phone {
  position: relative;
  z-index: 1;
  width: min(100%, 430px);
  min-height: 100vh;
  margin: 0 auto;
  padding: 70px 24px 126px;
}

.profile-header {
  display: grid;
  grid-template-columns: 78px minmax(0, 1fr) 54px;
  align-items: center;
  gap: 20px;
}

.profile-avatar {
  width: 78px;
  height: 78px;
  display: grid;
  place-items: center;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.76);
  border-radius: 50%;
  color: #2e241f;
  background:
    radial-gradient(circle at 35% 26%, rgba(255, 255, 255, 0.96), rgba(255, 248, 237, 0.86)),
    rgba(255, 250, 240, 0.88);
  box-shadow:
    0 14px 30px rgba(80, 50, 30, 0.13),
    inset 0 1px 0 rgba(255, 255, 255, 0.92);
}

.profile-avatar img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.profile-avatar span {
  font-size: 31px;
  font-weight: 900;
  line-height: 1;
}

.profile-copy {
  min-width: 0;
}

.profile-copy h1 {
  margin: 0;
  overflow: hidden;
  color: var(--text);
  font-size: clamp(26px, 7vw, 30px);
  font-weight: 950;
  line-height: 1.12;
  letter-spacing: 0;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.edit-btn {
  width: 54px;
  height: 54px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(255, 255, 255, 0.68);
  border-radius: 18px;
  color: #8b7d70;
  background: rgba(255, 250, 240, 0.84);
  box-shadow:
    0 12px 26px rgba(80, 50, 30, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(18px);
  -webkit-backdrop-filter: blur(18px);
  cursor: pointer;
  transition: transform 180ms ease, box-shadow 180ms ease, color 180ms ease;
}

.edit-btn svg {
  width: 29px;
  height: 29px;
  stroke-width: 2.25;
}

.glass-card {
  border: 1px solid var(--border);
  border-radius: 30px;
  background: var(--cream);
  box-shadow: var(--shadow);
  backdrop-filter: blur(20px) saturate(1.08);
  -webkit-backdrop-filter: blur(20px) saturate(1.08);
}

.edit-mask {
  position: fixed;
  inset: 0 0 calc(106px + var(--safe-bottom, 0px));
  z-index: 90;
  display: grid;
  align-items: center;
  padding: 24px;
  background:
    linear-gradient(180deg, rgba(46, 36, 31, 0.12), rgba(46, 36, 31, 0.34)),
    rgba(46, 36, 31, 0.2);
  backdrop-filter: blur(12px) saturate(1.04);
  -webkit-backdrop-filter: blur(12px) saturate(1.04);
}

.edit-dialog {
  width: min(100%, 430px);
  max-height: min(74dvh, 620px);
  display: flex;
  flex-direction: column;
  justify-self: center;
  overflow-y: auto;
  padding: 10px 22px 24px;
  border-radius: 30px;
  background:
    radial-gradient(circle at 15% 0%, rgba(255, 255, 255, 0.54), transparent 42%),
    rgba(255, 250, 240, 0.94);
  box-shadow:
    0 24px 70px rgba(80, 50, 30, 0.24),
    inset 0 1px 0 rgba(255, 255, 255, 0.88);
  overscroll-behavior: contain;
}

.sheet-handle {
  width: 48px;
  height: 5px;
  flex: 0 0 auto;
  margin: 2px auto 18px;
  border-radius: 999px;
  background: rgba(122, 106, 95, 0.22);
}

.edit-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.edit-header h2 {
  margin: 0;
  color: var(--text);
  font-size: 24px;
  font-weight: 900;
  line-height: 1.1;
}

.dialog-close {
  width: 40px;
  height: 40px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 14px;
  color: #8b7d70;
  background: rgba(255, 255, 255, 0.68);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.9);
  cursor: pointer;
}

.dialog-close svg {
  width: 20px;
  height: 20px;
  stroke-width: 2.2;
}

.edit-field {
  display: grid;
  gap: 10px;
  margin-top: 18px;
}

.edit-field span {
  color: #5d4b40;
  font-size: 15px;
  font-weight: 760;
  line-height: 1;
}

.edit-field input {
  width: 100%;
  min-height: 54px;
  padding: 0 16px;
  border: 1px solid rgba(120, 100, 80, 0.14);
  border-radius: 16px;
  outline: 0;
  color: var(--text);
  background: rgba(255, 255, 255, 0.64);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.82);
  font: inherit;
  font-size: 16px;
  font-weight: 620;
}

.edit-field input::placeholder {
  color: #ab9d92;
}

.avatar-uploader {
  min-height: 92px;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 16px;
  border: 1px dashed rgba(160, 120, 90, 0.28);
  border-radius: 20px;
  color: var(--text);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.56), rgba(255, 250, 242, 0.42));
  box-shadow:
    0 10px 24px rgba(80, 50, 30, 0.06),
    inset 0 1px 0 rgba(255, 255, 255, 0.82);
  cursor: pointer;
  transition: transform 180ms ease, border-color 180ms ease, background 180ms ease, opacity 180ms ease;
}

.avatar-uploader input {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  opacity: 0;
  pointer-events: none;
}

.avatar-uploader.uploading {
  opacity: 0.72;
  pointer-events: none;
}

.avatar-preview {
  width: 64px;
  height: 64px;
  display: grid;
  flex: 0 0 64px;
  place-items: center;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.78);
  border-radius: 22px;
  color: #2e241f;
  background:
    radial-gradient(circle at 35% 26%, rgba(255, 255, 255, 0.96), rgba(255, 248, 237, 0.86)),
    rgba(255, 250, 240, 0.88);
  box-shadow: 0 10px 20px rgba(80, 50, 30, 0.1);
  font-size: 25px;
  font-weight: 900;
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.avatar-upload-copy {
  min-width: 0;
  display: grid;
  gap: 7px;
}

.avatar-upload-copy strong {
  color: #2e241f;
  font-size: 17px;
  font-weight: 850;
  line-height: 1;
}

.gender-group {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.gender-chip {
  min-height: 48px;
  border: 1px solid rgba(120, 100, 80, 0.14);
  border-radius: 16px;
  color: #6f6054;
  background: rgba(255, 255, 255, 0.58);
  font-size: 15px;
  font-weight: 760;
  cursor: pointer;
  transition: transform 180ms ease, background 180ms ease, color 180ms ease;
}

.gender-chip.active {
  border-color: transparent;
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 12px 24px rgba(233, 86, 69, 0.18);
}

.dialog-actions {
  position: sticky;
  bottom: 0;
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin: 24px -2px 0;
  padding: 16px 0 0;
  border-top: 1px solid rgba(120, 100, 80, 0.1);
  background: transparent;
}

.secondary-btn,
.primary-btn {
  min-height: 52px;
  border: 0;
  border-radius: 16px;
  font-size: 16px;
  font-weight: 820;
  cursor: pointer;
  transition: transform 180ms ease, opacity 180ms ease;
}

.secondary-btn {
  color: #7a6a5f;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.74), rgba(255, 251, 245, 0.62));
  box-shadow:
    0 8px 18px rgba(80, 50, 30, 0.06),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
}

.primary-btn {
  color: #fff;
  background: linear-gradient(135deg, #f06152, #e9473a);
  box-shadow: 0 14px 26px rgba(233, 86, 69, 0.24);
}

.primary-btn:disabled {
  opacity: 0.62;
  cursor: not-allowed;
}

.stats-card {
  height: 96px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  align-items: center;
  margin-top: 28px;
  overflow: hidden;
}

.stat-item {
  position: relative;
  height: 100%;
  display: grid;
  place-items: center;
  align-content: center;
  gap: 9px;
  border: 0;
  color: var(--sub);
  background: transparent;
  cursor: pointer;
  transition: transform 180ms ease, background 180ms ease;
}

.stat-item + .stat-item::before {
  content: "";
  position: absolute;
  left: 0;
  top: 24px;
  width: 1px;
  height: 48px;
  background: rgba(120, 100, 80, 0.16);
}

.stat-value {
  display: block;
  color: var(--text);
  font-size: 30px;
  font-weight: 950;
  line-height: 1;
}

.stat-label {
  display: block;
  color: #76665a;
  font-size: 14px;
  font-weight: 760;
  line-height: 1;
}

.menu-group {
  overflow: hidden;
  margin-top: 27px;
}

.menu-group + .menu-group {
  margin-top: 26px;
}

.menu-item {
  position: relative;
  width: 100%;
  height: 74px;
  display: flex;
  align-items: center;
  padding: 0 23px;
  border: 0;
  color: var(--text);
  background: transparent;
  text-align: left;
  cursor: pointer;
  transition: transform 180ms ease, background 180ms ease, opacity 180ms ease;
}

.menu-item + .menu-item::before {
  content: "";
  position: absolute;
  left: 80px;
  right: 23px;
  top: 0;
  height: 1px;
  background: var(--line);
}

.icon-badge {
  width: 50px;
  height: 50px;
  flex: 0 0 auto;
  display: grid;
  place-items: center;
  border-radius: 50%;
  color: var(--item-color);
  background: var(--item-bg);
}

.icon-badge svg {
  width: 28px;
  height: 28px;
  stroke-width: 2.25;
}

.menu-label {
  margin-left: 18px;
  color: #2e241f;
  font-size: 21px;
  font-weight: 860;
  letter-spacing: 0;
}

.chevron {
  width: 24px;
  height: 24px;
  margin-left: auto;
  color: #8b7d70;
  stroke-width: 2.4;
}

.logout-btn {
  width: 100%;
  min-height: 42px;
  margin-top: 40px;
  border: 0;
  color: var(--coral);
  background: transparent;
  font-size: 21px;
  font-weight: 850;
  letter-spacing: 0;
  cursor: pointer;
  transition: transform 180ms ease, opacity 180ms ease;
}

.edit-btn:active,
.dialog-close:active,
.stat-item:active,
.menu-item:active,
.logout-btn:active,
.gender-chip:active,
.secondary-btn:active,
.primary-btn:active {
  transform: scale(0.98);
}

.menu-item:active,
.stat-item:active {
  background: rgba(255, 255, 255, 0.28);
}

@media (hover: hover) {
  .edit-btn:hover {
    color: var(--coral);
    transform: translateY(-1px);
    box-shadow:
      0 15px 30px rgba(80, 50, 30, 0.14),
      inset 0 1px 0 rgba(255, 255, 255, 0.94);
  }

  .dialog-close:hover,
  .avatar-uploader:hover,
  .gender-chip:hover,
  .secondary-btn:hover,
  .primary-btn:hover:not(:disabled) {
    transform: translateY(-1px);
  }

  .menu-item:hover,
  .stat-item:hover {
    background: rgba(255, 255, 255, 0.22);
  }

  .logout-btn:hover {
    opacity: 0.82;
  }
}

@media (max-width: 380px) {
  .edit-mask {
    inset: 0 0 calc(98px + var(--safe-bottom, 0px));
    padding: 16px;
  }

  .edit-dialog {
    max-height: 72dvh;
    padding-right: 18px;
    padding-left: 18px;
  }

  .edit-header h2 {
    font-size: 22px;
  }

  .dialog-actions {
    grid-template-columns: 1fr;
  }

  .profile-phone {
    padding-left: 20px;
    padding-right: 20px;
  }

  .profile-header {
    grid-template-columns: 72px minmax(0, 1fr) 50px;
    gap: 16px;
  }

  .profile-avatar {
    width: 72px;
    height: 72px;
  }

  .profile-copy h1 {
    font-size: 25px;
  }

  .edit-btn {
    width: 50px;
    height: 50px;
  }

  .menu-item {
    padding: 0 20px;
  }

  .menu-item + .menu-item::before {
    left: 76px;
    right: 20px;
  }

  .menu-label {
    font-size: 20px;
  }
}

@media (min-width: 431px) {
  .profile-phone {
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
