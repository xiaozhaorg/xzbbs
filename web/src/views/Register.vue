<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()
const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')

async function submit() {
  error.value = ''
  try {
    await auth.register(username.value, email.value, password.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.msg || '注册失败'
  }
}
</script>

<template>
  <div class="max-w-sm mx-auto mt-12">
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h1 class="text-xl font-bold text-center mb-6">注册</h1>
      <div v-if="error" class="text-red-500 text-sm mb-3">{{ error }}</div>
      <div class="space-y-3">
        <input
          v-model="username"
          class="w-full border rounded-lg px-3 py-2 text-sm"
          placeholder="用户名 (2-32字符)"
        />
        <input
          v-model="email"
          type="email"
          class="w-full border rounded-lg px-3 py-2 text-sm"
          placeholder="邮箱"
        />
        <input
          v-model="password"
          type="password"
          class="w-full border rounded-lg px-3 py-2 text-sm"
          placeholder="密码 (至少6位)"
          @keyup.enter="submit"
        />
        <button
          @click="submit"
          class="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          注册
        </button>
        <p class="text-center text-sm text-gray-500">
          已有账号？
          <RouterLink to="/login" class="text-blue-600">登录</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>
