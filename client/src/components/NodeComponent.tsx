'use client';

import { WorkflowNode } from '@/types/workflow';

interface NodeComponentProps {
  node: WorkflowNode;
  isSelected: boolean;
  onDragStart: (event: React.MouseEvent) => void;
  onClick: () => void;
}

const getNodeIcon = (type: string) => {
  const icons: Record<string, string> = {
    'schedule': '⏰',
    'webhook': '🔗',
    'postgres': '🐘',
    'mysql': '🐬',
    'mariadb': '🗃️',
    'tidb': '⚡',
    'neon': '🌟',
    'mongodb': '🍃',
    'snowflake': '❄️',
    'supabase': '⚡',
    'clickhouse': '📊',
    'hydra': '🐍',
    'rest-api': '🌐',
    'graphql': '📋',
    'ai-agent': '🤖',
    'transformer': '⚙️',
    'condition': '🔀',
    'loop': '🔄',
    'response': '📤',
    'error-handler': '⚠️',
  };
  return icons[type] || '📦';
};

const getNodeColor = (type: string) => {
  if (['schedule', 'webhook'].includes(type)) return 'bg-green-100 border-green-300 text-green-800';
  if (['postgres', 'mysql', 'mariadb', 'tidb', 'neon', 'mongodb', 'snowflake', 'supabase', 'clickhouse', 'hydra'].includes(type)) return 'bg-blue-100 border-blue-300 text-blue-800';
  if (['rest-api', 'graphql'].includes(type)) return 'bg-purple-100 border-purple-300 text-purple-800';
  return 'bg-gray-100 border-gray-300 text-gray-800';
};

const formatNodeTitle = (type: string) => {
  return type.split('-').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
};

export default function NodeComponent({ node, isSelected, onDragStart, onClick }: NodeComponentProps) {
  return (
    <div
      className={`absolute cursor-move select-none transition-all ${
        isSelected ? 'ring-2 ring-blue-500 shadow-lg' : 'hover:shadow-md'
      }`}
      style={{
        left: node.position.x,
        top: node.position.y,
        transform: isSelected ? 'scale(1.05)' : 'scale(1)',
      }}
      onMouseDown={onDragStart}
      onClick={(e) => {
        e.stopPropagation();
        onClick();
      }}
    >
      <div className={`
        min-w-[200px] p-4 rounded-lg border-2 bg-white shadow-sm
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
        <div className="absolute -right-2 top-1/2 transform -translate-y-1/2 w-4 h-4 bg-white border-2 border-gray-300 rounded-full hover:border-blue-500"></div>
        <div className="absolute -left-2 top-1/2 transform -translate-y-1/2 w-4 h-4 bg-white border-2 border-gray-300 rounded-full hover:border-blue-500"></div>
      </div>
    </div>
  );
}
