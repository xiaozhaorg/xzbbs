<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const username = ref('')
const email = ref('')
const signature = ref('')
const emailVerified = ref(false)
const verifyingEmail = ref(false)
const error = ref('')
const success = ref('')

async function save() {
  try {
    const { updateUser } = await import('@/api/bbs')
    const updates: any = {}
    if (username.value) updates.username = username.value
    if (email.value) updates.email = email.value
    if (signature.value !== undefined) updates.signature = signature.value
    await updateUser(auth.user!.id, updates)
    success.value = '已更新'
    await auth.init()
  } catch (e: any) {
    error.value = e.msg || '更新失败'
  }
}

// Initialize form values
onMounted(() => {
  if (auth.user) {
    username.value = auth.user.username
    email.value = auth.user.email
    signature.value = (auth.user as any).signature || ''
    emailVerified.value = (auth.user as any).email_verified || false
  }
})

async function verifyEmail() {
  verifyingEmail.value = true
  try {
    const { requestEmailVerify } = await import('@/api/bbs')
    const res = await requestEmailVerify()
    success.value = `验证码: ${(res.data as any).token}（生产环境会发送邮件）`
  } catch (e: any) {
    error.value = e.msg || '验证请求失败'
  } finally {
    verifyingEmail.value = false
  }
}
</script>

<template>
  <div class="max-w-lg mx-auto">
    <h1 class="text-xl font-bold mb-4">个人设置</h1>
    <div class="bg-white rounded-lg shadow-sm border p-6 space-y-4">
      <div v-if="error" class="text-red-500 text-sm">{{ error }}</div>
      <div v-if="success" class="text-green-500 text-sm">{{ success }}</div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
        <input v-model="username" class="w-full border rounded-lg px-3 py-2 text-sm" />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
        <div class="flex items-center gap-2">
          <input v-model="email" type="email" class="flex-1 border rounded-lg px-3 py-2 text-sm" />
          <span v-if="emailVerified" class="text-xs text-green-600 bg-green-50 px-2 py-1 rounded">✓ 已验证</span>
          <button v-else @click="verifyEmail" :disabled="verifyingEmail" class="text-xs text-blue-600 border border-blue-600 px-2 py-1 rounded hover:bg-blue-50 disabled:opacity-50">
            {{ verifyingEmail ? '发送中...' : '验证邮箱' }}
          </button>
        </div>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">签名</label>
        <input v-model="signature" class="w-full border rounded-lg px-3 py-2 text-sm" maxlength="255" placeholder="一句话介绍自己" />
      </div>
      <button @click="save" class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
        保存
      </button>
    </div>
  </div>
</template>
