import { apiClient } from './api-client';
import { 
  Workflow, 
  WorkflowNode, 
  Connection, 
  WorkflowBackend,
  convertBackendToFrontend,
  convertFrontendToBackendRequest
} from '@/types/workflow';

// Workflow API Service
export class WorkflowService {
  
  // Get all workflows for a team
  static async getWorkflows(): Promise<Workflow[]> { 
    try {
      const response = await apiClient.get<{ workflows: WorkflowBackend[] }>(`/api/v1/workflows`);
      
      // Extract workflows array from response object
      const workflowsArray = response?.workflows || [];
      
      // Convert backend workflows to frontend format, filtering out any invalid ones
      return workflowsArray
        .filter(workflow => workflow != null) // Filter out null/undefined workflows
        .map(convertBackendToFrontend);
    } catch (error) {
      console.error('Failed to fetch workflows:', error);
      throw new Error(error instanceof Error ? error.message : 'Failed to fetch workflows');
    }
  }

  // Get a specific workflow by workflow ID
  static async getWorkflow(workflowId: string): Promise<Workflow> { 
    try {
      const backendWorkflow = await apiClient.get<WorkflowBackend>(`/api/v1/workflows/${workflowId}`);
      
      if (!backendWorkflow) {
        throw new Error('Workflow not found or response is empty');
      }
      
      // Convert backend workflow to frontend format
      return convertBackendToFrontend(backendWorkflow);
    } catch (error) {
      console.error(`Failed to fetch workflow ${workflowId}:`, error);
      throw new Error(error instanceof Error ? error.message : 'Failed to fetch workflow');
    }
  }

  // Create a new workflow
  static async createWorkflow(workflow: Partial<Workflow>): Promise<Workflow> {
    try {
      const requestData = convertFrontendToBackendRequest(workflow);
      
      const backendWorkflow = await apiClient.post<WorkflowBackend>(`/api/v1/workflows`, requestData);
      
      if (!backendWorkflow) {
        throw new Error('Create workflow response is empty');
      }
      
      return convertBackendToFrontend(backendWorkflow);
    } catch (error) {
      console.error('Failed to create workflow:', error);
      throw new Error(error instanceof Error ? error.message : 'Failed to create workflow');
    }
  }

  // Update an existing workflow
  static async updateWorkflow(
    workflowId: string, 
    workflow: Partial<Workflow>
  ): Promise<Workflow> {
    try {
      const requestData = convertFrontendToBackendRequest(workflow);
      
      const backendWorkflow = await apiClient.put<WorkflowBackend>(
        `/api/v1/workflows/${workflowId}`, 
        requestData
      );
      
      if (!backendWorkflow) {
        throw new Error('Update workflow response is empty');
      }
      
      return convertBackendToFrontend(backendWorkflow);
    } catch (error) {
      console.error(`Failed to update workflow ${workflowId}:`, error);
      throw new Error(error instanceof Error ? error.message : 'Failed to update workflow');
    }
  }

  // Delete a workflow
  static async deleteWorkflow(workflowId: string): Promise<void> {
    await apiClient.delete<void>(`/api/v1/workflows/${workflowId}`);
  }

  // Execute a workflow
  static async executeWorkflow(
    workflowId: string, 
    input?: Record<string, unknown>
  ): Promise<{ executionId: string }> {
    return await apiClient.post<{ executionId: string }>(
      `/api/v1/workflows/${workflowId}/execute`, 
      { input }
    );
  }

  // Get workflow execution status
  static async getExecutionStatus(
    workflowId: string, 
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
    }>(`/api/v1/workflows/${workflowId}/executions/${executionId}`);
  }

  // Save workflow (nodes and connections)
  static async saveWorkflow(
    workflowId: string, 
    nodes: WorkflowNode[], 
    connections: Connection[]
  ): Promise<Workflow> {
    try {
      const requestData = convertFrontendToBackendRequest(
        { name: 'Updated Workflow' }, // Will be overridden by existing name
        nodes,
        connections
      );
      
      const backendWorkflow = await apiClient.put<WorkflowBackend>(
        `/api/v1/workflows/${workflowId}`,
        requestData
      );
      
      if (!backendWorkflow) {
        throw new Error('Save workflow response is empty');
      }
      
      return convertBackendToFrontend(backendWorkflow);
    } catch (error) {
      console.error(`Failed to save workflow ${workflowId}:`, error);
      throw new Error(error instanceof Error ? error.message : 'Failed to save workflow');
    }
  }

  // Activate/Deactivate workflow (if backend supports this)
  static async toggleWorkflowStatus(
    workflowId: string, 
    isActive: boolean
  ): Promise<Workflow> {
    const backendWorkflow = await apiClient.put<WorkflowBackend>(
      `/api/v1/workflows/${workflowId}/status`, 
      { isActive }
    );
    
    return convertBackendToFrontend(backendWorkflow);
  }

  // Test webhook endpoint (if backend supports this)
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
}
