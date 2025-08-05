import { apiClient } from './api-client';
import { 
  Flow, 
  FlowNode, 
  Connection, 
  FlowBackend,
  convertBackendToFrontend,
  convertFrontendToBackendRequest
} from '@/types/flow';

// Flow API Service (renamed from WorkflowService for consistency)
export class FlowService {
  
  // Get all flows for a team
  static async getFlows(): Promise<Flow[]> { 
    try {
      const response = await apiClient.get<{ flows: FlowBackend[] }>(`/api/v1/flows`);
      
      // Extract flows array from response object
      const flowsArray = response?.flows || [];
      
      // Convert backend flows to frontend format, filtering out any invalid ones
      return flowsArray
        .filter(flow => flow != null) // Filter out null/undefined flows
        .map(convertBackendToFrontend);
    } catch (error) {
      console.error('Failed to fetch flows:', error);
      throw new Error(error instanceof Error ? error.message : 'Failed to fetch flows');
    }
  }

  // Get a specific flow by flow ID
  static async getFlow(flowId: string): Promise<Flow> { 
    try {
      const backendFlow = await apiClient.get<FlowBackend>(`/api/v1/flows/${flowId}`);
      
      if (!backendFlow) {
        throw new Error('Flow not found or response is empty');
      }
      
      // Convert backend flow to frontend format
      return convertBackendToFrontend(backendFlow);
    } catch (error) {
      console.error(`Failed to fetch flow ${flowId}:`, error);
      throw new Error(error instanceof Error ? error.message : 'Failed to fetch flow');
    }
  }

  // Create a new flow
  static async createFlow(
    flow: Partial<Flow>, 
    nodes: FlowNode[] = [], 
    connections: Connection[] = []
  ): Promise<Flow> {
    try {
      const requestData = convertFrontendToBackendRequest(flow, nodes, connections);
      
      const backendFlow = await apiClient.post<FlowBackend>(`/api/v1/flows`, requestData);
      
      if (!backendFlow) {
        throw new Error('Create flow response is empty');
      }
      
      return convertBackendToFrontend(backendFlow);
    } catch (error) {
      console.error('Failed to create flow:', error);
      throw new Error(error instanceof Error ? error.message : 'Failed to create flow');
    }
  }

  // Update an existing flow
  static async updateFlow(
    flowId: string, 
    flow: Partial<Flow>,
    nodes: FlowNode[] = [],
    connections: Connection[] = []
  ): Promise<Flow> {
    try {
      const requestData = convertFrontendToBackendRequest(flow, nodes, connections);
      
      const backendFlow = await apiClient.put<FlowBackend>(
        `/api/v1/flows/${flowId}`, 
        requestData
      );
      
      if (!backendFlow) {
        throw new Error('Update flow response is empty');
      }
      
      return convertBackendToFrontend(backendFlow);
    } catch (error) {
      console.error(`Failed to update flow ${flowId}:`, error);
      throw new Error(error instanceof Error ? error.message : 'Failed to update flow');
    }
  }

  // Delete a flow
  static async deleteFlow(flowId: string): Promise<void> {
    await apiClient.delete<void>(`/api/v1/flows/${flowId}`);
  }

  // Execute a flow
  static async executeFlow(
    flowId: string, 
    input?: Record<string, unknown>
  ): Promise<{ executionId: string }> {
    return await apiClient.post<{ executionId: string }>(
      `/api/v1/flows/${flowId}/execute`, 
      { input }
    );
  }

  // Get flow execution status
  static async getExecutionStatus(
    flowId: string, 
    executionId: string
  ): Promise<{
    status: 'running' | 'completed' | 'failed' | 'cancelled';
    progress: number;
    result?: unknown;
    error?: string;
  }> {
    return await apiClient.get<{
      status: 'running' | 'completed' | 'failed' | 'cancelled';
      progress: number;
      result?: unknown;
      error?: string;
    }>(`/api/v1/flows/${flowId}/executions/${executionId}`);
  }

  // Save flow (nodes and connections) - this is an alias for updateFlow
  static async saveFlow(
    flowId: string, 
    nodes: FlowNode[], 
    connections: Connection[]
  ): Promise<Flow> {
    try {
      // Use updateFlow with empty flow object since we're just updating nodes/connections
      return await this.updateFlow(flowId, {}, nodes, connections);
    } catch (error) {
      console.error(`Failed to save flow ${flowId}:`, error);
      throw new Error(error instanceof Error ? error.message : 'Failed to save flow');
    }
  }

  // Test webhook endpoint 
  static async testWebhook(webhookUrl: string, payload: Record<string, unknown>): Promise<{
    success: boolean;
    response?: unknown;
    error?: string;
  }> {
    return await apiClient.post<{
      success: boolean;
      response?: unknown;
      error?: string;
    }>('/webhook/test', {
      url: webhookUrl,
      payload
    });
  }

  // Legacy method aliases for backward compatibility
  static async getWorkflows(): Promise<Flow[]> {
    return this.getFlows();
  }

  static async getWorkflow(workflowId: string): Promise<Flow> {
    return this.getFlow(workflowId);
  }

  static async createWorkflow(
    workflow: Partial<Flow>, 
    nodes: FlowNode[] = [], 
    connections: Connection[] = []
  ): Promise<Flow> {
    return this.createFlow(workflow, nodes, connections);
  }

  static async updateWorkflow(
    workflowId: string, 
    workflow: Partial<Flow>,
    nodes: FlowNode[] = [],
    connections: Connection[] = []
  ): Promise<Flow> {
    return this.updateFlow(workflowId, workflow, nodes, connections);
  }

  static async deleteWorkflow(workflowId: string): Promise<void> {
    return this.deleteFlow(workflowId);
  }

  static async executeWorkflow(
    workflowId: string, 
    input?: Record<string, unknown>
  ): Promise<{ executionId: string }> {
    return this.executeFlow(workflowId, input);
  }

  static async saveWorkflow(
    workflowId: string, 
    nodes: FlowNode[], 
    connections: Connection[]
  ): Promise<Flow> {
    return this.saveFlow(workflowId, nodes, connections);
  }
}

