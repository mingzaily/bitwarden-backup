<template>
  <div v-if="booting" class="min-h-screen grid place-items-center bg-ink text-muted">
    <div class="flex items-center gap-3 text-sm"><span class="status-pulse"></span>正在连接保险库…</div>
  </div>

  <div v-else-if="!authenticated" class="min-h-screen bg-ink px-6 py-10 text-main">
    <button class="icon-button theme-toggle fixed right-6 top-6 z-10" type="button" :aria-label="theme === 'dark' ? '切换到浅色主题' : '切换到深色主题'" :title="theme === 'dark' ? '切换到浅色主题' : '切换到深色主题'" :aria-pressed="theme === 'light'" @click="toggleTheme">
      <svg v-if="theme === 'dark'" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 3v1.5M12 19.5V21M4.64 4.64l1.06 1.06m12.6 12.6 1.06 1.06M3 12h1.5m16 0H21M4.64 19.36l1.06-1.06m12.6-12.6 1.06-1.06M16.5 12a4.5 4.5 0 1 1-9 0 4.5 4.5 0 0 1 9 0Z" /></svg>
      <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M20.3 15.5A8.5 8.5 0 0 1 8.5 3.7 8.5 8.5 0 1 0 20.3 15.5Z" /></svg>
    </button>
    <div class="mx-auto flex min-h-[calc(100vh-5rem)] max-w-5xl items-center justify-center">
      <div class="grid w-full overflow-hidden rounded-3xl border border-line bg-panel shadow-2xl lg:grid-cols-[1.05fr_0.95fr]">
        <section class="login-promo relative hidden overflow-hidden p-12 lg:block">
          <div class="absolute -right-24 -top-24 h-72 w-72 rounded-full bg-accent/15 blur-3xl"></div>
          <div class="absolute -bottom-32 -left-16 h-80 w-80 rounded-full bg-violet/15 blur-3xl"></div>
          <div class="relative flex h-full flex-col justify-between">
            <div>
              <div class="brand-lockup mb-10">
                <div class="brand-mark"><span class="brand-mark-core"></span></div>
                <div class="brand-copy">
                  <div class="brand-heading">
                    <p class="brand-name text-main">VAULT//SYNC</p>
                    <span class="brand-version">{{ appVersion }}</span>
                  </div>
                </div>
              </div>
              <p class="eyebrow">SECURE BACKUP CONTROL</p>
              <h1 class="mt-4 max-w-md text-4xl font-semibold leading-tight tracking-tight text-main">让每一次备份，都有迹可循。</h1>
              <p class="mt-5 max-w-sm text-sm leading-7 text-muted">集中管理 Bitwarden 源站、存储目标与备份任务。凭证加密保存，备份过程清晰可见。</p>
            </div>
            <div class="mt-16 flex items-center gap-3 text-xs text-subtle"><span class="h-2 w-2 rounded-full bg-accent shadow-[0_0_14px_rgba(30,215,96,.8)]"></span>LOCAL CONTROL PLANE · READY</div>
          </div>
        </section>

        <section class="flex items-center p-8 sm:p-12">
          <form class="w-full max-w-sm" @submit.prevent="handleLogin">
            <div class="mb-10 lg:hidden">
              <div class="brand-lockup mb-7"><div class="brand-mark"><span class="brand-mark-core"></span></div><div class="brand-copy"><div class="brand-heading"><p class="brand-name">VAULT//SYNC</p><span class="brand-version">{{ appVersion }}</span></div></div></div>
              <p class="eyebrow">SECURE BACKUP CONTROL</p>
            </div>
            <p class="eyebrow">WELCOME BACK</p>
            <h2 class="mt-3 text-3xl font-semibold tracking-tight text-main">进入控制台</h2>
            <p class="mt-3 text-sm leading-6 text-muted">输入管理员密码，继续管理你的备份工作区。</p>
            <label class="mt-8 block text-sm font-medium text-muted" for="admin-password">管理员密码</label>
            <div class="relative mt-2">
              <input id="admin-password" v-model="password" class="input pr-12" :type="showPassword ? 'text' : 'password'" autocomplete="current-password" autofocus placeholder="输入密码" />
              <button
                type="button"
                class="absolute inset-y-1 right-1 grid w-10 place-items-center rounded-lg text-subtle transition hover:bg-surface-hover hover:text-main focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent/50"
                :aria-label="showPassword ? '隐藏密码' : '显示密码'"
                :aria-pressed="showPassword"
                :title="showPassword ? '隐藏密码' : '显示密码'"
                @click="showPassword = !showPassword"
              >
                <svg v-if="showPassword" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.6" d="M3.5 12s3.1-5 8.5-5c1.35 0 2.57.3 3.63.73M20.5 12s-3.1 5-8.5 5c-1.35 0-2.57-.3-3.63-.73M9.88 9.88a3 3 0 1 0 4.24 4.24M4 4l16 16" /></svg>
                <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.6" d="M12 15.5a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7Z M3.5 12s3.1-5 8.5-5 8.5 5 8.5 5-3.1 5-8.5 5-8.5-5-8.5-5Z" /></svg>
              </button>
            </div>
            <p v-if="loginError" class="mt-3 rounded-xl border border-danger/30 bg-danger/10 px-3 py-2 text-sm text-danger">{{ loginError }}</p>
            <button class="btn-primary mt-6 w-full justify-center" :disabled="loggingIn || !password">
              <span v-if="loggingIn" class="spinner"></span>
              {{ loggingIn ? '验证中…' : '进入控制台' }}
            </button>
            <p class="mt-6 text-center text-xs text-subtle">会话有效期 12 小时 · 重启服务后自动失效</p>
          </form>
        </section>
      </div>
    </div>
  </div>

  <div v-else class="min-h-screen bg-ink text-main lg:flex">
    <aside class="hidden w-64 shrink-0 border-r border-line bg-panel px-4 py-6 lg:flex lg:flex-col">
      <div class="brand-lockup px-3"><div class="brand-mark"><span class="brand-mark-core"></span></div><div class="brand-copy"><div class="brand-heading"><p class="brand-name text-main">VAULT//SYNC</p><span class="brand-version">{{ appVersion }}</span></div><p class="brand-subtitle">backup control plane</p></div></div>
      <nav class="mt-12 space-y-1" aria-label="主导航">
        <template v-for="group in navGroups" :key="group.label || group.items[0].path">
          <div :class="['nav-group', group.label ? 'is-grouped' : '']">
            <div v-if="group.label" class="nav-group-heading">
              <span class="nav-group-label">{{ group.label }}</span>
              <span class="nav-group-count">{{ group.items.length }} 项</span>
            </div>
            <div :class="group.label ? 'nav-group-items' : ''">
              <router-link v-for="item in group.items" :key="item.path" :to="item.path" :class="['nav-link', group.label ? 'nav-child-link' : '']">
            <span class="nav-icon" aria-hidden="true"><svg v-if="item.icon === 'overview'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M4 4h6v6H4zM14 4h6v6h-6zM4 14h6v6H4zM14 14h6v6h-6z" /></svg><svg v-else-if="item.icon === 'server'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M5 5.5h14v5H5zM5 13.5h14v5H5zM8 8h.01M8 16h.01" /></svg><svg v-else-if="item.icon === 'target'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 3.5 19 7v5c0 4.1-2.6 7.4-7 8.5-4.4-1.1-7-4.4-7-8.5V7l7-3.5Z M9 12h6M12 9v6" /></svg><svg v-else-if="item.icon === 'task'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M7 4h10a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2ZM8 9h8M8 13h5M8 17h3" /></svg><svg v-else fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M5 5v14h14M8 16l3-4 3 2 4-5" /></svg></span>
            <span>{{ item.label }}</span><span v-if="isNavItemActive(item)" class="nav-active-dot"></span>
          </router-link>
            </div>
          </div>
        </template>
      </nav>
      <div class="mt-auto rounded-2xl border border-line bg-panel-soft/70 p-4"><p class="eyebrow">SYSTEM STATUS</p><div class="mt-3 flex items-center gap-2 text-sm text-muted"><span class="status-pulse"></span>运行正常</div><p class="mt-2 text-xs leading-5 text-subtle">凭证已加密 · API 会话安全</p></div>
    </aside>

    <main class="min-w-0 flex-1">
      <header class="sticky top-0 z-20 border-b border-line/80 bg-ink/90 px-5 py-3 backdrop-blur-xl sm:px-8 lg:px-10">
        <div class="mx-auto flex max-w-[1440px] items-center justify-between gap-4">
          <div class="brand-lockup lg:hidden"><div class="brand-mark"><span class="brand-mark-core"></span></div><div class="brand-copy"><div class="brand-heading"><p class="brand-name text-main">VAULT//SYNC</p><span class="brand-version">{{ appVersion }}</span></div></div></div>
          <div class="ml-auto flex items-center gap-2">
            <button class="icon-button theme-toggle" type="button" :aria-label="theme === 'dark' ? '切换到浅色主题' : '切换到深色主题'" :title="theme === 'dark' ? '切换到浅色主题' : '切换到深色主题'" :aria-pressed="theme === 'light'" @click="toggleTheme">
              <svg v-if="theme === 'dark'" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 3v1.5M12 19.5V21M4.64 4.64l1.06 1.06m12.6 12.6 1.06 1.06M3 12h1.5m16 0H21M4.64 19.36l1.06-1.06m12.6-12.6 1.06-1.06M16.5 12a4.5 4.5 0 1 1-9 0 4.5 4.5 0 0 1 9 0Z" /></svg>
              <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M20.3 15.5A8.5 8.5 0 0 1 8.5 3.7 8.5 8.5 0 1 0 20.3 15.5Z" /></svg>
            </button>
            <button class="btn-ghost" type="button" title="退出登录" @click="handleLogout"><svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M15 8V6a2 2 0 0 0-2-2H6a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h7a2 2 0 0 0 2-2v-2M10 12h10m0 0-3-3m3 3-3 3" /></svg><span class="hidden sm:inline">退出</span></button>
          </div>
        </div>
        <nav class="mx-auto mt-3 flex max-w-[1440px] gap-1 overflow-x-auto lg:hidden" aria-label="移动端导航"><router-link v-for="item in navItems" :key="item.path" :to="item.path" class="mobile-nav-link">{{ item.label }}</router-link></nav>
      </header>

      <div class="mx-auto max-w-[1440px] p-5 sm:p-8 lg:p-10"><router-view v-slot="{ Component }"><transition name="fade" mode="out-in"><component :is="Component" /></transition></router-view></div>
    </main>
  </div>

  <ToastContainer ref="toastRef" />
  <ConfirmModal v-if="authenticated" :visible="confirmState.visible" :title="confirmState.title" :message="confirmState.message" :type="confirmState.type" :confirm-text="confirmState.confirmText" :cancel-text="confirmState.cancelText" @confirm="handleConfirm" @cancel="handleCancel" />
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import ToastContainer from '@/components/ui/ToastContainer.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import { authApi, metaApi } from '@/api'
import { setToastInstance, useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'

const route = useRoute()
const toast = useToast()
const { state: confirmState, handleConfirm, handleCancel } = useConfirm()
const toastRef = ref(null)
const booting = ref(true)
const authenticated = ref(false)
const appVersion = ref('DEV')
const password = ref('')
const showPassword = ref(false)
const loggingIn = ref(false)
const loginError = ref('')
const theme = ref(document.documentElement.dataset.theme === 'light' ? 'light' : 'dark')

const navGroups = [
  { label: '', items: [{ path: '/overview', label: '总览', kicker: 'OVERVIEW', icon: 'overview' }] },
  { label: '', items: [{ path: '/tasks', label: '备份任务', kicker: 'BACKUP TASKS', icon: 'task' }] },
  {
    label: '备份资源',
    items: [
      { path: '/servers', label: 'Bitwarden 源站', kicker: 'BITWARDEN SOURCES', icon: 'server' },
      { path: '/destinations', label: '存储目标', kicker: 'STORAGE TARGETS', icon: 'target' }
    ]
  },
  { label: '', items: [{ path: '/logs', label: '运行记录', kicker: 'ACTIVITY', icon: 'log' }] }
]
const navItems = navGroups.flatMap(group => group.items)
const isNavItemActive = (item) => route.path === item.path || route.path.startsWith(`${item.path}/`)

const applyTheme = (nextTheme) => {
  theme.value = nextTheme === 'light' ? 'light' : 'dark'
  document.documentElement.dataset.theme = theme.value
  window.localStorage.setItem('vaultsync-theme', theme.value)
}

const toggleTheme = () => applyTheme(theme.value === 'dark' ? 'light' : 'dark')

const handleLogin = async () => {
  if (!password.value || loggingIn.value) return
  loggingIn.value = true
  loginError.value = ''
  try {
    await authApi.login(password.value)
    authenticated.value = true
    password.value = ''
    showPassword.value = false
  } catch (error) {
    loginError.value = error.message || '登录失败，请稍后重试'
  } finally {
    loggingIn.value = false
  }
}

const handleLogout = async () => {
  try {
    await authApi.logout()
  } catch {
    // The local UI should still be logged out when the server session expired.
  }
  authenticated.value = false
}

const handleAuthRequired = () => {
  if (authenticated.value) {
    authenticated.value = false
    toast.info('会话已失效，请重新登录')
  }
}

const loadAppMeta = async () => {
  try {
    const data = await metaApi.get()
    if (typeof data?.version === 'string' && data.version.trim()) {
      appVersion.value = data.version.trim()
    }
  } catch {
    // Keep DEV for local static builds or older backends without /api/meta.
  }
}

const restoreSession = async () => {
  try {
    await authApi.session()
    authenticated.value = true
  } catch {
    authenticated.value = false
  }
}

onMounted(async () => {
  setToastInstance(toastRef.value)
  window.addEventListener('auth:required', handleAuthRequired)
  applyTheme(theme.value)
  await Promise.all([loadAppMeta(), restoreSession()])
  booting.value = false
})

onBeforeUnmount(() => window.removeEventListener('auth:required', handleAuthRequired))
</script>
