import { apiClient } from './api-client';
import { 
  Workflow, 
  WorkflowNode, 
  Connection, 
  WorkflowBackend,
  ApiResponse,
  convertBackendToFrontend,
  convertFrontendToBackendRequest
} from '@/types/workflow';

// Workflow API Service
export class WorkflowService {
  
  // Get all workflows for a team
  static async getWorkflows(teamId: number): Promise<Workflow[]> { 
    const response = await apiClient.get<ApiResponse<WorkflowBackend[]>>(`/teams/${teamId}/workflow`);
    const backendWorkflows = response.data || [];
    
    // Convert backend workflows to frontend format
    return backendWorkflows.map(convertBackendToFrontend);
  }

  // Get a specific workflow by team ID and workflow ID
  static async getWorkflow(workflowId: string, teamId: number): Promise<Workflow> { 
    const response = await apiClient.get<ApiResponse<WorkflowBackend>>(`/teams/2324354/workflow/${workflowId}`); // todo: Replace with actual team ID
    const backendWorkflow = response.data;
    
    // Convert backend workflow to frontend format
    return convertBackendToFrontend(backendWorkflow);
  }

  // Create a new workflow
  static async createWorkflow(workflow: Partial<Workflow>, teamId: number = 1): Promise<Workflow> {
    const requestData = convertFrontendToBackendRequest(workflow);
    
    const response = await apiClient.post<ApiResponse<WorkflowBackend>>(`/teams/${teamId}/workflow`, requestData);
    const backendWorkflow = response.data;
    
    return convertBackendToFrontend(backendWorkflow);
  }

  // Update an existing workflow
  static async updateWorkflow(
    workflowId: string, 
    workflow: Partial<Workflow>,
    teamId: number = 1
  ): Promise<Workflow> {
    const requestData = convertFrontendToBackendRequest(workflow);
    
    const response = await apiClient.put<ApiResponse<WorkflowBackend>>(
      `/teams/${teamId}/workflow/${workflowId}`, 
      requestData
    );
    const backendWorkflow = response.data;
    
    return convertBackendToFrontend(backendWorkflow);
  }

  // Delete a workflow
  static async deleteWorkflow(workflowId: string, teamId: number = 1): Promise<void> {
    await apiClient.delete<ApiResponse<void>>(`/teams/${teamId}/workflow/${workflowId}`);
  }

  // Execute a workflow
  static async executeWorkflow(
    workflowId: string, 
    input?: Record<string, unknown>,
    teamId: number = 1
  ): Promise<{ executionId: string }> {
    const response = await apiClient.post<ApiResponse<{ executionId: string }>>(
      `/teams/${teamId}/workflow/${workflowId}/execute`, 
      { input }
    );
    return response.data;
  }

  // Get workflow execution status
  static async getExecutionStatus(
    workflowId: string, 
    executionId: string,
    teamId: number = 1
  ): Promise<{
    status: 'running' | 'completed' | 'failed' | 'cancelled';
    progress: number;
    result?: unknown;
    error?: string;
  }> {
    const response = await apiClient.get<ApiResponse<{
      status: 'running' | 'completed' | 'failed' | 'cancelled';
      progress: number;
      result?: unknown;
      error?: string;
    }>>(`/teams/${teamId}/workflow/${workflowId}/executions/${executionId}`);
    return response.data;
  }

  // Save workflow (nodes and connections)
  static async saveWorkflow(
    workflowId: string, 
    nodes: WorkflowNode[], 
    connections: Connection[],
    teamId: number = 1
  ): Promise<Workflow> {
    const requestData = convertFrontendToBackendRequest(
      { name: 'Updated Workflow' }, // Will be overridden by existing name
      nodes,
      connections
    );
    
    const response = await apiClient.put<ApiResponse<WorkflowBackend>>(
      `/teams/${teamId}/workflow/${workflowId}`,
      requestData
    );
    const backendWorkflow = response.data;
    
    return convertBackendToFrontend(backendWorkflow);
  }

  // Activate/Deactivate workflow (if backend supports this)
  static async toggleWorkflowStatus(
    workflowId: string, 
    isActive: boolean,
    teamId: number = 1
  ): Promise<Workflow> {
    const response = await apiClient.put<ApiResponse<WorkflowBackend>>(
      `/teams/${teamId}/workflow/${workflowId}/status`, 
      { isActive }
    );
    const backendWorkflow = response.data;
    
    return convertBackendToFrontend(backendWorkflow);
  }

  // Test webhook endpoint (if backend supports this)
  static async testWebhook(webhookUrl: string, payload: Record<string, unknown>): Promise<{
    success: boolean;
    response?: unknown;
    error?: string;
  }> {
    const response = await apiClient.post<ApiResponse<{
      success: boolean;
      response?: unknown;
      error?: string;
    }>>('/webhook/test', {
      url: webhookUrl,
      payload
    });
    return response.data;
  }
}
