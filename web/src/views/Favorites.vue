<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getFavorites, type Thread, type PageResult } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const threads = ref<PageResult<Thread> | null>(null)
const page = ref(1)
const loading = ref(true)

onMounted(() => {
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  loadFavorites()
})

async function loadFavorites(p = 1) {
  page.value = p
  loading.value = true
  try {
    const res = await getFavorites(p)
    threads.value = res.data as any
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="text-xl font-bold mb-4">我的收藏</h1>

    <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>

    <div v-else-if="!threads?.items?.length" class="text-center py-12 text-gray-500 bg-white rounded-lg">
      还没有收藏任何帖子
    </div>

    <div v-else class="space-y-2">
      <RouterLink
        v-for="t in threads.items"
        :key="t.id"
        :to="`/thread/${t.id}`"
        class="block bg-white rounded-lg shadow-sm border p-4 hover:shadow-md transition"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0">
            <h3 class="font-medium text-gray-800 truncate">{{ t.title }}</h3>
            <div class="text-xs text-gray-400 mt-1.5">
              {{ t.user?.username || '' }} · {{ new Date(t.created_at).toLocaleDateString() }}
            </div>
          </div>
          <div class="text-xs text-gray-400 ml-4 flex gap-3">
            <span>{{ t.views }} 浏览</span>
            <span>{{ t.posts }} 回复</span>
          </div>
        </div>
      </RouterLink>
    </div>

    <div v-if="threads && threads.pages > 1" class="flex justify-center gap-2 mt-6">
      <button
        v-for="p in threads.pages"
        :key="p"
        @click="loadFavorites(p)"
        :class="p === threads.page ? 'bg-blue-600 text-white' : 'bg-white border text-gray-600'"
        class="px-3 py-1 rounded text-sm"
      >
        {{ p }}
      </button>
    </div>
  </div>
</template>
