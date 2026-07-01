<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { createThread } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'
import BBEditor from '@/components/BBEditor.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const fid = Number(route.params.fid)
const title = ref('')
const content = ref('')
const submitting = ref(false)

async function submit() {
  if (!title.value.trim() || !content.value.trim()) return
  submitting.value = true
  try {
    await createThread(fid, title.value, content.value)
    router.push(`/forum/${fid}`)
  } finally {
    submitting.value = false
  }
}

function handleImageUpload(file: File) {
  // BBEditor image upload placeholder
  const reader = new FileReader()
  reader.onload = (e) => {
    content.value += `\n![${file.name}](${e.target?.result})\n`
  }
  reader.readAsDataURL(file)
}
</script>

<template>
  <div class="max-w-3xl mx-auto">
    <h1 class="text-xl font-bold mb-4">发布新帖</h1>
    <div class="bg-white rounded-lg shadow-sm border p-6 space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">标题</label>
        <input
          v-model="title"
          class="w-full border rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          placeholder="请输入标题"
          maxlength="200"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">内容</label>
        <BBEditor
          v-model="content"
          placeholder="写下你的帖子内容... (支持 Markdown)"
          :rows="10"
          auto-focus
          @submit="submit"
          @image-upload="handleImageUpload"
        >
          <template #smilies-picker="{ onSelect }">
            <SmileyPicker :show="true" @select="onSelect" />
          </template>
          <template #footer="{ charCount, wordCount }">
            <div class="px-3 py-1.5 bg-gray-50 border-t text-xs text-gray-400 flex justify-between">
              <span></span>
              <span>{{ charCount }} 字 · {{ wordCount }} 词</span>
            </div>
          </template>
        </BBEditor>
      </div>
      <div class="flex justify-between">
        <button @click="router.back()" class="px-4 py-2 text-gray-600 border rounded-lg hover:bg-gray-50">
          取消
        </button>
        <button
          @click="submit"
          :disabled="submitting || !title.trim() || !content.trim()"
          class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          {{ submitting ? '发布中...' : '发布' }}
        </button>
      </div>
    </div>
  </div>
</template>
