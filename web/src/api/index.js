const API_BASE = '/api'

let csrfToken = ''

const handleResponse = async (response) => {
  if (response.status === 204) return null
  if (!response.ok) {
    let message = `HTTP ${response.status}`
    try {
      const data = await response.json()
      message = data.error || message
    } catch {
      // Keep the generic HTTP message when the server did not return JSON.
    }
    throw new Error(message)
  }
  return response.json()
}

const request = async (path, options = {}) => {
  const method = (options.method || 'GET').toUpperCase()
  const headers = new Headers(options.headers || {})
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }
  if (!['GET', 'HEAD', 'OPTIONS'].includes(method) && csrfToken) {
    headers.set('X-CSRF-Token', csrfToken)
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    method,
    headers,
    credentials: 'same-origin'
  })

  if (response.status === 401 && !path.startsWith('/auth/login')) {
    csrfToken = ''
    window.dispatchEvent(new CustomEvent('auth:required'))
  }
  return handleResponse(response)
}

export const authApi = {
  login: async (password) => {
    const data = await request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ password })
    })
    csrfToken = data.csrf_token || ''
    return data
  },
  session: async () => {
    const data = await request('/auth/session')
    csrfToken = data.csrf_token || ''
    return data
  },
  logout: async () => {
    try {
      return await request('/auth/logout', { method: 'POST' })
    } finally {
      csrfToken = ''
    }
  }
}

export const overviewApi = {
  get: () => request('/overview')
}

const paginatedPath = (resource, params = {}) => {
  const query = new URLSearchParams()
  if (params.enabled !== undefined) query.append('enabled', params.enabled)
  if (params.task_id) query.append('task_id', params.task_id)
  if (params.page) query.append('page', params.page)
  if (params.page_size) query.append('page_size', params.page_size)
  const queryString = query.toString() ? `?${query.toString()}` : ''
  return `/${resource}${queryString}`
}

export const serversApi = {
  getAll: (params = {}) => request(paginatedPath('servers', params)),
  getById: (id) => request(`/servers/${id}`),
  create: (data) => request('/servers', { method: 'POST', body: JSON.stringify(data) }),
  update: (id, data) => request(`/servers/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  setEnabled: (id, enabled) => request(`/servers/${id}/enabled`, { method: 'PATCH', body: JSON.stringify({ enabled }) }),
  delete: (id) => request(`/servers/${id}`, { method: 'DELETE' })
}

export const destinationsApi = {
  getAll: (params = {}) => request(paginatedPath('destinations', params)),
  getById: (id) => request(`/destinations/${id}`),
  create: (data) => request('/destinations', { method: 'POST', body: JSON.stringify(data) }),
  update: (id, data) => request(`/destinations/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  setEnabled: (id, enabled) => request(`/destinations/${id}/enabled`, { method: 'PATCH', body: JSON.stringify({ enabled }) }),
  delete: (id) => request(`/destinations/${id}`, { method: 'DELETE' }),
  toggle: (id) => request(`/destinations/${id}/toggle`, { method: 'PATCH' })
}

export const tasksApi = {
  getAll: (params = {}) => request(paginatedPath('tasks', params)),
  getById: (id) => request(`/tasks/${id}`),
  create: (data) => request('/tasks', { method: 'POST', body: JSON.stringify(data) }),
  update: (id, data) => request(`/tasks/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  setEnabled: (id, enabled) => request(`/tasks/${id}/enabled`, { method: 'PATCH', body: JSON.stringify({ enabled }) }),
  delete: (id) => request(`/tasks/${id}`, { method: 'DELETE' }),
  execute: (id) => request(`/tasks/${id}/execute`, { method: 'POST' })
}

export const logsApi = {
  getAll: (params = {}) => request(paginatedPath('logs', params))
}
