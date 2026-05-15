import React, { useState, useRef, useEffect } from 'react'
import { MessageCircle, X, Send, Loader2, Bot, User } from 'lucide-react'
import { chatService } from '../../services/chat'

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
}

const ChatBot: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false)
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  useEffect(() => {
    if (isOpen && inputRef.current) {
      inputRef.current.focus()
    }
  }, [isOpen])

  const handleSend = async () => {
    const trimmed = input.trim()
    if (!trimmed || isLoading) return

    const userMsg: Message = {
      role: 'user',
      content: trimmed,
      timestamp: new Date(),
    }
    setMessages((prev) => [...prev, userMsg])
    setInput('')
    setIsLoading(true)

    try {
      const response = await chatService.sendMessage(trimmed)
      const assistantMsg: Message = {
        role: 'assistant',
        content: response.reply,
        timestamp: new Date(),
      }
      setMessages((prev) => [...prev, assistantMsg])
    } catch {
      const errorMsg: Message = {
        role: 'assistant',
        content: '抱歉，AI 服务暂时不可用，请稍后再试。',
        timestamp: new Date(),
      }
      setMessages((prev) => [...prev, errorMsg])
    } finally {
      setIsLoading(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }

  const renderContent = (content: string) => {
    const lines = content.split('\n')
    return lines.map((line, i) => {
      // 处理加粗
      const boldParts = line.split(/(\*\*[^*]+\*\*)/g)
      const rendered = boldParts.map((part, j) => {
        if (part.startsWith('**') && part.endsWith('**')) {
          return <strong key={j}>{part.slice(2, -2)}</strong>
        }
        return part
      })
      return (
        <p key={i} className={line === '' ? 'h-3' : ''}>
          {line === '' ? null : rendered}
        </p>
      )
    })
  }

  return (
    <>
      {/* 浮动按钮 */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className={`fixed bottom-6 right-6 z-50 p-4 rounded-full shadow-lg transition-all duration-300 hover:scale-105 ${
          isOpen
            ? 'bg-gray-200 text-gray-600 hover:bg-gray-300'
            : 'bg-gradient-to-r from-blue-500 to-purple-500 text-white hover:from-blue-600 hover:to-purple-600'
        }`}
        title={isOpen ? '关闭助手' : '饮食小助手'}
      >
        {isOpen ? <X size={24} /> : <MessageCircle size={24} />}
      </button>

      {/* 聊天面板 */}
      {isOpen && (
        <div className="fixed bottom-24 right-6 z-50 w-full max-w-md bg-white rounded-2xl shadow-2xl border flex flex-col overflow-hidden animate-in slide-in-from-bottom-4 duration-300"
          style={{ height: '560px', maxHeight: 'calc(100vh - 140px)' }}
        >
          {/* 头部 */}
          <div className="flex items-center justify-between px-5 py-4 bg-gradient-to-r from-blue-500 to-purple-500 text-white">
            <div className="flex items-center space-x-3">
              <div className="w-9 h-9 bg-white/20 rounded-full flex items-center justify-center">
                <Bot size={20} />
              </div>
              <div>
                <h3 className="font-semibold text-sm">饮食小助手</h3>
                <p className="text-xs text-white/70">AI 营养顾问</p>
              </div>
            </div>
            <button
              onClick={() => setIsOpen(false)}
              className="p-1.5 rounded-full hover:bg-white/20 transition"
            >
              <X size={18} />
            </button>
          </div>

          {/* 消息列表 */}
          <div className="flex-1 overflow-y-auto px-4 py-4 space-y-4 bg-gray-50">
            {messages.length === 0 && (
              <div className="flex flex-col items-center justify-center h-full text-gray-400 space-y-3">
                <Bot size={48} strokeWidth={1.5} />
                <div className="text-center">
                  <p className="font-medium text-gray-500">你好！我是饮食小助手 👋</p>
                  <p className="text-sm mt-1">可以根据你的饮食记录提供个性化建议</p>
                  <div className="mt-4 space-y-2 text-xs text-gray-400">
                    <p>试试问我：</p>
                    <button
                      onClick={() => setInput('帮我分析一下我最近的饮食结构')}
                      className="block w-full text-left px-3 py-1.5 bg-white rounded-lg border hover:border-blue-300 hover:text-blue-500 transition"
                    >
                      "帮我分析一下我最近的饮食结构"
                    </button>
                    <button
                      onClick={() => setInput('我最近在减脂，推荐一下今天晚餐吃什么')}
                      className="block w-full text-left px-3 py-1.5 bg-white rounded-lg border hover:border-blue-300 hover:text-blue-500 transition"
                    >
                      "我最近在减脂，推荐一下今天晚餐吃什么"
                    </button>
                    <button
                      onClick={() => setInput('鸡胸肉的热量大概是多少？')}
                      className="block w-full text-left px-3 py-1.5 bg-white rounded-lg border hover:border-blue-300 hover:text-blue-500 transition"
                    >
                      "鸡胸肉的热量大概是多少？"
                    </button>
                  </div>
                </div>
              </div>
            )}

            {messages.map((msg, idx) => (
              <div
                key={idx}
                className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
              >
                <div className={`flex items-start space-x-2 max-w-[85%] ${msg.role === 'user' ? 'flex-row-reverse space-x-reverse' : ''}`}>
                  <div
                    className={`w-7 h-7 rounded-full flex items-center justify-center flex-shrink-0 ${
                      msg.role === 'user'
                        ? 'bg-gradient-to-r from-blue-500 to-purple-500'
                        : 'bg-green-500'
                    }`}
                  >
                    {msg.role === 'user' ? (
                      <User size={14} className="text-white" />
                    ) : (
                      <Bot size={14} className="text-white" />
                    )}
                  </div>
                  <div
                    className={`px-4 py-2.5 rounded-2xl text-sm leading-relaxed ${
                      msg.role === 'user'
                        ? 'bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-tr-md'
                        : 'bg-white border text-gray-700 rounded-tl-md shadow-sm'
                    }`}
                  >
                    <div className={msg.role === 'user' ? 'text-white' : 'text-gray-700'}>
                      {renderContent(msg.content)}
                    </div>
                    <span
                      className={`text-xs mt-1 block ${
                        msg.role === 'user' ? 'text-white/60' : 'text-gray-400'
                      }`}
                    >
                      {formatTime(msg.timestamp)}
                    </span>
                  </div>
                </div>
              </div>
            ))}

            {/* 加载动画 */}
            {isLoading && (
              <div className="flex justify-start">
                <div className="flex items-start space-x-2">
                  <div className="w-7 h-7 rounded-full bg-green-500 flex items-center justify-center flex-shrink-0">
                    <Bot size={14} className="text-white" />
                  </div>
                  <div className="px-4 py-3 bg-white border rounded-2xl rounded-tl-md shadow-sm">
                    <div className="flex space-x-1.5">
                      <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                      <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
                      <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
                    </div>
                  </div>
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {/* 输入区域 */}
          <div className="px-4 py-3 bg-white border-t">
            <div className="flex items-end space-x-2">
              <textarea
                ref={inputRef}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="输入你的问题..."
                rows={1}
                className="flex-1 resize-none rounded-xl border border-gray-200 px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent max-h-32"
                style={{ minHeight: '40px' }}
                disabled={isLoading}
                onInput={(e) => {
                  const target = e.target as HTMLTextAreaElement
                  target.style.height = 'auto'
                  target.style.height = Math.min(target.scrollHeight, 128) + 'px'
                }}
              />
              <button
                onClick={handleSend}
                disabled={!input.trim() || isLoading}
                className="p-2.5 rounded-xl bg-gradient-to-r from-blue-500 to-purple-500 text-white disabled:opacity-40 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-600 transition flex-shrink-0"
              >
                {isLoading ? (
                  <Loader2 size={18} className="animate-spin" />
                ) : (
                  <Send size={18} />
                )}
              </button>
            </div>
            <p className="text-xs text-gray-400 mt-1.5 text-center">
              按 Enter 发送，Shift+Enter 换行
            </p>
          </div>
        </div>
      )}
    </>
  )
}

export default ChatBot
