import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  { path: '/', name: 'Home', component: () => import('@/views/Home.vue') },
  { path: '/forum/:id', name: 'Forum', component: () => import('@/views/Forum.vue') },
  { path: '/thread/:id', name: 'Thread', component: () => import('@/views/Thread.vue') },
  { path: '/new-thread/:fid', name: 'NewThread', component: () => import('@/views/NewThread.vue'), meta: { requiresAuth: true } },
  { path: '/user/:id', name: 'User', component: () => import('@/views/User.vue') },
  { path: '/profile', name: 'Profile', component: () => import('@/views/Profile.vue'), meta: { requiresAuth: true } },
  { path: '/login', name: 'Login', component: () => import('@/views/Login.vue'), meta: { guestOnly: true } },
  { path: '/register', name: 'Register', component: () => import('@/views/Register.vue'), meta: { guestOnly: true } },
  { path: '/admin', name: 'Admin', component: () => import('@/views/admin/Dashboard.vue'), meta: { requiresAuth: true, requiresAdmin: true } },
  { path: '/search', name: 'Search', component: () => import('@/views/Search.vue') },
  { path: '/favorites', name: 'Favorites', component: () => import('@/views/Favorites.vue'), meta: { requiresAuth: true } },
  { path: '/notifications', name: 'Notifications', component: () => import('@/views/Notifications.vue'), meta: { requiresAuth: true } },
  { path: '/pm', name: 'PM', component: () => import('@/views/PM.vue'), meta: { requiresAuth: true } },
  { path: '/pm/:otherId', name: 'PMConversation', component: () => import('@/views/PMConversation.vue'), meta: { requiresAuth: true } },
  { path: '/hot', name: 'HotThreads', component: () => import('@/views/HotThreads.vue') },
  { path: '/credits', name: 'Credits', component: () => import('@/views/Credits.vue'), meta: { requiresAuth: true } },
  { path: '/posts/:id/edits', name: 'PostEdits', component: () => import('@/views/PostEdits.vue'), meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  // Init auth state
  if (!auth.user && auth.token) {
    await auth.init()
  }

  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return '/login'
  }
  if (to.meta.guestOnly && auth.isLoggedIn) {
    return '/'
  }
  if (to.meta.requiresAdmin && auth.user?.group_id !== 1) {
    return '/'
  }
})

export default router
