<template>
  <div class="space-y-6">
    <div class="page-header">
      <div>
        <p class="eyebrow">OVERVIEW</p>
        <h2 class="page-title">工作区总览</h2>
        <p class="page-subtitle">从一个页面掌握备份资源、备份任务和最近运行状态。</p>
      </div>
      <router-link class="btn-primary" to="/tasks">
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M12 5v14m7-7H5" /></svg>
        创建备份任务
      </router-link>
    </div>

    <div v-if="loading" class="loading-state"><div><div class="spinner mx-auto text-accent"></div><p class="mt-3 text-sm">正在读取工作区…</p></div></div>
    <template v-else>
      <section class="overview-metric-grid" aria-label="资源概览">
        <article v-for="metric in metrics" :key="metric.label" class="metric-card">
          <div :class="['metric-icon', `metric-icon-${metric.tone}`]" aria-hidden="true">
            <svg v-if="metric.icon === 'server'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M5 5.5h14v5H5zM5 13.5h14v5H5zM8 8h.01M8 16h.01" /></svg>
            <svg v-else-if="metric.icon === 'target'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 3.5 19 7v5c0 4.1-2.6 7.4-7 8.5-4.4-1.1-7-4.4-7-8.5V7l7-3.5Z M9 12h6M12 9v6" /></svg>
            <svg v-else-if="metric.icon === 'task'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M7 4h10a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2ZM8 9h8M8 13h5M8 17h3" /></svg>
            <svg v-else fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 7v5l3 2m6-2a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" /></svg>
          </div>
          <div class="metric-copy">
            <p class="metric-label">{{ metric.label }}</p>
            <p class="metric-value">{{ metric.value }}</p>
            <p class="metric-detail">{{ metric.detail }}</p>
          </div>
        </article>
      </section>

      <div class="overview-grid">
        <section class="surface overview-panel">
          <div class="overview-panel-header">
            <div>
              <p class="eyebrow">BACKUP TASKS</p>
              <h3 class="overview-panel-title">最近备份任务</h3>
            </div>
            <router-link class="overview-panel-link" to="/tasks">查看全部</router-link>
          </div>
          <div v-if="overview.recent_tasks.length" class="overview-list">
            <router-link v-for="task in overview.recent_tasks" :key="task.id" class="overview-list-item" to="/tasks">
              <span :class="['overview-list-icon', !task.enabled ? 'is-muted' : '']" aria-hidden="true"><svg fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M7 4h10a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2ZM8 9h8M8 13h5M8 17h3" /></svg></span>
              <span class="overview-list-copy"><span class="overview-list-title">{{ task.name }}</span><span class="overview-list-meta">{{ task.source_server_name || '未关联源站' }} · {{ task.destination_count }} 个目标 · {{ task.cron_expression || '手动触发' }}</span></span>
              <span :class="['status-badge', task.enabled ? 'status-success' : 'status-neutral']">{{ task.enabled ? '已启用' : '已停用' }}</span>
            </router-link>
          </div>
          <div v-else class="overview-empty"><p>还没有备份任务</p><router-link to="/tasks">创建第一个任务</router-link></div>
        </section>

        <section class="surface overview-panel">
          <div class="overview-panel-header">
            <div>
              <p class="eyebrow">ACTIVITY</p>
              <h3 class="overview-panel-title">最近运行记录</h3>
            </div>
            <router-link class="overview-panel-link" to="/logs">查看全部</router-link>
          </div>
          <div v-if="overview.recent_logs.length" class="overview-list">
            <router-link v-for="log in overview.recent_logs" :key="log.id" class="overview-list-item" to="/logs">
              <span :class="['overview-list-icon', log.status === 'failed' ? 'is-danger' : log.status === 'running' ? 'is-info' : '']" aria-hidden="true"><svg v-if="log.status === 'success'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="m5 12.5 4 4L19 7" /></svg><svg v-else-if="log.status === 'failed'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M6 18 18 6M6 6l12 12" /></svg><svg v-else fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.7" d="M12 7v5l3 2m6-2a9 9 0 1 1-18 0 9 9 0 0 1-18 0Z" /></svg></span>
              <span class="overview-list-copy"><span class="overview-list-title">{{ log.task_name || '未知任务' }}</span><span class="overview-list-meta">{{ statusLabel(log.status) }} · {{ formatTime(log.created_at) }}<template v-if="log.message"> · {{ compactMessage(log.message) }}</template></span></span>
              <span :class="statusClass(log.status)">{{ statusLabel(log.status) }}</span>
            </router-link>
          </div>
          <div v-else class="overview-empty"><p>还没有运行记录</p><router-link to="/tasks">执行一个备份任务</router-link></div>
        </section>
      </div>

      <section class="surface overview-actions-panel">
        <div>
          <p class="eyebrow">QUICK ACTIONS</p>
          <h3 class="overview-panel-title">继续配置你的备份工作区</h3>
          <p class="mt-2 text-sm text-muted">先连接源站和存储目标，再用备份任务把两者编排起来。</p>
        </div>
        <div class="overview-actions">
          <router-link class="quick-action" to="/servers"><span class="quick-action-label">添加 Bitwarden 源站</span><span class="quick-action-arrow">→</span></router-link>
          <router-link class="quick-action" to="/destinations"><span class="quick-action-label">配置存储目标</span><span class="quick-action-arrow">→</span></router-link>
          <router-link class="quick-action" to="/logs"><span class="quick-action-label">检查运行记录</span><span class="quick-action-arrow">→</span></router-link>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { overviewApi } from '@/api'
