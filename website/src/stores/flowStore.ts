import { create } from 'zustand'
import { Flow, FlowNode, FlowEdge, NodeData, NodeConfig, ExecutionStatus, NodeType } from '@/types'
import { flowService } from '@/services/flowService'

interface FlowState {
  flows: Flow[]
  currentFlow: Flow | null
  selectedNode: NodeData | null
  isRunning: boolean
  executionHistory: Record<string, ExecutionStatus>
  isLoading: boolean
  
  // Actions
  loadFlows: () => Promise<void>
  setCurrentFlow: (flow: Flow | null) => void
  createFlow: (flow: Partial<Flow>) => Promise<Flow>
  updateFlow: (flowId: string, updates: Partial<Flow>) => Promise<void>
  deleteFlow: (flowId: string) => Promise<void>
  saveFlow: () => Promise<void>
  addNode: (node: FlowNode) => void
  updateNode: (nodeId: string, updates: Partial<NodeData>) => void
  deleteNode: (nodeId: string) => void
  addEdge: (edge: FlowEdge) => void
  deleteEdge: (edgeId: string) => void
  setSelectedNode: (node: NodeData | null) => void
  updateNodeConfig: (nodeId: string, config: NodeConfig) => void
  setNodeStatus: (nodeId: string, status: ExecutionStatus, output?: any, error?: string) => void
  setIsRunning: (isRunning: boolean) => void
  resetNodeStatuses: () => void
}

export const useFlowStore = create<FlowState>((set, get) => ({
  flows: [],
  currentFlow: null,
  selectedNode: null,
  isRunning: false,
  executionHistory: {},
  isLoading: false,

  loadFlows: async () => {
    set({ isLoading: true })
    try {
      const flows = await flowService.getFlows()
      set({ flows: flows || [], isLoading: false })
    } catch (error) {
      console.error('Failed to load flows:', error)
      set({ flows: [], isLoading: false })
    }
  },

  setCurrentFlow: (flow) => set({ currentFlow: flow }),

  createFlow: async (flowData) => {
    try {
      const flow = await flowService.createFlow(flowData)
      set((state) => ({
        flows: [...state.flows, flow],
        currentFlow: flow
      }))
      return flow
    } catch (error) {
      console.error('Failed to create flow:', error)
      throw error
    }
  },

  updateFlow: async (flowId, updates) => {
    try {
      const updatedFlow = await flowService.updateFlow(flowId, updates)
      set((state) => ({
        flows: state.flows.map(f => f.id === flowId ? updatedFlow : f),
        currentFlow: state.currentFlow?.id === flowId ? updatedFlow : state.currentFlow
      }))
    } catch (error) {
      console.error('Failed to update flow:', error)
      throw error
    }
  },

  deleteFlow: async (flowId) => {
    try {
      await flowService.deleteFlow(flowId)
      set((state) => ({
        flows: state.flows.filter(f => f.id !== flowId),
        currentFlow: state.currentFlow?.id === flowId ? null : state.currentFlow
      }))
    } catch (error) {
      console.error('Failed to delete flow:', error)
      throw error
    }
  },

  saveFlow: async () => {
    const { currentFlow } = get()
    if (!currentFlow) return

    try {
      await get().updateFlow(currentFlow.id, currentFlow)
    } catch (error) {
      console.error('Failed to save flow:', error)
      throw error
    }
  },

  addNode: (node) => set((state) => {
    if (!state.currentFlow) return state
    // Ensure nodes array exists
    const currentNodes = state.currentFlow.nodes || []
    const updatedFlow = {
      ...state.currentFlow,
      nodes: [...currentNodes, node]
    }
    return {
      currentFlow: updatedFlow,
      flows: state.flows.map(f => f.id === updatedFlow.id ? updatedFlow : f)
    }
  }),

  updateNode: (nodeId, updates) => set((state) => {
    if (!state.currentFlow) return state
    const updatedFlow = {
      ...state.currentFlow,
      nodes: state.currentFlow.nodes.map(n => 
        n.id === nodeId ? { ...n, data: { ...n.data, ...updates } } : n
      )
    }
    return {
      currentFlow: updatedFlow,
      flows: state.flows.map(f => f.id === updatedFlow.id ? updatedFlow : f)
    }
  }),

  deleteNode: (nodeId) => set((state) => {
    if (!state.currentFlow) return state
    const updatedFlow = {
      ...state.currentFlow,
      nodes: state.currentFlow.nodes.filter(n => n.id !== nodeId),
      edges: state.currentFlow.edges.filter(e => e.source !== nodeId && e.target !== nodeId)
    }
    return {
      currentFlow: updatedFlow,
      flows: state.flows.map(f => f.id === updatedFlow.id ? updatedFlow : f)
    }
  }),

  addEdge: (edge) => set((state) => {
    if (!state.currentFlow) return state
    const updatedFlow = {
      ...state.currentFlow,
      edges: [...state.currentFlow.edges, edge]
    }
    return {
      currentFlow: updatedFlow,
      flows: state.flows.map(f => f.id === updatedFlow.id ? updatedFlow : f)
    }
  }),

  deleteEdge: (edgeId) => set((state) => {
    if (!state.currentFlow) return state
    const updatedFlow = {
      ...state.currentFlow,
      edges: state.currentFlow.edges.filter(e => e.id !== edgeId)
    }
    return {
      currentFlow: updatedFlow,
      flows: state.flows.map(f => f.id === updatedFlow.id ? updatedFlow : f)
    }
  }),

  setSelectedNode: (node) => set({ selectedNode: node }),

  updateNodeConfig: (nodeId, config) => set((state) => {
    if (!state.currentFlow) return state
    const updatedFlow = {
      ...state.currentFlow,
      nodes: state.currentFlow.nodes.map(n => 
        n.id === nodeId ? { ...n, data: { ...n.data, config } } : n
      )
    }
    return {
      currentFlow: updatedFlow,
      flows: state.flows.map(f => f.id === updatedFlow.id ? updatedFlow : f)
    }
  }),

  setNodeStatus: (nodeId, status, output, error) => set((state) => {
    if (!state.currentFlow) return state
    const updatedFlow = {
      ...state.currentFlow,
      nodes: state.currentFlow.nodes.map(n => 
        n.id === nodeId 
          ? { ...n, data: { ...n.data, status, output, error } } 
          : n
      )
    }
    return {
      currentFlow: updatedFlow,
      flows: state.flows.map(f => f.id === updatedFlow.id ? updatedFlow : f),
      executionHistory: { ...state.executionHistory, [nodeId]: status }
    }
  }),

  setIsRunning: (isRunning) => set({ isRunning }),

  resetNodeStatuses: () => set((state) => {
    if (!state.currentFlow) return state
    const updatedFlow = {
      ...state.currentFlow,
      nodes: state.currentFlow.nodes.map(n => ({
        ...n,
        data: { ...n.data, status: undefined, output: undefined, error: undefined }
      }))
    }
    return {
      currentFlow: updatedFlow,
      flows: state.flows.map(f => f.id === updatedFlow.id ? updatedFlow : f),
      executionHistory: {}
    }
  }),
}))
