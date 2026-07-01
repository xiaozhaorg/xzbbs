<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getPMList, getPMConversation, sendPM, markPMRead, type PrivateMessage } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const conversations = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  if (!auth.isLoggedIn) { router.push('/login'); return }
  try {
    const res = await getPMList()
    conversations.value = res.data || []
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <h1 class="text-xl font-bold mb-4">私信</h1>
    <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>
    <div v-else-if="!conversations.length" class="text-center py-12 text-gray-500 bg-white rounded-lg">
      暂无私信对话
    </div>
    <div v-else class="space-y-2">
      <div v-for="conv in conversations" :key="conv.other_id" class="bg-white rounded-lg shadow-sm border p-4 flex items-center justify-between">
        <div>
          <div class="font-medium text-gray-800">用户 #{{ conv.other_id }}</div>
          <div class="text-sm text-gray-500 line-clamp-1">{{ conv.last_message || '' }}</div>
        </div>
        <RouterLink :to="`/pm/${conv.other_id}`" class="text-blue-600 text-sm hover:underline">查看</RouterLink>
      </div>
    </div>
  </div>
</template>
