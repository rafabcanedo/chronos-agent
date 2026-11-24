import { useRef, useState } from 'react'
import { ArrowUpIcon } from 'lucide-react'
import { useChatHistory } from '@/hooks/useChatHistory'
import { ChatEmptyState } from '@/components/chronos/chat-empty-state'
import { ChatHeader } from '@/components/chat/header'
import { Card, CardContent } from '@/components/ui/card'
import { Bubble, BubbleContent } from '@/components/ui/bubble'
import { InputGroup, InputGroupAddon, InputGroupButton, InputGroupTextarea } from '@/components/ui/input-group'
import { Message, MessageContent } from '@/components/ui/message'
import { MessageScroller, MessageScrollerButton, MessageScrollerContent, MessageScrollerItem, MessageScrollerProvider, MessageScrollerViewport } from '@/components/ui/message-scroller'

export default function Chat() {
  const { messages, addMessage } = useChatHistory()
  const [text, setText] = useState('')
  const formRef = useRef<HTMLFormElement>(null)

  const canSend = text.trim().length > 0

  function onSubmit(e: { preventDefault(): void }) {
    e.preventDefault()
    if (!canSend) return
    addMessage({ role: 'user', content: text.trim() })
    setText('')
  }

  return (
    <MessageScrollerProvider>
      <div className="h-screen p-6">
        <Card className="w-full h-full shadow-md">
          <CardContent className="flex flex-col h-full p-0">
            <ChatHeader />
            <div className="flex-1 min-h-0">
              {messages.length === 0 ? (
                <ChatEmptyState />
              ) : (
                <MessageScroller>
                  <MessageScrollerViewport>
                    <MessageScrollerContent className="gap-4 px-6 py-4">
                      {messages.map((msg) => (
                        <MessageScrollerItem
                          key={msg.id}
                          messageId={msg.id}
                          scrollAnchor={msg.role === 'user'}
                        >
                          {msg.role === 'user' ? (
                            <Message align="end">
                              <MessageContent>
                                <Bubble align="end" variant="muted">
                                  <BubbleContent>{msg.content}</BubbleContent>
                                </Bubble>
                              </MessageContent>
                            </Message>
                          ) : (
                            <Message>
                              <MessageContent>
                                <p>{msg.content}</p>
                              </MessageContent>
                            </Message>
                          )}
                        </MessageScrollerItem>
                      ))}
                    </MessageScrollerContent>
                  </MessageScrollerViewport>
                  <MessageScrollerButton />
                </MessageScroller>
              )}
            </div>
            <div className="p-4 border-t">
              <form ref={formRef} onSubmit={onSubmit}>
                <InputGroup>
                  <InputGroupTextarea
                    value={text}
                    onChange={(e) => setText(e.target.value)}
                    placeholder="Tell me what you paid or spent..."
                    onKeyDown={(e) => {
                      if (e.key === 'Enter' && !e.shiftKey) {
                        e.preventDefault()
                        formRef.current?.requestSubmit()
                      }
                    }}
                  />
                  <InputGroupAddon align="block-end">
                    <InputGroupButton
                      type="submit"
                      size="icon-sm"
                      variant="default"
                      disabled={!canSend}
                      className="ml-auto"
                    >
                      <ArrowUpIcon />
                    </InputGroupButton>
                  </InputGroupAddon>
                </InputGroup>
              </form>
            </div>
          </CardContent>
        </Card>
      </div>
    </MessageScrollerProvider>
  )
}
