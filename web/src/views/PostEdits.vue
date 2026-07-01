<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getPostEdits } from '@/api/bbs'

const route = useRoute()
const edits = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await getPostEdits(Number(route.params.id))
    edits.value = res.data || []
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="max-w-3xl mx-auto">
    <h1 class="text-xl font-bold mb-4">编辑历史</h1>
    <div v-if="loading" class="text-center py-8 text-gray-500">加载中...</div>
    <div v-else-if="!edits.length" class="text-center py-8 text-gray-500 bg-white rounded-lg">暂无编辑记录</div>
    <div v-else class="space-y-3">
      <div v-for="edit in edits" :key="edit.id" class="bg-white rounded-lg shadow-sm border p-4">
        <div class="text-xs text-gray-400 mb-2">{{ new Date(edit.created_at).toLocaleString() }} · 编辑者 #{{ edit.user_id }}</div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <div class="text-xs text-gray-500 mb-1">修改前</div>
            <div class="text-sm text-gray-600 bg-gray-50 p-2 rounded line-clamp-3">{{ edit.old_content }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-500 mb-1">修改后</div>
            <div class="text-sm text-gray-600 bg-gray-50 p-2 rounded line-clamp-3">{{ edit.new_content }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
