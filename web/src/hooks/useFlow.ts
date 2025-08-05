import { useState, useEffect } from 'react';
import { Flow, FlowNode, Connection } from '@/types/flow';
import { FlowService } from '@/services/flow-service';

// Hook for managing flow state
export function useFlow(flowId?: string) {
  const [flow, setFlow] = useState<Flow | null>(null);
  const [nodes, setNodes] = useState<FlowNode[]>([]);
  const [connections, setConnections] = useState<Connection[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load flow data
  const loadFlow = async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      const flowData = await FlowService.getFlow(id);
      setFlow(flowData);
      setNodes(flowData.nodes || []);
      setConnections(flowData.connections || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load flow');
      setNodes([]);
      setConnections([]);
    } finally {
      setLoading(false);
    }
  };

  // Save flow changes
  const saveFlow = async () => {
    if (!flow) return;
    
    setLoading(true);
    setError(null);
    try {
      const updatedFlow = await FlowService.saveFlow(
        flow.uid || flow.id, // Use UUID first, fallback to id
        nodes,
        connections
      );
      setFlow(updatedFlow);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save flow');
    } finally {
      setLoading(false);
    }
  };

  // Create new flow
  const createFlow = async (name: string, description?: string) => {
    setLoading(true);
    setError(null);
    try {
      const newFlow = await FlowService.createFlow({
        name,
        description,
        nodes: [],
        connections: []
      });
      setFlow(newFlow);
      setNodes([]);
      setConnections([]);
      return newFlow;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create flow');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // Execute flow
  const executeFlow = async (input?: Record<string, unknown>) => {
    if (!flow) return;
    
    setLoading(true);
    setError(null);
    try {
      const result = await FlowService.executeFlow(flow.id, input);
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
    if (flowId) {
      loadFlow(flowId);
    }
  }, [flowId]);

  return {
    flow,
    nodes,
    setNodes,
    connections,
    setConnections,
    loading,
    error,
    loadFlow,
    saveFlow,
    createFlow,
    executeFlow,
  };
}

// Hook for managing flow list
export function useFlowList() {
  const [flows, setFlows] = useState<Flow[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadFlows = async () => {
    setLoading(true);
    setError(null);
    try {
      const flowList = await FlowService.getFlows();
      setFlows(flowList);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load flows');
    } finally {
      setLoading(false);
    }
  };

  const deleteFlow = async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      await FlowService.deleteFlow(id);
      setFlows(prev => prev.filter(w => w.uid !== id && w.id !== id)); // Filter by both uid and id
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete flow');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadFlows();
  }, []);

  return {
    flows,
    loading,
    error,
    loadFlows,
    deleteFlow,
  };
}
