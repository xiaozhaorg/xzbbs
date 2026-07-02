<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getThread, createPost, updatePost, toggleFavorite, checkFavorite, type Thread as ThreadType, type Post as PostType } from '@/api/bbs'
import { useAuthStore } from '@/stores/auth'
import MarkdownIt from 'markdown-it'
import BBEditor from '@/components/BBEditor.vue'
import SmileyPicker from '@/components/SmileyPicker.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const threadId = Number(route.params.id)
const thread = ref<ThreadType | null>(null)
const firstPost = ref<PostType | null>(null)
const replies = ref<PostType[]>([])
const content = ref('')
const submitting = ref(false)
const loading = ref(true)
const isFavorited = ref(false)
const favoriting = ref(false)

// Edit state
const editingPostId = ref<number | null>(null)
const editContent = ref('')
const editingThreadTitle = ref(false)
const editTitleText = ref('')

const md = MarkdownIt({ html: false, linkify: true })

const levelTitle = (level: number) => {
  const titles = ['', '新手', '学徒', '入门', '小有所成', '登堂入室', '炉火纯青', '出类拔萃', '技压群雄', '一代宗师', '武林至尊']
  return titles[Math.min(level, 10)] || ''
}

const levelColor = (level: number) => {
  if (level >= 8) return 'text-purple-600'
  if (level >= 5) return 'text-blue-600'
  if (level >= 3) return 'text-green-600'
  return 'text-gray-400'
}

const canEdit = (post: PostType) => {
  if (!auth.isLoggedIn) return false
  return auth.userId === post.user_id || (auth.user && (auth.user as any).group_id <= 5)
}

onMounted(async () => {
  try {
    const res = await getThread(threadId)
    const data = res.data as { thread: ThreadType; first_post: PostType; replies: { items: PostType[] } }
    thread.value = data.thread
    firstPost.value = data.first_post
    replies.value = data.replies.items

    if (auth.isLoggedIn) {
      try {
        const favRes = await checkFavorite(threadId)
        isFavorited.value = (favRes.data as any).favorited
      } catch { /* ignore */ }
    }
  } finally {
    loading.value = false
  }
})

async function toggleFav() {
  if (!auth.isLoggedIn) return
  favoriting.value = true
  try {
    const res = await toggleFavorite(threadId)
    isFavorited.value = (res.data as any).favorited
  } finally {
    favoriting.value = false
  }
}

async function submitReply() {
  if (!content.value.trim() || submitting.value) return
  submitting.value = true
  try {
    await createPost(threadId, { content: content.value })
    content.value = ''
    await reload()
  } finally {
    submitting.value = false
  }
}

async function startEditPost(post: PostType) {
  editingPostId.value = post.id
  editContent.value = post.content
}

async function cancelEdit() {
  editingPostId.value = null
  editContent.value = ''
}

async function saveEdit(postId: number) {
  if (!editContent.value.trim()) return
  try {
    await updatePost(postId, { content: editContent.value })
    editingPostId.value = null
    await reload()
  } catch (e) {
    alert('编辑失败')
  }
}

function startEditTitle() {
  editTitleText.value = thread.value?.title || ''
  editingThreadTitle.value = true
}

async function saveTitle() {
  if (!editTitleText.value.trim()) return
  try {
    const { updateThread } = await import('@/api/bbs')
    await updateThread(threadId, { title: editTitleText.value })
    editingThreadTitle.value = false
    if (thread.value) thread.value.title = editTitleText.value
  } catch (e) {
    alert('修改失败')
  }
}

async function reload() {
  const res = await getThread(threadId)
  const data = res.data as { thread: ThreadType; first_post: PostType; replies: { items: PostType[] } }
  thread.value = data.thread
  firstPost.value = data.first_post
  replies.value = data.replies.items
}

function renderMd(text: string) {
  return md.render(text || '')
}

function handleImageUpload(file: File) {
  const reader = new FileReader()
  reader.onload = (e) => {
    content.value += `\n![${file.name}](${e.target?.result})\n`
  }
  reader.readAsDataURL(file)
}
</script>

