import React, { useState } from 'react'
import { Outlet, Link, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import ChatBot from '../chat/ChatBot'
import {
  Home,
  Utensils,
  BarChart3,
  Calendar,
  PlusCircle,
  LogOut,
  Menu,
  X,
  User,
} from 'lucide-react'

const Layout: React.FC = () => {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const navigation = [
    { name: '首页', path: '/', icon: Home },
    { name: '我的记录', path: '/meals', icon: Utensils },
    { name: '添加记录', path: '/meals/new', icon: PlusCircle },
    { name: '范围查询', path: '/stats', icon: BarChart3 },
    { name: '日历视图', path: '/calendar', icon: Calendar },
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 导航栏 */}
      <nav className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <button
                onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                className="md:hidden p-2 rounded-md text-gray-500 hover:text-gray-700"
              >
                {isMobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
              </button>
              <div className="flex items-center ml-4 md:ml-0">
                <Utensils className="h-8 w-8 text-blue-500" />
                <span className="ml-2 text-xl font-bold text-gray-800">放纵餐计划</span>
              </div>
            </div>

            <div className="hidden md:flex items-center space-x-4">
              <div className="flex items-center space-x-2 text-gray-600">
                <User size={18} />
                <span>{user?.name || user?.username}</span>
              </div>
              <button
                onClick={handleLogout}
                className="flex items-center space-x-2 px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition"
              >
                <LogOut size={18} />
                <span>退出</span>
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* 移动端菜单 */}
      {isMobileMenuOpen && (
        <div className="md:hidden bg-white border-b">
          <div className="px-2 pt-2 pb-3 space-y-1">
            {navigation.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                onClick={() => setIsMobileMenuOpen(false)}
                className={`flex items-center px-3 py-2 rounded-md text-base font-medium ${
                  location.pathname === item.path
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-700 hover:bg-gray-50'
                }`}
              >
                <item.icon className="mr-3 h-5 w-5" />
                {item.name}
              </Link>
            ))}
            <div className="px-3 py-2 text-sm text-gray-500">
              <div className="flex items-center">
                <User size={16} className="mr-2" />
                {user?.name || user?.username}
              </div>
            </div>
            <button
              onClick={handleLogout}
              className="w-full flex items-center px-3 py-2 text-base font-medium text-gray-700 hover:bg-gray-50 rounded-md"
            >
              <LogOut className="mr-3 h-5 w-5" />
              退出登录
            </button>
          </div>
        </div>
      )}

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="md:flex">
          {/* 侧边栏（桌面端） */}
          <div className="hidden md:block w-64 flex-shrink-0 mr-8">
            <div className="bg-white rounded-xl shadow-sm p-4">
              <div className="mb-6">
                <h2 className="text-lg font-semibold text-gray-800 mb-4">导航</h2>
                <nav className="space-y-2">
                  {navigation.map((item) => (
                    <Link
                      key={item.path}
                      to={item.path}
                      className={`flex items-center px-3 py-2 rounded-lg transition ${
                        location.pathname === item.path
                          ? 'bg-blue-50 text-blue-600 font-medium'
                          : 'text-gray-600 hover:bg-gray-50'
                      }`}
                    >
                      <item.icon className="mr-3 h-5 w-5" />
                      {item.name}
                    </Link>
                  ))}
                </nav>
              </div>

              <div className="pt-4 border-t">
                <div className="flex items-center px-3 py-2 text-gray-600">
                  <User className="mr-3 h-5 w-5" />
                  <div>
                    <p className="font-medium">{user?.name || user?.username}</p>
                    <p className="text-sm text-gray-500">欢迎回来</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* 主内容区 */}
          <div className="flex-1">
            <div className="bg-white rounded-xl shadow-sm p-6">
              <Outlet />
            </div>
          </div>
        </div>
      </div>

      <ChatBot />
    </div>
  )
}

export default Layout