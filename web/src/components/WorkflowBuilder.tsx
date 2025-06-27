'use client';

import { useState } from 'react';

import PropertiesPanel from './PropertiesPanel';
import { WorkflowNode, TriggerType, EventType } from '@/types/workflow';
import Header from './Header';
import Sidebar from './Sidebar';
import Canvas from './Canvas';

export default function WorkflowBuilder() {
  const [nodes, setNodes] = useState<WorkflowNode[]>([]);
  const [selectedNode, setSelectedNode] = useState<WorkflowNode | null>(null);

  const addNode = (type: EventType, triggerType?: TriggerType) => {
    const newNode: WorkflowNode = {
      id: Date.now().toString(),
      type,
      triggerType,
      position: { x: 100, y: 100 },
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
      <Header />
      <div className="flex flex-1">
        <Sidebar onAddNode={addNode} />
        <Canvas 
          nodes={nodes} 
          setNodes={setNodes}
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
