'use client';

import { FlowNode, NodeType } from '@/types/flow';
import { useDraggable } from '@dnd-kit/core';

interface NodeComponentProps {
  node: FlowNode;
  isSelected: boolean;
  isConnecting: boolean;
  onStartConnection: (nodeId: string, position: { x: number; y: number }) => void;
  onCompleteConnection: (targetNodeId: string) => void;
  onClick: () => void;
}

const getNodeIcon = (type: string) => {
  const icons: Record<string, string> = {
    // Virtual/Local Actions
    [NodeType.TRANSFORMER]: '⚙️',
    
    // APIs
    [NodeType.RESTAPI]: '🌐',
    [NodeType.GRAPHQL]: '📋',
    
    // Cache/Messaging
    [NodeType.REDIS]: '�',
    [NodeType.UPSTASH]: '⚡',
    
    // Databases
    [NodeType.MYSQL]: '🐬',
    [NodeType.MARIADB]: '🗃️',
    [NodeType.POSTGRESQL]: '🐘',
    [NodeType.MONGODB]: '🍃',
    [NodeType.TIDB]: '⚡',
    [NodeType.ELASTICSEARCH]: '🔍',
    [NodeType.SUPABASEDB]: '⚡',
    [NodeType.FIREBASE]: '🔥',
    [NodeType.CLICKHOUSE]: '📊',
    [NodeType.MSSQL]: '🗄️',
    [NodeType.DYNAMODB]: '🟡',
    [NodeType.SNOWFLAKE]: '❄️',
    [NodeType.COUCHDB]: '🛋️',
    [NodeType.ORACLE]: '🔷',
    [NodeType.ORACLE_9I]: '🔷',
    [NodeType.NEON]: '🌟',
    [NodeType.HYDRA]: '🐍',
    
    // Storage
    [NodeType.S3]: '☁️',
    
    // Communication
    [NodeType.SMTP]: '📧',
    
    // AI/ML
    [NodeType.HUGGINGFACE]: '🤗',
    [NodeType.HFENDPOINT]: '🤖',
    [NodeType.AI_AGENT]: '🤖',
    
    // External Services
    [NodeType.GOOGLESHEETS]: '📊',
    [NodeType.AIRTABLE]: '📋',
    [NodeType.APPWRITE]: '📱',
    
    // Flow Control
    [NodeType.TRIGGER]: '🎯',
    [NodeType.SERVER_SIDE_TRANSFORMER]: '⚙️',
    [NodeType.CONDITION]: '🔀',
    [NodeType.WEBHOOK_RESPONSE]: '📤',
    [NodeType.WF_DRIVE]: '💾',
    
    // Legacy - keeping for backward compatibility
    [NodeType.SCHEDULE]: '⏰',
    [NodeType.WEBHOOK]: '🔗',
    [NodeType.REST_API]: '🌐',
    [NodeType.LOOP]: '🔄',
    [NodeType.RESPONSE]: '📤',
    [NodeType.ERROR_HANDLER]: '⚠️',
  };
  return icons[type] || '📦';
};

const getNodeColor = (type: string) => {
  // Triggers and Flow Control
  if ([NodeType.SCHEDULE, NodeType.WEBHOOK, NodeType.TRIGGER].includes(type as NodeType)) 
    return 'bg-green-100 border-green-300 text-green-800';
  
  // Databases
  if ([NodeType.POSTGRESQL, NodeType.MYSQL, NodeType.MARIADB, NodeType.TIDB, NodeType.NEON, 
       NodeType.MONGODB, NodeType.SNOWFLAKE, NodeType.SUPABASEDB, NodeType.CLICKHOUSE, 
       NodeType.HYDRA, NodeType.ELASTICSEARCH, NodeType.FIREBASE, NodeType.MSSQL, 
       NodeType.DYNAMODB, NodeType.COUCHDB, NodeType.ORACLE, NodeType.ORACLE_9I].includes(type as NodeType)) 
    return 'bg-blue-100 border-blue-300 text-blue-800';
  
  // APIs
  if ([NodeType.REST_API, NodeType.RESTAPI, NodeType.GRAPHQL].includes(type as NodeType)) 
    return 'bg-purple-100 border-purple-300 text-purple-800';
  
  // AI/ML
  if ([NodeType.AI_AGENT, NodeType.HUGGINGFACE, NodeType.HFENDPOINT].includes(type as NodeType)) 
    return 'bg-pink-100 border-pink-300 text-pink-800';
  
  // Storage and Services
  if ([NodeType.S3, NodeType.SMTP, NodeType.GOOGLESHEETS, NodeType.AIRTABLE, NodeType.APPWRITE, NodeType.WF_DRIVE].includes(type as NodeType)) 
    return 'bg-yellow-100 border-yellow-300 text-yellow-800';
  
  // Cache and Messaging
  if ([NodeType.REDIS, NodeType.UPSTASH].includes(type as NodeType)) 
    return 'bg-red-100 border-red-300 text-red-800';
  
  // Flow Actions
  if ([NodeType.TRANSFORMER, NodeType.SERVER_SIDE_TRANSFORMER, NodeType.CONDITION, 
       NodeType.WEBHOOK_RESPONSE, NodeType.LOOP, NodeType.RESPONSE, NodeType.ERROR_HANDLER].includes(type as NodeType)) 
    return 'bg-indigo-100 border-indigo-300 text-indigo-800';
  
  return 'bg-gray-100 border-gray-300 text-gray-800';
};

