// API 服务层
const API_BASE = '/api'

const handleResponse = async (response) => {
  if (!response.ok) {
    // 尝试解析后端返回的错误信息
    try {
      const data = await response.json()
      throw new Error(data.error || `HTTP ${response.status}`)
    } catch (e) {
      if (e.message) throw e
      throw new Error(`HTTP ${response.status}`)
    }
  }
  return response.json()
}

// Servers API
export const serversApi = {
  getAll: (params = {}) => {
    const query = new URLSearchParams()
    if (params.enabled !== undefined) query.append('enabled', params.enabled)
    if (params.page) query.append('page', params.page)
    if (params.page_size) query.append('page_size', params.page_size)
    const queryString = query.toString() ? `?${query.toString()}` : ''
    return fetch(`${API_BASE}/servers${queryString}`).then(handleResponse)
  },
  getById: (id) => fetch(`${API_BASE}/servers/${id}`).then(handleResponse),
  create: (data) => fetch(`${API_BASE}/servers`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }).then(handleResponse),
  update: (id, data) => fetch(`${API_BASE}/servers/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }).then(handleResponse),
  delete: (id) => fetch(`${API_BASE}/servers/${id}`, {
    method: 'DELETE'
  }).then(handleResponse)
}

// Destinations API
export const destinationsApi = {
  getAll: (params = {}) => {
    const query = new URLSearchParams()
    if (params.page) query.append('page', params.page)
    if (params.page_size) query.append('page_size', params.page_size)
    const queryString = query.toString() ? `?${query.toString()}` : ''
    return fetch(`${API_BASE}/destinations${queryString}`).then(handleResponse)
  },
  getById: (id) => fetch(`${API_BASE}/destinations/${id}`).then(handleResponse),
  create: (data) => fetch(`${API_BASE}/destinations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }).then(handleResponse),
  update: (id, data) => fetch(`${API_BASE}/destinations/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }).then(handleResponse),
  delete: (id) => fetch(`${API_BASE}/destinations/${id}`, {
    method: 'DELETE'
  }).then(handleResponse)
}

// Tasks API
export const tasksApi = {
  getAll: (params = {}) => {
    const query = new URLSearchParams()
    if (params.page) query.append('page', params.page)
    if (params.page_size) query.append('page_size', params.page_size)
    const queryString = query.toString() ? `?${query.toString()}` : ''
    return fetch(`${API_BASE}/tasks${queryString}`).then(handleResponse)
  },
  getById: (id) => fetch(`${API_BASE}/tasks/${id}`).then(handleResponse),
  create: (data) => fetch(`${API_BASE}/tasks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }).then(handleResponse),
  update: (id, data) => fetch(`${API_BASE}/tasks/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  }).then(handleResponse),
  delete: (id) => fetch(`${API_BASE}/tasks/${id}`, {
    method: 'DELETE'
  }).then(handleResponse),
  execute: (id) => fetch(`${API_BASE}/tasks/${id}/execute`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' }
  }).then(handleResponse)
}

// Logs API
export const logsApi = {
  getAll: (params = {}) => {
    const query = new URLSearchParams()
    if (params.task_id) query.append('task_id', params.task_id)
    if (params.page) query.append('page', params.page)
    if (params.page_size) query.append('page_size', params.page_size)
    const queryString = query.toString() ? `?${query.toString()}` : ''
    return fetch(`${API_BASE}/logs${queryString}`).then(handleResponse)
  }
}
