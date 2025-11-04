import React, { useCallback, useEffect, useRef, useState } from 'react'
import ReactFlow, {
  Node,
  Edge,
  addEdge,
  Connection,
  useNodesState,
  useEdgesState,
  Background,
  Controls,
  MiniMap,
  Panel,
} from 'reactflow'
import 'reactflow/dist/style.css'
import { CustomNode, nodeTypes } from './nodes/CustomNode'
import { NodeSettings } from './nodes/NodeSettings'
import { useFlowStore } from '@/stores/flowStore'
import { Button } from './ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog'
import { Play, Square, Plus, Save, Trash2 } from 'lucide-react'
import { NodeType, NodeData, ExecutionStatus } from '@/types'
import { flowService } from '@/services/flowService'
import { websocketService, WebSocketMessage } from '@/services/websocketService'
import { useNavigate } from 'react-router-dom'

interface FlowEditorProps {
  flowId: string
}

export const FlowEditor: React.FC<FlowEditorProps> = ({ flowId }) => {
  const {
    currentFlow,
    selectedNode,
    isRunning,
    setCurrentFlow,
    setSelectedNode,
    updateNodeConfig,
    setNodeStatus,
    setIsRunning,
    resetNodeStatuses,
    addEdge: addFlowEdge,
    addNode: addFlowNode,
    saveFlow,
    deleteFlow,
    deleteNode,
  } = useFlowStore()
  const navigate = useNavigate()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const [nodes, setNodes, onNodesChange] = useNodesState([])
  const [edges, setEdges, onEdgesChange] = useEdgesState([])
  const [settingsOpen, setSettingsOpen] = useState(false)
  const testRunIdRef = useRef<string | null>(null)
  const isInitializedRef = useRef(false)

  // Load flow when flowId changes
  useEffect(() => {
    if (flowId && (!currentFlow || currentFlow.id !== flowId)) {
      flowService.getFlow(flowId)
        .then((flow) => {
          setCurrentFlow(flow)
        })
        .catch((error) => {
          console.error('Failed to load flow:', error)
          navigate('/')
        })
    }
  }, [flowId, currentFlow?.id, setCurrentFlow, navigate])

  // Track previous flow ID to detect when flow actually changes
  const prevFlowIdRef = useRef<string | null>(null)
  
  // Sync store nodes/edges with ReactFlow (only when flow ID actually changes)
  useEffect(() => {
    const flowId = currentFlow?.id
    if (flowId && flowId !== prevFlowIdRef.current) {
      const flowNodes = (currentFlow.nodes || []) as Node[]
      const flowEdges = (currentFlow.edges || []) as Edge[]
      setNodes(flowNodes)
      setEdges(flowEdges)
      prevFlowIdRef.current = flowId
      isInitializedRef.current = true
    }
  }, [currentFlow?.id]) // Only sync when flow ID changes

  // Sync ReactFlow changes back to store (debounced)
  // This runs when nodes/edges change, but we use a ref to prevent infinite loops
  const isSyncingRef = useRef(false)
  
  useEffect(() => {
    if (!currentFlow || !isInitializedRef.current || isSyncingRef.current) return

    const timeoutId = setTimeout(() => {
      isSyncingRef.current = true
      
      // Update store with current ReactFlow state
      const updatedFlow = {
        ...currentFlow,
        nodes: nodes as any,
        edges: edges as any,
      }
      
      // Update store directly
      useFlowStore.setState((state) => ({
        currentFlow: updatedFlow,
        flows: state.flows.map(f => f.id === updatedFlow.id ? updatedFlow : f)
      }))

      // Reset sync flag after a short delay
      setTimeout(() => {
        isSyncingRef.current = false
      }, 100)

      // Auto-save after updating store
      setTimeout(() => {
        useFlowStore.getState().saveFlow().catch(console.error)
      }, 500)
    }, 1500) // Debounce to avoid too many updates

    return () => clearTimeout(timeoutId)
  }, [nodes.length, edges.length]) // Only depend on lengths to avoid re-running on every change

  const onConnect = useCallback(
    (params: Connection) => {
      const newEdge = {
        id: `edge-${params.source}-${params.target}`,
        source: params.source!,
        target: params.target!,
        sourceHandle: params.sourceHandle,
        targetHandle: params.targetHandle,
      }
      addFlowEdge(newEdge)
      setEdges((eds) => addEdge(params, eds))
    },
    [addFlowEdge, setEdges]
  )

  const onNodeClick = useCallback(
    (_: React.MouseEvent, node: Node) => {
      setSelectedNode(node.data as NodeData)
      setSettingsOpen(true)
    },
    [setSelectedNode]
  )

  const handleSaveConfig = async (config: any) => {
    if (selectedNode) {
      updateNodeConfig(selectedNode.id, config)
      setNodes((nds) =>
        nds.map((n) =>
          n.id === selectedNode.id
            ? { ...n, data: { ...n.data, config } }
            : n
        )
      )
      await saveFlow()
    }
  }

  const handleDeleteNode = async () => {
    if (selectedNode) {
      // Remove from ReactFlow state
      setNodes((nds) => nds.filter((n) => n.id !== selectedNode.id))
      // Remove connected edges
      setEdges((eds) => 
        eds.filter((e) => e.source !== selectedNode.id && e.target !== selectedNode.id)
      )
      // Remove from store
      deleteNode(selectedNode.id)
      // Save to backend
      await saveFlow()
      // Close settings and clear selection
      setSettingsOpen(false)
      setSelectedNode(null)
    }
  }

  const handleRunFlow = async () => {
    if (!currentFlow) return

    resetNodeStatuses()
    setIsRunning(true)

    try {
      // Save flow first to ensure latest version is on server
      await saveFlow()

      // Trigger flow execution on backend
      // Note: Backend should return test run ID in response for WebSocket connection
      // For now, we'll poll for the latest test run
      await flowService.runFlow(currentFlow.id)

      // Poll for test run ID (in production, backend should return it immediately)
      // This is a workaround - ideally the backend would return test run ID
      setTimeout(async () => {
        try {
          const testRuns = await flowService.getTestRuns(currentFlow.id)
          if (testRuns.length > 0) {
            const latestRun = testRuns[0]
            testRunIdRef.current = latestRun.id
            // WebSocket will connect in the useEffect hook
          }
        } catch (error) {
          console.error('Failed to get test run:', error)
        }
      }, 1000)
      
    } catch (error: any) {
      console.error('Failed to run flow:', error)
      setIsRunning(false)
    }
  }

  const handleStopFlow = () => {
    setIsRunning(false)
    websocketService.disconnect()
    testRunIdRef.current = null
  }

  // Set up WebSocket listener for execution updates
  useEffect(() => {
    if (!isRunning || !testRunIdRef.current) return

    const handleMessage = (message: WebSocketMessage) => {
      if (message.type === 'node_update' && message.nodeId) {
        const status = message.status as ExecutionStatus
        setNodeStatus(message.nodeId!, status, message.output, message.error)
        setNodes((nds) =>
          nds.map((n) =>
            n.id === message.nodeId
              ? {
                  ...n,
                  data: {
                    ...n.data,
                    status,
                    output: message.output,
                    error: message.error,
                  },
                }
              : n
          )
        )
      } else if (message.type === 'test_run_complete') {
        setIsRunning(false)
        websocketService.disconnect()
        testRunIdRef.current = null
      }
    }

    websocketService.connect(testRunIdRef.current, handleMessage)

    return () => {
      websocketService.disconnect()
    }
  }, [isRunning, setNodeStatus, setIsRunning, setNodes])

  const addNode = (type: NodeType, position: { x: number; y: number }) => {
    if (!currentFlow) return

    const newNodeId = `node-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    const nodeLabels: Record<NodeType, string> = {
      [NodeType.API]: 'API Call',
      [NodeType.VERIFICATION]: 'Verification',
      [NodeType.MOCK]: 'Mock Response',
      [NodeType.REPORT]: 'Report',
      [NodeType.EVENT_TRIGGER]: 'Event Trigger',
    }

    const newNode: Node<NodeData> = {
      id: newNodeId,
      type: 'custom',
      position,
      data: {
        id: newNodeId,
        type,
        label: nodeLabels[type],
        config: {},
      },
    }

    // Add to ReactFlow state first (immediate UI update)
    setNodes((nds) => [...nds, newNode])
    
    // Then update store (will trigger sync effect which will save)
    addFlowNode(newNode as any)
  }

  const handleSave = async () => {
    try {
      await saveFlow()
      alert('Flow saved successfully!')
    } catch (error) {
      console.error('Failed to save flow:', error)
      alert('Failed to save flow')
    }
  }

  const handleDelete = async () => {
    if (!currentFlow) return

    try {
      await deleteFlow(currentFlow.id)
      setDeleteDialogOpen(false)
      navigate('/')
    } catch (error) {
      console.error('Failed to delete flow:', error)
      alert('Failed to delete flow')
    }
  }

  if (!currentFlow) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center">
          <p className="text-muted-foreground">Loading flow...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="w-full h-screen relative">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onNodeClick={onNodeClick}
        nodeTypes={nodeTypes}
        fitView
      >
        <Background />
        <Controls />
        <MiniMap />
        
        <Panel position="top-left" className="bg-white p-2 rounded shadow-lg">
          <div className="flex gap-2">
            <Button
              onClick={handleSave}
              variant="outline"
              size="sm"
            >
              <Save className="h-4 w-4 mr-2" />
              Save
            </Button>
            <Button
              onClick={handleRunFlow}
              disabled={isRunning || !currentFlow}
              size="sm"
            >
              <Play className="h-4 w-4 mr-2" />
              Run Flow
            </Button>
            {isRunning && (
              <Button onClick={handleStopFlow} variant="destructive" size="sm">
                <Square className="h-4 w-4 mr-2" />
                Stop
              </Button>
            )}
            <Button
              onClick={() => setDeleteDialogOpen(true)}
              variant="destructive"
              size="sm"
              disabled={!currentFlow}
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Delete
            </Button>
          </div>
        </Panel>

        <Panel position="top-right" className="bg-white p-2 rounded shadow-lg">
          <div className="flex flex-col gap-1">
            {Object.values(NodeType).map((type) => (
              <Button
                key={type}
                onClick={() => {
                  const centerX = window.innerWidth / 2
                  const centerY = window.innerHeight / 2
                  addNode(type, { x: centerX, y: centerY })
                }}
                variant="outline"
                size="sm"
                className="text-xs"
              >
                <Plus className="h-3 w-3 mr-1" />
                {type.replace('_', ' ')}
              </Button>
            ))}
          </div>
        </Panel>
      </ReactFlow>

      <NodeSettings
        node={selectedNode}
        open={settingsOpen}
        onClose={() => {
          setSettingsOpen(false)
          setSelectedNode(null)
        }}
        onSave={handleSaveConfig}
        onDelete={handleDeleteNode}
        flowId={currentFlow?.id}
      />

      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Flow</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete "{currentFlow?.name}"? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleDelete}
            >
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
