import api from './api'

export interface Meal {
  ID: number
  user_id: number
  food_name: string
  meal_type: string
  meal_date: string
  image_path?: string
  calories: number
  protein_g: number
  fat_g: number
  carbs_g: number
  created_at?: string
  updated_at?: string
}

export interface CreateMealRequest {
  food_name: string
  meal_type: number
  meal_date: string
  image?: File
  calories?: number
  protein_g?: number
  fat_g?: number
  carbs_g?: number
}

export interface UpdateMealRequest {
  food_name?: string
  meal_type?: number
  meal_date?: string
  image?: File
  calories?: number
  protein_g?: number
  fat_g?: number
  carbs_g?: number
}

export interface FoodStatistic {
  food_name: string
  count: number
}

export interface NutritionSummary {
  calories: number
  protein_g: number
  fat_g: number
  carbs_g: number
}

export interface DateRangeResult {
  records: Meal[]
  statistics: FoodStatistic[]
  nutrition_summary: NutritionSummary
}

export const mealService = {
  getMeals: async (): Promise<Meal[]> => {
    const response = await api.get<{ data: Meal[] }>('/api/v1/meals')
    return response.data.data
  },

  createMeal: async (data: CreateMealRequest): Promise<Meal> => {
    const formData = new FormData()
    formData.append('food_name', data.food_name)
    formData.append('meal_type', data.meal_type.toString())
    formData.append('meal_date', data.meal_date)
    if (data.image) formData.append('image', data.image)
    if (data.calories) formData.append('calories', data.calories.toString())
    if (data.protein_g) formData.append('protein_g', data.protein_g.toString())
    if (data.fat_g) formData.append('fat_g', data.fat_g.toString())
    if (data.carbs_g) formData.append('carbs_g', data.carbs_g.toString())

    const response = await api.post<{ data: Meal }>('/api/v1/meals', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return response.data.data
  },

  updateMeal: async (id: number, data: UpdateMealRequest): Promise<Meal> => {
    const formData = new FormData()
    if (data.food_name) formData.append('food_name', data.food_name)
    if (data.meal_type) formData.append('meal_type', data.meal_type.toString())
    if (data.meal_date) formData.append('meal_date', data.meal_date)
    if (data.image) formData.append('image', data.image)
    if (data.calories !== undefined) formData.append('calories', data.calories.toString())
    if (data.protein_g !== undefined) formData.append('protein_g', data.protein_g.toString())
    if (data.fat_g !== undefined) formData.append('fat_g', data.fat_g.toString())
    if (data.carbs_g !== undefined) formData.append('carbs_g', data.carbs_g.toString())

    const response = await api.put<{ data: Meal }>(`/api/v1/meals/${id}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return response.data.data
  },

  deleteMeal: async (id: number): Promise<void> => {
    await api.delete(`/api/v1/meals/${id}`)
  },

  getMealsByDateRange: async (startDate: string, endDate: string): Promise<DateRangeResult> => {
    const response = await api.get<{ data: DateRangeResult }>('/api/v1/meals/range', {
      params: { start_date: startDate, end_date: endDate },
    })
    return response.data.data
  },

  exportCSV: async (startDate: string, endDate: string): Promise<void> => {
    const token = localStorage.getItem('token')
    const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
    const response = await fetch(
      `${baseUrl}/api/v1/meals/export?start_date=${startDate}&end_date=${endDate}`,
      { headers: { Authorization: `Bearer ${token}` } }
    )
    if (!response.ok) throw new Error('导出失败')

    const blob = await response.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `meals_export_${startDate}_${endDate}.csv`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  },

  getImageUrl: (filename: string): string => {
    return `/images/${filename}`
  },
}
