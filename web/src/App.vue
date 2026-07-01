<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getUnreadCount } from '@/api/bbs'

const auth = useAuthStore()
const unreadCount = ref(0)

onMounted(async () => {
  if (auth.isLoggedIn) {
    try {
      const res = await getUnreadCount()
      unreadCount.value = (res.data as any).count
    } catch { /* ignore */ }
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm border-b sticky top-0 z-50">
      <div class="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
        <RouterLink to="/" class="text-xl font-bold text-blue-600 hover:text-blue-700">
          XzBBS
        </RouterLink>
        <div class="flex items-center gap-4 text-sm">
          <RouterLink to="/" class="text-gray-600 hover:text-blue-600">首页</RouterLink>
          <RouterLink to="/hot" class="text-red-500 hover:text-red-600">🔥 热门</RouterLink>
          <template v-if="auth.isLoggedIn">
            <RouterLink to="/favorites" class="text-gray-600 hover:text-blue-600">收藏</RouterLink>
            <RouterLink to="/pm" class="text-gray-600 hover:text-blue-600">私信</RouterLink>
            <RouterLink to="/notifications" class="text-gray-600 hover:text-blue-600 relative">
              通知
              <span v-if="unreadCount > 0" class="absolute -top-1 -right-2 bg-red-500 text-white text-xs rounded-full w-4 h-4 flex items-center justify-center">{{ unreadCount }}</span>
            </RouterLink>
            <RouterLink to="/profile" class="text-gray-600 hover:text-blue-600">个人中心</RouterLink>
            <RouterLink to="/credits" class="text-gray-600 hover:text-blue-600">💎 {{ auth.user?.credits || 0 }}</RouterLink>
            <span class="text-gray-400">|</span>
            <button @click="auth.logout()" class="text-gray-500 hover:text-red-600">退出</button>
          </template>
          <template v-else>
            <RouterLink to="/login" class="text-gray-600 hover:text-blue-600">登录</RouterLink>
            <RouterLink to="/register" class="text-blue-600 hover:text-blue-700">注册</RouterLink>
          </template>
        </div>
      </div>
    </nav>
    <main class="max-w-6xl mx-auto px-4 py-6">
      <RouterView />
    </main>
  </div>
</template>