const formatNodeTitle = (type: string | undefined) => {
  if (!type) return '';
  return type.split('-').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
};

export default function NodeComponent({ 
  node, 
  isSelected, 
  isConnecting, 
  onStartConnection, 
  onCompleteConnection, 
  onClick 
}: NodeComponentProps) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: node.id,
  });

  // During dragging, use transform. After dragging, the position is updated and transform is reset
  const finalStyle = {
    left: node.position.x,
    top: node.position.y,
    transform: transform ? `translate3d(${transform.x}px, ${transform.y}px, 0)` : 'translate3d(0px, 0px, 0)',
    // Remove transition during dragging to prevent conflicts
    transition: isDragging ? 'none' : 'none',
  };

  const handleConnectionClick = (event: React.MouseEvent, isOutput: boolean) => {
    event.stopPropagation();
    
    if (isConnecting) {
      // Complete connection
      onCompleteConnection(node.id);
    } else if (isOutput) {
      // Start connection from output
      const position = {
        x: node.position.x + 200, // Right side of node
        y: node.position.y + 30,  // Middle of node
      };
      onStartConnection(node.id, position);
    }
  };

  return (
    <div
      ref={setNodeRef}
      style={finalStyle}
      className={`absolute select-none ${
        isSelected ? 'ring-2 ring-blue-500 shadow-lg z-20' : 'hover:shadow-md z-10'
      } ${isDragging ? 'opacity-80 z-30' : ''}`}
      {...listeners}
      {...attributes}
      onClick={(e) => {
        e.stopPropagation();
        onClick();
      }}
    >
      <div className={`
        min-w-[200px] p-4 rounded-lg border-2 bg-white shadow-sm cursor-grab active:cursor-grabbing
        ${getNodeColor(node.type)}
      `}>
        <div className="flex items-center space-x-3">
          <span className="text-2xl">{getNodeIcon(node.type)}</span>
          <div>
            <h3 className="font-medium text-sm">{formatNodeTitle(node.type)}</h3>
            {node.triggerType && (
              <p className="text-xs opacity-75">Trigger: {formatNodeTitle(node.triggerType)}</p>
            )}
          </div>
        </div>
        
        {/* Connection points */}
        {/* Output connection point (right side) */}
        <div 
          className={`absolute -right-2 top-1/2 transform -translate-y-1/2 w-4 h-4 bg-white border-2 rounded-full cursor-pointer transition-colors ${
            isConnecting 
              ? 'border-gray-300 hover:border-gray-400' 
              : 'border-blue-300 hover:border-blue-500'
          }`}
          onClick={(e) => handleConnectionClick(e, true)}
          title="Output - Click to start connection"
        />
        
        {/* Input connection point (left side) */}
        <div 
          className={`absolute -left-2 top-1/2 transform -translate-y-1/2 w-4 h-4 bg-white border-2 rounded-full cursor-pointer transition-colors ${
            isConnecting 
              ? 'border-green-300 hover:border-green-500' 
              : 'border-gray-300 hover:border-gray-400'
          }`}
          onClick={(e) => handleConnectionClick(e, false)}
          title="Input - Click to complete connection"
        />
      </div>
    </div>
  );
}
