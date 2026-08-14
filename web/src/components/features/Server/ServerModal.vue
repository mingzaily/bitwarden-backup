<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self.stop>
      <div class="modal-panel" role="dialog" aria-modal="true" :aria-label="server ? '编辑源站' : '新建源站'">
        <div class="modal-header">
          <div>
            <h3 class="modal-title">{{ server ? '编辑源站' : '新建源站' }}</h3>
            <p class="modal-subtitle">{{ server ? '更新连接信息；敏感字段留空会保留当前凭证。' : '配置 Bitwarden 源站和用于备份的访问凭证。' }}</p>
          </div>
          <button class="icon-button" type="button" aria-label="关闭" @click="$emit('close')">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="m6 6 12 12M18 6 6 18" /></svg>
          </button>
        </div>

        <form id="server-form" class="modal-body grid gap-5" @submit.prevent="handleSubmit">
          <section class="form-section">
            <div class="form-section-heading">
              <h4 class="form-section-title">基础信息</h4>
              <p class="form-section-description">给连接起一个容易识别的名称，并选择 Bitwarden 源站类型。</p>
            </div>
            <div class="field">
              <label class="field-label" for="server-name">源站名称</label>
              <input id="server-name" v-model.trim="formData.name" class="input" type="text" required placeholder="例如：生产环境 Bitwarden" />
            </div>
            <TabSelector v-model="formData.server_type" :options="serverTypes" label="服务器类型" />
            <TabSelector v-if="formData.server_type === 'official'" v-model="formData.region" :options="officialRegions" label="服务器区域" />
            <div v-if="formData.server_type === 'self-hosted'" class="field">
              <label class="field-label" for="server-url">源站地址</label>
              <input id="server-url" v-model="formData.url" class="input" type="url" required placeholder="https://example.com" />
              <p class="field-hint">填写自建 Bitwarden 服务器的完整地址。</p>
            </div>
          </section>

          <section class="form-section">
            <div class="form-section-heading">
              <h4 class="form-section-title">访问凭证</h4>
              <p class="form-section-description">这些凭证仅用于调用 Bitwarden CLI；编辑时只填写需要替换的字段。</p>
            </div>
            <div class="form-grid">
              <div class="field form-grid-span">
                <label class="field-label" for="client-id">Client ID</label>
                <input id="client-id" v-model="formData.client_id" class="input" type="text" required autocomplete="off" placeholder="user.xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" />
              </div>
              <div class="field">
                <div class="field-label-row">
                  <label class="field-label" for="client-secret">Client Secret</label>
                  <span v-if="server" class="field-label-note">可选</span>
                </div>
                <input id="client-secret" v-model="formData.client_secret" class="input" type="password" :required="!server" autocomplete="new-password" :placeholder="server ? '留空保持原值' : '输入 Client Secret'" />
                <p v-if="server" class="field-hint">留空表示不修改当前凭证。</p>
              </div>
              <div class="field">
                <div class="field-label-row">
                  <label class="field-label" for="master-password">Master Password</label>
                </div>
                <input id="master-password" v-model="formData.master_password" class="input" type="password" :required="!server" autocomplete="new-password" :placeholder="server ? '留空保持原值' : '输入 Master Password'" />
                <p class="field-hint">{{ server ? '留空表示不修改当前密码；填写后将替换。' : '备份时通过 Bitwarden CLI 解锁 Bitwarden 保险库。' }}</p>
              </div>
            </div>
          </section>
        </form>

        <div class="modal-footer">
          <button type="button" class="btn-secondary" @click="$emit('close')">取消</button>
          <button form="server-form" type="submit" class="btn-primary" :disabled="loading">
            <span v-if="loading" class="spinner"></span>{{ loading ? '保存中…' : server ? '保存更改' : '保存源站' }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'
import { serversApi } from '@/api'
import { useToast } from '@/composables/useToast'
import TabSelector from '@/components/ui/TabSelector.vue'

const props = defineProps({ server: Object })
const emit = defineEmits(['close', 'saved'])
const toast = useToast()

const serverTypes = [
  { label: '官方服务器', value: 'official' },
  { label: '自建服务器', value: 'self-hosted' }
]

const officialRegions = [
  { label: 'bitwarden.com', value: 'https://vault.bitwarden.com' },
  { label: 'bitwarden.eu', value: 'https://vault.bitwarden.eu' }
]

const getOfficialRegion = (rawURL) => {
  try {
    const parsed = new URL(String(rawURL || '').trim())
    const hostname = parsed.hostname.toLowerCase().replace(/\.$/, '')
    if (parsed.protocol !== 'https:') return ''
    if (hostname === 'vault.bitwarden.com') return 'https://vault.bitwarden.com'
    if (hostname === 'vault.bitwarden.eu') return 'https://vault.bitwarden.eu'
    return ''
  } catch {
    return ''
  }
}

const formData = ref({
  name: '',
  server_type: 'official',
  region: 'https://vault.bitwarden.com',
  url: '',
  client_id: '',
  client_secret: '',
  master_password: ''
})
const loading = ref(false)

watch([() => formData.value.server_type, () => formData.value.region], ([type, region]) => {
  if (type === 'official') formData.value.url = region
  else if (!props.server || formData.value.url === 'https://vault.bitwarden.com' || formData.value.url === 'https://vault.bitwarden.eu') formData.value.url = ''
})

watch(() => props.server, (newServer) => {
  if (newServer) {
    const serverUrl = newServer.server_url || newServer.url || ''
    const officialRegion = getOfficialRegion(serverUrl)
    const isOfficial = newServer.is_official === true || Boolean(officialRegion)
    formData.value = { ...newServer, url: serverUrl, server_type: isOfficial ? 'official' : 'self-hosted', region: isOfficial ? (officialRegion || 'https://vault.bitwarden.com') : 'https://vault.bitwarden.com' }
  } else {
    formData.value = { name: '', server_type: 'official', region: 'https://vault.bitwarden.com', url: 'https://vault.bitwarden.com', client_id: '', client_secret: '', master_password: '' }
  }
}, { immediate: true })

const handleSubmit = async () => {
  loading.value = true
  try {
    const submitData = {
      name: formData.value.name,
      server_url: formData.value.url,
      client_id: formData.value.client_id,
      client_secret: formData.value.client_secret,
      master_password: formData.value.master_password
    }
    if (props.server?.id) {
      await serversApi.update(props.server.id, submitData)
      toast.success('源站已更新')
    } else {
      await serversApi.create(submitData)
      toast.success('源站已创建')
    }
    emit('saved')
  } catch (error) {
    console.error('Failed to save server:', error)
    toast.error(error.message || '保存源站失败')
  } finally {
    loading.value = false
  }
}
</script>
