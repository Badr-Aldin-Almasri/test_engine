import React, { useState, useEffect } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogDescription } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { NodeData, NodeType, NodeConfig } from '@/types'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Trash2, Play } from 'lucide-react'
import { flowService } from '@/services/flowService'

interface NodeSettingsProps {
  node: NodeData | null
  open: boolean
  onClose: () => void
  onSave: (config: NodeConfig) => void
  onDelete?: () => void
  flowId?: string
}

export const NodeSettings: React.FC<NodeSettingsProps> = ({ node, open, onClose, onSave, onDelete, flowId }) => {
  const [config, setConfig] = useState<NodeConfig>(node?.config || {})
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [isRunning, setIsRunning] = useState(false)
  const [output, setOutput] = useState<any>(node?.output || null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (node) {
      setConfig(node.config || {})
      setOutput(node.output || null)
      setError(node.error || null)
    }
  }, [node])

  const handleSave = () => {
    onSave(config)
    onClose()
  }

  const handleDelete = () => {
    if (onDelete) {
      onDelete()
      setShowDeleteConfirm(false)
      onClose()
    }
  }

  const handleRun = async () => {
    if (!flowId || !node || node.type !== NodeType.API) return

    // Validate required fields
    const urlValue = typeof config.url === 'string' ? config.url.trim() : ''
    const methodValue = typeof config.method === 'string' ? config.method.trim() : (config.method || 'GET')
    
    if (urlValue === '' || methodValue === '') {
      setError('URL and Method are required to run the API call')
      return
    }

    setIsRunning(true)
    setError(null)
    
    try {
      // Send the current config to the backend so it uses the latest configuration
      // even if not saved yet
      const result = await flowService.executeNode(flowId, node!.id, config)
      setOutput(result.output)
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Failed to execute API call')
      setOutput(null)
    } finally {
      setIsRunning(false)
    }
  }

  if (!node) return null

  const renderAPISettings = () => {
    // Check if Run button should be enabled
    const urlValue = typeof config.url === 'string' ? config.url.trim() : ''
    const methodValue = typeof config.method === 'string' ? config.method.trim() : (config.method || 'GET')
    const canRun = urlValue !== '' && methodValue !== ''
    
    return (
    <div className="space-y-4">
      <div className="flex items-center justify-between mb-4">
        <Label htmlFor="method" className="text-base font-semibold">API Configuration</Label>
        {flowId && (
          <Button
            onClick={handleRun}
            disabled={isRunning || !canRun}
            size="sm"
          >
            <Play className="h-4 w-4 mr-2" />
            {isRunning ? 'Running...' : 'Run'}
          </Button>
        )}
      </div>
      <div>
        <Label htmlFor="method">Method</Label>
        <Select
          value={config.method || 'GET'}
          onValueChange={(value) => setConfig({ ...config, method: value })}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="GET">GET</SelectItem>
            <SelectItem value="POST">POST</SelectItem>
            <SelectItem value="PUT">PUT</SelectItem>
            <SelectItem value="DELETE">DELETE</SelectItem>
            <SelectItem value="PATCH">PATCH</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div>
        <Label htmlFor="url">URL</Label>
        <Input
          id="url"
          value={config.url || ''}
          onChange={(e) => setConfig({ ...config, url: e.target.value })}
          placeholder="https://api.example.com/endpoint"
        />
      </div>
      <div>
        <Label htmlFor="headers">Headers (JSON)</Label>
        <Textarea
          id="headers"
          value={config.headers ? JSON.stringify(config.headers, null, 2) : ''}
          onChange={(e) => {
            try {
              const headers = JSON.parse(e.target.value)
              setConfig({ ...config, headers })
            } catch {
              // Invalid JSON, keep as is
            }
          }}
          placeholder='{"Content-Type": "application/json"}'
          className="font-mono text-xs"
        />
      </div>
      <div>
        <Label htmlFor="body">Body</Label>
        <Textarea
          id="body"
          value={config.body || ''}
          onChange={(e) => setConfig({ ...config, body: e.target.value })}
          placeholder='{"key": "value"}'
          className="font-mono text-xs"
        />
      </div>
    </div>
    )
  }

  const renderVerificationSettings = () => (
    <div className="space-y-4">
      <div>
        <Label htmlFor="assertionType">Assertion Type</Label>
        <Select
          value={config.assertionType || 'equals'}
          onValueChange={(value) => setConfig({ ...config, assertionType: value as any })}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="equals">Equals</SelectItem>
            <SelectItem value="contains">Contains</SelectItem>
            <SelectItem value="regex">Regex Match</SelectItem>
            <SelectItem value="custom">Custom Script</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div>
        <Label htmlFor="expected">Expected Value</Label>
        <Textarea
          id="expected"
          value={config.expected ? JSON.stringify(config.expected, null, 2) : ''}
          onChange={(e) => {
            try {
              const expected = JSON.parse(e.target.value)
              setConfig({ ...config, expected })
            } catch {
              // Invalid JSON
            }
          }}
          placeholder='{"status": 200}'
          className="font-mono text-xs"
        />
      </div>
      {config.assertionType === 'custom' && (
        <div>
          <Label htmlFor="customScript">Custom Script</Label>
          <Textarea
            id="customScript"
            value={config.customScript || ''}
            onChange={(e) => setConfig({ ...config, customScript: e.target.value })}
            placeholder="// JavaScript verification code"
            className="font-mono text-xs"
            rows={10}
          />
        </div>
      )}
    </div>
  )

  const renderMockSettings = () => (
    <div className="space-y-4">
      <div>
        <Label htmlFor="mockResponse">Mock Response (JSON)</Label>
        <Textarea
          id="mockResponse"
          value={config.mockResponse ? JSON.stringify(config.mockResponse, null, 2) : ''}
          onChange={(e) => {
            try {
              const mockResponse = JSON.parse(e.target.value)
              setConfig({ ...config, mockResponse })
            } catch {
              // Invalid JSON
            }
          }}
          placeholder='{"status": 200, "data": {...}}'
          className="font-mono text-xs"
          rows={10}
        />
      </div>
      <div>
        <Label htmlFor="mockDelay">Delay (ms)</Label>
        <Input
          id="mockDelay"
          type="number"
          value={config.mockDelay || 0}
          onChange={(e) => setConfig({ ...config, mockDelay: parseInt(e.target.value) || 0 })}
        />
      </div>
    </div>
  )

  const renderReportSettings = () => (
    <div className="space-y-4">
      <div>
        <Label htmlFor="reportName">Report Name</Label>
        <Input
          id="reportName"
          value={config.reportName || ''}
          onChange={(e) => setConfig({ ...config, reportName: e.target.value })}
          placeholder="Test Report Name"
        />
      </div>
    </div>
  )

  const renderEventTriggerSettings = () => (
    <div className="space-y-4">
      <div>
        <Label htmlFor="triggerEvent">Event Name</Label>
        <Input
          id="triggerEvent"
          value={config.triggerEvent || ''}
          onChange={(e) => setConfig({ ...config, triggerEvent: e.target.value })}
          placeholder="user.created"
        />
      </div>
      <div>
        <Label htmlFor="triggerCondition">Condition (JSON)</Label>
        <Textarea
          id="triggerCondition"
          value={config.triggerCondition || ''}
          onChange={(e) => setConfig({ ...config, triggerCondition: e.target.value })}
          placeholder='{"userId": "123"}'
          className="font-mono text-xs"
        />
      </div>
    </div>
  )

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Configure {node.label}</DialogTitle>
        </DialogHeader>
        
        <Tabs defaultValue="config" className="w-full">
          <TabsList>
            <TabsTrigger value="config">Configuration</TabsTrigger>
            <TabsTrigger value="output">Output</TabsTrigger>
          </TabsList>
          
          <TabsContent value="config" className="mt-4">
            {node.type === NodeType.API && renderAPISettings()}
            {node.type === NodeType.VERIFICATION && renderVerificationSettings()}
            {node.type === NodeType.MOCK && renderMockSettings()}
            {node.type === NodeType.REPORT && renderReportSettings()}
            {node.type === NodeType.EVENT_TRIGGER && renderEventTriggerSettings()}
          </TabsContent>
          
          <TabsContent value="output" className="mt-4">
            {error && (
              <div className="bg-red-50 border border-red-200 p-4 rounded mb-4">
                <p className="text-sm text-red-800 font-semibold mb-2">Error:</p>
                <p className="text-xs text-red-700">{error}</p>
              </div>
            )}
            {output && (
              <div className="bg-gray-50 p-4 rounded border">
                <pre className="text-xs overflow-auto max-h-96">
                  {JSON.stringify(output, null, 2)}
                </pre>
              </div>
            )}
            {!output && !error && (
              <div className="text-sm text-muted-foreground">
                {node.type === NodeType.API
                  ? 'Click "Run" to execute this API call and see the output here.'
                  : 'No output yet. Run the flow to see results.'}
              </div>
            )}
          </TabsContent>
        </Tabs>

        <DialogFooter className="flex justify-between">
          <div>
            {onDelete && (
              <Button
                variant="destructive"
                onClick={() => setShowDeleteConfirm(true)}
                className="mr-2"
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Delete Node
              </Button>
            )}
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={onClose}>Cancel</Button>
            <Button onClick={handleSave}>Save</Button>
          </div>
        </DialogFooter>
      </DialogContent>

      {showDeleteConfirm && (
        <Dialog open={showDeleteConfirm} onOpenChange={setShowDeleteConfirm}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Delete Node</DialogTitle>
              <DialogDescription>
                Are you sure you want to delete "{node.label}"? This action cannot be undone and will also remove all connections to this node.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button variant="outline" onClick={() => setShowDeleteConfirm(false)}>
                Cancel
              </Button>
              <Button variant="destructive" onClick={handleDelete}>
                Delete
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      )}
    </Dialog>
  )
}

