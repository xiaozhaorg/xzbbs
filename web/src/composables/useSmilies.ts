import { ref } from 'vue'

export interface Smiley {
  id: number
  code: string
  image: string
  sort: number
}

// Default smilies (matches server-side defaults)
export const defaultSmilies: Smiley[] = [
  { id: 1, code: ':)', image: '/uploads/smilies/1.png', sort: 0 },
  { id: 2, code: ':(', image: '/uploads/smilies/2.png', sort: 0 },
  { id: 3, code: ':D', image: '/uploads/smilies/3.png', sort: 0 },
  { id: 4, code: ':o', image: '/uploads/smilies/4.png', sort: 0 },
  { id: 5, code: ':P', image: '/uploads/smilies/5.png', sort: 0 },
  { id: 6, code: ':love:', image: '/uploads/smilies/6.png', sort: 0 },
  { id: 7, code: ':think:', image: '/uploads/smilies/7.png', sort: 0 },
  { id: 8, code: ':angry:', image: '/uploads/smilies/8.png', sort: 0 },
  { id: 9, code: ':cool:', image: '/uploads/smilies/9.png', sort: 0 },
  { id: 10, code: ':cry:', image: '/uploads/smilies/10.png', sort: 0 },
  { id: 11, code: ':ok:', image: '/uploads/smilies/11.png', sort: 0 },
  { id: 12, code: ':no:', image: '/uploads/smilies/12.png', sort: 0 },
]

export function useSmilies() {
  const smilies = ref<Smiley[]>([])
  const loaded = ref(false)

  async function load() {
    if (loaded.value) return
    try {
      const res = await fetch('/api/smilies')
      if (res.ok) {
        const json = await res.json()
        smilies.value = json.data || defaultSmilies
      } else {
        smilies.value = defaultSmilies
      }
    } catch {
      smilies.value = defaultSmilies
    }
    loaded.value = true
  }

  return { smilies, defaultSmilies, load }
}
