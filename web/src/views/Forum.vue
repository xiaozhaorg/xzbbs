<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getForumThreads, type Thread, type PageResult } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const forumId = Number(route.params.id)
const threads = ref<PageResult<Thread> | null>(null)
const order = ref('reply')
const loading = ref(true)

onMounted(async () => {
  await loadThreads()
})

async function loadThreads(page = 1) {
  loading.value = true
  try {
    const res = await getForumThreads(forumId, order.value, page)
    threads.value = res.data as PageResult<Thread>
  } finally {
    loading.value = false
  }
}

function changeOrder(o: string) {
  order.value = o
  loadThreads()
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-4">
      <div class="flex gap-2">
        <button
          @click="changeOrder('reply')"
          :class="order === 'reply' ? 'bg-blue-600 text-white' : 'bg-white text-gray-600 border'"
          class="px-3 py-1 rounded text-sm"
        >
          最新回复
        </button>
        <button
          @click="changeOrder('created')"
          :class="order === 'created' ? 'bg-blue-600 text-white' : 'bg-white text-gray-600 border'"
          class="px-3 py-1 rounded text-sm"
        >
          最新发布
        </button>
      </div>
      <button
        v-if="auth.isLoggedIn"
        @click="router.push(`/new-thread/${forumId}`)"
        class="px-4 py-1.5 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-700"
      >
        发帖
      </button>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>

    <div v-else-if="!threads?.items?.length" class="text-center py-12 text-gray-500 bg-white rounded-lg">
      暂无帖子
    </div>

    <div v-else class="space-y-2">
      <RouterLink
        v-for="t in threads.items"
        :key="t.id"
        :to="`/thread/${t.id}`"
        class="block bg-white rounded-lg shadow-sm border p-4 hover:shadow-md transition"
        :class="{ 'border-blue-400 bg-blue-50': t.is_top }"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <span v-if="t.is_top" class="text-xs bg-red-500 text-white px-1.5 py-0.5 rounded">置顶</span>
              <h3 class="font-medium text-gray-800 truncate">{{ t.title }}</h3>
            </div>
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

    <!-- Pagination -->
    <div v-if="threads && threads.pages > 1" class="flex justify-center gap-2 mt-6">
      <button
        v-for="p in threads.pages"
        :key="p"
        @click="loadThreads(p)"
        :class="p === threads.page ? 'bg-blue-600 text-white' : 'bg-white border text-gray-600'"
        class="px-3 py-1 rounded text-sm"
      >
        {{ p }}
      </button>
    </div>
  </div>
</template>
