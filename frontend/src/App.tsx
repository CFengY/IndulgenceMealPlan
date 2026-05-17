import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import PrivateRoute from './components/common/PrivateRoute'
import Layout from './components/layout/Layout'

// 页面组件
import Login from './components/auth/Login'
import Register from './components/auth/Register'
import Home from './components/Home'
import MealList from './components/meals/MealList'
import MealForm from './components/meals/MealForm'
import Stats from './components/stats/Stats'
import Calendar from './components/stats/Calendar'
import Timeline from './components/posts/Timeline'

function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          {/* 公开路由 */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          {/* 需要认证的路由 */}
          <Route path="/" element={
            <PrivateRoute>
              <Layout />
            </PrivateRoute>
          }>
            <Route index element={<Home />} />
            <Route path="meals" element={<MealList />} />
            <Route path="meals/new" element={<MealForm />} />
            <Route path="meals/edit/:id" element={<MealForm />} />
            <Route path="stats" element={<Stats />} />
            <Route path="calendar" element={<Calendar />} />
            <Route path="posts" element={<Timeline />} />
          </Route>

          {/* 默认重定向 */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </Router>
  )
}

export default App