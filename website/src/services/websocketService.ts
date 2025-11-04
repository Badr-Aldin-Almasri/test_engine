import { ExecutionStatus } from '@/types'

export interface WebSocketMessage {
  type: 'node_update' | 'test_run_complete'
  testRunId: string
  nodeId?: string
  status?: string
  output?: any
  error?: string
  duration?: number
}

export type WebSocketCallback = (message: WebSocketMessage) => void

class WebSocketService {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  private callbacks: Set<WebSocketCallback> = new Set()

  connect(testRunId: string, onMessage: WebSocketCallback) {
    const token = localStorage.getItem('auth_token')
    if (!token) {
      console.error('No auth token available')
      return
    }

    // Close existing connection if any
    if (this.ws) {
      this.ws.close()
    }

    const wsUrl = `ws://localhost:8080/api/ws?testRunId=${testRunId}`
    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.reconnectAttempts = 0
    }

    this.ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)
        onMessage(message)
        // Also notify all registered callbacks
        this.callbacks.forEach((callback) => callback(message))
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
      this.attemptReconnect(testRunId, onMessage)
    }
  }

  private attemptReconnect(testRunId: string, onMessage: WebSocketCallback) {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      setTimeout(() => {
        console.log(`Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
        this.connect(testRunId, onMessage)
      }, this.reconnectDelay * this.reconnectAttempts)
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.callbacks.clear()
  }

  subscribe(callback: WebSocketCallback) {
    this.callbacks.add(callback)
    return () => {
      this.callbacks.delete(callback)
    }
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
}

export const websocketService = new WebSocketService()

