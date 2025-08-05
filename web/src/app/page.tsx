"use client"

import Link from 'next/link'
import { useOptionalAuth } from '@/hooks/useAuth'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import { AppLayout } from '@/components/layout/AppLayout'

export default function Home() {
  const { isAuthenticated, isLoading } = useOptionalAuth()
  const router = useRouter()

  useEffect(() => {
    if (isAuthenticated) {
      router.push('/dashboard')
    }
  }, [isAuthenticated, router])

  if (isLoading) {
    return (
      <AppLayout showAuthButtons={false}>
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-indigo-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading...</p>
          </div>
        </div>
      </AppLayout>
    )
  }

  if (isAuthenticated) {
    return null // Will redirect via useEffect
  }

  return (
    <AppLayout>
      {/* Hero Section */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center">
          <div className="text-8xl mb-8">⚡</div>
          <h2 className="text-4xl font-bold text-gray-900 mb-4">
            Build Powerful Automation Flows
          </h2>
          <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
            Create visual flows that connect your favorite tools and services. 
            Drag and drop components, configure triggers, and automate your processes with ease.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/auth/register"
              className="bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-4 rounded-lg font-medium transition-colors text-lg"
            >
              Get Started Free
            </Link>
            <Link
              href="/auth/login"
              className="border border-gray-300 hover:border-gray-400 text-gray-700 px-8 py-4 rounded-lg font-medium transition-colors text-lg"
            >
              Sign In
            </Link>
          </div>
        </div>

        {/* Features */}
        <div className="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="text-center">
            <div className="text-4xl mb-4">🎯</div>
            <h3 className="text-xl font-semibold text-gray-900 mb-2">Visual Builder</h3>
            <p className="text-gray-600">
              Drag and drop components to build flows visually. No coding required.
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
    </AppLayout>
  )
}
