import { defineStore } from 'pinia'

interface AdminState {
  users: any[]
  posts: any[]
  news: any[]
  clubs: any[]
  loading: boolean
  stats: {
    totalUsers: number
    totalPosts: number
  }
}

export const useAdminStore = defineStore('admin', {
  state: (): AdminState => ({
    users: [],
    posts: [],
    news: [],
    clubs: [],
    loading: false,
    stats: {
      totalUsers: 0,
      totalPosts: 0
    }
  }),
  actions: {
    setUsers(users: any[]) {
      this.users = users
    },
    setPosts(posts: any[]) {
      this.posts = posts
    },
    setLoading(status: boolean) {
      this.loading = status
    }
  }
})
