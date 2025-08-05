'use client';

import { FlowNode, NodeType } from '@/types/flow';

interface PropertiesPanelProps {
  node: FlowNode | null;
  onUpdateNode: (nodeId: string, data: Record<string, unknown>) => void;
}

export default function PropertiesPanel({ node, onUpdateNode }: PropertiesPanelProps) {
  if (!node) {
    return (
      <div className="w-80 bg-white border-l border-gray-200 p-6">
        <div className="text-center text-gray-500">
          <div className="text-4xl mb-4">📋</div>
          <h3 className="text-lg font-medium mb-2">Properties</h3>
          <p className="text-sm">Select a node to configure its properties</p>
        </div>
      </div>
    );
  }

  const formatTitle = (type: string) => {
    if (!type || typeof type !== 'string') {
      return 'Unknown Node';
    }
    return type.split('-').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
  };

  const renderDatabaseConfig = () => (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Connection String</label>
        <input
          type="text"
          placeholder="postgresql://user:password@host:port/database"
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-600 text-gray-900"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Query</label>
        <textarea
          rows={4}
          placeholder="SELECT * FROM users WHERE created_at > NOW() - INTERVAL '1 hour'"
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-600 text-gray-900"
        />
      </div>
    </div>
  );

  const renderAPIConfig = () => (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">URL</label>
        <input
          type="url"
          placeholder="https://api.example.com/endpoint"
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-600 text-gray-900"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Method</label>
        <select className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white text-gray-900">
          <option>GET</option>
          <option>POST</option>
          <option>PUT</option>
          <option>DELETE</option>
        </select>
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Headers</label>
        <textarea
          rows={3}
          placeholder='{"Authorization": "Bearer token", "Content-Type": "application/json"}'
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-600 text-gray-900"
        />
      </div>
    </div>
  );

  const renderWebhookConfig = () => (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Webhook URL</label>
        <div className="flex items-center space-x-2">
          <input
            type="text"
            value={`https://your-domain.com/webhook/${node.id}`}
            readOnly
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-600 text-sm"
          />
          <button 
            className="px-3 py-2 text-sm bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            onClick={() => navigator.clipboard.writeText(`https://your-domain.com/webhook/${node.id}`)}
          >
            Copy
          </button>
        </div>
        <p className="text-xs text-gray-500 mt-1">This URL will trigger your flow when called</p>
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">HTTP Method</label>
        <select className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white text-gray-900">
          <option>POST</option>
          <option>GET</option>
          <option>PUT</option>
          <option>PATCH</option>
        </select>
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Authentication</label>
        <select className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white text-gray-900">
          <option>None</option>
          <option>API Key</option>
          <option>Bearer Token</option>
          <option>Basic Auth</option>
        </select>
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Expected Content Type</label>
        <select className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white text-gray-900">
          <option>application/json</option>
          <option>application/x-www-form-urlencoded</option>
          <option>text/plain</option>
          <option>application/xml</option>
        </select>
      </div>
    </div>
  );

  const renderScheduleConfig = () => (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Schedule Type</label>
        <select className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white text-gray-900">
          <option>Every minute</option>
          <option>Every hour</option>
          <option>Daily</option>
          <option>Weekly</option>
          <option>Custom cron</option>
        </select>
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Cron Expression</label>
        <input
          type="text"
          placeholder="0 0 * * *"
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-600 text-gray-900"
        />
      </div>
    </div>
  );

  const renderTransformerConfig = () => (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">JavaScript Code</label>
        <textarea
          rows={8}
          placeholder={`// Transform the input data
function transform(input) {
  return {
    ...input,
    processedAt: new Date().toISOString()
  };
}`}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm placeholder-gray-600 text-gray-900"
        />
      </div>
    </div>
  );

  const renderConditionConfig = () => (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Condition</label>
        <select className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white text-gray-900">
          <option>Equals</option>
          <option>Not equals</option>
          <option>Greater than</option>
          <option>Less than</option>
          <option>Contains</option>
          <option>Custom expression</option>
        </select>
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Value</label>
        <input
          type="text"
          placeholder="Enter value to compare"
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent placeholder-gray-600 text-gray-900"
        />
      </div>
    </div>
  );

  const renderConfig = () => {
    switch (node.type) {
      case NodeType.WEBHOOK:
      case NodeType.TRIGGER:
        return renderWebhookConfig();
      case NodeType.SCHEDULE:
        return renderScheduleConfig();
      case NodeType.POSTGRESQL:
      case NodeType.MYSQL:
      case NodeType.MARIADB:
      case NodeType.TIDB:
      case NodeType.NEON:
      case NodeType.MONGODB:
      case NodeType.SNOWFLAKE:
      case NodeType.SUPABASEDB:
      case NodeType.CLICKHOUSE:
      case NodeType.HYDRA:
      case NodeType.MSSQL:
      case NodeType.ORACLE:
      case NodeType.ORACLE_9I:
      case NodeType.ELASTICSEARCH:
      case NodeType.FIREBASE:
      case NodeType.DYNAMODB:
      case NodeType.COUCHDB:
        return renderDatabaseConfig();
      case NodeType.REST_API:
      case NodeType.RESTAPI:
      case NodeType.GRAPHQL:
        return renderAPIConfig();
      case NodeType.TRANSFORMER:
      case NodeType.SERVER_SIDE_TRANSFORMER:
        return renderTransformerConfig();
      case NodeType.CONDITION:
        return renderConditionConfig();
      case NodeType.AI_AGENT:
      case NodeType.HUGGINGFACE:
      case NodeType.HFENDPOINT:
        return (
          <div className="text-sm text-gray-500">
            AI/ML configuration options coming soon...
          </div>
        );
      case NodeType.S3:
      case NodeType.REDIS:
      case NodeType.UPSTASH:
      case NodeType.WF_DRIVE:
        return (
          <div className="text-sm text-gray-500">
            Storage configuration options coming soon...
          </div>
        );
      case NodeType.SMTP:
      case NodeType.WEBHOOK_RESPONSE:
        return (
          <div className="text-sm text-gray-500">
            Communication configuration options coming soon...
          </div>
        );
      case NodeType.GOOGLESHEETS:
      case NodeType.AIRTABLE:
      case NodeType.APPWRITE:
        return (
          <div className="text-sm text-gray-500">
            External service configuration options coming soon...
          </div>
        );
      default:
        return (
          <div className="text-sm text-gray-500">
            Configuration options for {formatTitle(node.type)} coming soon...
          </div>
        );
    }
  };

  return (
    <div className="w-80 bg-white border-l border-gray-200 flex flex-col">
      <div className="p-6 border-b border-gray-200">
        <h3 className="text-lg font-semibold text-gray-900">{formatTitle(node.type)}</h3>
        <p className="text-sm text-gray-500 mt-1">Configure this node</p>
      </div>
      
      <div className="flex-1 p-6 overflow-y-auto">
        {renderConfig()}
      </div>
      
      <div className="p-6 border-t border-gray-200">
        <button
          onClick={() => onUpdateNode(node.id, {})}
          className="w-full px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
        >
          Save Changes
        </button>
      </div>
    </div>
  );
}
