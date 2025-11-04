import { Flow, FlowNode, FlowEdge, NodeType, ExecutionStatus, NodeData } from '@/types'

interface ExecutionContext {
  nodeOutputs: Map<string, any>
  nodeStatuses: Map<string, ExecutionStatus>
}

export class MockExecutionEngine {
  private context: ExecutionContext

  constructor() {
    this.context = {
      nodeOutputs: new Map(),
      nodeStatuses: new Map(),
    }
  }

  async executeFlow(
    flow: Flow,
    onNodeUpdate: (nodeId: string, status: ExecutionStatus, output?: any, error?: string) => void
  ): Promise<void> {
    this.context.nodeOutputs.clear()
    this.context.nodeStatuses.clear()

    // Build dependency graph
    const nodeMap = new Map(flow.nodes.map(n => [n.id, n]))
    const incomingEdges = new Map<string, FlowEdge[]>()
    
    flow.edges.forEach(edge => {
      if (!incomingEdges.has(edge.target)) {
        incomingEdges.set(edge.target, [])
      }
      incomingEdges.get(edge.target)!.push(edge)
    })

    // Find root nodes (nodes with no incoming edges)
    const rootNodes = flow.nodes.filter(node => !incomingEdges.has(node.id))

    // Execute nodes in topological order
    const executed = new Set<string>()
    const queue = [...rootNodes.map(n => n.id)]

    while (queue.length > 0) {
      const nodeId = queue.shift()!
      if (executed.has(nodeId)) continue

      const node = nodeMap.get(nodeId)!
      const dependencies = incomingEdges.get(nodeId) || []
      
      // Check if all dependencies are executed
      const allDepsExecuted = dependencies.every(edge => executed.has(edge.source))
      if (!allDepsExecuted) {
        // Re-add to queue later
        queue.push(nodeId)
        continue
      }

      // Execute node
      await this.executeNode(node, onNodeUpdate)
      executed.add(nodeId)

      // Add dependent nodes to queue
      flow.edges
        .filter(edge => edge.source === nodeId)
        .forEach(edge => {
          if (!queue.includes(edge.target) && !executed.has(edge.target)) {
            queue.push(edge.target)
          }
        })
    }
  }

  private async executeNode(
    node: FlowNode,
    onUpdate: (nodeId: string, status: ExecutionStatus, output?: any, error?: string) => void
  ): Promise<void> {
    onUpdate(node.id, ExecutionStatus.RUNNING)

    try {
      await new Promise(resolve => setTimeout(resolve, 500 + Math.random() * 1000))

      let output: any

      switch (node.data.type) {
        case NodeType.API:
          output = await this.executeAPINode(node.data)
          break
        case NodeType.MOCK:
          output = await this.executeMockNode(node.data)
          break
        case NodeType.VERIFICATION:
          output = await this.executeVerificationNode(node.data)
          break
        case NodeType.REPORT:
          output = await this.executeReportNode(node.data)
          break
        case NodeType.EVENT_TRIGGER:
          output = await this.executeEventTriggerNode(node.data)
          break
        default:
          throw new Error(`Unknown node type: ${node.data.type}`)
      }

      this.context.nodeOutputs.set(node.id, output)
      this.context.nodeStatuses.set(node.id, ExecutionStatus.SUCCESS)
      onUpdate(node.id, ExecutionStatus.SUCCESS, output)
    } catch (error: any) {
      this.context.nodeStatuses.set(node.id, ExecutionStatus.FAILED)
      onUpdate(node.id, ExecutionStatus.FAILED, undefined, error.message)
    }
  }

  private async executeAPINode(data: NodeData): Promise<any> {
    const { method = 'GET', url = '', headers = {}, body } = data.config

    // Mock API response
    const mockResponse = {
      status: 200,
      statusText: 'OK',
      headers: { 'content-type': 'application/json' },
      data: {
        id: Math.floor(Math.random() * 1000),
        message: `Successfully called ${method} ${url}`,
        timestamp: new Date().toISOString(),
        ...(body ? { requestBody: JSON.parse(body) } : {}),
      }
    }

    return mockResponse
  }

  private async executeMockNode(data: NodeData): Promise<any> {
    const { mockResponse, mockDelay = 0 } = data.config

    if (mockDelay > 0) {
      await new Promise(resolve => setTimeout(resolve, mockDelay))
    }

    return mockResponse || { status: 200, data: { mock: true } }
  }

  private async executeVerificationNode(data: NodeData): Promise<any> {
    const { expected, assertionType = 'equals' } = data.config
    
    // In real implementation, this would compare with previous node outputs
    // For mock, we'll simulate verification
    const previousOutput = Array.from(this.context.nodeOutputs.values())[0] || {}

    let passed = false
    switch (assertionType) {
      case 'equals':
        passed = JSON.stringify(previousOutput) === JSON.stringify(expected)
        break
      case 'contains':
        passed = JSON.stringify(previousOutput).includes(JSON.stringify(expected))
        break
      case 'regex':
        // Mock regex match
        passed = true
        break
      case 'custom':
        // Mock custom script execution
        passed = true
        break
    }

    if (!passed) {
      throw new Error(`Verification failed: Expected ${JSON.stringify(expected)}, got ${JSON.stringify(previousOutput)}`)
    }

    return {
      passed: true,
      expected,
      actual: previousOutput,
    }
  }

  private async executeReportNode(data: NodeData): Promise<any> {
    const { reportName } = data.config

    const report = {
      name: reportName || 'Test Report',
      timestamp: new Date().toISOString(),
      nodes: Array.from(this.context.nodeStatuses.entries()).map(([id, status]) => ({
        nodeId: id,
        status,
      })),
    }

    return report
  }

  private async executeEventTriggerNode(data: NodeData): Promise<any> {
    const { triggerEvent, triggerCondition } = data.config

    return {
      event: triggerEvent,
      condition: triggerCondition,
      triggered: true,
      timestamp: new Date().toISOString(),
    }
  }
}

