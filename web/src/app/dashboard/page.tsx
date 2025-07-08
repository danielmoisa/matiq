"use client"

import { useSession } from 'next-auth/react'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import Link from 'next/link'
import { AppLayout } from '@/components/layout/AppLayout'

export default function DashboardPage() {
  const { data: session, status } = useSession()
  const router = useRouter()

  useEffect(() => {
    if (status === 'unauthenticated') {
      router.push('/auth/login')
    }
  }, [status, router])

  if (status === 'loading') {
    return (
      <AppLayout showAuthButtons={false}>
        <div className="min-h-[80vh] flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-indigo-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading...</p>
          </div>
        </div>
      </AppLayout>
    )
  }

  if (!session) {
    return null // Will redirect via useEffect
  }

  return (
    <AppLayout>
      {/* Main content */}
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="border-4 border-dashed border-gray-200 rounded-lg p-8">
            <div className="text-center">
              <h2 className="text-3xl font-extrabold text-gray-900 mb-4">
                Welcome to your Dashboard
              </h2>
              <p className="text-gray-600 mb-8">
                You are successfully authenticated with Keycloak through NextAuth.js
              </p>
              
              {/* User info card */}
              <div className="bg-white shadow rounded-lg p-6 max-w-md mx-auto">
                <h3 className="text-lg font-medium text-gray-900 mb-4">User Information</h3>
                <div className="space-y-2 text-left">
                  <p><strong>ID:</strong> {session.user.id}</p>
                  <p><strong>Username:</strong> {session.user.username}</p>
                  <p><strong>Email:</strong> {session.user.email}</p>
                  <p><strong>Enabled:</strong> {session.user.enabled ? 'Yes' : 'No'}</p>
                  <p><strong>Roles:</strong> {session.user.roles.join(', ')}</p>
                </div>
              </div>

              {/* Session info */}
              <div className="bg-white shadow rounded-lg p-6 max-w-md mx-auto mt-6">
                <h3 className="text-lg font-medium text-gray-900 mb-4">Session Information</h3>
                <div className="space-y-2 text-left">
                  <p><strong>Token Type:</strong> {session.tokenType}</p>
                  <p><strong>Expires At:</strong> {new Date(session.expiresAt).toLocaleString()}</p>
                  {session.error && (
                    <p className="text-red-600"><strong>Error:</strong> {session.error}</p>
                  )}
                </div>
              </div>

              <div className="mt-8">
                <Link
                  href="/workflows"
                  className="bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded"
                >
                  Go to Workflows
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AppLayout>
  )
}
