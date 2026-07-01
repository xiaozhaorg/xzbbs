<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getUserThreads, type Thread, type PageResult } from '@/api/bbs'

const route = useRoute()
const threads = ref<PageResult<Thread> | null>(null)
const page = ref(1)
const loading = ref(true)

async function load(p = 1) {
  page.value = p
  loading.value = true
  try {
    const res = await getUserThreads(Number(route.params.id), p)
    threads.value = res.data as PageResult<Thread>
  } finally {
    loading.value = false
  }
}

onMounted(() => load())
</script>

<template>
  <div>
    <h1 class="text-xl font-bold mb-4">发布的主题</h1>
    <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>
    <div v-else-if="!threads?.items?.length" class="text-center py-12 text-gray-500 bg-white rounded-lg">暂无主题</div>
    <div v-else class="space-y-2">
      <RouterLink v-for="t in threads.items" :key="t.id" :to="`/thread/${t.id}`" class="block bg-white rounded-lg shadow-sm border p-4 hover:shadow-md transition">
        <h3 class="font-medium text-gray-800">{{ t.title }}</h3>
        <div class="text-xs text-gray-400 mt-1">{{ new Date(t.created_at).toLocaleDateString() }} · {{ t.views }} 浏览 · {{ t.posts }} 回复</div>
      </RouterLink>
    </div>
    <div v-if="threads && threads.pages > 1" class="flex justify-center gap-2 mt-6">
      <button v-for="p in threads.pages" :key="p" @click="load(p)" :class="p === threads.page ? 'bg-blue-600 text-white' : 'bg-white border text-gray-600'" class="px-3 py-1 rounded text-sm">{{ p }}</button>
    </div>
  </div>
</template>
