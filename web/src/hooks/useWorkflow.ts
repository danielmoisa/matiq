import { useState, useEffect } from 'react';
import { Workflow, WorkflowNode, Connection } from '@/types/workflow';
import { WorkflowService } from '@/services/workflow-service';

// Hook for managing flow state
export function useWorkflow(workflowId?: string) {
  const [flow, setWorkflow] = useState<Workflow | null>(null);
  const [nodes, setNodes] = useState<WorkflowNode[]>([]);
  const [connections, setConnections] = useState<Connection[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load flow data
  const loadWorkflow = async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      const workflowData = await WorkflowService.getWorkflow(id);
      setWorkflow(workflowData);
      setNodes(workflowData.nodes || []);
      setConnections(workflowData.connections || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load flow');
      setNodes([]);
      setConnections([]);
    } finally {
      setLoading(false);
    }
  };

  // Save flow changes
  const saveWorkflow = async () => {
    if (!flow) return;
    
    setLoading(true);
    setError(null);
    try {
      const updatedWorkflow = await WorkflowService.saveWorkflow(
        flow.uid || flow.id, // Use UUID first, fallback to id
        nodes,
        connections
      );
      setWorkflow(updatedWorkflow);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save flow');
    } finally {
      setLoading(false);
    }
  };

  // Create new flow
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
      setError(err instanceof Error ? err.message : 'Failed to create flow');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // Execute flow
  const executeWorkflow = async (input?: Record<string, unknown>) => {
    if (!flow) return;
    
    setLoading(true);
    setError(null);
    try {
      const result = await WorkflowService.executeWorkflow(flow.id, input);
      return result;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to execute flow');
      throw err;
    } finally {
      setLoading(false);
    }
  };


  // Load flow on mount if ID provided
  useEffect(() => {
    if (workflowId) {
      loadWorkflow(workflowId);
    }
  }, [workflowId]);

  return {
    flow,
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
  };
}

// Hook for managing flow list
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
      setWorkflows(prev => prev.filter(w => w.uid !== id && w.id !== id)); // Filter by both uid and id
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete flow');
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
