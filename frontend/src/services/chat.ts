const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const chatService = {
  sendMessageStream: async (
    message: string,
    onChunk: (chunk: string) => void,
    onDone: () => void,
    onError: (error: string) => void
  ): Promise<void> => {
    const token = localStorage.getItem('token')

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/chat`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ message }),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: '请求失败' }))
        onError(errorData.error || '请求失败')
        return
      }

      const reader = response.body?.getReader()
      if (!reader) {
        onError('浏览器不支持流式读取')
        return
      }

      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6)
            if (data === '[DONE]') {
              onDone()
              return
            }
            try {
              const parsed = JSON.parse(data)
              if (parsed.content) {
                onChunk(parsed.content)
              }
            } catch {
              // 忽略 JSON 解析错误
            }
          }
        }
      }

      onDone()
    } catch {
      onError('网络请求失败，请稍后再试')
    }
  },
}
