// Flow and Node Types
export enum NodeType {
  API = 'api',
  VERIFICATION = 'verification',
  MOCK = 'mock',
  REPORT = 'report',
  EVENT_TRIGGER = 'event_trigger',
}

export enum ExecutionStatus {
  PENDING = 'pending',
  RUNNING = 'running',
  SUCCESS = 'success',
  FAILED = 'failed',
  SKIPPED = 'skipped',
}

export interface NodeData {
  id: string
  type: NodeType
  label: string
  status?: ExecutionStatus
  config: NodeConfig
  output?: any
  error?: string
}

export interface NodeConfig {
  // API Node
  method?: string
  url?: string
  headers?: Record<string, string>
  body?: string
  
  // Verification Node
  expected?: any
  assertionType?: 'equals' | 'contains' | 'regex' | 'custom'
  customScript?: string
  
  // Mock Node
  mockResponse?: any
  mockDelay?: number
  
  // Report Node
  reportName?: string
  
  // Event Trigger Node
  triggerCondition?: string
  triggerEvent?: string
}

export interface FlowNode {
  id: string
  type: string
  position: { x: number; y: number }
  data: NodeData
}

export interface FlowEdge {
  id: string
  source: string
  target: string
  sourceHandle?: string
  targetHandle?: string
}

export interface Flow {
  id: string
  name: string
  description?: string
  tags?: string[]
  nodes: FlowNode[]
  edges: FlowEdge[]
  createdAt: string
  updatedAt: string
}

export interface TestRun {
  id: string
  flow_id: string
  flow_name?: string
  status: ExecutionStatus
  started_at: string
  completed_at?: string
  duration_ms?: number
  node_results: Record<string, {
    status: ExecutionStatus
    output?: any
    error?: string
    duration: number
  }>
  error?: string
  created_at: string
}

