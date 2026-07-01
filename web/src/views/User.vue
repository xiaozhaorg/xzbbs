<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getUser, type User } from '@/api/bbs'

const route = useRoute()
const user = ref<User | null>(null)
const loading = ref(true)

const levelTitle = (level: number) => {
  const titles = ['', '新手', '学徒', '入门', '小有所成', '登堂入室', '炉火纯青', '出类拔萃', '技压群雄', '一代宗师', '武林至尊']
  return titles[Math.min(level, 10)] || '未知'
}

const levelColor = (level: number) => {
  if (level >= 8) return 'text-purple-600 bg-purple-50'
  if (level >= 5) return 'text-blue-600 bg-blue-50'
  if (level >= 3) return 'text-green-600 bg-green-50'
  return 'text-gray-500 bg-gray-50'
}

onMounted(async () => {
  try {
    const res = await getUser(Number(route.params.id))
    user.value = res.data as User
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>
  <div v-else-if="user" class="max-w-2xl mx-auto">
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex items-center gap-4">
        <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center text-blue-600 font-bold text-2xl">
          {{ user.username[0].toUpperCase() }}
        </div>
        <div>
          <h1 class="text-xl font-bold">{{ user.username }}</h1>
          <p class="text-sm text-gray-500">ID: {{ user.id }} · 组: {{ user.group_id }}</p>
        </div>
      </div>
      <div class="flex items-center gap-2 mt-3">
        <span :class="levelColor(user.level || 1)" class="text-xs px-2 py-0.5 rounded-full font-medium">
          Lv.{{ user.level || 1 }} {{ levelTitle(user.level || 1) }}
        </span>
      </div>
      <div class="grid grid-cols-3 gap-4 mt-4">
        <div class="text-center p-3 bg-gray-50 rounded-lg">
          <div class="text-lg font-bold">{{ user.threads }}</div>
          <div class="text-sm text-gray-500">主题</div>
        </div>
        <div class="text-center p-3 bg-gray-50 rounded-lg">
          <div class="text-lg font-bold">{{ user.posts }}</div>
          <div class="text-sm text-gray-500">回帖</div>
        </div>
        <div class="text-center p-3 bg-gray-50 rounded-lg">
          <div class="text-lg font-bold">{{ user.credits }}</div>
          <div class="text-sm text-gray-500">积分</div>
        </div>
      </div>
      <div v-if="user.signature" class="mt-4 p-3 bg-yellow-50 border border-yellow-100 rounded-lg text-sm text-gray-600">
        {{ user.signature }}
      </div>
      <div class="flex gap-3 mt-4">
        <RouterLink :to="`/user/${user.id}/threads`" class="text-sm text-blue-600 hover:underline">查看主题</RouterLink>
        <RouterLink :to="`/user/${user.id}/posts`" class="text-sm text-blue-600 hover:underline">查看回帖</RouterLink>
      </div>
      <div class="text-sm text-gray-400 mt-4">
        注册时间: {{ new Date(user.created_at).toLocaleDateString() }}
      </div>
    </div>
  </div>
  <div v-else class="text-center py-12 text-gray-500">用户不存在</div>
</template>
