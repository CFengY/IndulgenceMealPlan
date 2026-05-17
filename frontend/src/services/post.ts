import api from './api'

export interface Post {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  user_id: number
  content: string
  images?: string
  user?: {
    ID: number
    CreatedAt: string
    UpdatedAt: string
    DeletedAt: string | null
    username: string
  }
}

export const postService = {
  getTimeline: async (): Promise<Post[]> => {
    const response = await api.get<{ data: Post[] }>('/api/v1/posts')
    return response.data.data
  },

  createPost: async (content: string, images?: FileList): Promise<Post> => {
    const formData = new FormData()
    formData.append('content', content)
    if (images) {
      Array.from(images).forEach((file) => {
        formData.append('images', file)
      })
    }
    const response = await api.post<{ data: Post }>('/api/v1/posts', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return response.data.data
  },

  deletePost: async (id: number): Promise<void> => {
    await api.delete(`/api/v1/posts/${id}`)
  },

  getImageUrl: (filename: string): string => {
    return `/images/${filename}`
  },
}
