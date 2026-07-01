<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const account = ref('')
const password = ref('')
const error = ref('')

async function submit() {
  error.value = ''
  try {
    await auth.login(account.value, password.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.msg || '登录失败'
  }
}
</script>

<template>
  <div class="max-w-sm mx-auto mt-12">
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h1 class="text-xl font-bold text-center mb-6">登录</h1>
      <div v-if="error" class="text-red-500 text-sm mb-3">{{ error }}</div>
      <div class="space-y-3">
        <input
          v-model="account"
          class="w-full border rounded-lg px-3 py-2 text-sm"
          placeholder="邮箱或用户名"
          @keyup.enter="submit"
        />
        <input
          v-model="password"
          type="password"
          class="w-full border rounded-lg px-3 py-2 text-sm"
          placeholder="密码"
          @keyup.enter="submit"
        />
        <button
          @click="submit"
          class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          登录
        </button>
        <p class="text-center text-sm text-gray-500">
          还没有账号？
          <RouterLink to="/register" class="text-blue-600">注册</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>
