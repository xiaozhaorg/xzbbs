import api from '@/api'

export interface User {
  id: number
  username: string
  email: string
  group_id: number
  avatar: string
  threads: number
  posts: number
  credits: number
  level?: number
  signature?: string
  email_verified?: boolean
  created_at: string
}

export interface Forum {
  id: number
  name: string
  description: string
  icon: string
  sort_order: number
  threads: number
  posts: number
  created_at: string
}

export interface Thread {
  id: number
  forum_id: number
  user_id: number
  title: string
  is_top: number
  is_closed: boolean
  views: number
  posts: number
  last_reply_at: string | null
  created_at: string
  user?: { username: string; avatar: string }
  forum?: { name: string }
}

export interface Post {
  id: number
  thread_id: number
  user_id: number
  is_first: boolean
  content: string
  content_type: number
  created_at: string
  user?: { username: string; avatar: string }
}

export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  pages: number
}

// Smilies
export interface Smiley {
  id: number
  code: string
  image: string
  sort: number
}

// Auth
export const register = (data: { username: string; email: string; password: string }) =>
  api.post('/auth/register', data)

export const login = (data: { account: string; password: string }) =>
  api.post('/auth/login', data)

export const getMe = () => api.get('/auth/me')

// Forums
export const getForums = () => api.get('/forums')

export const getForum = (id: number) => api.get(`/forums/${id}`)

// Threads
export const getThread = (id: number) => api.get(`/threads/${id}`)
export const createThread = (data: { forum_id: number; title: string; content: string }) =>
  api.post('/threads', data)
export const updateThread = (id: number, data: { title?: string }) =>
  api.put(`/threads/${id}`, data)
export const deleteThread = (id: number) => api.delete(`/threads/${id}`)
export const getForumThreads = (forumId: number, order = 'reply', page = 1) =>
  api.get(`/forums/${forumId}/threads?order=${order}&page=${page}`)
export const listThreads = (forumId: number, order = 'reply', page = 1, pageSize = 20) =>
  api.get(`/forums/${forumId}/threads?order=${order}&page=${page}&page_size=${pageSize}`)

// Posts
export const createPost = (threadId: number, data: { content: string }) =>
  api.post(`/threads/${threadId}/posts`, data)
export const updatePost = (id: number, data: { content: string }) =>
  api.put(`/posts/${id}`, data)
export const deletePost = (id: number) => api.delete(`/posts/${id}`)

// Upload
export const uploadFile = (file: File) => {
  const form = new FormData()
  form.append('file', file)
  return api.post('/attachments', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

// Users
export const getUser = (id: number) => api.get(`/users/${id}`)
export const updateUser = (id: number, data: { username?: string; email?: string }) =>
  api.put(`/users/${id}`, data)

// Admin
export const getAdminStats = () => api.get('/admin/stats')
export const getAdminGroups = () => api.get('/admin/groups')
export const listAdminUsers = (page = 1, search = '') =>
  api.get(`/admin/users?page=${page}&search=${search}`)
export const listIPBans = (page = 1) => api.get(`/admin/ip-bans?page=${page}`)
export const banIP = (ip: string, reason: string, hours: number) =>
  api.post('/admin/ip-bans', { ip, reason, expire_hours: hours })
export const unbanIP = (id: number) => api.delete(`/admin/ip-bans/${id}`)

// Search
export const searchThreads = (q: string, forumId?: number, page = 1) =>
  api.get(`/search?q=${encodeURIComponent(q)}&forum_id=${forumId || 0}&page=${page}`)

// Favorites
export const toggleFavorite = (threadId: number) =>
  api.post(`/favorites/threads/${threadId}`)
export const getFavorites = (page = 1) =>
  api.get(`/favorites/threads?page=${page}`)
export const checkFavorite = (threadId: number) =>
  api.get(`/favorites/threads/${threadId}/check`)

// Notifications
export const getNotifications = (unreadOnly = false, page = 1) =>
  api.get(`/notifications?unread=${unreadOnly ? 1 : 0}&page=${page}`)
export const getUnreadCount = () =>
  api.get('/notifications/unread-count')
export const markNotificationsRead = (ids: number[]) =>
  api.post('/notifications/read', { ids })
export const markAllNotificationsRead = () =>
  api.post('/notifications/read-all')

// Online users
export const getOnlineUsers = () =>
  api.get('/online')

// Private Messages
export const sendPM = (receiverId: number, content: string) =>
  api.post('/pms', { receiver_id: receiverId, content })
export const getPMConversation = (otherId: number, page = 1) =>
  api.get(`/pms/conversations/${otherId}?page=${page}`)
export const getPMList = () =>
  api.get('/pms')
export const getPMUnreadCount = () =>
  api.get('/pms/unread-count')
export const markPMRead = (otherId: number) =>
  api.post(`/pms/read/${otherId}`)

// Smilies
export const getSmilies = () =>
  api.get('/smilies')

// Email verification
export const requestEmailVerify = () => api.post('/email/verify-request')
export const confirmEmailVerify = (token: string) => api.post('/email/verify-confirm', { token })

// Post edits
export const getPostEdits = (postId: number) =>
  api.get(`/posts/${postId}/edits`)

// Credits
export const getCreditLogs = (page = 1) =>
  api.get(`/credits/logs?page=${page}`)

// User threads/posts
export const getUserThreads = (userId: number, page = 1) =>
  api.get(`/users/${userId}/threads?page=${page}`)
export const getUserPosts = (userId: number, page = 1) =>
  api.get(`/users/${userId}/posts?page=${page}`)
