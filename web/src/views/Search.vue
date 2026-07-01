<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { searchThreads, type Thread } from '@/api/bbs'

const route = useRoute()
const router = useRouter()
const query = ref((route.query.q as string) || '')
const threads = ref<Thread[]>([])
const total = ref(0)
const page = ref(1)
const pages = ref(1)
const loading = ref(false)

function doSearch(p = 1) {
  page.value = p
  loading.value = true
  searchThreads(query.value, undefined, p)
    .then(res => {
      const data = res.data as any
      threads.value = data.items
      total.value = data.total
      pages.value = data.pages
    })
    .finally(() => { loading.value = false })
}

onMounted(() => {
  if (query.value) doSearch()
})
</script>

<template>
  <div>
    <div class="flex items-center gap-2 mb-4">
      <button @click="router.push('/')" class="text-gray-500 hover:text-blue-600">← 返回首页</button>
    </div>
    <h1 class="text-xl font-bold mb-2">搜索结果</h1>
    <p class="text-sm text-gray-500 mb-4">"{{ query }}" 共 {{ total }} 条结果</p>

    <div class="flex gap-2 mb-4">
      <input
        v-model="query"
        @keyup.enter="doSearch(1)"
        class="flex-1 border rounded-lg px-4 py-2 text-sm focus:ring-2 focus:ring-blue-500"
        placeholder="搜索帖子..."
      />
      <button @click="doSearch(1)" class="px-6 py-2 bg-gray-800 text-white rounded-lg text-sm hover:bg-gray-900">搜索</button>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-500">搜索中...</div>
    <div v-else-if="threads.length === 0" class="text-center py-8 text-gray-500">未找到相关帖子</div>
    <div v-else class="space-y-2">
      <RouterLink
        v-for="t in threads"
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
              {{ t.user?.username || '' }} · {{ t.forum?.name || '' }} · {{ new Date(t.created_at).toLocaleDateString() }}
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
    <div v-if="pages > 1" class="flex justify-center gap-2 mt-6">
      <button
        v-for="p in pages"
        :key="p"
        @click="doSearch(p)"
        :class="p === page ? 'bg-blue-600 text-white' : 'bg-white border text-gray-600'"
        class="px-3 py-1 rounded text-sm"
      >
        {{ p }}
      </button>
    </div>
  </div>
</template>
