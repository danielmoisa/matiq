'use client';

import { useState, useEffect } from 'react';
import { useWorkflow } from '@/hooks/useWorkflow';
import PropertiesPanel from './PropertiesPanel';
import { WorkflowNode, TriggerType, EventType } from '@/types/workflow';
import Header from './Header';
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
    saveWorkflow,
    createWorkflow
  } = useWorkflow(workflowId);

  const [selectedNode, setSelectedNode] = useState<WorkflowNode | null>(null);
  const [autoSave, setAutoSave] = useState(true);

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
    <div className="h-screen flex flex-col">
      <Header 
        workflowName={workflow?.name || 'Untitled Workflow'}
        onWorkflowNameChange={(name) => {
          // TODO: Update workflow name through API
          console.log('Updating workflow name to:', name);
        }}
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
