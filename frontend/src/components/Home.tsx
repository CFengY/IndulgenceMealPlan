import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Utensils, BarChart3, Calendar, PlusCircle, Coffee, Pizza } from 'lucide-react'
import { mealService } from '../services/meal'
import type { Meal, FoodStatistic } from '../services/meal'

const Home: React.FC = () => {
  const [recentMeals, setRecentMeals] = useState<Meal[]>([])
  const [stats, setStats] = useState<FoodStatistic[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchData = async () => {
      try {
        // 获取最近7天的数据用于热门食物统计
        const endDate = new Date()
        const startDate = new Date()
        startDate.setDate(endDate.getDate() - 7)

        // 使用本地日期格式 YYYY-MM-DD
        const formatDateToYYYYMMDD = (date: Date) => {
          const year = date.getFullYear()
          const month = String(date.getMonth() + 1).padStart(2, '0')
          const day = String(date.getDate()).padStart(2, '0')
          return `${year}-${month}-${day}`
        }

        const endDateStr = formatDateToYYYYMMDD(endDate)
        const startDateStr = formatDateToYYYYMMDD(startDate)

        const result = await mealService.getMealsByDateRange(startDateStr, endDateStr)

        // 获取当天日期
        const todayStr = formatDateToYYYYMMDD(endDate)

        // 筛选出当天的记录并按创建时间降序排序
        const todayMeals = result.records
          .filter(meal => {
            // 处理日期格式，meal.meal_date 可能是 "2026-05-06T00:00:00Z" 或 "2026-05-06" 格式
            let mealDateStr: string
            if (typeof meal.meal_date === 'string') {
              if (meal.meal_date.includes('T')) {
                mealDateStr = meal.meal_date.split('T')[0]
              } else {
                mealDateStr = meal.meal_date
              }
            } else {
              // 如果是其他类型，转换为字符串
              mealDateStr = String(meal.meal_date).split('T')[0]
            }
            console.log('Checking meal date:', mealDateStr, 'vs today:', todayStr, 'match:', mealDateStr === todayStr)
            return mealDateStr === todayStr
          })
          .sort((a, b) => {
            // 按创建时间降序排序，如果没有创建时间则按ID降序
            if (a.created_at && b.created_at) {
              return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
            }
            return b.ID - a.ID
          })

        console.log('Today meals found:', todayMeals.length, todayMeals)
        setRecentMeals(todayMeals)

        // 取前3个热门食物
        const topStats = result.statistics.slice(0, 3)
        setStats(topStats)
      } catch (error) {
        console.error('Failed to fetch home data:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    )
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', weekday: 'short' })
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800 mb-2">欢迎回来！</h1>
        <p className="text-gray-600">记录你的放纵餐，享受美食生活</p>
      </div>

      {/* 快速操作卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <Link
          to="/meals/new"
          className="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-xl p-6 shadow-lg hover:shadow-xl transition"
        >
          <div className="flex items-center">
            <PlusCircle className="h-8 w-8 mr-3" />
            <div>
              <h3 className="text-xl font-bold">添加记录</h3>
              <p className="opacity-90">记录新的一餐</p>
            </div>
          </div>
        </Link>

        <Link
          to="/meals"
          className="bg-white border rounded-xl p-6 shadow-sm hover:shadow-md transition"
        >
          <div className="flex items-center">
            <Utensils className="h-8 w);-8 mr-3 text-blue-500" />
            <div>
              <h3 className="text-xl font-bold text-gray-800">查看记录</h3>
              <p className="text-gray-600">查看所有用餐记录</p>
            </div>
          </div>
        </Link>

        <Link
          to="/stats"
          className="bg-white border rounded-xl p-6 shadow-sm hover:shadow-md transition"
        >
          <div className="flex items-center">
            <BarChart3 className="h-8 w-8 mr-3 text-green-500" />
            <div>
              <h3 className="text-xl font-bold text-gray-800">统计分析</h3>
              <p className="text-gray-600">查看食物频次统计</p>
            </div>
          </div>
        </Link>

        <Link
          to="/calendar"
          className="bg-white border rounded-xl p-6 shadow-sm hover:shadow-md transition"
        >
          <div className="flex items-center">
            <Calendar className="h-8 w-8 mr-3 text-orange-500" />
            <div>
              <h3 className="text-xl font-bold text-gray-800">日历视图</h3>
              <p className="text-gray-600">按日期查看记录</p>
            </div>
          </div>
        </Link>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* 最近记录 */}
        <div className="bg-white border rounded-xl p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-bold text-gray-800">最近记录</h2>
            <Link to="/meals" className="text-blue-600 hover:text-blue-800 text-sm font-medium">
              查看全部 →
            </Link>
          </div>
          <div className="space-y-4">
            {recentMeals.length === 0 ? (
              <div className="text-center py-6 text-gray-500">
                今天还没有记录，开始记录你的第一餐吧！
              </div>
            ) : (
              recentMeals.map((meal) => (
              <div key={meal.ID} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center">
                  {meal.meal_type === '早餐' && <Coffee className="h-5 w-5 mr-3 text-yellow-500" />}
                  {meal.meal_type === '午餐' && <Utensils className="h-5 w-5 mr-3 text-blue-500" />}
                  {meal.meal_type === '晚餐' && <Pizza className="h-5 w-5 mr-3 text-purple-500" />}
                  <div>
                    <h4 className="font-medium text-gray-800">{meal.food_name}</h4>
                    <p className="text-sm text-gray-500">{formatDate(meal.meal_date)} · {meal.meal_type}</p>
                  </div>
                </div>
                <span className="px-3 py-1 bg-blue-100 text-blue-600 text-sm font-medium rounded-full">
                  {meal.meal_type}
                </span>
              </div>
            )))}
          </div>
        </div>

        {/* 热门食物 */}
        <div className="bg-white border rounded-xl p-6">
          <h2 className="text-xl font-bold text-gray-800 mb-6">近一周热门食物</h2>
          <div className="space-y-4">
            {stats.length === 0 ? (
              <div className="text-center py-6 text-gray-500">
                近一周暂无统计数据
              </div>
            ) : (
              stats.map((stat, index) => (
                <div key={stat.food_name} className="flex items-center justify-between">
                  <div className="flex items-center">
                    <div className="w-8 h-8 flex items-center justify-center bg-gray-100 rounded-lg mr-3">
                      <span className="font-bold text-gray-700">{index + 1}</span>
                    </div>
                    <span className="font-medium text-gray-800">{stat.food_name}</span>
                  </div>
                  <div className="flex items-center">
                    <div className="w-32 bg-gray-200 rounded-full h-2 mr-3">
                      <div
                        className="bg-gradient-to-r from-blue-500 to-purple-500 h-2 rounded-full"
                        style={{ width: `${(stat.count / (stats[0]?.count || 1)) * 100}%` }}
                      ></div>
                    </div>
                    <span className="text-gray-600 font-medium">{stat.count} 次</span>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default Home