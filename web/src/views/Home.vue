<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getForums, type Forum } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const forums = ref<Forum[]>([])
const loading = ref(true)
const searchQuery = ref('')
const isSearching = ref(false)

onMounted(async () => {
  try {
    const res = await getForums()
    forums.value = res.data as Forum[]
  } finally {
    loading.value = false
  }
})

function doSearch() {
  if (!searchQuery.value.trim()) return
  isSearching.value = true
  router.push(`/search?q=${encodeURIComponent(searchQuery.value.trim())}`)
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div class="flex gap-3 items-center">
        <h1 class="text-2xl font-bold text-gray-800">版块列表</h1>
        <RouterLink to="/hot" class="text-sm text-red-500 hover:text-red-600">🔥 热门</RouterLink>
      </div>
      <div class="flex gap-2">
        <button
          v-if="auth.isLoggedIn"
          @click="router.push(`/new-thread/${forums[0]?.id}`)"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          发帖
        </button>
        <button
          v-else
          @click="router.push('/login')"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          登录
        </button>
      </div>
    </div>

    <!-- Search bar -->
    <div class="mb-6">
      <div class="flex gap-2">
        <input
          v-model="searchQuery"
          @keyup.enter="doSearch"
          class="flex-1 border rounded-lg px-4 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          placeholder="搜索帖子..."
        />
        <button @click="doSearch" class="px-6 py-2 bg-gray-800 text-white rounded-lg hover:bg-gray-900 text-sm">
          搜索
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>

    <div v-else-if="forums.length === 0" class="text-center py-12 text-gray-500">
      暂无版块
    </div>

    <div v-else class="space-y-3">
      <RouterLink
        v-for="f in forums"
        :key="f.id"
        :to="`/forum/${f.id}`"
        class="block bg-white rounded-lg shadow-sm border p-5 hover:shadow-md transition"
      >
        <div class="flex items-start justify-between">
          <div>
            <h2 class="text-lg font-semibold text-blue-600">{{ f.name }}</h2>
            <p class="text-gray-500 text-sm mt-1 line-clamp-1">{{ f.description || '暂无描述' }}</p>
          </div>
          <div class="text-right text-sm text-gray-400 ml-4 whitespace-nowrap">
            <div>{{ f.threads }} 主题</div>
            <div>{{ f.posts }} 回复</div>
          </div>
        </div>
      </RouterLink>
    </div>
  </div>
</template>
