'use client';

import { useParams } from 'next/navigation';
import { useFlow } from '@/hooks/useFlow';
import FlowBuilder from '@/components/FlowBuilder';
import Link from 'next/link';
import { AppLayout } from '@/components/layout/AppLayout';

export default function FlowPage() {
  const params = useParams();
  const flowId = params.id as string;
  const { flow, loading, error } = useFlow(flowId);

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
              href="/flows"
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors inline-block"
            >
              Back to Flows
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
            <h3 className="text-xl font-medium text-gray-900 mb-2">Flow not found</h3>
            <p className="text-gray-600 mb-6">The flow you&apos;re looking for doesn&apos;t exist or has been deleted.</p>
            <Link
              href="/flows"
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors inline-block"
            >
              Back to Flows
            </Link>
          </div>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className="h-[calc(100vh-4rem)]">
        <FlowBuilder flowId={flowId} />
      </div>
    </AppLayout>
  );
}
