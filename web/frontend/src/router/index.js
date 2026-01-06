import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/servers'
  },
  {
    path: '/servers',
    name: 'Servers',
    component: () => import('@/views/Servers.vue')
  },
  {
    path: '/destinations',
    name: 'Destinations',
    component: () => import('@/views/Destinations.vue')
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: () => import('@/views/Tasks.vue')
  },
  {
    path: '/logs',
    name: 'Logs',
    component: () => import('@/views/Logs.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
