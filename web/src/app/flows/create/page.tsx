'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { FlowService } from '@/services/flow-service';
import { FlowNode, Connection, EventType, TriggerType } from '@/types/flow';
import { AppLayout } from '@/components/layout/AppLayout';
import Sidebar from '@/components/Sidebar';
import Canvas from '@/components/Canvas';
import PropertiesPanel from '@/components/PropertiesPanel';

export default function CreateFlowPage() {
  const router = useRouter();
  const [flowName, setFlowName] = useState('');
  const [flowDescription, setFlowDescription] = useState('');
  const [nodes, setNodes] = useState<FlowNode[]>([]);
  const [connections, setConnections] = useState<Connection[]>([]);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showNameModal, setShowNameModal] = useState(false);

  const handleSaveFlow = async () => {
    if (!flowName.trim()) {
      setShowNameModal(true);
      return;
    }

    setSaving(true);
    setError(null);
    
    try {
      // Create the flow payload using the template from nodes and connections
      const flowData = {
        name: flowName.trim(),
        description: flowDescription.trim(),
        nodes,
        connections,
        status: 'draft' as const,
        isActive: false,
        triggerMode: '1', // Default trigger mode
        actionType: 'restapi' // Default flow type
      };

      // Create the flow
      const createdFlow = await FlowService.createFlow(flowData, nodes, connections);

      // Redirect to the created flow's edit page
      router.push(`/flows/${createdFlow.uid}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create flow');
    } finally {
      setSaving(false);
    }
  };

  const handlePublishFlow = async () => {
    if (!flowName.trim()) {
      setShowNameModal(true);
      return;
    }

    setSaving(true);
    setError(null);
    
    try {
      // Create the flow payload with active status
      const flowData = {
        name: flowName.trim(),
        description: flowDescription.trim(),
        nodes,
        connections,
        status: 'active' as const,
        isActive: true,
        triggerMode: '1',
        actionType: 'restapi'
      };

      // Create and publish the flow
      const createdFlow = await FlowService.createFlow(flowData, nodes, connections);

      // Redirect to the created flow's page
      router.push(`/flows/${createdFlow.uid}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create and publish flow');
    } finally {
      setSaving(false);
    }
  };

  const handleNameModalSubmit = (name: string, description: string) => {
    setFlowName(name);
    setFlowDescription(description);
    setShowNameModal(false);
    // The save will be triggered again automatically
    setTimeout(() => handleSaveFlow(), 100);
  };

  return (
    <AppLayout>
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center space-x-4">
              <button
                onClick={() => router.back()}
                className="text-gray-500 hover:text-gray-700"
              >
                ← Back
              </button>
              <div>
                <h1 className="text-2xl font-bold text-gray-900">
                  {flowName || 'New Flow'}
                </h1>
                {flowDescription && (
                  <p className="text-gray-600 text-sm mt-1">{flowDescription}</p>
                )}
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              <button
                onClick={() => setShowNameModal(true)}
                className="px-4 py-2 text-gray-700 hover:text-gray-900 border border-gray-300 rounded-md font-medium transition-colors"
              >
                {flowName ? 'Edit Details' : 'Set Name'}
              </button>
              <button
                onClick={handleSaveFlow}
                disabled={saving || nodes.length === 0}
                className="px-4 py-2 bg-gray-600 hover:bg-gray-700 text-white rounded-md font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {saving ? 'Saving...' : 'Save Draft'}
              </button>
              <button
                onClick={handlePublishFlow}
                disabled={saving || nodes.length === 0}
                className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {saving ? 'Publishing...' : 'Save & Publish'}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3">
          <div className="flex">
            <div className="flex-shrink-0">⚠️</div>
            <div className="ml-3">
              <h3 className="text-sm font-medium">Error</h3>
              <p className="text-sm mt-1">{error}</p>
            </div>
          </div>
        </div>
      )}

      {/* Flow Builder */}
      <div className="flex-1">
        <CreateFlowBuilder
          nodes={nodes}
          setNodes={setNodes}
          connections={connections}
          setConnections={setConnections}
          flowName={flowName}
          onSave={handleSaveFlow}
          onPublish={handlePublishFlow}
          saving={saving}
        />
      </div>

      {/* Name/Description Modal */}
      {showNameModal && (
        <FlowNameModal
          initialName={flowName}
          initialDescription={flowDescription}
          onSubmit={handleNameModalSubmit}
          onCancel={() => setShowNameModal(false)}
        />
      )}
    </AppLayout>
  );
}

// Custom flow builder for creation (without loading existing flow)
function CreateFlowBuilder({
  nodes,
  setNodes,
  connections,
  setConnections,
  flowName,
  onSave,
  onPublish,
  saving
}: {
  nodes: FlowNode[];
  setNodes: (nodes: FlowNode[]) => void;
  connections: Connection[];
  setConnections: (connections: Connection[]) => void;
  flowName: string;
  onSave: () => void;
  onPublish: () => void;
  saving: boolean;
}) {
  const [selectedNode, setSelectedNode] = useState<FlowNode | null>(null);

  const addNode = (type: EventType, triggerType?: TriggerType) => {
    const newNode: FlowNode = {
      id: Date.now().toString(),
      type,
      triggerType,
      position: { 
        x: 100 + (nodes.length * 250),
        y: 100 + (nodes.length % 3) * 150
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
    <div className="h-full flex flex-col">
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

// Modal for setting flow name and description
function FlowNameModal({
  initialName,
  initialDescription,
  onSubmit,
  onCancel
}: {
  initialName: string;
  initialDescription: string;
  onSubmit: (name: string, description: string) => void;
  onCancel: () => void;
}) {
  const [name, setName] = useState(initialName);
  const [description, setDescription] = useState(initialDescription);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (name.trim()) {
      onSubmit(name.trim(), description.trim());
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg max-w-md w-full">
        <form onSubmit={handleSubmit}>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              {initialName ? 'Edit Flow Details' : 'Set Flow Details'}
            </h3>
            
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Flow Name *
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Enter flow name"
                required
                autoFocus
              />
            </div>
            
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Description
              </label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Describe what this flow does"
              />
            </div>
          </div>
          
          <div className="px-6 py-4 bg-gray-50 flex justify-end space-x-3 rounded-b-lg">
            <button
              type="button"
              onClick={onCancel}
              className="px-4 py-2 text-gray-700 hover:text-gray-900 font-medium"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-medium transition-colors disabled:opacity-50"
              disabled={!name.trim()}
            >
              {initialName ? 'Update' : 'Set Details'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
