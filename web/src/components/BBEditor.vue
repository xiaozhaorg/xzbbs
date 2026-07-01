<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'

const props = defineModel<string>({ required: true })
const props_ = defineProps<{
  placeholder?: string
  rows?: number
  readonly?: boolean
  showPreview?: boolean
  autoFocus?: boolean
  serverUpload?: boolean
}>()

const emit = defineEmits<{
  submit: []
  imageUpload: [file: File]
  uploadDone: [result: { url: string }]
}>()

const textarea = ref<HTMLTextAreaElement>()
const showPreview = ref(false)
const showSmilies = ref(false)
const previewHtml = ref('')
const charCount = computed(() => props.value.length)
const wordCount = computed(() => props.value.trim() ? props.value.trim().split(/\s+/).length : 0)
const uploading = ref(false)

// Insert markdown syntax at cursor
function insert(prefix: string, suffix = '') {
  const el = textarea.value
  if (!el) return
  const start = el.selectionStart
  const end = el.selectionEnd
  const text = props.value
  const selected = text.slice(start, end)
  const before = text.slice(0, start)
  const after = text.slice(end)

  // Find line start for block-level elements
  const lineStart = before.lastIndexOf('\n') + 1
  const lineText = before.slice(lineStart)

  if (['**', '~~', '`'].includes(prefix) && !suffix) {
    // Inline: wrap selection
    props.value = before + prefix + selected + prefix + after
  } else if (['> ', '- ', '1. ', '```'].includes(prefix) && !suffix) {
    // Block: insert at line start
    props.value = before.slice(0, lineStart) + prefix + lineText + after
  } else {
    props.value = before + prefix + selected + suffix + after
  }

  nextTick(() => {
    el.focus()
    const newPos = lineStart + prefix.length + selected.length
    el.setSelectionRange(newPos, newPos)
  })
}

function insertLink() {
  const url = prompt('请输入链接地址:')
  if (url) insert('[', '](' + url + ')')
}

function insertImage() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = async () => {
    const file = input.files?.[0]
    if (!file) return

    if (props_.serverUpload) {
      uploading.value = true
      try {
        const { uploadFile } = await import('@/api/bbs')
        const res = await uploadFile(file)
        const url = (res.data as any)?.url || ''
        const pos = textarea.value?.selectionStart ?? props.value.length
        const before = props.value.slice(0, pos)
        const after = props.value.slice(pos)
        props.value = before + `\n![${file.name}](${url})\n` + after
        emit('uploadDone', { url })
      } catch {
        alert('图片上传失败')
      } finally {
        uploading.value = false
      }
    } else {
      emit('imageUpload', file)
    }
  }
  input.click()
}

function insertSmiley(code: string) {
  props.value += ' ' + code + ' '
  nextTick(() => textarea.value?.focus())
}

function updatePreview() {
  // Import markdown-it dynamically to avoid circular deps
  import('markdown-it').then(m => {
    const md = m.default({ html: false, linkify: true, typographer: true })
    previewHtml.value = md.render(props.value || '')
  }).catch(() => {
    previewHtml.value = '<p>预览暂不可用</p>'
  })
}

function togglePreview() {
  showPreview.value = !showPreview.value
  if (showPreview.value) updatePreview()
}

const toolbar = [
  { icon: 'B', title: '粗体 (Ctrl+B)', action: () => insert('**') },
  { icon: 'I', title: '斜体 (Ctrl+I)', action: () => insert('*') },
  { icon: 'S', title: '删除线', action: () => insert('~~') },
  { icon: '<>', title: '代码', action: () => insert('`') },
  { icon: 'H', title: '标题', action: () => insert('## ') },
  { icon: '🔗', title: '链接', action: insertLink },
  { icon: '❝', title: '引用', action: () => insert('> ') },
  { icon: '•', title: '列表', action: () => insert('- ') },
  { icon: '1.', title: '编号', action: () => insert('1. ') },
  { icon: '---', title: '分割线', action: () => insert('\n---\n') },
  { icon: '🖼️', title: '图片', action: insertImage },
  { icon: '👀', title: '预览', action: togglePreview },
]

function onKeydown(e: KeyboardEvent) {
  if (e.ctrlKey || e.metaKey) {
    if (e.key === 'b') { e.preventDefault(); insert('**') }
    if (e.key === 'i') { e.preventDefault(); insert('*') }
    if (e.key === 'k') { e.preventDefault(); insertLink() }
  }
  if (e.key === 'Tab') {
    e.preventDefault()
    insert('  ')
  }
}

onMounted(() => {
  if (props_.autoFocus) textarea.value?.focus()
})
</script>

<template>
  <div class="bb-editor border rounded-lg overflow-hidden">
    <!-- Toolbar -->
    <div class="flex flex-wrap items-center gap-1 px-2 py-1.5 bg-gray-50 border-b text-xs">
      <button v-for="btn in toolbar" :key="btn.title"
        type="button" @click="btn.action"
        :title="btn.title"
        class="w-7 h-7 flex items-center justify-center rounded hover:bg-gray-200 text-gray-600 font-medium transition-colors"
        :class="btn.icon === 'B' ? 'font-bold' : btn.icon === 'I' ? 'italic' : ''">
        {{ btn.icon }}
      </button>
      <span class="flex-1"></span>
      <button type="button" @click="showSmilies = !showSmilies"
        class="px-2 py-1 rounded hover:bg-gray-200 text-gray-500">
        😀 表情
      </button>
      <span class="text-gray-300">|</span>
      <span v-if="uploading" class="text-blue-500">上传中...</span>
      <span v-else class="text-gray-400">{{ charCount }}字</span>
    </div>

    <!-- Smilies overlay -->
    <div v-if="showSmilies" class="relative">
      <div class="absolute z-10 top-1 left-0">
        <slot name="smilies-picker" :onSelect="insertSmiley" />
      </div>
    </div>

    <!-- Editor / Preview -->
    <div class="flex">
      <div v-show="!showPreview" class="flex-1">
        <textarea
          ref="textarea"
          v-model="props.value"
          :readonly="props_.readonly"
          :placeholder="props_.placeholder || '支持 Markdown 格式...'"
          :rows="props_.rows || 8"
          class="w-full p-3 text-sm resize-y focus:outline-none focus:ring-0 border-0"
          style="min-height: 120px"
          @keydown="onKeydown"
        ></textarea>
      </div>
      <div v-if="showPreview" class="flex-1 p-3 bg-white min-h-[120px] prose prose-sm max-w-none"
        v-html="previewHtml">
      </div>
    </div>

    <slot name="footer" :charCount :wordCount :showPreview>
      <div v-if="!$slots.footer" class="px-3 py-1.5 bg-gray-50 border-t text-xs text-gray-400 flex justify-between">
        <span>Ctrl+B 粗体 · Ctrl+I 斜体 · Tab 缩进</span>
        <span>{{ wordCount }} 词</span>
      </div>
    </slot>
  </div>
</template>
