'use client';

import { useRef, useState } from 'react';
import { WorkflowNode } from '@/types/workflow';
import NodeComponent from './NodeComponent';

interface CanvasProps {
  nodes: WorkflowNode[];
  setNodes: (nodes: WorkflowNode[]) => void;
  selectedNode: WorkflowNode | null;
  setSelectedNode: (node: WorkflowNode | null) => void;
}

export default function Canvas({ nodes, setNodes, selectedNode, setSelectedNode }: CanvasProps) {
  const canvasRef = useRef<HTMLDivElement>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 });

  const handleNodeDragStart = (nodeId: string, event: React.MouseEvent) => {
    const node = nodes.find(n => n.id === nodeId);
    if (!node) return;

    setIsDragging(true);
    setSelectedNode(node);
    
    const rect = event.currentTarget.getBoundingClientRect();
    setDragOffset({
      x: event.clientX - rect.left,
      y: event.clientY - rect.top,
    });
  };

  const handleMouseMove = (event: React.MouseEvent) => {
    if (!isDragging || !selectedNode) return;

    const canvasRect = canvasRef.current?.getBoundingClientRect();
    if (!canvasRect) return;

    const newX = event.clientX - canvasRect.left - dragOffset.x;
    const newY = event.clientY - canvasRect.top - dragOffset.y;

    setNodes(nodes.map(node => 
      node.id === selectedNode.id 
        ? { ...node, position: { x: newX, y: newY } }
        : node
    ));
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const handleCanvasClick = (event: React.MouseEvent) => {
    if (event.target === canvasRef.current) {
      setSelectedNode(null);
    }
  };

  return (
    <div className="flex-1 relative overflow-hidden bg-gray-100">
      <div
        ref={canvasRef}
        className="w-full h-full relative"
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onClick={handleCanvasClick}
      >
        {/* Grid pattern */}
        <div 
          className="absolute inset-0 opacity-20"
          style={{
            backgroundImage: `
              linear-gradient(rgba(0,0,0,0.1) 1px, transparent 1px),
              linear-gradient(90deg, rgba(0,0,0,0.1) 1px, transparent 1px)
            `,
            backgroundSize: '20px 20px',
          }}
        />

        {/* Nodes */}
        {nodes.map((node) => (
          <NodeComponent
            key={node.id}
            node={node}
            isSelected={selectedNode?.id === node.id}
            onDragStart={(event: React.MouseEvent) => handleNodeDragStart(node.id, event)}
            onClick={() => setSelectedNode(node)}
          />
        ))}

        {/* Instructions when empty */}
        {nodes.length === 0 && (
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="text-center text-gray-500">
              <div className="text-6xl mb-4">⚡</div>
              <h3 className="text-xl font-medium mb-2">Build Your First Workflow</h3>
              <p className="text-sm">Start by adding a trigger from the sidebar</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
