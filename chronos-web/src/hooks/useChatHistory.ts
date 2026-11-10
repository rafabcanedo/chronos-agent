import { useState, useEffect } from 'react'
import { CHAT_STORAGE_KEY } from '@/contants'
import type { ChatMessage } from '@/types'

export function useChatHistory() {
  const [messages, setMessages] = useState<ChatMessage[]>([])

  useEffect(() => {
    const saved = sessionStorage.getItem(CHAT_STORAGE_KEY)
    if (!saved) return
    try {
      setMessages(JSON.parse(saved))
    } catch {}
  }, [])

  useEffect(() => {
    sessionStorage.setItem(CHAT_STORAGE_KEY, JSON.stringify(messages))
  }, [messages])

  function addMessage(message: Omit<ChatMessage, 'id'>) {
    setMessages(prev => [...prev, { ...message, id: crypto.randomUUID() }])
  }

  return { messages, addMessage }
}
