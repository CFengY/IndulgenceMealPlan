import React, { useState, useEffect, useRef } from 'react'
import { useAuth } from '../../contexts/AuthContext'
import { postService, type Post } from '../../services/post'
import { ImagePlus, X, Trash2, Send, Camera } from 'lucide-react'

function timeAgo(dateStr: string): string {
  const now = Date.now()
  const then = new Date(dateStr).getTime()
  const diff = now - then
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}天前`
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

const Timeline: React.FC = () => {
  const { user } = useAuth()
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [content, setContent] = useState('')
  const [images, setImages] = useState<FileList | null>(null)
  const [previewUrls, setPreviewUrls] = useState<string[]>([])
  const [submitting, setSubmitting] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const fetchTimeline = async () => {
    try {
      setLoading(true)
      const data = await postService.getTimeline()
      setPosts(data)
    } catch (err: any) {
      setError(err.response?.data?.error || '获取动态列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchTimeline()
  }, [])

  useEffect(() => {
    return () => {
      previewUrls.forEach((url) => URL.revokeObjectURL(url))
    }
  }, [previewUrls])

  const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (!files || files.length === 0) return

    const newPreviews: string[] = []
    for (let i = 0; i < files.length; i++) {
      newPreviews.push(URL.createObjectURL(files[i]))
    }

    previewUrls.forEach((url) => URL.revokeObjectURL(url))
    setImages(files)
    setPreviewUrls(newPreviews)
  }

  const removeImage = (index: number) => {
    URL.revokeObjectURL(previewUrls[index])
    const newPreviews = previewUrls.filter((_, i) => i !== index)

    if (newPreviews.length === 0) {
      setImages(null)
      setPreviewUrls([])
      if (fileInputRef.current) fileInputRef.current.value = ''
      return
    }

    const dt = new DataTransfer()
    previewUrls.forEach((_, i) => {
      if (i !== index && images) dt.items.add(images[i])
    })
    setImages(dt.files)
    setPreviewUrls(newPreviews)
  }

  const handleSubmit = async () => {
    const trimmed = content.trim()
    if (!trimmed || submitting) return

    try {
      setSubmitting(true)
      await postService.createPost(trimmed, images || undefined)
      setContent('')
      setImages(null)
      setPreviewUrls([])
      if (fileInputRef.current) fileInputRef.current.value = ''
      await fetchTimeline()
    } catch (err: any) {
      setError(err.response?.data?.error || '发布失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async (postId: number) => {
    if (!window.confirm('确定要删除这条动态吗？')) return
    try {
      await postService.deletePost(postId)
      setPosts((prev) => prev.filter((p) => p.ID !== postId))
    } catch (err: any) {
      setError(err.response?.data?.error || '删除失败')
    }
  }

  const renderImages = (imageStr: string) => {
    const imageList = imageStr.split(',').filter(Boolean)
    if (imageList.length === 0) return null

    const gridClass =
      imageList.length === 1
        ? 'grid-cols-1'
        : imageList.length === 2
          ? 'grid-cols-2'
          : 'grid-cols-3'

    return (
      <div className={`grid ${gridClass} gap-2 mt-3`}>
        {imageList.map((filename, idx) => (
          <img
            key={idx}
            src={postService.getImageUrl(filename.trim())}
            alt={`图片 ${idx + 1}`}
            className="w-full h-48 object-cover rounded-lg cursor-pointer hover:opacity-90 transition"
            loading="lazy"
            onClick={() => {
              window.open(postService.getImageUrl(filename.trim()), '_blank')
            }}
          />
        ))}
      </div>
    )
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-800 dark:text-gray-100">动态</h1>
        <p className="text-gray-600 dark:text-gray-400 dark:text-gray-400">分享你的美食时刻和心情</p>
      </div>

      {/* 发布区 */}
      <div className="bg-gray-50 dark:bg-gray-900 rounded-xl p-4 mb-8 border">
        <div className="flex items-start space-x-3">
          <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center flex-shrink-0">
            <span className="text-white font-medium text-sm">
              {(user?.username || '?')[0].toUpperCase()}
            </span>
          </div>
          <div className="flex-1">
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="分享你的美食时刻..."
              rows={3}
              className="w-full resize-none rounded-xl border border-gray-200 dark:border-gray-700 px-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              disabled={submitting}
            />

            {/* 图片预览 */}
            {previewUrls.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-2">
                {previewUrls.map((url, idx) => (
                  <div key={idx} className="relative">
                    <img
                      src={url}
                      alt={`预览 ${idx + 1}`}
                      className="w-20 h-20 object-cover rounded-lg"
                    />
                    <button
                      onClick={() => removeImage(idx)}
                      className="absolute -top-1.5 -right-1.5 p-0.5 bg-red-500 text-white rounded-full hover:bg-red-600"
                    >
                      <X size={14} />
                    </button>
                  </div>
                ))}
              </div>
            )}

            <div className="flex items-center justify-between mt-3">
              <button
                onClick={() => fileInputRef.current?.click()}
                className="flex items-center space-x-1.5 px-3 py-1.5 text-sm text-gray-500 dark:text-gray-400 hover:text-blue-500 hover:bg-blue-50 rounded-lg transition"
                disabled={submitting}
              >
                <Camera size={18} />
                <span>图片</span>
              </button>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/jpeg,image/png,image/gif,image/webp"
                multiple
                className="hidden"
                onChange={handleImageSelect}
              />

              <button
                onClick={handleSubmit}
                disabled={!content.trim() || submitting}
                className="px-5 py-2 bg-gradient-to-r from-blue-500 to-purple-500 text-white text-sm font-medium rounded-lg disabled:opacity-40 disabled:cursor-not-allowed hover:from-blue-600 hover:to-purple-600 transition flex items-center space-x-2"
              >
                {submitting ? (
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                ) : (
                  <Send size={16} />
                )}
                <span>发布</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* 错误提示 */}
      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-600 flex items-center justify-between">
          <span>{error}</span>
          <button onClick={() => setError('')} className="text-red-400 hover:text-red-600">
            <X size={18} />
          </button>
        </div>
      )}

      {/* 加载态 */}
      {loading && (
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500" />
        </div>
      )}

      {/* 空态 */}
      {!loading && !error && posts.length === 0 && (
        <div className="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
          <Camera size={64} strokeWidth={1} />
          <p className="mt-4 text-lg font-medium text-gray-500 dark:text-gray-400">还没有动态</p>
          <p className="text-sm mt-1">发布第一条动态吧!</p>
        </div>
      )}

      {/* 时间线 */}
      {!loading && posts.length > 0 && (
        <div className="space-y-6">
          {posts.map((post) => (
            <div key={post.ID} className="bg-white dark:bg-gray-800 rounded-xl border p-5">
              <div className="flex items-start space-x-3">
                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center flex-shrink-0">
                  <span className="text-white font-medium text-sm">
                    {(post.user?.username || '?')[0].toUpperCase()}
                  </span>
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <span className="font-semibold text-gray-800 dark:text-gray-100">
                        {post.user?.username || '未知用户'}
                      </span>
                      <span className="text-xs text-gray-400 dark:text-gray-500">
                        {timeAgo(post.CreatedAt)}
                      </span>
                    </div>
                    {user && post.user_id === user.id && (
                      <button
                        onClick={() => handleDelete(post.ID)}
                        className="p-1.5 text-gray-400 dark:text-gray-500 hover:text-red-500 hover:bg-red-50 rounded-lg transition"
                        title="删除"
                      >
                        <Trash2 size={16} />
                      </button>
                    )}
                  </div>
                  <p className="mt-2 text-gray-700 dark:text-gray-200 whitespace-pre-wrap break-words">
                    {post.content}
                  </p>
                  {post.images && renderImages(post.images)}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default Timeline
