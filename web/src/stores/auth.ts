import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as bbsApi from '@/api/bbs'
import type { User } from '@/api/bbs'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)

  async function init() {
    if (!token.value) return
    try {
      const res: any = await bbsApi.getMe()
      user.value = res.data as User
    } catch {
      token.value = ''
      localStorage.removeItem('token')
    }
  }

  async function login(account: string, password: string) {
    const res: any = await bbsApi.login({ account, password })
    const data = res.data as { user: User; token: string }
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
  }

  async function register(username: string, email: string, password: string) {
    const res: any = await bbsApi.register({ username, email, password })
    const data = res.data as { user: User; token: string }
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  return { user, token, isLoggedIn, init, login, register, logout }
})
