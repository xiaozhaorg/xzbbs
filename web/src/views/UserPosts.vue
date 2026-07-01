<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getUserPosts, type Post, type PageResult } from '@/api/bbs'

const route = useRoute()
const posts = ref<PageResult<Post> | null>(null)
const page = ref(1)
const loading = ref(true)

async function load(p = 1) {
  page.value = p
  loading.value = true
  try {
    const res = await getUserPosts(Number(route.params.id), p)
    posts.value = res.data as PageResult<Post>
  } finally {
    loading.value = false
  }
}

onMounted(() => load())
</script>

<template>
  <div>
    <h1 class="text-xl font-bold mb-4">回复的帖子</h1>
    <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>
    <div v-else-if="!posts?.items?.length" class="text-center py-12 text-gray-500 bg-white rounded-lg">暂无回复</div>
    <div v-else class="space-y-2">
      <div v-for="p in posts.items" :key="p.id" class="bg-white rounded-lg shadow-sm border p-4">
        <p class="text-gray-700 text-sm line-clamp-2">{{ p.content }}</p>
        <div class="text-xs text-gray-400 mt-2">{{ new Date(p.created_at).toLocaleDateString() }}</div>
      </div>
    </div>
    <div v-if="posts && posts.pages > 1" class="flex justify-center gap-2 mt-6">
      <button v-for="p in posts.pages" :key="p" @click="load(p)" :class="p === posts.page ? 'bg-blue-600 text-white' : 'bg-white border text-gray-600'" class="px-3 py-1 rounded text-sm">{{ p }}</button>
    </div>
  </div>
</template>
