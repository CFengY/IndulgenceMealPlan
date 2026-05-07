import React, { useState, useEffect } from 'react'
import { mealService } from '../../services/meal'
import type { DateRangeResult } from '../../services/meal'
import { BarChart3, Calendar as CalendarIcon, TrendingUp, PieChart } from 'lucide-react'

const Stats: React.FC = () => {
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')
  const [result, setResult] = useState<DateRangeResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  // 设置默认日期范围（最近30天）
  useEffect(() => {
    const end = new Date()
    const start = new Date()
    start.setDate(end.getDate() - 30)

    setEndDate(end.toISOString().split('T')[0])
    setStartDate(start.toISOString().split('T')[0])
  }, [])

  const fetchStats = async () => {
    if (!startDate || !endDate) {
      setError('请选择日期范围')
      return
    }

    try {
      setLoading(true)
      const data = await mealService.getMealsByDateRange(startDate, endDate)
      setResult(data)
      setError('')
    } catch (err: any) {
      setError(err.response?.data?.error || '获取统计数据失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (startDate && endDate) {
      fetchStats()
    }
  }, [startDate, endDate])

  const getMealTypeCount = () => {
    if (!result) return { 早餐: 0, 午餐: 0, 晚餐: 0 }

    const counts = { 早餐: 0, 午餐: 0, 晚餐: 0 }
    result.records.forEach((record) => {
      if (record.meal_type in counts) {
        counts[record.meal_type as keyof typeof counts]++
      }
    })
    return counts
  }

  const mealTypeCounts = getMealTypeCount()
  const totalMeals = result ? result.records.length : 0
  const topFoods = result ? result.statistics.slice(0, 5) : []

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-800">统计分析</h1>
        <p className="text-gray-600">查看指定日期范围内的用餐统计</p>
      </div>

      {/* 日期选择器 */}
      <div className="bg-white border rounded-xl p-6 mb-8">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex-1">
            <label className="block text-sm font-medium text-gray-700 mb-2">开始日期</label>
            <input
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
            />
          </div>
          <div className="flex-1">
            <label className="block text-sm font-medium text-gray-700 mb-2">结束日期</label>
            <input
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
            />
          </div>
          <div className="flex items-end">
            <button
              onClick={fetchStats}
              disabled={loading}
              className="bg-gradient-to-r from-blue-500 to-purple-500 text-white px-6 py-3 rounded-lg font-medium hover:opacity-90 transition disabled:opacity-50"
            >
              {loading ? '查询中...' : '更新查询'}
            </button>
          </div>
        </div>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-600">
          {error}
        </div>
      )}

      {loading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
        </div>
      ) : result ? (
        <>
          {/* 概览卡片 */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <div className="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-xl p-6">
              <div className="flex items-center">
                <BarChart3 className="h-8 w-8 mr-3" />
                <div>
                  <p className="text-sm opacity-90">总记录数</p>
                  <p className="text-3xl font-bold">{totalMeals}</p>
                </div>
              </div>
            </div>

            <div className="bg-white border rounded-xl p-6">
              <div className="flex items-center">
                <CalendarIcon className="h-8 w-8 mr-3 text-blue-500" />
                <div>
                  <p className="text-sm text-gray-600">早餐次数</p>
                  <p className="text-3xl font-bold text-gray-800">{mealTypeCounts.早餐}</p>
                </div>
              </div>
            </div>

            <div className="bg-white border rounded-xl p-6">
              <div className="flex items-center">
                <CalendarIcon className="h-8 w-8 mr-3 text-green-500" />
                <div>
                  <p className="text-sm text-gray-600">午餐次数</p>
                  <p className="text-3xl font-bold text-gray-800">{mealTypeCounts.午餐}</p>
                </div>
              </div>
            </div>

            <div className="bg-white border rounded-xl p-6">
              <div className="flex items-center">
                <CalendarIcon className="h-8 w-8 mr-3 text-purple-500" />
                <div>
                  <p className="text-sm text-gray-600">晚餐次数</p>
                  <p className="text-3xl font-bold text-gray-800">{mealTypeCounts.晚餐}</p>
                </div>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* 热门食物排行 */}
            <div className="bg-white border rounded-xl p-6">
              <div className="flex items-center mb-6">
                <TrendingUp className="h-6 w-6 mr-2 text-blue-500" />
                <h2 className="text-xl font-bold text-gray-800">热门食物排行</h2>
              </div>
              <div className="space-y-4">
                {topFoods.map((food, index) => (
                  <div key={food.food_name} className="flex items-center">
                    <div className={`w-8 h-8 flex items-center justify-center rounded-lg mr-3 ${
                      index === 0 ? 'bg-yellow-100 text-yellow-800' :
                      index === 1 ? 'bg-gray-100 text-gray-800' :
                      index === 2 ? 'bg-orange-100 text-orange-800' :
                      'bg-gray-100 text-gray-600'
                    }`}>
                      <span className="font-bold">{index + 1}</span>
                    </div>
                    <div className="flex-1">
                      <h4 className="font-medium text-gray-800">{food.food_name}</h4>
                      <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                        <div
                          className="bg-gradient-to-r from-blue-500 to-purple-500 h-2 rounded-full"
                          style={{ width: `${(food.count / (topFoods[0]?.count || 1)) * 100}%` }}
                        ></div>
                      </div>
                    </div>
                    <div className="ml-4 text-right">
                      <p className="text-2xl font-bold text-gray-800">{food.count}</p>
                      <p className="text-sm text-gray-500">次</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* 餐类分布 */}
            <div className="bg-white border rounded-xl p-6">
              <div className="flex items-center mb-6">
                <PieChart className="h-6 w-6 mr-2 text-green-500" />
                <h2 className="text-xl font-bold text-gray-800">餐类分布</h2>
              </div>
              <div className="space-y-4">
                {Object.entries(mealTypeCounts).map(([type, count]) => (
                  <div key={type} className="flex items-center">
                    <div className={`w-12 h-12 rounded-full flex items-center justify-center mr-4 ${
                      type === '早餐' ? 'bg-yellow-100' :
                      type === '午餐' ? 'bg-blue-100' :
                      'bg-purple-100'
                    }`}>
                      <span className={`font-bold ${
                        type === '早餐' ? 'text-yellow-600' :
                        type === '午餐' ? 'text-blue-600' :
                        'text-purple-600'
                      }`}>
                        {type === '早餐' ? '早' : type === '午餐' ? '午' : '晚'}
                      </span>
                    </div>
                    <div className="flex-1">
                      <h4 className="font-medium text-gray-800">{type}</h4>
                      <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                        <div
                          className={`h-2 rounded-full ${
                            type === '早餐' ? 'bg-yellow-500' :
                            type === '午餐' ? 'bg-blue-500' :
                            'bg-purple-500'
                          }`}
                          style={{ width: totalMeals > 0 ? `${(count / totalMeals) * 100}%` : '0%' }}
                        ></div>
                      </div>
                    </div>
                    <div className="ml-4 text-right">
                      <p className="text-2xl font-bold text-gray-800">{count}</p>
                      <p className="text-sm text-gray-500">
                        {totalMeals > 0 ? `${Math.round((count / totalMeals) * 100)}%` : '0%'}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* 记录列表 */}
          {result.records.length > 0 && (
            <div className="mt-8">
              <h2 className="text-xl font-bold text-gray-800 mb-4">详细记录</h2>
              <div className="bg-white border rounded-xl overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="py-3 px-6 text-left text-sm font-medium text-gray-700">日期</th>
                        <th className="py-3 px-6 text-left text-sm font-medium text-gray-700">食物名称</th>
                        <th className="py-3 px-6 text-left text-sm font-medium text-gray-700">餐类</th>
                        <th className="py-3 px-6 text-left text-sm font-medium text-gray-700">图片</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200">
                      {result.records.map((record) => (
                        <tr key={record.id} className="hover:bg-gray-50">
                          <td className="py-4 px-6 text-gray-800">{record.meal_date}</td>
                          <td className="py-4 px-6 font-medium text-gray-800">{record.food_name}</td>
                          <td className="py-4 px-6">
                            <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                              record.meal_type === '早餐' ? 'bg-yellow-100 text-yellow-800' :
                              record.meal_type === '午餐' ? 'bg-blue-100 text-blue-800' :
                              'bg-purple-100 text-purple-800'
                            }`}>
                              {record.meal_type}
                            </span>
                          </td>
                          <td className="py-4 px-6">
                            {record.image_path && (
                              <img
                                src={mealService.getImageUrl(record.image_path)}
                                alt={record.food_name}
                                className="w-16 h-16 object-cover rounded-lg"
                              />
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          )}
        </>
      ) : null}
    </div>
  )
}

export default Stats