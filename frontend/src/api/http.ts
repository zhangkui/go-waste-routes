import axios from 'axios'

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 20000,
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => response,
  async (error) => {
    const status = error.response?.status
    if (status === 401) {
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        try {
          const refresh = await axios.post('/api/v1/auth/refresh', { refresh_token: refreshToken })
          localStorage.setItem('access_token', refresh.data.data.access_token)
          localStorage.setItem('refresh_token', refresh.data.data.refresh_token)
          error.config.headers.Authorization = `Bearer ${refresh.data.data.access_token}`
          return axios.request(error.config)
        } catch {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          window.location.href = '/login'
        }
      }
    }
    return Promise.reject(error)
  },
)

export default http
