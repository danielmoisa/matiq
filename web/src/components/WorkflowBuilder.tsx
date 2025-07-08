'use client';

import { useState } from 'react';
import { useWorkflow } from '@/hooks/useWorkflow';
import PropertiesPanel from './PropertiesPanel';
import { WorkflowNode, TriggerType, EventType } from '@/types/workflow';
import WorkflowHeader from './WorkflowHeader';
import Sidebar from './Sidebar';
import Canvas from './Canvas';

interface WorkflowBuilderProps {
  workflowId?: string;
}

export default function WorkflowBuilder({ workflowId }: WorkflowBuilderProps) {
  const {
    workflow,
    nodes,
    setNodes,
    connections,
    setConnections,
    loading,
    error,
    saveWorkflow
  } = useWorkflow(workflowId);

  const [selectedNode, setSelectedNode] = useState<WorkflowNode | null>(null);

  const addNode = (type: EventType, triggerType?: TriggerType) => {
    const newNode: WorkflowNode = {
      id: Date.now().toString(),
      type,
      triggerType,
      position: { 
        x: 100 + (nodes.length * 250), // Stagger nodes horizontally
        y: 100 + (nodes.length % 3) * 150 // Stagger nodes vertically every 3 nodes
      },
      data: {},
    };
    setNodes([...nodes, newNode]);
  };

  const updateNode = (nodeId: string, data: Record<string, unknown>) => {
    setNodes(nodes.map(node => 
      node.id === nodeId 
        ? { ...node, data: { ...node.data, ...data } }
        : node
    ));
  };


  return (
    <div className="h-full flex flex-col">
      {error && (
        <div className="bg-red-50 border-l-4 border-red-400 p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-sm text-red-700">
                <strong>Error:</strong> {error}
              </p>
            </div>
          </div>
        </div>
      )}
      <WorkflowHeader 
        workflowName={workflow?.name}
        onSave={saveWorkflow}
        onPublish={() => {
          // TODO: Implement publish functionality
          console.log('Publishing workflow...');
        }}
        saving={loading}
      />
      <div className="flex flex-1">
        <Sidebar onAddNode={addNode} />
        <Canvas 
          nodes={nodes} 
          setNodes={setNodes}
          connections={connections}
          setConnections={setConnections}
          selectedNode={selectedNode}
          setSelectedNode={setSelectedNode}
        />
        <PropertiesPanel 
          node={selectedNode}
          onUpdateNode={updateNode}
        />
      </div>
    </div>
  );
}
