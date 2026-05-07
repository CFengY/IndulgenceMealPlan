import api from './api'

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  name?: string
}

export interface AuthResponse {
  token: string
  user: {
    id: number
    username: string
    name?: string
  }
}

export const authService = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const formData = new URLSearchParams()
    formData.append('username', data.username)
    formData.append('password', data.password)

    const response = await api.post('/api/v1/auth/login', formData, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    })

    // 后端返回格式: { data: { userid, username, token } }
    const backendData = response.data.data
    const authResponse: AuthResponse = {
      token: backendData.token,
      user: {
        id: backendData.userid,
        username: backendData.username,
      },
    }

    if (authResponse.token) {
      localStorage.setItem('token', authResponse.token)
    }
    return authResponse
  },

  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const formData = new URLSearchParams()
    formData.append('username', data.username)
    formData.append('password', data.password)
    if (data.name) {
      formData.append('name', data.name)
    }

    const response = await api.post('/api/v1/auth/register', formData, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    })

    // 注册应该返回相同的格式
    const backendData = response.data.data || response.data
    const authResponse: AuthResponse = {
      token: backendData.token,
      user: {
        id: backendData.userid,
        username: backendData.username,
        name: backendData.name,
      },
    }

    if (authResponse.token) {
      localStorage.setItem('token', authResponse.token)
    }
    return authResponse
  },

  logout: async (): Promise<void> => {
    try {
      await api.post('/api/v1/logout')
    } catch (error) {
      // 即使登出 API 失败也清除本地 token
    } finally {
      localStorage.removeItem('token')
    }
  },

  getCurrentUser: () => {
    const token = localStorage.getItem('token')
    if (!token) return null
    // 简单解析 JWT payload（仅用于显示）
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return payload
    } catch {
      return null
    }
  },
}