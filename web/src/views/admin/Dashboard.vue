<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getAdminStats, getAdminGroups, listAdminUsers, listIPBans, banIP as banIPApi, unbanIP as unbanIPApi, getSmilies } from '@/api/bbs'
import type { User, Forum } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

onMounted(() => {
  if (!auth.isLoggedIn) {
    router.push('/login')
  }
})

// Stats
const stats = ref({ users: 0, threads: 0, posts: 0, today_posts: 0 })

// Tabs
type Tab = 'stats' | 'users' | 'forums' | 'online' | 'ipbans' | 'smilies'
const tab = ref<Tab>('stats')

// Users
const users = ref<User[]>([])
const usersTotal = ref(0)
const usersPage = ref(1)

// Groups (placeholder)
const groups = ref<any[]>([])

// Forums
const forums = ref<Forum[]>([])

// Online users
const onlineUsers = ref<any[]>([])

// IP bans
const ipBans = ref<any[]>([])
const banIP = ref('')
const banReason = ref('')
const banExpire = ref(0)

// Smilies
const smilies = ref<any[]>([])

onMounted(async () => {
  try {
    const sr = await getAdminStats()
    stats.value = sr.data
    const gr = await getAdminGroups()
    groups.value = gr.data
    await loadUsers()
    loadOnlineUsers()
    loadIPBans()
    loadSmilies()
  } catch (e) {
    // not admin
    router.push('/')
  }
})

async function loadIPBans() {
  try {
    const res = await listIPBans(1)
    ipBans.value = (res.data as any)?.items || []
  } catch { /* ignore */ }
}

async function doBanIP() {
  if (!banIP.value || !banReason.value) return
  await banIPApi(banIP.value, banReason.value, banExpire.value)
  banIP.value = ''; banReason.value = ''; banExpire.value = 0
  loadIPBans()
}

async function doUnban(id: number) {
  await unbanIPApi(id)
  loadIPBans()
}

async function loadSmilies() {
  try {
    const res = await getSmilies()
    smilies.value = res.data || []
  } catch { /* ignore */ }
}

async function loadOnlineUsers() {
  try {
    const res = await getOnlineUsers()
    onlineUsers.value = res.data
  } catch { /* ignore */ }
}

