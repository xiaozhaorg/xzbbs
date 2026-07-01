<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getNotifications, markAllNotificationsRead, markNotificationsRead, type Notification } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const notifications = ref<Notification[]>([])
const loading = ref(true)

onMounted(async () => {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  await loadNotifications()
})

async function loadNotifications() {
  loading.value = true
  try {
    const res = await getNotifications(false, 1)
    const data = res.data as any
    notifications.value = data.items
  } finally {
    loading.value = false
  }
}

async function markRead(ids: number[]) {
  await markNotificationsRead(ids)
  await loadNotifications()
}

async function markAll() {
  await markAllNotificationsRead()
  await loadNotifications()
}

function notifText(n: Notification): string {
  const typeText = { 0: '回复了你的帖子', 1: '在帖子中@了你', 2: '系统通知' }
  return typeText[n.type as keyof typeof typeText] || '新通知'
}

function timeAgo(dateStr: string): string {
  const d = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((now.getTime() - d.getTime()) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  if (diff < 604800) return `${Math.floor(diff / 86400)}天前`
  return d.toLocaleDateString()
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-xl font-bold">通知</h1>
      <button @click="markAll" class="text-sm text-blue-600 hover:underline">全部已读</button>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>

    <div v-else-if="notifications.length === 0" class="text-center py-12 text-gray-500 bg-white rounded-lg">
      暂无通知
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="n in notifications"
        :key="n.id"
        class="bg-white rounded-lg shadow-sm border p-4 flex items-start gap-3"
        :class="{ 'bg-blue-50': !n.is_read }"
      >
        <div class="w-2 h-2 rounded-full mt-2 flex-shrink-0" :class="n.is_read ? 'bg-gray-300' : 'bg-blue-500'"></div>
        <div class="flex-1 min-w-0">
          <p class="text-sm text-gray-700">
            <span class="font-medium">{{ n.message }}</span>
          </p>
          <div class="flex items-center gap-2 mt-1">
            <span class="text-xs text-gray-400">{{ timeAgo(n.created_at) }}</span>
          </div>
        </div>
        <button
          v-if="n.thread_id"
          @click="router.push(`/thread/${n.thread_id}`)"
          class="text-xs text-blue-600 hover:underline whitespace-nowrap"
        >
          查看帖子
        </button>
      </div>
    </div>
  </div>
</template>
