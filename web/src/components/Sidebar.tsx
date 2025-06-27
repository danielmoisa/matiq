'use client';

import { useState } from 'react';
import { EventType, TriggerType } from '@/types/workflow';

interface SidebarProps {
  onAddNode: (type: EventType, triggerType?: TriggerType) => void;
}

const triggers = [
  { type: 'schedule' as TriggerType, label: 'Schedule', icon: '⏰', description: 'Run on a schedule' },
  { type: 'webhook' as TriggerType, label: 'Webhook', icon: '🔗', description: 'Trigger via HTTP request' },
];

const databases = [
  { type: 'postgres' as EventType, label: 'PostgreSQL', icon: '🐘' },
  { type: 'mysql' as EventType, label: 'MySQL', icon: '🐬' },
  { type: 'mariadb' as EventType, label: 'MariaDB', icon: '🗃️' },
  { type: 'tidb' as EventType, label: 'TiDB', icon: '⚡' },
  { type: 'neon' as EventType, label: 'Neon', icon: '🌟' },
  { type: 'mongodb' as EventType, label: 'MongoDB', icon: '🍃' },
  { type: 'snowflake' as EventType, label: 'Snowflake', icon: '❄️' },
  { type: 'supabase' as EventType, label: 'Supabase', icon: '⚡' },
  { type: 'clickhouse' as EventType, label: 'ClickHouse', icon: '📊' },
  { type: 'hydra' as EventType, label: 'Hydra', icon: '🐍' },
];

const apis = [
  { type: 'rest-api' as EventType, label: 'REST API', icon: '🌐' },
  { type: 'graphql' as EventType, label: 'GraphQL', icon: '📋' },
];

const actions = [
  { type: 'ai-agent' as EventType, label: 'AI Agent', icon: '🤖' },
  { type: 'transformer' as EventType, label: 'Transformer', icon: '⚙️' },
  { type: 'condition' as EventType, label: 'Condition', icon: '🔀' },
  { type: 'loop' as EventType, label: 'Loop', icon: '🔄' },
  { type: 'response' as EventType, label: 'Response', icon: '📤' },
  { type: 'error-handler' as EventType, label: 'Error Handler', icon: '⚠️' },
];

export default function Sidebar({ onAddNode }: SidebarProps) {
  const [activeTab, setActiveTab] = useState<'triggers' | 'databases' | 'apis' | 'actions'>('triggers');

  const renderNodeList = (nodes: { type: EventType | TriggerType; label: string; icon: string; description?: string }[], isTrigger = false) => (
    <div className="space-y-2">
      {nodes.map((node) => (
        <button
          key={node.type}
          onClick={() => isTrigger 
            ? onAddNode('schedule' as EventType, node.type as TriggerType) 
            : onAddNode(node.type as EventType)
          }
          className="w-full p-3 text-left bg-white border border-gray-200 rounded-lg hover:border-blue-300 hover:shadow-sm transition-all group"
        >
          <div className="flex items-center space-x-3">
            <span className="text-lg">{node.icon}</span>
            <div>
              <div className="font-medium text-gray-900 group-hover:text-blue-600">{node.label}</div>
              {node.description && (
                <div className="text-xs text-gray-500">{node.description}</div>
              )}
            </div>
          </div>
        </button>
      ))}
    </div>
  );

  return (
    <aside className="w-80 bg-gray-50 border-r border-gray-200 flex flex-col">
      <div className="p-4 border-b border-gray-200">
        <h2 className="text-lg font-semibold text-gray-900">Components</h2>
        <p className="text-sm text-gray-500 mt-1">Drag and drop to build your workflow</p>
      </div>

      <div className="flex border-b border-gray-200">
        {[
          { key: 'triggers', label: 'Triggers' },
          { key: 'databases', label: 'Databases' },
          { key: 'apis', label: 'APIs' },
          { key: 'actions', label: 'Actions' },
        ].map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key as typeof activeTab)}
            className={`flex-1 px-3 py-2 text-sm font-medium border-b-2 transition-colors ${
              activeTab === tab.key
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div className="flex-1 p-4 overflow-y-auto">
        {activeTab === 'triggers' && renderNodeList(triggers, true)}
        {activeTab === 'databases' && renderNodeList(databases)}
        {activeTab === 'apis' && renderNodeList(apis)}
        {activeTab === 'actions' && renderNodeList(actions)}
      </div>
    </aside>
  );
}