async function loadUsers(page = 1) {
  usersPage.value = page
  const res = await listAdminUsers(page)
  users.value = res.data.items
  usersTotal.value = res.data.total
}
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">管理后台</h1>

    <div class="flex gap-2 mb-6 border-b overflow-x-auto">
      <button
        v-for="t in ['stats', 'users', 'forums', 'online', 'ipbans', 'smilies'] as Tab[]"
        :key="t"
        @click="tab = t"
        :class="tab === t ? 'border-b-2 border-blue-600 text-blue-600' : 'text-gray-500'"
        class="px-4 py-2 text-sm font-medium whitespace-nowrap"
      >
        {{ t === 'stats' ? '📊 统计' : t === 'users' ? '👤 用户' : t === 'forums' ? '📁 版块' : t === 'online' ? '🟢 在线' : t === 'ipbans' ? '🚫 IP封禁' : '😀 表情' }}
      </button>
    </div>

    <!-- Stats -->
    <div v-if="tab === 'stats'" class="grid grid-cols-4 gap-4">
      <div class="bg-white rounded-lg border p-4">
        <div class="text-2xl font-bold">{{ stats.users }}</div>
        <div class="text-sm text-gray-500">用户总数</div>
      </div>
      <div class="bg-white rounded-lg border p-4">
        <div class="text-2xl font-bold">{{ stats.threads }}</div>
        <div class="text-sm text-gray-500">主题总数</div>
      </div>
      <div class="bg-white rounded-lg border p-4">
        <div class="text-2xl font-bold">{{ stats.posts }}</div>
        <div class="text-sm text-gray-500">帖子总数</div>
      </div>
      <div class="bg-white rounded-lg border p-4">
        <div class="text-2xl font-bold text-green-600">{{ stats.today_posts }}</div>
        <div class="text-sm text-gray-500">今日帖子</div>
      </div>
    </div>

    <!-- Users -->
    <div v-if="tab === 'users'" class="bg-white rounded-lg border overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-gray-50 border-b">
          <tr>
            <th class="px-4 py-2 text-left">ID</th>
            <th class="px-4 py-2 text-left">用户名</th>
            <th class="px-4 py-2 text-left">邮箱</th>
            <th class="px-4 py-2 text-left">主题</th>
            <th class="px-4 py-2 text-left">回帖</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id" class="border-b hover:bg-gray-50">
            <td class="px-4 py-2">{{ u.id }}</td>
            <td class="px-4 py-2">{{ u.username }}</td>
            <td class="px-4 py-2 text-gray-500">{{ u.email }}</td>
            <td class="px-4 py-2">{{ u.threads }}</td>
            <td class="px-4 py-2">{{ u.posts }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="usersTotal > 20" class="flex justify-center gap-2 p-3">
        <button
          v-for="p in Math.ceil(usersTotal / 20)"
          :key="p"
          @click="loadUsers(p)"
          :class="p === usersPage ? 'bg-blue-600 text-white' : 'bg-gray-100'"
          class="px-3 py-1 rounded text-sm"
        >
          {{ p }}
        </button>
      </div>
    </div>

    <!-- Forums -->
    <div v-if="tab === 'forums'" class="bg-white rounded-lg border p-4 text-gray-500">
      版块管理（开发中...）
    </div>

    <!-- Online Users -->
    <div v-if="tab === 'online'" class="bg-white rounded-lg border overflow-hidden">
      <h3 class="px-4 py-2 bg-gray-50 border-b font-medium text-sm">在线用户</h3>
      <table class="w-full text-sm">
        <thead class="bg-gray-50 border-b">
          <tr>
            <th class="px-4 py-2 text-left">用户</th>
            <th class="px-4 py-2 text-left">IP</th>
            <th class="px-4 py-2 text-left">最后活跃</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in onlineUsers" :key="u.id" class="border-b hover:bg-gray-50">
            <td class="px-4 py-2">{{ u.username }}</td>
            <td class="px-4 py-2 text-gray-500">{{ u.ip }}</td>
            <td class="px-4 py-2 text-gray-500">{{ new Date(u.last_active).toLocaleTimeString() }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- IP Bans -->
    <div v-if="tab === 'ipbans'" class="bg-white rounded-lg border overflow-hidden">
      <div class="p-4 border-b flex gap-2">
        <input v-model="banIP" placeholder="IP地址" class="border rounded px-3 py-1.5 text-sm flex-1" />
        <input v-model="banReason" placeholder="封禁原因" class="border rounded px-3 py-1.5 text-sm flex-1" />
        <input v-model.number="banExpire" type="number" placeholder="小时(0=永久)" class="border rounded px-3 py-1.5 text-sm w-32" />
        <button @click="doBanIP" class="px-4 py-1.5 bg-red-600 text-white rounded text-sm hover:bg-red-700">封禁</button>
      </div>
      <table class="w-full text-sm">
        <thead class="bg-gray-50 border-b">
          <tr>
            <th class="px-4 py-2 text-left">IP</th>
            <th class="px-4 py-2 text-left">原因</th>
            <th class="px-4 py-2 text-left">到期时间</th>
            <th class="px-4 py-2 text-left">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ban in ipBans" :key="ban.id" class="border-b hover:bg-gray-50">
            <td class="px-4 py-2">{{ ban.ip }}</td>
            <td class="px-4 py-2 text-gray-500">{{ ban.reason }}</td>
            <td class="px-4 py-2 text-gray-500 text-xs">{{ ban.expire_at ? new Date(ban.expire_at).toLocaleString() : '永久' }}</td>
            <td class="px-4 py-2">
              <button @click="doUnban(ban.id)" class="text-red-600 text-xs hover:underline">解封</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Smilies -->
    <div v-if="tab === 'smilies'" class="bg-white rounded-lg border p-4">
      <div class="grid grid-cols-6 gap-3">
        <div v-for="s in smilies" :key="s.id" class="text-center p-2 border rounded">
          <div class="text-2xl">{{ s.code }}</div>
          <div class="text-xs text-gray-400">{{ s.image }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
