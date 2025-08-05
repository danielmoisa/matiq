'use client';

import Link from 'next/link';

interface HeaderProps {
  flowName?: string;
  onSave?: () => void;
  onPublish?: () => void;
  saving?: boolean;
}

export default function Header({ 
  flowName,
  onSave,
  onPublish,
  saving = false
}: HeaderProps) {

  return (
    <header className="bg-white border-b border-gray-200 px-6 py-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <Link href="/" className="flex items-center space-x-2 hover:opacity-80 transition-opacity">
              <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">WF</span>
              </div>
              <h1 className="text-xl font-semibold text-gray-900">Builder</h1>
            </Link>
          </div>
          <div className="h-6 w-px bg-gray-300" />
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
            {saving ? 'Saving...' : 'Save Draft'}
          </button>
          <button 
            onClick={onPublish}
            disabled={saving}
            className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors disabled:opacity-50"
          >
            Publish
          </button>
          <Link
            href="/flows"
            className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
            title="Back to flows"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </Link>
        </div>
      </div>
    </header>
  );
}
