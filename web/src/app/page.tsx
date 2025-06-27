import Link from 'next/link';

export default function Home() {
  return (
    <main className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Workflow Builder</h1>
              <p className="text-gray-600 mt-1">Design, build, and manage your automation workflows</p>
            </div>
            <Link
              href="/workflows"
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium transition-colors"
            >
              View Workflows
            </Link>
          </div>
        </div>
      </div>

      {/* Hero Section */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center">
          <div className="text-8xl mb-8">⚡</div>
          <h2 className="text-4xl font-bold text-gray-900 mb-4">
            Build Powerful Automation Workflows
          </h2>
          <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
            Create visual workflows that connect your favorite tools and services. 
            Drag and drop components, configure triggers, and automate your processes with ease.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/workflows"
              className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-4 rounded-lg font-medium transition-colors text-lg"
            >
              Get Started
            </Link>
            <Link
              href="/workflows"
              className="border border-gray-300 hover:border-gray-400 text-gray-700 px-8 py-4 rounded-lg font-medium transition-colors text-lg"
            >
              Browse Workflows
            </Link>
          </div>
        </div>

        {/* Features */}
        <div className="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="text-center">
            <div className="text-4xl mb-4">🎯</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">Visual Builder</h3>
            <p className="text-gray-600">
              Drag and drop components to build workflows visually. No coding required.
            </p>
          </div>
          
          <div className="text-center">
            <div className="text-4xl mb-4">🔗</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">Connect Everything</h3>
            <p className="text-gray-600">
              Integrate with databases, APIs, webhooks, and more. Connect your entire tech stack.
            </p>
          </div>
          
          <div className="text-center">
            <div className="text-4xl mb-4">⚙️</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">Powerful Automation</h3>
            <p className="text-gray-600">
              Schedule tasks, trigger on events, and create complex conditional logic with ease.
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
