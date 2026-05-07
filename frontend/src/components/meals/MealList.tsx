import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { mealService } from '../../services/meal'
import type { Meal } from '../../services/meal'
import { Edit2, Trash2, Camera, Coffee, Utensils, Moon } from 'lucide-react'

const MealList: React.FC = () => {
  const [meals, setMeals] = useState<Meal[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const fetchMeals = async () => {
    try {
      setLoading(true)
      const data = await mealService.getMeals()
      console.log('Fetched meals data:', data)
      if (data.length > 0) {
        console.log('First meal object:', data[0])
        console.log('Meal ID:', data[0].ID)
      }
      setMeals(data)
      setError('')
    } catch (err: any) {
      setError(err.response?.data?.error || '获取记录失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchMeals()
  }, [])

  const handleDelete = async (id: number) => {
    if (!window.confirm('确定要删除这条记录吗？')) return

    try {
      await mealService.deleteMeal(id)
      setMeals(meals.filter((meal) => meal.ID !== id))
    } catch (err: any) {
      alert(err.response?.data?.error || '删除失败')
    }
  }

  const getMealTypeIcon = (mealType: string) => {
    switch (mealType) {
      case '早餐':
        return <Coffee className="h-5 w-5 text-yellow-500" />
      case '午餐':
        return <Utensils className="h-5 w-5 text-blue-500" />
      case '晚餐':
        return <Moon className="h-5 w-5 text-purple-500" />
      default:
        return <Utensils className="h-5 w-5 text-gray-500" />
    }
  }

  const getMealTypeColor = (mealType: string) => {
    switch (mealType) {
      case '早餐':
        return 'bg-yellow-100 text-yellow-800'
      case '午餐':
        return 'bg-blue-100 text-blue-800'
      case '晚餐':
        return 'bg-purple-100 text-purple-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', weekday: 'short' })
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
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-800">我的用餐记录</h1>
          <p className="text-gray-600">记录你的每一餐，享受美食生活</p>
        </div>
        <Link
          to="/meals/new"
          className="bg-gradient-to-r from-blue-500 to-purple-500 text-white px-6 py-3 rounded-lg font-medium hover:opacity-90 transition"
        >
          添加新记录
        </Link>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-600">
          {error}
        </div>
      )}

      {meals.length === 0 ? (
        <div className="text-center py-12">
          <div className="w-24 h-24 mx-auto bg-gray-100 rounded-full flex items-center justify-center mb-4">
            <Utensils className="h-12 w-12 text-gray-400" />
          </div>
          <h3 className="text-xl font-medium text-gray-700 mb-2">还没有记录</h3>
          <p className="text-gray-500 mb-6">开始记录你的第一餐吧！</p>
          <Link
            to="/meals/new"
            className="inline-flex items-center bg-gradient-to-r from-blue-500 to-purple-500 text-white px-6 py-3 rounded-lg font-medium hover:opacity-90 transition"
          >
            <Camera className="mr-2 h-5 w-5" />
            添加第一餐
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {meals.map((meal) => (
            <div
              key={meal.ID}
              className="bg-white border rounded-xl p-6 hover:shadow-md transition"
            >
              <div className="flex flex-col md:flex-row md:items-center justify-between">
                <div className="flex-1">
                  <div className="flex items-center mb-2">
                    {getMealTypeIcon(meal.meal_type)}
                    <span className={`ml-2 px-3 py-1 rounded-full text-sm font-medium ${getMealTypeColor(meal.meal_type)}`}>
                      {meal.meal_type}
                    </span>
                    <span className="ml-4 text-gray-500 text-sm">
                      {formatDate(meal.meal_date)}
                    </span>
                  </div>
                  <h3 className="text-xl font-bold text-gray-800 mb-2">{meal.food_name}</h3>
                  {meal.image_path && (
                    <div className="mt-4">
                      <img
                        src={mealService.getImageUrl(meal.image_path)}
                        alt={meal.food_name}
                        className="w-32 h-32 object-cover rounded-lg"
                      />
                    </div>
                  )}
                </div>

                <div className="mt-4 md:mt-0 flex space-x-2">
                  <Link
                    to={`/meals/edit/${meal.ID}`}
                    className="flex items-center px-4 py-2 bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition"
                  >
                    <Edit2 className="h-4 w-4 mr-2" />
                    编辑
                  </Link>
                  <button
                    onClick={() => handleDelete(meal.ID)}
                    className="flex items-center px-4 py-2 bg-red-50 text-red-600 rounded-lg hover:bg-red-100 transition"
                  >
                    <Trash2 className="h-4 w-4 mr-2" />
                    删除
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default MealList