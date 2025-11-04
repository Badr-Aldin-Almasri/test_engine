import api from './api'
import { Flow, TestRun } from '@/types'

export const flowService = {
  async getFlows(): Promise<Flow[]> {
    const response = await api.get<Flow[]>('/flows')
    return response.data
  },

  async getFlow(id: string): Promise<Flow> {
    const response = await api.get<Flow>(`/flows/${id}`)
    return response.data
  },

  async createFlow(flow: Partial<Flow>): Promise<Flow> {
    const response = await api.post<Flow>('/flows', flow)
    return response.data
  },

  async updateFlow(id: string, flow: Partial<Flow>): Promise<Flow> {
    const response = await api.put<Flow>(`/flows/${id}`, flow)
    return response.data
  },

  async deleteFlow(id: string): Promise<void> {
    await api.delete(`/flows/${id}`)
  },

  async runFlow(id: string): Promise<{ message: string; flow_id: string }> {
    const response = await api.post(`/flows/${id}/run`)
    return response.data
  },

  async getTestRuns(flowId: string): Promise<TestRun[]> {
    const response = await api.get<TestRun[]>(`/flows/${flowId}/test-runs`)
    return response.data
  },

  async getTestRun(id: string): Promise<TestRun> {
    const response = await api.get<TestRun>(`/test-runs/${id}`)
    return response.data
  },

  async executeNode(flowId: string, nodeId: string, config?: any): Promise<any> {
    const response = await api.post(`/nodes/${flowId}/${nodeId}/execute`, config ? { config } : {})
    return response.data
  },
}

