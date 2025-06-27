'use client';

import { useRef, useState} from 'react';
import { WorkflowNode, Connection } from '@/types/workflow';
import NodeComponent from './NodeComponent';
import {
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import { restrictToWindowEdges } from '@dnd-kit/modifiers';

interface CanvasProps {
  nodes: WorkflowNode[];
  setNodes: (nodes: WorkflowNode[]) => void;
  connections: Connection[];
  setConnections: (connections: Connection[]) => void;
  selectedNode: WorkflowNode | null;
  setSelectedNode: (node: WorkflowNode | null) => void;
}

export default function Canvas({ 
  nodes, 
  setNodes, 
  connections, 
  setConnections, 
  selectedNode, 
  setSelectedNode 
}: CanvasProps) {
  const canvasRef = useRef<HTMLDivElement>(null);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [isConnecting, setIsConnecting] = useState(false);
  const [connectionStart, setConnectionStart] = useState<{ nodeId: string; position: { x: number; y: number } } | null>(null);
  const [tempConnection, setTempConnection] = useState<{ start: { x: number; y: number }; end: { x: number; y: number } } | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    })
  );

  const handleDragStart = (event: DragStartEvent) => {
    setActiveId(event.active.id as string);
    const node = nodes.find(n => n.id === event.active.id);
    if (node) {
      setSelectedNode(node);
    }
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, delta } = event;
    
    if (delta.x !== 0 || delta.y !== 0) {
      setNodes(nodes.map(node => 
        node.id === active.id 
          ? { ...node, position: { x: node.position.x + delta.x, y: node.position.y + delta.y } }
          : node
      ));
    }
    
    setActiveId(null);
  };

  const handleStartConnection = (nodeId: string, position: { x: number; y: number }) => {
    setIsConnecting(true);
    setConnectionStart({ nodeId, position });
    
    // Calculate the actual position relative to the canvas
    const sourceNode = nodes.find(n => n.id === nodeId);
    if (sourceNode) {
      const startPos = {
        x: sourceNode.position.x + 200, // Right edge of node
        y: sourceNode.position.y + 40,  // Middle of node
      };
      setTempConnection({ start: startPos, end: startPos });
    }
  };

  const handleMouseMove = (event: React.MouseEvent) => {
    if (isConnecting && tempConnection && canvasRef.current) {
      const rect = canvasRef.current.getBoundingClientRect();
      const end = {
        x: event.clientX - rect.left,
        y: event.clientY - rect.top,
      };
      setTempConnection({ ...tempConnection, end });
    }
  };

  const handleCompleteConnection = (targetNodeId: string) => {
    if (connectionStart && connectionStart.nodeId !== targetNodeId) {
      const newConnection: Connection = {
        id: `${connectionStart.nodeId}-${targetNodeId}`,
        sourceId: connectionStart.nodeId,
        targetId: targetNodeId,
      };
      
      setConnections([...connections, newConnection]);
    }
    
    setIsConnecting(false);
    setConnectionStart(null);
    setTempConnection(null);
  };

  const handleCanvasClick = (event: React.MouseEvent) => {
    if (isConnecting) {
      setIsConnecting(false);
      setConnectionStart(null);
      setTempConnection(null);
    } else if (event.target === canvasRef.current) {
      setSelectedNode(null);
    }
  };

  return (
    <div className="flex-1 relative overflow-hidden bg-gray-100">
      <DndContext
        sensors={sensors}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        modifiers={[restrictToWindowEdges]}
      >
        <div
          ref={canvasRef}
          className="w-full h-full relative"
          onMouseMove={handleMouseMove}
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

          {/* Connection lines */}
          <svg className="absolute inset-0 pointer-events-none z-10" style={{ width: '100%', height: '100%' }}>
            {/* Arrow marker definitions */}
            <defs>
              <marker
                id="arrowhead"
                markerWidth="12"
                markerHeight="8"
                refX="11"
                refY="4"
                orient="auto"
                markerUnits="strokeWidth"
              >
                <polygon
                  points="0 0, 12 4, 0 8"
                  fill="#3b82f6"
                />
              </marker>
              <marker
                id="arrowhead-temp"
                markerWidth="12"
                markerHeight="8"
                refX="11"
                refY="4"
                orient="auto"
                markerUnits="strokeWidth"
              >
                <polygon
                  points="0 0, 12 4, 0 8"
                  fill="#6b7280"
                />
              </marker>
            </defs>

            {/* Rendered connections */}
            {connections.map((connection) => {
              const sourceNode = nodes.find(n => n.id === connection.sourceId);
              const targetNode = nodes.find(n => n.id === connection.targetId);
              
              if (!sourceNode || !targetNode) return null;
              
              // Calculate connection points (from right edge of source to left edge of target)
              const sourceX = sourceNode.position.x + 200; // Right edge of source node
              const sourceY = sourceNode.position.y + 40; // Middle of source node
              const targetX = targetNode.position.x; // Left edge of target node
              const targetY = targetNode.position.y + 40; // Middle of target node
              
              // Create a curved path for better visual appeal
              const controlPoint1X = sourceX + Math.min(100, Math.abs(targetX - sourceX) / 3);
              const controlPoint2X = targetX - Math.min(100, Math.abs(targetX - sourceX) / 3);
              
              const pathData = `M ${sourceX} ${sourceY} C ${controlPoint1X} ${sourceY}, ${controlPoint2X} ${targetY}, ${targetX} ${targetY}`;
              
              return (
                <g key={connection.id}>
                  {/* Connection path */}
                  <path
                    d={pathData}
                    fill="none"
                    stroke="#3b82f6"
                    strokeWidth="2"
                    markerEnd="url(#arrowhead)"
                    className="hover:stroke-blue-600 transition-colors"
                  />
                  {/* Connection hit area for selection/deletion */}
                  <path
                    d={pathData}
                    fill="none"
                    stroke="transparent"
                    strokeWidth="10"
                    className="cursor-pointer pointer-events-auto"
                    onClick={(e) => {
                      e.stopPropagation();
                      // Remove connection on click
                      setConnections(connections.filter(c => c.id !== connection.id));
                    }}
                  />
                </g>
              );
            })}
            
            {/* Temporary connection line while connecting */}
            {tempConnection && (
              <line
                x1={tempConnection.start.x}
                y1={tempConnection.start.y}
                x2={tempConnection.end.x}
                y2={tempConnection.end.y}
                stroke="#6b7280"
                strokeWidth="2"
                strokeDasharray="5 5"
                opacity="0.8"
                markerEnd="url(#arrowhead-temp)"
              />
            )}
          </svg>

          {/* Nodes */}
          {nodes.map((node) => (
            <NodeComponent
              key={node.id}
              node={node}
              isSelected={selectedNode?.id === node.id}
              isConnecting={isConnecting}
              onStartConnection={handleStartConnection}
              onCompleteConnection={handleCompleteConnection}
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
        
        <DragOverlay>
          {activeId ? (
            <div className="bg-white rounded-lg shadow-lg border-2 border-blue-500 p-3 min-w-[200px] opacity-80">
              <div className="font-medium">{nodes.find(n => n.id === activeId)?.type}</div>
            </div>
          ) : null}
        </DragOverlay>
      </DndContext>
    </div>
  );
}
