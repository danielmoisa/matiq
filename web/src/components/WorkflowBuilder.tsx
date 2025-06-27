'use client';

import { useState } from 'react';

import PropertiesPanel from './PropertiesPanel';
import { WorkflowNode, Connection, TriggerType, EventType } from '@/types/workflow';
import Header from './Header';
import Sidebar from './Sidebar';
import Canvas from './Canvas';

export default function WorkflowBuilder() {
  const [nodes, setNodes] = useState<WorkflowNode[]>([]);
  const [connections, setConnections] = useState<Connection[]>([]);
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
    <div className="h-screen flex flex-col">
      <Header />
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
