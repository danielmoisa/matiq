'use client';

import { useState } from 'react';

export default function Header() {
  const [workflowName, setWorkflowName] = useState('Untitled Workflow');

  return (
    <header className="bg-white border-b border-gray-200 px-6 py-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">WF</span>
            </div>
            <h1 className="text-xl font-semibold text-gray-900">Builder</h1>
          </div>
          <div className="h-6 w-px bg-gray-300" />
          <input
            type="text"
            value={workflowName}
            onChange={(e) => setWorkflowName(e.target.value)}
            className="bg-transparent text-lg font-medium text-gray-900 border-none outline-none focus:bg-gray-50 px-2 py-1 rounded"
          />
        </div>
        
        <div className="flex items-center space-x-3">
          <button className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors">
            Save Draft
          </button>
          <button className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors">
            Publish
          </button>
          <button className="p-2 text-gray-400 hover:text-gray-600 transition-colors">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
            </svg>
          </button>
        </div>
      </div>
    </header>
  );
}
