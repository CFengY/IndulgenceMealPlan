import api from './api'

export interface ChatResponse {
  reply: string
}

export const chatService = {
  sendMessage: async (message: string): Promise<ChatResponse> => {
    const response = await api.post<ChatResponse>('/api/v1/chat', { message })
    return response.data
  },
}