import { useToast } from '@/composables/useToast'

const toast = useToast()
const loading = ref(true)
const overview = ref({
  servers: { total: 0, enabled: 0 },
  destinations: { total: 0, enabled: 0 },
  tasks: { total: 0, enabled: 0, scheduled: 0 },
  logs: { total: 0, success_24h: 0, failed_24h: 0, running_24h: 0 },
  recent_tasks: [],
  recent_logs: []
})

const metrics = computed(() => [
  { label: '活跃源站', value: overview.value.servers.enabled, detail: `共 ${overview.value.servers.total} 个已配置`, icon: 'server', tone: 'accent' },
  { label: '可用存储目标', value: overview.value.destinations.enabled, detail: `共 ${overview.value.destinations.total} 个已配置`, icon: 'target', tone: 'info' },
  { label: '启用备份任务', value: overview.value.tasks.enabled, detail: `${overview.value.tasks.scheduled} 个正在自动调度`, icon: 'task', tone: 'violet' },
  { label: '24 小时失败', value: overview.value.logs.failed_24h, detail: `${overview.value.logs.success_24h} 次成功 · ${overview.value.logs.running_24h} 次运行中`, icon: 'log', tone: overview.value.logs.failed_24h ? 'danger' : 'accent' }
])

const statusLabel = (status) => ({ success: '成功', failed: '失败', running: '运行中' }[status] || status || '未知')
const statusClass = (status) => ({ success: 'status-badge status-success', failed: 'status-badge status-danger', running: 'status-badge status-info' }[status] || 'status-badge status-neutral')
const formatTime = (time) => time ? new Date(time).toLocaleString('zh-CN', { month: 'numeric', day: 'numeric', hour: '2-digit', minute: '2-digit' }) : 'N/A'
const compactMessage = (message) => {
  const text = String(message).replace(/\s+/g, ' ').trim()
  return text.length > 38 ? `${text.slice(0, 38)}…` : text
}

const loadOverview = async () => {
  loading.value = true
  try {
    const data = await overviewApi.get()
    overview.value = {
      ...overview.value,
      ...data,
      recent_tasks: data.recent_tasks || [],
      recent_logs: data.recent_logs || []
    }
  } catch (error) {
    console.error('Failed to load overview:', error)
    toast.error('加载工作区总览失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadOverview)
</script>
