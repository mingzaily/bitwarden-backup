// API 服务层
import { useToast } from '@/composables/useToast'

const API_BASE = '/api'

const handleResponse = async (response) => {
  if (!response.ok) {
    const toast = useToast()
    toast.error(`请求失败: ${response.statusText}`)
    throw new Error(`HTTP ${response.status}`)
  }
  return response.json()
}

// Servers API
export const serversApi = {
  getAll: () => fetch(`${API_BASE}/servers`).then(handleResponse),
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
  getAll: () => fetch(`${API_BASE}/destinations`).then(handleResponse),
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
  getAll: () => fetch(`${API_BASE}/tasks`).then(handleResponse),
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
  getAll: () => fetch(`${API_BASE}/logs`).then(handleResponse)
}
