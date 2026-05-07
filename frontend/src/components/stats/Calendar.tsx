import React, { useState, useEffect } from 'react'
import { mealService } from '../../services/meal'
import type { Meal } from '../../services/meal'
import { ChevronLeft, ChevronRight, Calendar as CalendarIcon } from 'lucide-react'

const Calendar: React.FC = () => {
  const [currentDate, setCurrentDate] = useState(new Date())
  const [meals, setMeals] = useState<Meal[]>([])
  const [loading, setLoading] = useState(false)

  // 获取当前月份的所有餐食记录
  useEffect(() => {
    const fetchMealsForMonth = async () => {
      setLoading(true)
      try {
        // 使用本地日期格式，避免时区问题
        const formatDateToYYYYMMDD = (date: Date) => {
          const year = date.getFullYear()
          const month = String(date.getMonth() + 1).padStart(2, '0')
          const day = String(date.getDate()).padStart(2, '0')
          return `${year}-${month}-${day}`
        }

        const startDate = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1)
        const endDate = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0)

        const startDateStr = formatDateToYYYYMMDD(startDate)
        const endDateStr = formatDateToYYYYMMDD(endDate)

        console.log('Calendar fetching data for:', startDateStr, 'to', endDateStr)

        const result = await mealService.getMealsByDateRange(startDateStr, endDateStr)
        console.log('Calendar data:', result)
        setMeals(result.records)
      } catch (error) {
        console.error('Failed to fetch meals:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchMealsForMonth()
  }, [currentDate])

  const getDaysInMonth = (year: number, month: number) => {
    return new Date(year, month + 1, 0).getDate()
  }

  const getFirstDayOfMonth = (year: number, month: number) => {
    return new Date(year, month, 1).getDay()
  }

  const prevMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1))
  }

  const nextMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1))
  }

  const year = currentDate.getFullYear()
  const month = currentDate.getMonth()
  const daysInMonth = getDaysInMonth(year, month)
  const firstDayOfMonth = getFirstDayOfMonth(year, month)

  const monthNames = [
    '一月', '二月', '三月', '四月', '五月', '六月',
    '七月', '八月', '九月', '十月', '十一月', '十二月'
  ]

  const dayNames = ['日', '一', '二', '三', '四', '五', '六']

  // 获取某一天的餐食记录
  const getMealsForDay = (day: number) => {
    const dateStr = `${year}-${(month + 1).toString().padStart(2, '0')}-${day.toString().padStart(2, '0')}`
    return meals.filter(meal => {
      // meal.meal_date 可能是 "2026-05-06T00:00:00Z" 或 "2026-05-06" 格式
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
      return mealDateStr === dateStr
    })
  }

  // 生成日历网格
  const calendarDays = []
  for (let i = 0; i < firstDayOfMonth; i++) {
    calendarDays.push(null)
  }
  for (let day = 1; day <= daysInMonth; day++) {
    calendarDays.push(day)
  }

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-800">日历视图</h1>
        <p className="text-gray-600">按日期查看你的用餐记录</p>
      </div>

      {/* 日历控件 */}
      <div className="bg-white border rounded-xl p-6 mb-8">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center">
            <CalendarIcon className="h-6 w-6 mr-2 text-blue-500" />
            <h2 className="text-xl font-bold text-gray-800">
              {year}年 {monthNames[month]}
            </h2>
          </div>
          <div className="flex space-x-2">
            <button
              onClick={prevMonth}
              className="p-2 rounded-lg hover:bg-gray-100 transition"
            >
              <ChevronLeft className="h-5 w-5" />
            </button>
            <button
              onClick={() => setCurrentDate(new Date())}
              className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition"
            >
              今天
            </button>
            <button
              onClick={nextMonth}
              className="p-2 rounded-lg hover:bg-gray-100 transition"
            >
              <ChevronRight className="h-5 w-5" />
            </button>
          </div>
        </div>

        {/* 日历网格 */}
        <div className="grid grid-cols-7 gap-2">
          {/* 星期标题 */}
          {dayNames.map(day => (
            <div key={day} className="text-center font-medium text-gray-500 py-2">
              {day}
            </div>
          ))}

          {/* 日期格子 */}
          {calendarDays.map((day, index) => {
            if (day === null) {
              return <div key={`empty-${index}`} className="h-32"></div>
            }

            const dayMeals = getMealsForDay(day)
            const isToday = new Date().getDate() === day &&
              new Date().getMonth() === month &&
              new Date().getFullYear() === year

            return (
              <div
                key={day}
                className={`min-h-32 border rounded-lg p-2 ${
                  isToday ? 'bg-blue-50 border-blue-200' : 'border-gray-200'
                }`}
              >
                <div className="flex justify-between items-start mb-1">
                  <span className={`font-medium ${isToday ? 'text-blue-600' : 'text-gray-700'}`}>
                    {day}
                  </span>
                  {dayMeals.length > 0 && (
                    <span className="text-xs bg-blue-100 text-blue-600 px-2 py-1 rounded-full">
                      {dayMeals.length} 餐
                    </span>
                  )}
                </div>

                <div className="space-y-1 max-h-24 overflow-y-auto">
                  {dayMeals.slice(0, 3).map(meal => (
                    <div
                      key={meal.ID}
                      className={`text-xs p-1 rounded truncate ${
                        meal.meal_type === '早餐' ? 'bg-yellow-100 text-yellow-800' :
                        meal.meal_type === '午餐' ? 'bg-blue-100 text-blue-800' :
                        'bg-purple-100 text-purple-800'
                      }`}
                      title={`${meal.food_name} (${meal.meal_type})`}
                    >
                      {meal.food_name}
                    </div>
                  ))}
                  {dayMeals.length > 3 && (
                    <div className="text-xs text-gray-500 text-center">
                      +{dayMeals.length - 3} 更多
                    </div>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      </div>

      {/* 当月统计 */}
      <div className="bg-white border rounded-xl p-6">
        <h3 className="text-lg font-medium text-gray-800 mb-4">本月统计</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-sm text-gray-600">总记录数</p>
            <p className="text-2xl font-bold text-gray-800">{meals.length}</p>
          </div>
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-sm text-gray-600">早餐次数</p>
            <p className="text-2xl font-bold text-gray-800">
              {meals.filter(m => m.meal_type === '早餐').length}
            </p>
          </div>
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-sm text-gray-600">最常吃</p>
            <p className="text-2xl font-bold text-gray-800">
              {(() => {
                const foodCounts: Record<string, number> = {}
                meals.forEach(meal => {
                  foodCounts[meal.food_name] = (foodCounts[meal.food_name] || 0) + 1
                })
                const topFood = Object.entries(foodCounts)
                  .sort((a, b) => b[1] - a[1])[0]
                return topFood ? topFood[0] : '暂无'
              })()}
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Calendar