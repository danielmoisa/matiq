'use client';

import Link from 'next/link';

interface FlowHeaderProps {
  flowName?: string;
  onSave?: () => void;
  onPublish?: () => void;
  saving?: boolean;
}

export default function FlowHeader({ 
  flowName,
  onSave,
  onPublish,
  saving = false
}: FlowHeaderProps) {

  return (
    <div className="bg-white border-b border-gray-200 px-6 py-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2 text-sm text-gray-500">
            <Link href="/" className="hover:text-gray-700">Home</Link>
            <span>›</span>
            <Link href="/flows" className="hover:text-gray-700">Flows</Link>
            <span>›</span>
            <span>Edit</span>
          </div>
          <div className="h-6 w-px bg-gray-300" />
          <h2 className="text-lg font-medium text-gray-900">
            {flowName || 'Untitled Flow'}
          </h2>
        </div>
        
        <div className="flex items-center space-x-3">
          <button 
            onClick={onSave}
            disabled={saving}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors disabled:opacity-50"
          >
            {saving ? 'Saving...' : 'Save'}
          </button>
          <button 
            onClick={onPublish}
            className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
          >
            Publish
          </button>
        </div>
      </div>
    </div>
  );
}
