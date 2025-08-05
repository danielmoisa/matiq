'use client';

import { useParams } from 'next/navigation';
import { useWorkflow } from '@/hooks/useWorkflow';
import WorkflowBuilder from '@/components/WorkflowBuilder';
import Link from 'next/link';
import { AppLayout } from '@/components/layout/AppLayout';

export default function WorkflowPage() {
  const params = useParams();
  const workflowId = params.id as string;
  const { flow, loading, error } = useWorkflow(workflowId);

  if (loading) {
    return (
      <AppLayout>
        <div className="min-h-[calc(100vh-4rem)] bg-gray-50 flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600">Loading flow...</p>
          </div>
        </div>
      </AppLayout>
    );
  }

  if (error) {
    return (
      <AppLayout>
        <div className="min-h-[calc(100vh-4rem)] bg-gray-50 flex items-center justify-center">
          <div className="text-center max-w-md">
            <div className="text-6xl mb-4">⚠️</div>
            <h3 className="text-xl font-medium text-gray-900 mb-2">Failed to load flow</h3>
            <p className="text-gray-600 mb-6">{error}</p>
            <Link
              href="/workflows"
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors inline-block"
            >
              Back to Workflows
            </Link>
          </div>
        </div>
      </AppLayout>
    );
  }

  if (!flow) {
    return (
      <AppLayout>
        <div className="min-h-[calc(100vh-4rem)] bg-gray-50 flex items-center justify-center">
          <div className="text-center max-w-md">
            <div className="text-6xl mb-4">❓</div>
            <h3 className="text-xl font-medium text-gray-900 mb-2">Workflow not found</h3>
            <p className="text-gray-600 mb-6">The flow you&apos;re looking for doesn&apos;t exist or has been deleted.</p>
            <Link
              href="/workflows"
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors inline-block"
            >
              Back to Workflows
            </Link>
          </div>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className="h-[calc(100vh-4rem)]">
        <WorkflowBuilder workflowId={workflowId} />
      </div>
    </AppLayout>
  );
}
