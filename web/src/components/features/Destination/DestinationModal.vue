<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self.stop>
      <div class="modal-panel modal-panel-wide" role="dialog" aria-modal="true" :aria-label="destination ? '编辑存储目标' : '新建存储目标'">
        <div class="modal-header">
          <div>
            <h3 class="modal-title">{{ destination ? '编辑存储目标' : '新建存储目标' }}</h3>
            <p class="modal-subtitle">{{ destination ? '更新存储方式或策略；敏感字段留空会保留当前值。' : '选择存储方式，再补充连接参数和备份保留策略。' }}</p>
          </div>
          <button class="icon-button" type="button" aria-label="关闭" @click="$emit('close')">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m6 6 12 12M18 6 6 18" /></svg>
          </button>
        </div>

        <form id="destination-form" class="modal-body grid gap-5" @submit.prevent="handleSubmit">
          <section class="form-section">
            <div class="form-section-heading">
              <h4 class="form-section-title">目标信息</h4>
              <p class="form-section-description">名称会出现在任务流程和运行记录中。</p>
            </div>
            <div class="field">
              <label class="field-label" for="destination-name">存储目标名称</label>
              <input id="destination-name" v-model.trim="formData.name" class="input" type="text" required placeholder="例如：异地 WebDAV" />
            </div>
            <TabSelector v-model="formData.type" :options="storageTypes" label="存储类型" />
            <div v-if="formData.type === 'server'" class="field destination-primary-field">
              <CustomSelect
                v-model="formData.target_server_id"
                :options="serverOptions"
                label="目标服务器"
                placeholder="请选择目标服务器"
                empty-text="暂无可用源站，请先创建源站"
              />
              <p class="field-hint">备份文件会导入这个已配置的 Bitwarden 目标服务器。</p>
            </div>
          </section>

          <section v-if="formData.type !== 'server'" class="form-section">
            <div class="form-section-heading">
              <h4 class="form-section-title">连接配置</h4>
              <p class="form-section-description">只显示当前存储类型需要的字段，密码编辑时留空表示保持原值。</p>
            </div>

            <div v-if="formData.type === 'local'" class="field">
              <label class="field-label" for="local-path">本地路径</label>
              <input id="local-path" v-model="formData.local_path" class="input" type="text" required placeholder="/app/backups" />
              <p class="field-hint">Docker 部署建议填写 <code>/app/backups</code>，并将宿主机目录挂载到此容器路径。</p>
              <p class="field-hint mt-1 text-muted">例如 <code>-v /data/backups:/app/backups</code>；非 Docker 部署请填写运行服务所在机器的绝对路径。</p>
            </div>

            <div v-else-if="formData.type === 'webdav'" class="grid gap-4">
              <div class="field">
                <label class="field-label" for="webdav-url">WebDAV URL</label>
                <input id="webdav-url" v-model="formData.webdav_url" class="input" type="url" required placeholder="https://dav.example.com" />
              </div>
              <div class="form-grid">
                <div class="field">
                  <label class="field-label" for="webdav-username">用户名</label>
                  <input id="webdav-username" v-model="formData.webdav_username" class="input" type="text" autocomplete="username" />
                </div>
                <div class="field">
                  <label class="field-label" for="webdav-password">密码</label>
                  <input id="webdav-password" v-model="formData.webdav_password" class="input" type="password" :required="!destination" autocomplete="new-password" :placeholder="destination ? '留空保持原值' : '输入 WebDAV 密码'" />
                  <p v-if="destination" class="field-hint">留空表示不修改当前密码。</p>
                </div>
              </div>
              <div class="field">
                <label class="field-label" for="webdav-path">存储路径 <span>可选</span></label>
                <input id="webdav-path" v-model="formData.webdav_path" class="input" type="text" placeholder="/bitwarden-backup" />
                <p class="field-hint">留空会使用默认路径 <code>/bitwarden-backup</code>。</p>
                <p class="field-hint mt-1 text-muted">首次备份时会自动创建缺失目录。</p>
              </div>
            </div>

            <div v-else-if="formData.type === 's3'" class="grid gap-4">
              <div class="form-grid">
                <div class="field">
                  <label class="field-label" for="s3-endpoint">Endpoint</label>
                  <input id="s3-endpoint" v-model="formData.s3_endpoint" class="input" type="url" required placeholder="https://s3.amazonaws.com" />
                  <p class="field-hint">填写完整地址，例如 <code>https://s3.amazonaws.com</code>。</p>
                </div>
                <div class="field">
                  <label class="field-label" for="s3-region">区域 <span>Region</span></label>
                  <input id="s3-region" v-model="formData.s3_region" class="input" type="text" required placeholder="us-east-1" />
                </div>
              </div>
              <div class="field">
                <label class="field-label" for="s3-bucket">Bucket 名称</label>
                <input id="s3-bucket" v-model="formData.s3_bucket" class="input" type="text" required />
              </div>
              <div class="form-grid">
                <div class="field">
                  <label class="field-label" for="s3-access-key">Access Key <span v-if="destination">可选</span></label>
                  <input id="s3-access-key" v-model="formData.s3_access_key" class="input" type="text" autocomplete="off" />
                </div>
                <div class="field">
                  <label class="field-label" for="s3-secret-key">Secret Key</label>
                  <input id="s3-secret-key" v-model="formData.s3_secret_key" class="input" type="password" :required="!destination" autocomplete="new-password" :placeholder="destination ? '留空保持原值' : '输入 Secret Key'" />
                  <p v-if="destination" class="field-hint">留空表示不修改当前密钥。</p>
                </div>
              </div>
              <div class="field">
                <label class="field-label" for="s3-path">存储路径 <span>可选</span></label>
                <input id="s3-path" v-model="formData.s3_path" class="input" type="text" placeholder="/bitwarden-backup" />
                <p class="field-hint">留空会使用默认路径 <code>/bitwarden-backup</code>。</p>
              </div>
            </div>

          </section>

          <section v-if="['local', 'webdav', 's3'].includes(formData.type)" class="form-section">
            <div class="form-section-heading">
              <h4 class="form-section-title">安全与保留</h4>
              <p class="form-section-description">把备份文件保护和清理策略放在一起，保存前可以清楚确认影响范围。</p>
            </div>
            <div class="surface-muted flex items-center justify-between gap-4 p-3">
              <div>
                <p class="text-sm font-semibold text-main">加密备份文件</p>
                <p class="mt-1 text-xs text-muted">使用密码保护导出的备份文件。</p>
              </div>
              <ToggleButton v-model="formData.encrypted" label="启用" aria-label="加密备份文件" />
            </div>
            <div v-if="formData.encrypted" class="field">
              <label class="field-label" for="encryption-password">加密密码</label>
              <input id="encryption-password" v-model="formData.encryption_password" class="input" type="password" :required="!destination || !destination.encrypted" autocomplete="new-password" :placeholder="destination ? '留空保持原值' : '请输入加密密码'" />
              <p class="field-hint">{{ destination ? '留空表示不修改。' : '解密备份文件时需要使用相同密码。' }}</p>
            </div>
            <div class="surface-muted flex items-center justify-between gap-4 p-3">
              <div>
                <p class="text-sm font-semibold text-main">限制保留数量</p>
                <p class="mt-1 text-xs text-muted">超过数量后自动删除最旧的备份文件。</p>
              </div>
              <ToggleButton v-model="retentionEnabled" label="启用" aria-label="限制保留数量" />
            </div>
            <div v-if="retentionEnabled" class="field">
              <label class="field-label" for="max-backup-count">最多保留份数</label>
              <div class="relative">
                <input id="max-backup-count" v-model.number="formData.max_backup_count" class="input pr-12" type="number" min="1" placeholder="5" />
                <span class="pointer-events-none absolute inset-y-0 right-3 flex items-center text-xs font-semibold text-muted">份</span>
              </div>
              <p class="field-hint text-warning">超过限制时会自动删除最旧的备份文件。</p>
            </div>
            <p v-else class="field-hint">当前保留所有历史备份文件，不限制数量。</p>
          </section>
        </form>

        <div class="modal-footer">
          <button type="button" class="btn-secondary" @click="$emit('close')">取消</button>
          <button form="destination-form" type="submit" class="btn-primary" :disabled="loading">
            <span v-if="loading" class="spinner"></span>{{ loading ? '保存中…' : destination ? '保存更改' : '保存存储目标' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { destinationsApi, serversApi } from '@/api'
import { useToast } from '@/composables/useToast'
import TabSelector from '@/components/ui/TabSelector.vue'
import ToggleButton from '@/components/ui/ToggleButton.vue'
import CustomSelect from '@/components/ui/CustomSelect.vue'

const props = defineProps({ destination: Object })
const emit = defineEmits(['close', 'saved'])
const toast = useToast()

const storageTypes = [
  { label: '本地存储', value: 'local' },
  { label: 'WebDAV', value: 'webdav' },
  { label: 'S3', value: 's3' },
  { label: '服务器', value: 'server' }
]

const servers = ref([])
const emptyForm = () => ({
  name: '', type: 'local', local_path: '', webdav_url: '', webdav_username: '', webdav_password: '', webdav_path: '',
  s3_endpoint: '', s3_region: '', s3_bucket: '', s3_access_key: '', s3_secret_key: '', s3_path: '', target_server_id: '',
  enabled: true, encrypted: false, encryption_password: '', max_backup_count: 5
})
const formData = ref(emptyForm())
const loading = ref(false)
const retentionEnabled = ref(false)
const serverOptions = computed(() => {
  const currentID = Number(formData.value.target_server_id || 0)
  return servers.value
    .filter(server => server.enabled || Number(server.id) === currentID)
    .map(server => ({
      label: server.name,
      value: server.id,
      description: `${server.server_url}${server.enabled ? '' : ' · 已停用'}`
    }))
})

watch(() => props.destination, (newDestination) => {
  if (newDestination) {
    formData.value = {
      ...newDestination,
      local_path: newDestination.local_path || (newDestination.type === 'local' ? newDestination.path : ''),
      target_server_id: newDestination.target_server_id || '',
      webdav_password: '',
      s3_access_key: '',
      s3_secret_key: '',
      encrypted: newDestination.encrypted || false,
      encryption_password: newDestination.encryption_password || '',
      max_backup_count: newDestination.max_backup_count || 5
    }
    retentionEnabled.value = Boolean(newDestination.max_backup_count && newDestination.max_backup_count > 0)
  } else {
    formData.value = emptyForm()
    retentionEnabled.value = false
  }
}, { immediate: true })

watch(() => formData.value.type, (type) => {
  if (type === 'server') {
    formData.value.encrypted = false
    formData.value.encryption_password = ''
    retentionEnabled.value = false
  } else {
    formData.value.target_server_id = ''
  }
})

const loadServers = async () => {
  try {
    const res = await serversApi.getAll({ page: 1, page_size: 1000 })
    servers.value = res.data || []
  } catch (error) {
    console.error('Failed to load servers:', error)
  }
}

onMounted(loadServers)

const buildSubmitData = () => {
  const current = formData.value
  const data = {
    name: current.name.trim(),
    type: current.type,
    enabled: current.enabled,
    encrypted: current.type === 'server' ? false : Boolean(current.encrypted),
    max_backup_count: current.type === 'server' ? 0 : retentionEnabled.value ? Number(current.max_backup_count) || 5 : 0
  }

  if (current.type === 'local') {
    data.local_path = current.local_path.trim()
  } else if (current.type === 'webdav') {
    data.webdav_url = current.webdav_url.trim()
    data.webdav_username = current.webdav_username
    data.webdav_path = current.webdav_path.trim() || '/bitwarden-backup'
    if (current.webdav_password) data.webdav_password = current.webdav_password
  } else if (current.type === 's3') {
    data.s3_endpoint = current.s3_endpoint.trim()
    data.s3_region = current.s3_region.trim()
    data.s3_bucket = current.s3_bucket.trim()
    data.s3_path = current.s3_path.trim() || '/bitwarden-backup'
    if (current.s3_access_key) data.s3_access_key = current.s3_access_key
    if (current.s3_secret_key) data.s3_secret_key = current.s3_secret_key
  } else if (current.type === 'server') {
    data.target_server_id = Number(current.target_server_id)
  }

  if (data.encrypted && current.encryption_password) data.encryption_password = current.encryption_password
  if (!props.destination?.id) delete data.enabled
  return data
}

const handleSubmit = async () => {
  if (!formData.value.name.trim()) {
    toast.error('请输入存储目标名称')
    return
  }
  if (formData.value.type === 'server' && !formData.value.target_server_id) {
    toast.error('请选择目标服务器')
    return
  }
  if (formData.value.encrypted && !formData.value.encryption_password && !props.destination?.encrypted) {
    toast.error('启用加密时必须设置加密密码')
    return
  }
  if (retentionEnabled.value && Number(formData.value.max_backup_count) < 1) {
    toast.error('最多保留份数必须大于 0')
    return
  }

  loading.value = true
  try {
    const data = buildSubmitData()

    if (props.destination?.id) {
      await destinationsApi.update(props.destination.id, data)
      toast.success('存储目标已更新')
    } else {
      await destinationsApi.create(data)
      toast.success('存储目标已创建')
    }
    emit('saved')
  } catch (error) {
    console.error('保存失败:', error)
    toast.error(error.message || '保存存储目标失败')
  } finally {
    loading.value = false
  }
}
</script>
