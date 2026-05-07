import api from './api'

export interface Meal {
  ID: number
  user_id: number
  food_name: string
  meal_type: string
  meal_date: string
  image_path?: string
  created_at?: string
  updated_at?: string
}

export interface CreateMealRequest {
  food_name: string
  meal_type: number // 1-早餐 2-午餐 3-晚餐
  meal_date: string // YYYY-MM-DD
  image?: File
}

export interface UpdateMealRequest {
  food_name?: string
  meal_type?: number
  meal_date?: string
  image?: File
}

export interface FoodStatistic {
  food_name: string
  count: number
}

export interface DateRangeResult {
  records: Meal[]
  statistics: FoodStatistic[]
}

export const mealService = {
  // 获取当前用户的所有记录
  getMeals: async (): Promise<Meal[]> => {
    const response = await api.get<{ data: Meal[] }>('/api/v1/meals')
    return response.data.data
  },

  // 创建记录（支持图片上传）
  createMeal: async (data: CreateMealRequest): Promise<Meal> => {
    const formData = new FormData()
    formData.append('food_name', data.food_name)
    formData.append('meal_type', data.meal_type.toString())
    formData.append('meal_date', data.meal_date)
    if (data.image) {
      formData.append('image', data.image)
    }

    const response = await api.post<{ data: Meal }>('/api/v1/meals', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data.data
  },

  // 更新记录
  updateMeal: async (id: number, data: UpdateMealRequest): Promise<Meal> => {
    const formData = new FormData()
    if (data.food_name) formData.append('food_name', data.food_name)
    if (data.meal_type) formData.append('meal_type', data.meal_type.toString())
    if (data.meal_date) formData.append('meal_date', data.meal_date)
    if (data.image) formData.append('image', data.image)

    const response = await api.put<{ data: Meal }>(`/api/v1/meals/${id}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data.data
  },

  // 删除记录
  deleteMeal: async (id: number): Promise<void> => {
    await api.delete(`/api/v1/meals/${id}`)
  },

  // 日期范围查询
  getMealsByDateRange: async (startDate: string, endDate: string): Promise<DateRangeResult> => {
    const response = await api.get<{ data: DateRangeResult }>('/api/v1/meals/range', {
      params: { start_date: startDate, end_date: endDate },
    })
    return response.data.data
  },

  // 获取图片完整 URL
  getImageUrl: (filename: string): string => {
    return `/images/${filename}`
  },
}