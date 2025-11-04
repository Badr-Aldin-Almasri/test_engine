import React from 'react'
import { Handle, Position } from 'reactflow'
import { NodeData, ExecutionStatus } from '@/types'
import { Badge } from '@/components/ui/badge'
import { 
  Cloud, 
  CheckCircle2, 
  Loader2, 
  XCircle, 
  FileJson,
  Zap,
  AlertCircle
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface CustomNodeProps {
  data: NodeData
}

const getNodeIcon = (type: string) => {
  switch (type) {
    case 'api':
      return <Cloud className="h-4 w-4" />
    case 'verification':
      return <CheckCircle2 className="h-4 w-4" />
    case 'mock':
      return <FileJson className="h-4 w-4" />
    case 'report':
      return <AlertCircle className="h-4 w-4" />
    case 'event_trigger':
      return <Zap className="h-4 w-4" />
    default:
      return <Cloud className="h-4 w-4" />
  }
}

const getStatusIcon = (status?: ExecutionStatus) => {
  switch (status) {
    case ExecutionStatus.RUNNING:
      return <Loader2 className="h-3 w-3 animate-spin text-blue-500" />
    case ExecutionStatus.SUCCESS:
      return <CheckCircle2 className="h-3 w-3 text-green-500" />
    case ExecutionStatus.FAILED:
      return <XCircle className="h-3 w-3 text-red-500" />
    default:
      return null
  }
}

const getNodeColor = (type: string) => {
  switch (type) {
    case 'api':
      return 'border-blue-500 bg-blue-50'
    case 'verification':
      return 'border-green-500 bg-green-50'
    case 'mock':
      return 'border-purple-500 bg-purple-50'
    case 'report':
      return 'border-orange-500 bg-orange-50'
    case 'event_trigger':
      return 'border-yellow-500 bg-yellow-50'
    default:
      return 'border-gray-500 bg-gray-50'
  }
}

export const CustomNode: React.FC<CustomNodeProps> = ({ data }) => {
  return (
    <div className={cn(
      "px-4 py-2 shadow-md rounded-lg border-2 bg-white min-w-[150px]",
      getNodeColor(data.type)
    )}>
      <Handle type="target" position={Position.Top} className="w-3 h-3" />
      
      <div className="flex items-center gap-2">
        {getNodeIcon(data.type)}
        <div className="flex-1">
          <div className="text-sm font-semibold">{data.label}</div>
          <div className="text-xs text-gray-500 capitalize">{data.type.replace('_', ' ')}</div>
        </div>
        {getStatusIcon(data.status)}
      </div>
      
      {data.error && (
        <div className="mt-2 text-xs text-red-600 bg-red-50 p-1 rounded">
          {data.error}
        </div>
      )}
      
      <Handle type="source" position={Position.Bottom} className="w-3 h-3" />
    </div>
  )
}

export const nodeTypes = {
  custom: CustomNode,
}

