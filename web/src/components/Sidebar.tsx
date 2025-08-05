'use client';

import { useState } from 'react';
import { EventType, TriggerType, NodeType } from '@/types/flow';

interface SidebarProps {
  onAddNode: (type: EventType, triggerType?: TriggerType) => void;
}

const triggers = [
  { type: NodeType.SCHEDULE, label: 'Schedule', icon: '⏰', description: 'Run on a schedule' },
  { type: NodeType.WEBHOOK, label: 'Webhook', icon: '🔗', description: 'Trigger via HTTP request' },
  { type: NodeType.TRIGGER, label: 'Trigger', icon: '🎯', description: 'Generic trigger' },
];

const databases = [
  { type: NodeType.POSTGRESQL, label: 'PostgreSQL', icon: '🐘' },
  { type: NodeType.MYSQL, label: 'MySQL', icon: '🐬' },
  { type: NodeType.MARIADB, label: 'MariaDB', icon: '🗃️' },
  { type: NodeType.TIDB, label: 'TiDB', icon: '⚡' },
  { type: NodeType.NEON, label: 'Neon', icon: '🌟' },
  { type: NodeType.MONGODB, label: 'MongoDB', icon: '🍃' },
  { type: NodeType.SNOWFLAKE, label: 'Snowflake', icon: '❄️' },
  { type: NodeType.SUPABASEDB, label: 'Supabase DB', icon: '⚡' },
  { type: NodeType.CLICKHOUSE, label: 'ClickHouse', icon: '📊' },
  { type: NodeType.HYDRA, label: 'Hydra', icon: '🐍' },
  { type: NodeType.MSSQL, label: 'MS SQL Server', icon: '🗄️' },
  { type: NodeType.ORACLE, label: 'Oracle', icon: '🔷' },
  { type: NodeType.ORACLE_9I, label: 'Oracle 9i', icon: '🔷' },
  { type: NodeType.ELASTICSEARCH, label: 'Elasticsearch', icon: '🔍' },
  { type: NodeType.FIREBASE, label: 'Firebase', icon: '🔥' },
  { type: NodeType.DYNAMODB, label: 'DynamoDB', icon: '🟡' },
  { type: NodeType.COUCHDB, label: 'CouchDB', icon: '🛋️' },
];

const apis = [
  { type: NodeType.REST_API, label: 'REST API', icon: '🌐' },
  { type: NodeType.RESTAPI, label: 'REST API (Backend)', icon: '🌐' },
  { type: NodeType.GRAPHQL, label: 'GraphQL', icon: '📋' },
];

const storage = [
  { type: NodeType.S3, label: 'Amazon S3', icon: '☁️' },
  { type: NodeType.REDIS, label: 'Redis', icon: '🔴' },
  { type: NodeType.UPSTASH, label: 'Upstash', icon: '⚡' },
  { type: NodeType.WF_DRIVE, label: 'Workflow Drive', icon: '💾' },
];

const communication = [
  { type: NodeType.SMTP, label: 'SMTP Email', icon: '📧' },
  { type: NodeType.WEBHOOK_RESPONSE, label: 'Webhook Response', icon: '📤' },
];

const ai = [
  { type: NodeType.AI_AGENT, label: 'AI Agent', icon: '🤖' },
  { type: NodeType.HUGGINGFACE, label: 'Hugging Face', icon: '🤗' },
  { type: NodeType.HFENDPOINT, label: 'HF Endpoint', icon: '🤖' },
];

const external = [
  { type: NodeType.GOOGLESHEETS, label: 'Google Sheets', icon: '📊' },
  { type: NodeType.AIRTABLE, label: 'Airtable', icon: '📋' },
  { type: NodeType.APPWRITE, label: 'Appwrite', icon: '📱' },
];

const actions = [
  { type: NodeType.TRANSFORMER, label: 'Transformer', icon: '⚙️' },
  { type: NodeType.SERVER_SIDE_TRANSFORMER, label: 'Server Transformer', icon: '⚙️' },
  { type: NodeType.CONDITION, label: 'Condition', icon: '🔀' },
  { type: NodeType.LOOP, label: 'Loop', icon: '🔄' },
  { type: NodeType.RESPONSE, label: 'Response', icon: '📤' },
  { type: NodeType.ERROR_HANDLER, label: 'Error Handler', icon: '⚠️' },
];

export default function Sidebar({ onAddNode }: SidebarProps) {
  const [activeTab, setActiveTab] = useState<'triggers' | 'databases' | 'apis' | 'storage' | 'communication' | 'ai' | 'external' | 'actions'>('triggers');

  const renderNodeList = (nodes: { type: EventType; label: string; icon: string; description?: string }[], isTrigger = false) => (
    <div className="space-y-2">
      {nodes.map((node) => (
        <button
          key={node.type}
          onClick={() => isTrigger 
            ? onAddNode(node.type, node.type as TriggerType) 
            : onAddNode(node.type)
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
        <p className="text-sm text-gray-500 mt-1">Drag and drop to build your flow</p>
      </div>

      <div className="flex flex-wrap border-b border-gray-200">
        {[
          { key: 'triggers', label: 'Triggers' },
          { key: 'databases', label: 'Databases' },
          { key: 'apis', label: 'APIs' },
          { key: 'storage', label: 'Storage' },
          { key: 'communication', label: 'Comm' },
          { key: 'ai', label: 'AI/ML' },
          { key: 'external', label: 'External' },
          { key: 'actions', label: 'Actions' },
        ].map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key as typeof activeTab)}
            className={`px-2 py-2 text-xs font-medium border-b-2 transition-colors ${
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
        {activeTab === 'storage' && renderNodeList(storage)}
        {activeTab === 'communication' && renderNodeList(communication)}
        {activeTab === 'ai' && renderNodeList(ai)}
        {activeTab === 'external' && renderNodeList(external)}
        {activeTab === 'actions' && renderNodeList(actions)}
      </div>
    </aside>
  );
}
