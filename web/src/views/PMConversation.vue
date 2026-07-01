<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPMConversation, sendPM, markPMRead, type PrivateMessage } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const otherId = Number(route.params.otherId)
const messages = ref<PrivateMessage[]>([])
const newMsg = ref('')
const sending = ref(false)
const loading = ref(true)

async function load() {
  loading.value = true
  try {
    const res = await getPMConversation(otherId, 1)
    messages.value = (res.data as any).items || []
    await markPMRead(otherId)
  } finally {
    loading.value = false
  }
}

async function send() {
  if (!newMsg.value.trim() || sending.value) return
  sending.value = true
  try {
    await sendPM(otherId, newMsg.value.trim())
    newMsg.value = ''
    await load()
  } finally {
    sending.value = false
  }
}

onMounted(() => { if (auth.isLoggedIn) { load() } else { router.push('/login') } })
</script>

<template>
  <div class="max-w-2xl mx-auto">
    <div class="flex items-center gap-2 mb-4">
      <button @click="router.push('/pm')" class="text-gray-500 hover:text-blue-600 text-sm">← 返回私信列表</button>
      <h1 class="text-lg font-bold">与用户 #{{ otherId }} 的对话</h1>
    </div>

    <div v-if="loading" class="text-center py-8 text-gray-500">加载中...</div>

    <div v-else-if="!messages.length" class="text-center py-8 text-gray-500 bg-white rounded-lg">
      暂无消息
    </div>

    <div v-else class="space-y-3 mb-4">
      <div v-for="msg in messages" :key="msg.id" class="p-3 rounded-lg"
        :class="msg.sender_id === auth.user?.id ? 'bg-blue-50 ml-8' : 'bg-gray-100 mr-8'">
        <p class="text-sm text-gray-700">{{ msg.content }}</p>
        <div class="text-xs text-gray-400 mt-1">{{ new Date(msg.created_at).toLocaleString() }}</div>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow-sm border p-3">
      <textarea v-model="newMsg" rows="2" class="w-full border rounded-lg p-2 text-sm resize-none" placeholder="输入私信..." @keydown.enter.exact.prevent="send"></textarea>
      <div class="flex justify-end mt-2">
        <button @click="send" :disabled="sending || !newMsg.trim()" class="px-4 py-1.5 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:opacity-50">
          {{ sending ? '发送中...' : '发送' }}
        </button>
      </div>
    </div>
  </div>
</template>
