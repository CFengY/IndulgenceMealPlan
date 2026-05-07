import React, { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import { authService } from '../services/auth'

interface User {
  id: number
  username: string
  name?: string
}

interface AuthContextType {
  user: User | null
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  register: (username: string, password: string, name?: string) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    // 检查本地是否有 token，并尝试获取用户信息
    const token = localStorage.getItem('token')
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        console.log('Token payload:', payload) // 调试
        setUser({
          id: payload.UserId || payload.sub || payload.id,
          username: payload.username,
          name: payload.name,
        })
      } catch (error) {
        console.error('Failed to parse token:', error)
        // 不立即移除token，可能token仍然有效，只是解析有问题
        // localStorage.removeItem('token')
      }
    }
    setIsLoading(false)
  }, [])

  const login = async (username: string, password: string) => {
    const response = await authService.login({ username, password })
    console.log('Login response:', response) // 调试
    // 优先使用响应中的用户信息
    setUser({
      id: response.user.id,
      username: response.user.username,
      name: response.user.name,
    })
    // 同时尝试解析token以验证
    try {
      const payload = JSON.parse(atob(response.token.split('.')[1]))
      console.log('Token payload from login:', payload)
    } catch (error) {
      console.error('Failed to parse token after login:', error)
    }
  }

  const register = async (username: string, password: string, name?: string) => {
    const response = await authService.register({ username, password, name })
    const payload = JSON.parse(atob(response.token.split('.')[1]))
    setUser({
      id: payload.sub || payload.id,
      username: payload.username,
      name: payload.name,
    })
  }

  const logout = async () => {
    await authService.logout()
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, isLoading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}