import { defineStore } from 'pinia'
import http from '@/api/http'

type User = {
  id: number
  username: string
  display_name: string
  email: string
  phone: string
  permissions: Array<{ code: string }>
}

type LoginPayload = { username: string; password: string }

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    ready: false,
  }),
  getters: {
    permissions(state) {
      return state.user?.permissions?.map((item) => item.code) ?? []
    },
  },
  actions: {
    async login(payload: LoginPayload) {
      const response = await http.post('/auth/login', payload)
      localStorage.setItem('access_token', response.data.data.access_token)
      localStorage.setItem('refresh_token', response.data.data.refresh_token)
      this.user = response.data.data.user
    },
    async register(payload: User & { password: string }) {
      await http.post('/auth/register', payload)
    },
    async loadMe() {
      try {
        const response = await http.get('/users/me')
        this.user = response.data.data
      } finally {
        this.ready = true
      }
    },
    logout() {
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        void http.post('/auth/logout', { refresh_token: refreshToken })
      }
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      this.user = null
    },
  },
})
