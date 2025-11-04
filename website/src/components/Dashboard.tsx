import React, { useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card'
import { Button } from './ui/button'
import { Badge } from './ui/badge'
import { useFlowStore } from '@/stores/flowStore'
import { ExecutionStatus } from '@/types'
import { Play, Plus, Calendar, Clock, LogOut } from 'lucide-react'
import { formatDistanceToNow } from 'date-fns'
import { useAuth } from '@/contexts/AuthContext'
import { flowService } from '@/services/flowService'
import { useNavigate } from 'react-router-dom'

export const Dashboard: React.FC = () => {
  const { flows, isLoading, currentFlow, setCurrentFlow, loadFlows, createFlow } = useFlowStore()
  const { user, logout } = useAuth()
  const [testRuns, setTestRuns] = React.useState<any[]>([])
  const navigate = useNavigate()

  // Ensure flows is always an array
  const safeFlows = flows || []
  // Ensure testRuns is always an array
  const safeTestRuns = testRuns || []

  useEffect(() => {
    loadFlows()
  }, [loadFlows])

  const createNewFlow = async () => {
    try {
      const newFlow = await createFlow({
        name: `New Flow ${safeFlows.length + 1}`,
        description: 'A new test flow',
        tags: [],
        nodes: [],
        edges: [],
      })
      setCurrentFlow(newFlow)
      navigate(`/flow/${newFlow.id}/editor`)
    } catch (error) {
      console.error('Failed to create flow:', error)
    }
  }

  const loadTestRuns = async (flowId: string) => {
    try {
      const runs = await flowService.getTestRuns(flowId)
      setTestRuns(runs || [])
    } catch (error) {
      console.error('Failed to load test runs:', error)
      setTestRuns([])
    }
  }

  const getStatusBadge = (status: ExecutionStatus) => {
    switch (status) {
      case ExecutionStatus.SUCCESS:
        return <Badge variant="success">Passed</Badge>
      case ExecutionStatus.FAILED:
        return <Badge variant="destructive">Failed</Badge>
      case ExecutionStatus.RUNNING:
        return <Badge variant="secondary">Running</Badge>
      default:
        return <Badge variant="outline">Pending</Badge>
    }
  }

  return (
    <div className="container mx-auto p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Test Flow Dashboard</h1>
          <p className="text-muted-foreground mt-1">
            Manage your API test flows and view execution history
          </p>
          {user && (
            <p className="text-sm text-muted-foreground mt-1">
              Welcome, {user.name} ({user.email})
            </p>
          )}
        </div>
        <div className="flex gap-2">
          <Button onClick={createNewFlow}>
            <Plus className="h-4 w-4 mr-2" />
            New Flow
          </Button>
          <Button onClick={logout} variant="outline">
            <LogOut className="h-4 w-4 mr-2" />
            Logout
          </Button>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Flows</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{safeFlows.length}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Runs</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{safeTestRuns.length}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {safeTestRuns.length > 0
                ? Math.round(
                    (safeTestRuns.filter((r) => r && r.status === ExecutionStatus.SUCCESS).length /
                      safeTestRuns.length) *
                      100
                  )
                : 0}
              %
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Flows</CardTitle>
            <CardDescription>Your test flow definitions</CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8 text-muted-foreground">
                <p>Loading flows...</p>
              </div>
            ) : safeFlows.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <p>No flows yet. Create your first flow to get started.</p>
                <Button onClick={createNewFlow} className="mt-4" variant="outline">
                  <Plus className="h-4 w-4 mr-2" />
                  Create Flow
                </Button>
              </div>
            ) : (
              <div className="space-y-2">
                {safeFlows.map((flow) => (
                  <Card
                    key={flow.id}
                    className="cursor-pointer hover:bg-accent"
                    onClick={() => {
                      setCurrentFlow(flow)
                      loadTestRuns(flow.id)
                      navigate(`/flow/${flow.id}/editor`)
                    }}
                  >
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold">{flow.name}</h3>
                          <p className="text-sm text-muted-foreground">
                            {flow.nodes?.length || 0} nodes • {flow.edges?.length || 0} connections
                          </p>
                        </div>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={(e) => {
                            e.stopPropagation()
                            setCurrentFlow(flow)
                            navigate(`/flow/${flow.id}/editor`)
                          }}
                        >
                          <Play className="h-4 w-4" />
                        </Button>
                      </div>
                      {flow.tags && flow.tags.length > 0 && (
                        <div className="flex gap-1 mt-2">
                          {flow.tags.map((tag) => (
                            <Badge key={tag} variant="outline" className="text-xs">
                              {tag}
                            </Badge>
                          ))}
                        </div>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Recent Test Runs</CardTitle>
            <CardDescription>Execution history and results</CardDescription>
          </CardHeader>
          <CardContent>
            {safeTestRuns.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                <p>No test runs yet. Run a flow to see results here.</p>
              </div>
            ) : (
              <div className="space-y-2">
                {safeTestRuns.slice(0, 10).filter(run => run !== null).map((run) => (
                  <Card key={run.id}>
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between mb-2">
                        <h3 className="font-semibold">{run.flow_name || 'Test Run'}</h3>
                        {getStatusBadge(run.status)}
                      </div>
                      <div className="flex items-center gap-4 text-sm text-muted-foreground">
                        <div className="flex items-center gap-1">
                          <Calendar className="h-3 w-3" />
                          {formatDistanceToNow(new Date(run.started_at), {
                            addSuffix: true,
                          })}
                        </div>
                        {run.duration_ms && (
                          <div className="flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            {(run.duration_ms / 1000).toFixed(2)}s
                          </div>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
