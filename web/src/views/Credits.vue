<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { getCreditLogs } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const logs = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const loading = ref(true)

async function load(p = 1) {
  loading.value = true
  page.value = p
  try {
    const res = await getCreditLogs(p)
    logs.value = (res.data as any)?.items || []
    total.value = (res.data as any)?.total || 0
  } finally {
    loading.value = false
  }
}

const pages = computed(() => total.value > 0 ? Math.ceil(total.value / 20) : 0)

onMounted(() => load())
</script>

<template>
  <div class="max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold mb-4">💎 积分明细</h1>
    <div class="bg-white rounded-lg shadow-sm border p-4 mb-4">
      <div class="text-sm text-gray-500">当前积分</div>
      <div class="text-3xl font-bold text-blue-600">{{ auth.user?.credits || 0 }}</div>
    </div>
    <div v-if="loading" class="text-center py-8 text-gray-500">加载中...</div>
    <div v-else-if="!logs.length" class="text-center py-8 text-gray-500 bg-white rounded-lg">暂无记录</div>
    <div v-else class="space-y-2">
      <div v-for="log in logs" :key="log.id" class="bg-white rounded-lg border p-3 flex justify-between items-center">
        <div>
          <div class="text-sm text-gray-700">{{ log.reason }}</div>
          <div class="text-xs text-gray-400">{{ new Date(log.created_at).toLocaleString() }}</div>
        </div>
        <div :class="log.amount > 0 ? 'text-green-600' : 'text-red-600'" class="font-bold text-sm">
          {{ log.amount > 0 ? '+' : '' }}{{ log.amount }}
        </div>
      </div>
    </div>
    <div v-if="pages > 1" class="flex justify-center gap-2 mt-4">
      <button v-for="p in pages" :key="p" @click="load(p)" :class="p === page ? 'bg-blue-600 text-white' : 'bg-white border'" class="px-3 py-1 rounded text-sm">{{ p }}</button>
    </div>
  </div>
</template>
