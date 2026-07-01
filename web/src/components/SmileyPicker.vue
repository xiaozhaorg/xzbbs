<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getSmilies, type Smiley } from '@/api/bbs'

const show = defineModel<boolean>('show', { default: false })
const smilies = ref<Smiley[]>([])

// Expose insert for backward compatibility (direct textarea insertion)
function insert(code: string) {
  emit('select', code)
}

const emit = defineEmits<{
  select: [code: string]
}>()

onMounted(async () => {
  const res = await getSmilies()
  smilies.value = res.data || []
})

// Watch for show changes and load if needed
watch(show, (v) => {
  if (v && !smilies.value.length) {
    getSmilies().then(res => { smilies.value = res.data || [] })
  }
})
</script>

<template>
  <div v-if="show" class="bg-white border rounded-lg shadow-lg p-3 z-50">
    <div class="grid grid-cols-8 gap-1">
      <button
        v-for="s in smilies"
        :key="s.id"
        @click="insert(s.code)"
        class="w-8 h-8 flex items-center justify-center hover:bg-gray-100 rounded text-lg cursor-pointer"
        :title="s.code"
      >
        {{ s.code }}
      </button>
    </div>
  </div>
</template>