<template>
  <div v-if="loading" class="text-center py-12 text-gray-500">加载中...</div>

  <div v-else-if="thread" class="space-y-4">
    <!-- Breadcrumb -->
    <div class="flex items-center justify-between">
      <div class="text-sm text-gray-500">
        <RouterLink to="/" class="hover:text-blue-600">首页</RouterLink>
        <span class="mx-2">/</span>
        <RouterLink :to="`/forum/${thread.forum?.id}`" class="hover:text-blue-600">{{ thread.forum?.name || '版块' }}</RouterLink>
      </div>
      <button
        v-if="auth.isLoggedIn"
        @click="toggleFav"
        :class="isFavorited ? 'text-yellow-500' : 'text-gray-400 hover:text-yellow-500'"
        class="text-sm flex items-center gap-1"
      >
        <span>{{ isFavorited ? '★ 已收藏' : '☆ 收藏' }}</span>
      </button>
    </div>

    <!-- First Post -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex items-center justify-between mb-2">
        <h1 class="text-xl font-bold text-gray-800">
          <span v-if="thread.is_top" class="text-red-500 mr-1">[置顶]</span>
          <span v-if="!editingThreadTitle">{{ thread.title }}</span>
          <input v-else v-model="editTitleText" class="border rounded px-2 py-1 text-lg w-full" @keyup.enter="saveTitle" @keyup.escape="editingThreadTitle = false" />
        </h1>
        <div v-if="canEdit(firstPost!) && !editingThreadTitle">
          <button @click="startEditTitle" class="text-xs text-gray-400 hover:text-blue-600">编辑标题</button>
        </div>
        <div v-if="editingThreadTitle" class="flex gap-2">
          <button @click="saveTitle" class="text-xs px-2 py-1 bg-blue-600 text-white rounded">保存</button>
          <button @click="editingThreadTitle = false" class="text-xs px-2 py-1 bg-gray-200 rounded">取消</button>
        </div>
      </div>

      <div v-if="editingPostId !== firstPost?.id" class="flex items-center gap-3 mb-4 pb-3 border-b">
        <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center text-blue-600 font-bold text-sm">
          {{ firstPost?.user?.username?.[0]?.toUpperCase() || '?' }}
        </div>
        <div class="text-sm">
          <div class="font-medium text-gray-700 flex items-center gap-1">
            {{ firstPost?.user?.username || '未知' }}
            <span v-if="firstPost?.user?.level" :class="levelColor(firstPost.user.level)" class="text-xs">{{ levelTitle(firstPost.user.level) }}</span>
          </div>
          <div class="text-xs text-gray-400">{{ new Date(firstPost?.created_at || '').toLocaleString() }}</div>
        </div>
        <button v-if="canEdit(firstPost!)" @click="startEditPost(firstPost!)" class="text-xs text-gray-400 hover:text-blue-600 ml-auto">编辑</button>
      </div>

      <!-- Edit form for first post -->
      <div v-if="editingPostId === firstPost?.id" class="mb-4">
        <div class="text-xs text-gray-400 mb-1">编辑帖子</div>
        <BBEditor v-model="editContent" :rows="6" @image-upload="(f: File) => { editContent += '\n![](' + URL.createObjectURL(f) + ')\n' }">
          <template #footer="{ charCount }">
            <div class="px-3 py-1.5 bg-gray-50 border-t text-xs text-gray-400 flex justify-between">
              <span></span>
              <span>{{ charCount }} 字</span>
            </div>
          </template>
        </BBEditor>
        <div class="flex gap-2 mt-2">
          <button @click="saveEdit(firstPost!.id)" class="px-4 py-1.5 bg-blue-600 text-white rounded text-sm">保存</button>
          <button @click="cancelEdit" class="px-4 py-1.5 bg-gray-200 rounded text-sm">取消</button>
        </div>
      </div>
      <div v-else class="prose prose-sm max-w-none text-gray-700" v-html="renderMd(firstPost?.content || '')"></div>
    </div>

    <!-- Replies -->
    <div v-if="replies.length" class="space-y-3">
      <h3 class="text-sm font-medium text-gray-500">共 {{ thread.posts }} 条回复</h3>
      <div v-for="p in replies" :key="p.id" class="bg-white rounded-lg shadow-sm border p-4">
        <div v-if="editingPostId !== p.id" class="flex items-center justify-between mb-2">
          <div class="flex items-center gap-2">
            <div class="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center text-gray-600 font-bold text-xs">
              {{ p.user?.username?.[0]?.toUpperCase() || '?' }}
            </div>
            <span class="text-sm font-medium">{{ p.user?.username || '未知' }}</span>
            <span v-if="p.user?.level" :class="levelColor(p.user.level)" class="text-xs">{{ levelTitle(p.user.level) }}</span>
            <span class="text-xs text-gray-400">{{ new Date(p.created_at).toLocaleString() }}</span>
          </div>
          <div class="flex gap-3">
            <RouterLink :to="`/posts/${p.id}/edits`" class="text-xs text-gray-400 hover:text-blue-600">编辑历史</RouterLink>
            <button v-if="canEdit(p)" @click="startEditPost(p)" class="text-xs text-gray-400 hover:text-blue-600">编辑</button>
          </div>
        </div>
        <!-- Edit form -->
        <div v-if="editingPostId === p.id" class="mb-2">
          <div class="text-xs text-gray-400 mb-1">编辑回复</div>
          <BBEditor v-model="editContent" :rows="4" @image-upload="(f: File) => { editContent += '\n![](' + URL.createObjectURL(f) + ')\n' }">
            <template #footer="{ charCount }">
              <div class="px-3 py-1.5 bg-gray-50 border-t text-xs text-gray-400 flex justify-between">
                <span></span>
                <span>{{ charCount }} 字</span>
              </div>
            </template>
          </BBEditor>
          <div class="flex gap-2 mt-2">
            <button @click="saveEdit(p.id)" class="px-4 py-1.5 bg-blue-600 text-white rounded text-sm">保存</button>
            <button @click="cancelEdit" class="px-4 py-1.5 bg-gray-200 rounded text-sm">取消</button>
          </div>
        </div>
        <div v-else class="prose prose-sm max-w-none text-gray-700 ml-8" v-html="renderMd(p.content)"></div>
      </div>
    </div>

    <!-- Reply Form -->
    <div v-if="auth.isLoggedIn && !thread.is_closed" class="bg-white rounded-lg shadow-sm border p-4">
      <h3 class="text-sm font-medium text-gray-700 mb-2">回复</h3>
      <BBEditor
        v-model="content"
        placeholder="写下你的回复... (支持 Markdown)"
        :rows="4"
        @submit="submitReply"
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

    <div v-else-if="thread.is_closed" class="text-center py-4 text-gray-500 bg-white rounded-lg">
      该主题已关闭，无法回复
    </div>
    <div v-else class="text-center py-4">
      <RouterLink to="/login" class="text-blue-600 hover:underline">登录后回复</RouterLink>
    </div>
  </div>
</template>
