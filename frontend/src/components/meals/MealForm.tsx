import React, { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { mealService } from '../../services/meal'
import type { CreateMealRequest, UpdateMealRequest } from '../../services/meal'
import { Camera, Upload, ArrowLeft } from 'lucide-react'

const MealForm: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const isEditMode = !!id
  const navigate = useNavigate()

  const [foodName, setFoodName] = useState('')
  const [mealType, setMealType] = useState('2') // 默认午餐
  const [mealDate, setMealDate] = useState('')
  const [image, setImage] = useState<File | null>(null)
  const [previewUrl, setPreviewUrl] = useState<string>('')
  const [calories, setCalories] = useState('')
  const [proteinG, setProteinG] = useState('')
  const [fatG, setFatG] = useState('')
  const [carbsG, setCarbsG] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [existingImage, setExistingImage] = useState<string>('')

  useEffect(() => {
    // 设置默认日期为今天
    const today = new Date().toISOString().split('T')[0]
    setMealDate(today)

    // 如果是编辑模式，加载现有数据
    if (isEditMode) {
      fetchMeal()
    }
  }, [id])

  const fetchMeal = async () => {
    try {
      const meal = await mealService.getMeals().then((meals) =>
        meals.find((m) => m.ID === parseInt(id!))
      )
      if (meal) {
        setFoodName(meal.food_name)
        setMealType(meal.meal_type === '早餐' ? '1' : meal.meal_type === '午餐' ? '2' : '3')
        setMealDate(meal.meal_date)
        if (meal.image_path) setExistingImage(mealService.getImageUrl(meal.image_path))
        if (meal.calories) setCalories(meal.calories.toString())
        if (meal.protein_g) setProteinG(meal.protein_g.toString())
        if (meal.fat_g) setFatG(meal.fat_g.toString())
        if (meal.carbs_g) setCarbsG(meal.carbs_g.toString())
      }
    } catch (err) {
      setError('加载记录失败')
    }
  }

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setImage(file)
      const reader = new FileReader()
      reader.onloadend = () => {
        setPreviewUrl(reader.result as string)
      }
      reader.readAsDataURL(file)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      if (isEditMode) {
        const updateData: UpdateMealRequest = {
          food_name: foodName,
          meal_type: parseInt(mealType),
          meal_date: mealDate,
          calories: calories ? parseFloat(calories) : undefined,
          protein_g: proteinG ? parseFloat(proteinG) : undefined,
          fat_g: fatG ? parseFloat(fatG) : undefined,
          carbs_g: carbsG ? parseFloat(carbsG) : undefined,
        }
        if (image) updateData.image = image
        await mealService.updateMeal(parseInt(id!), updateData)
      } else {
        const createData: CreateMealRequest = {
          food_name: foodName,
          meal_type: parseInt(mealType),
          meal_date: mealDate,
          calories: calories ? parseFloat(calories) : undefined,
          protein_g: proteinG ? parseFloat(proteinG) : undefined,
          fat_g: fatG ? parseFloat(fatG) : undefined,
          carbs_g: carbsG ? parseFloat(carbsG) : undefined,
        }
        if (image) createData.image = image
        await mealService.createMeal(createData)
      }
      navigate('/meals')
    } catch (err: any) {
      setError(err.response?.data?.error || '提交失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <div className="mb-8">
        <button
          onClick={() => navigate('/meals')}
          className="flex items-center text-gray-600 dark:text-gray-400 dark:text-gray-400 hover:text-gray-800 dark:text-gray-100 mb-4"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          返回列表
        </button>
        <h1 className="text-2xl font-bold text-gray-800 dark:text-gray-100">
          {isEditMode ? '编辑用餐记录' : '添加用餐记录'}
        </h1>
        <p className="text-gray-600 dark:text-gray-400 dark:text-gray-400">记录你的美食时刻</p>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-600">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-8">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* 表单区域 */}
          <div className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                食物名称 *
              </label>
              <input
                type="text"
                value={foodName}
                onChange={(e) => setFoodName(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                placeholder="例如：炸鸡、奶茶、沙拉"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                餐类 *
              </label>
              <div className="grid grid-cols-3 gap-3">
                {[
                  { value: '1', label: '早餐', color: 'bg-yellow-100 text-yellow-800' },
                  { value: '2', label: '午餐', color: 'bg-blue-100 text-blue-800' },
                  { value: '3', label: '晚餐', color: 'bg-purple-100 text-purple-800' },
                ].map((option) => (
                  <button
                    key={option.value}
                    type="button"
                    onClick={() => setMealType(option.value)}
                    className={`py-3 rounded-lg text-center font-medium transition ${
                      mealType === option.value
                        ? `${option.color} ring-2 ring-offset-2 ring-current`
                        : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400 dark:text-gray-400 hover:bg-gray-200'
                    }`}
                  >
                    {option.label}
                  </button>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                用餐日期 *
              </label>
              <input
                type="date"
                value={mealDate}
                onChange={(e) => setMealDate(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                营养信息（可选）
              </label>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">热量 (kcal)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={calories}
                    onChange={(e) => setCalories(e.target.value)}
                    className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                    placeholder="例如：350"
                  />
                </div>
                <div>
                  <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">蛋白质 (g)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={proteinG}
                    onChange={(e) => setProteinG(e.target.value)}
                    className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                    placeholder="例如：30"
                  />
                </div>
                <div>
                  <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">脂肪 (g)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={fatG}
                    onChange={(e) => setFatG(e.target.value)}
                    className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                    placeholder="例如：10"
                  />
                </div>
                <div>
                  <label className="block text-xs text-gray-500 dark:text-gray-400 mb-1">碳水 (g)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={carbsG}
                    onChange={(e) => setCarbsG(e.target.value)}
                    className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
                    placeholder="例如：45"
                  />
                </div>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                食物图片（可选）
              </label>
              <div className="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-6 text-center hover:border-blue-400 transition">
                <input
                  type="file"
                  id="image-upload"
                  accept="image/*"
                  onChange={handleImageChange}
                  className="hidden"
                />
                <label htmlFor="image-upload" className="cursor-pointer">
                  {previewUrl || existingImage ? (
                    <div className="flex flex-col items-center">
                      <img
                        src={previewUrl || existingImage}
                        alt="预览"
                        className="w-32 h-32 object-cover rounded-lg mb-4"
                      />
                      <span className="text-blue-600 hover:text-blue-800 font-medium">
                        更换图片
                      </span>
                    </div>
                  ) : (
                    <div className="flex flex-col items-center">
                      <Camera className="h-12 w-12 text-gray-400 dark:text-gray-500 mb-3" />
                      <span className="text-gray-600 dark:text-gray-400 dark:text-gray-400">点击上传图片</span>
                      <span className="text-sm text-gray-500 dark:text-gray-400 mt-1">支持 JPG, PNG, GIF 格式</span>
                    </div>
                  )}
                </label>
              </div>
            </div>
          </div>

          {/* 预览区域 */}
          <div>
            <div className="bg-gray-50 dark:bg-gray-900 rounded-xl p-6 sticky top-6">
              <h3 className="text-lg font-medium text-gray-800 dark:text-gray-100 mb-4">预览</h3>
              <div className="bg-white dark:bg-gray-800 rounded-lg p-6 shadow-sm">
                <div className="flex items-center mb-4">
                  <div className={`w-12 h-12 rounded-full flex items-center justify-center ${
                    mealType === '1' ? 'bg-yellow-100' :
                    mealType === '2' ? 'bg-blue-100' :
                    'bg-purple-100'
                  }`}>
                    <span className={`font-bold ${
                      mealType === '1' ? 'text-yellow-600' :
                      mealType === '2' ? 'text-blue-600' :
                      'text-purple-600'
                    }`}>
                      {mealType === '1' ? '早' : mealType === '2' ? '午' : '晚'}
                    </span>
                  </div>
                  <div className="ml-4">
                    <h4 className="font-bold text-gray-800 dark:text-gray-100">
                      {foodName || '食物名称'}
                    </h4>
                    <p className="text-gray-500 dark:text-gray-400 text-sm">
                      {mealDate || '选择日期'} ·
                      {mealType === '1' ? '早餐' : mealType === '2' ? '午餐' : '晚餐'}
                    </p>
                  </div>
                </div>

                {(previewUrl || existingImage) && (
                  <div className="mt-4">
                    <img
                      src={previewUrl || existingImage}
                      alt="预览"
                      className="w-full h-48 object-cover rounded-lg"
                    />
                  </div>
                )}

                <div className="mt-6 pt-6 border-t">
                  <p className="text-gray-600 dark:text-gray-400 dark:text-gray-400 text-sm">
                    记录你的美食时刻，分享健康饮食生活
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="flex justify-end space-x-4 pt-6 border-t">
          <button
            type="button"
            onClick={() => navigate('/meals')}
            className="px-6 py-3 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 dark:bg-gray-900 transition"
          >
            取消
          </button>
          <button
            type="submit"
            disabled={loading}
            className="flex items-center bg-gradient-to-r from-blue-500 to-purple-500 text-white px-8 py-3 rounded-lg font-medium hover:opacity-90 transition disabled:opacity-50"
          >
            {loading ? (
              '提交中...'
            ) : (
              <>
                <Upload className="h-5 w-5 mr-2" />
                {isEditMode ? '更新记录' : '保存记录'}
              </>
            )}
          </button>
        </div>
      </form>
    </div>
  )
}

export default MealForm