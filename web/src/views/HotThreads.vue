<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import api from '@/api'

const threads = ref<any>(null)
const loading = ref(true)

async function load(page = 1) {
  loading.value = true
  try {
    const res = await api.get('/threads', { params: { order: 'views', page } })
    threads.value = res.data
  } finally {
    loading.value = false
  }
}

onMounted(() => load())
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">🔥 热门帖子</h1>
    <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>
    <div v-else-if="!threads?.items?.length" class="text-center py-12 text-gray-500 bg-white rounded-lg">暂无帖子</div>
    <div v-else class="space-y-2">
      <RouterLink v-for="t in threads.items" :key="t.id" :to="`/thread/${t.id}`" class="block bg-white rounded-lg shadow-sm border p-4 hover:shadow-md transition">
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0">
            <h3 class="font-medium text-gray-800">{{ t.title }}</h3>
            <div class="text-xs text-gray-400 mt-1">{{ t.user?.username || '' }} · {{ t.forum?.name || '' }} · {{ new Date(t.created_at).toLocaleDateString() }}</div>
          </div>
          <div class="text-xs text-gray-400 ml-4 flex gap-3">
            <span>{{ t.views }} 浏览</span>
            <span>{{ t.posts }} 回复</span>
          </div>
        </div>
      </RouterLink>
    </div>
    <div v-if="threads && threads.pages > 1" class="flex justify-center gap-2 mt-6">
      <button v-for="p in threads.pages" :key="p" @click="load(p)" :class="p === threads.page ? 'bg-blue-600 text-white' : 'bg-white border text-gray-600'" class="px-3 py-1 rounded text-sm">{{ p }}</button>
    </div>
  </div>
</template>
