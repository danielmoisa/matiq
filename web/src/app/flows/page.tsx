'use client';

import { useState } from 'react';
import Link from 'next/link';

import { FlowService } from '@/services/flow-service';
import { AppLayout } from '@/components/layout/AppLayout';
import { useFlowList } from '@/hooks/useFlow';

export default function FlowsPage() {
  const { flows, loading, error, loadFlows, deleteFlow } = useFlowList();
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [createLoading, setCreateLoading] = useState(false);

  const handleCreateFlow = async (name: string, description: string) => {
    setCreateLoading(true);
    try {
      await FlowService.createFlow({ name, description }, [], []);
      setShowCreateModal(false);
      loadFlows(); // Refresh the list
    } catch (error) {
      console.error('Failed to create flow:', error);
    } finally {
      setCreateLoading(false);
    }
  };

  const handleDeleteFlow = async (id: string, name: string) => {
    if (confirm(`Are you sure you want to delete "${name}"?`)) {
      try {
        await deleteFlow(id);
      } catch (error) {
        console.error('Failed to delete flow:', error);
      }
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'bg-green-100 text-green-800';
      case 'draft': return 'bg-gray-100 text-gray-800';
      case 'paused': return 'bg-yellow-100 text-yellow-800';
      case 'error': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };


  if (loading && flows.length === 0) {
    return (
      <AppLayout>
        <div className="min-h-[80vh] flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading flows...</p>
          </div>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
            <div className="flex">
              <div className="flex-shrink-0">⚠️</div>
              <div className="ml-3">
                <h3 className="text-sm font-medium">Error loading flows</h3>
                <p className="text-sm mt-1">{error}</p>
              </div>
            </div>
          </div>
        )}

        {!loading && (
          <div>
            {/* Header with Create Button */}
            <div className="flex justify-between items-center mb-6">
              <div>
                <h1 className="text-2xl font-bold text-gray-900">Flows</h1>
                <p className="text-gray-600 mt-1">Manage and monitor your automation flows</p>
              </div>
              <Link
                href="/flows/create"
                className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors inline-flex items-center"
              >
                <span className="mr-2">+</span>
                Create Flow
              </Link>
            </div>

            {/* Flows grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {flows.length === 0 ? (
                <div className="col-span-full text-center py-12">
                  <div className="text-6xl mb-4">⚡</div>
                  <h3 className="text-xl font-medium text-gray-900 mb-2">No flows yet</h3>
                  <p className="text-gray-600 mb-6">Get started by creating your first automation flow</p>
                  <Link
                    href="/flows/create"
                    className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors inline-block"
                  >
                    Create Your First Flow
                  </Link>
                </div>
              ) : (
                flows.map((flow) => (
                  <div key={flow.uid} className="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow">
                    <div className="p-6">
                      <div className="flex justify-between items-start mb-4">
                        <h3 className="text-lg font-semibold text-gray-900 truncate">{flow.name}</h3>
                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(flow.status)}`}>
                          {flow.status}
                        </span>
                      </div>
                      
                      {flow.description && (
                        <p className="text-gray-600 text-sm mb-4 line-clamp-2">{flow.description}</p>
                      )}
                      
                      <div className="flex items-center text-sm text-gray-500 mb-4">
                        <span className="mr-4">📊 {flow.nodes.length} nodes</span>
                        <span>🔗 {flow.connections.length} connections</span>
                      </div>
                      
                      <div className="text-xs text-gray-500 mb-4">
                        <div>Created: {formatDate(flow.createdAt)}</div>
                        <div>Updated: {formatDate(flow.updatedAt)}</div>
                      </div>
                      
                      <div className="flex space-x-2">
                        <Link
                          href={`/flows/${flow.uid}`}
                          className="flex-1 bg-blue-600 hover:bg-blue-700 text-white text-center py-2 rounded-md text-sm font-medium transition-colors"
                        >
                          Open
                        </Link>
                        <button
                          onClick={() => handleDeleteFlow(flow.uid || flow.id, flow.name)}
                          className="px-3 py-2 border border-red-300 text-red-700 hover:bg-red-50 rounded-md text-sm font-medium transition-colors"
                        >
                          Delete
                        </button>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        )}
      </div>

      {/* Create Flow Modal */}
      {showCreateModal && (
        <CreateFlowModal
          onCreate={handleCreateFlow}
          onCancel={() => setShowCreateModal(false)}
          loading={createLoading}
        />
      )}
    </AppLayout>
  );
}

// Create Flow Modal Component
function CreateFlowModal({ 
  onCreate, 
  onCancel, 
  loading 
}: { 
  onCreate: (name: string, description: string) => void;
  onCancel: () => void;
  loading: boolean;
}) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (name.trim()) {
      onCreate(name.trim(), description.trim());
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg max-w-md w-full">
        <form onSubmit={handleSubmit}>
          <div className="p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Create New Flow</h3>
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
                disabled={loading}
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
                disabled={loading}
              />
            </div>
          </div>
          
          <div className="px-6 py-4 bg-gray-50 flex justify-end space-x-3 rounded-b-lg">
            <button
              type="button"
              onClick={onCancel}
              className="px-4 py-2 text-gray-700 hover:text-gray-900 font-medium"
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-medium transition-colors disabled:opacity-50"
              disabled={loading || !name.trim()}
            >
              {loading ? 'Creating...' : 'Create Flow'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
