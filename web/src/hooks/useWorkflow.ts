import { useState, useEffect } from 'react';
import { Workflow, WorkflowNode, Connection } from '@/types/workflow';
import { WorkflowService } from '@/services/workflow-service';

// Hook for managing workflow state
export function useWorkflow(workflowId?: string) {
  const [workflow, setWorkflow] = useState<Workflow | null>(null);
  const [nodes, setNodes] = useState<WorkflowNode[]>([]);
  const [connections, setConnections] = useState<Connection[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load workflow data
  const loadWorkflow = async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      const workflowData = await WorkflowService.getWorkflow(id);
      setWorkflow(workflowData);
      setNodes(workflowData.nodes);
      setConnections(workflowData.connections);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load workflow');
    } finally {
      setLoading(false);
    }
  };

  // Save workflow changes
  const saveWorkflow = async () => {
    if (!workflow) return;
    
    setLoading(true);
    setError(null);
    try {
      const updatedWorkflow = await WorkflowService.saveWorkflow(
        workflow.id,
        nodes,
        connections
      );
      setWorkflow(updatedWorkflow);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save workflow');
    } finally {
      setLoading(false);
    }
  };

  // Create new workflow
  const createWorkflow = async (name: string, description?: string) => {
    setLoading(true);
    setError(null);
    try {
      const newWorkflow = await WorkflowService.createWorkflow({
        name,
        description,
        nodes: [],
        connections: []
      });
      setWorkflow(newWorkflow);
      setNodes([]);
      setConnections([]);
      return newWorkflow;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create workflow');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // Execute workflow
  const executeWorkflow = async (input?: Record<string, unknown>) => {
    if (!workflow) return;
    
    setLoading(true);
    setError(null);
    try {
      const result = await WorkflowService.executeWorkflow(workflow.id, input);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to execute workflow');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // Toggle workflow active status
  const toggleActive = async () => {
    if (!workflow) return;
    
    setLoading(true);
    setError(null);
    try {
      const updatedWorkflow = await WorkflowService.toggleWorkflowStatus(
        workflow.id,
        !workflow.isActive
      );
      setWorkflow(updatedWorkflow);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update workflow status');
    } finally {
      setLoading(false);
    }
  };

  // Load workflow on mount if ID provided
  useEffect(() => {
    if (workflowId) {
      loadWorkflow(workflowId);
    }
  }, [workflowId]);

  return {
    workflow,
    nodes,
    setNodes,
    connections,
    setConnections,
    loading,
    error,
    loadWorkflow,
    saveWorkflow,
    createWorkflow,
    executeWorkflow,
    toggleActive,
  };
}

// Hook for managing workflow list
export function useWorkflowList() {
  const [workflows, setWorkflows] = useState<Workflow[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadWorkflows = async () => {
    setLoading(true);
    setError(null);
    try {
      const workflowList = await WorkflowService.getWorkflows();
      setWorkflows(workflowList);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load workflows');
    } finally {
      setLoading(false);
    }
  };

  const deleteWorkflow = async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      await WorkflowService.deleteWorkflow(id);
      setWorkflows(prev => prev.filter(w => w.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete workflow');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadWorkflows();
  }, []);

  return {
    workflows,
    loading,
    error,
    loadWorkflows,
    deleteWorkflow,
  };
}
