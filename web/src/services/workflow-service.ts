import { apiClient } from './api-client';
import { Workflow, WorkflowNode, Connection } from '@/types/workflow';

// API Response types
interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

interface ApiError {
  success: false;
  error: string;
  message?: string;
}

// Workflow API Service
export class WorkflowService {
  
  // Get all workflows
  static async getWorkflows(): Promise<Workflow[]> {
    const response = await apiClient.get<ApiResponse<Workflow[]>>('/workflows');
    return response.data;
  }

  // Get a specific workflow by ID
  static async getWorkflow(id: string): Promise<Workflow> {
    const response = await apiClient.get<ApiResponse<Workflow>>(`/workflows/${id}`);
    return response.data;
  }

  // Create a new workflow
  static async createWorkflow(workflow: Partial<Workflow>): Promise<Workflow> {
    const response = await apiClient.post<ApiResponse<Workflow>>('/workflows', workflow);
    return response.data;
  }

  // Update an existing workflow
  static async updateWorkflow(id: string, workflow: Partial<Workflow>): Promise<Workflow> {
    const response = await apiClient.put<ApiResponse<Workflow>>(`/workflows/${id}`, workflow);
    return response.data;
  }

  // Delete a workflow
  static async deleteWorkflow(id: string): Promise<void> {
    await apiClient.delete<ApiResponse<void>>(`/workflows/${id}`);
  }

  // Execute a workflow
  static async executeWorkflow(id: string, input?: Record<string, unknown>): Promise<{ executionId: string }> {
    const response = await apiClient.post<ApiResponse<{ executionId: string }>>(`/workflows/${id}/execute`, { input });
    return response.data;
  }

  // Get workflow execution status
  static async getExecutionStatus(workflowId: string, executionId: string): Promise<{
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
    }>>(`/workflows/${workflowId}/executions/${executionId}`);
    return response.data;
  }

  // Save workflow (nodes and connections)
  static async saveWorkflow(
    id: string, 
    nodes: WorkflowNode[], 
    connections: Connection[]
  ): Promise<Workflow> {
    const response = await apiClient.put<ApiResponse<Workflow>>(`/workflows/${id}`, {
      nodes,
      connections,
      updatedAt: new Date().toISOString()
    });
    return response.data;
  }

  // Activate/Deactivate workflow
  static async toggleWorkflowStatus(id: string, isActive: boolean): Promise<Workflow> {
    const response = await apiClient.put<ApiResponse<Workflow>>(`/workflows/${id}/status`, {
      isActive
    });
    return response.data;
  }

  // Test webhook endpoint
  static async testWebhook(webhookUrl: string, payload: Record<string, unknown>): Promise<{
    success: boolean;
    response?: unknown;
    error?: string;
  }> {
    const response = await apiClient.post<ApiResponse<{
      success: boolean;
      response?: unknown;
      error?: string;
    }>>('/webhooks/test', {
      url: webhookUrl,
      payload
    });
    return response.data;
  }

  // Get workflow templates
  static async getTemplates(): Promise<Array<{
    id: string;
    name: string;
    description: string;
    category: string;
    nodes: WorkflowNode[];
    connections: Connection[];
  }>> {
    const response = await apiClient.get<ApiResponse<Array<{
      id: string;
      name: string;
      description: string;
      category: string;
      nodes: WorkflowNode[];
      connections: Connection[];
    }>>>('/workflows/templates');
    return response.data;
  }

  // Create workflow from template
  static async createFromTemplate(templateId: string, name: string): Promise<Workflow> {
    const response = await apiClient.post<ApiResponse<Workflow>>('/workflows/from-template', {
      templateId,
      name
    });
    return response.data;
  }
}
