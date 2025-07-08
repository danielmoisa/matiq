"use client"

import { useRequireAuth } from '@/hooks/useAuth'
import { authService } from '@/services/authService'
import { AppLayout } from '@/components/layout/AppLayout'
import { useState, useEffect } from 'react'

interface UserProfile {
  id: string
  username: string
  email: string
  enabled: boolean
  roles: string[]
}

export default function ProfilePage() {
  const { session, isLoading: authLoading } = useRequireAuth()
  const [profile, setProfile] = useState<UserProfile | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchProfile = async () => {
      if (!session) return

      try {
        const response = await authService.getProfile()
        if (response.success && response.data) {
          setProfile(response.data)
        } else {
          setError(response.message || 'Failed to load profile')
        }
      } catch (err) {
        console.error('Profile fetch error:', err)
        setError('Failed to load profile')
      } finally {
        setIsLoading(false)
      }
    }

    if (session && !authLoading) {
      fetchProfile()
    }
  }, [session, authLoading])

  if (authLoading || isLoading) {
    return (
      <AppLayout>
        <div className="min-h-[80vh] flex items-center justify-center">
          <div className="text-center">
            <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-indigo-600 mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading profile...</p>
          </div>
        </div>
      </AppLayout>
    )
  }

  if (error) {
    return (
      <AppLayout>
        <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          <div className="px-4 py-6 sm:px-0">
            <div className="bg-red-50 border border-red-200 rounded-md p-4">
              <div className="text-red-800">
                <h3 className="text-lg font-medium">Error</h3>
                <p className="mt-2">{error}</p>
              </div>
            </div>
          </div>
        </div>
      </AppLayout>
    )
  }

  return (
    <AppLayout>
      <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg leading-6 font-medium text-gray-900 mb-6">
                User Profile
              </h3>
              
              {profile && (
                <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      User ID
                    </label>
                    <div className="mt-1 text-sm text-gray-900 bg-gray-50 rounded-md px-3 py-2">
                      {profile.id}
                    </div>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Username
                    </label>
                    <div className="mt-1 text-sm text-gray-900 bg-gray-50 rounded-md px-3 py-2">
                      {profile.username}
                    </div>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Email
                    </label>
                    <div className="mt-1 text-sm text-gray-900 bg-gray-50 rounded-md px-3 py-2">
                      {profile.email}
                    </div>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700">
                      Status
                    </label>
                    <div className="mt-1">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        profile.enabled 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {profile.enabled ? 'Enabled' : 'Disabled'}
                      </span>
                    </div>
                  </div>
                  
                  <div className="sm:col-span-2">
                    <label className="block text-sm font-medium text-gray-700">
                      Roles
                    </label>
                    <div className="mt-1">
                      <div className="flex flex-wrap gap-2">
                        {profile.roles.map((role, index) => (
                          <span
                            key={index}
                            className="inline-flex px-3 py-1 text-sm font-medium rounded-full bg-indigo-100 text-indigo-800"
                          >
                            {role}
                          </span>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>
              )}
              
              {session && (
                <div className="mt-8 pt-6 border-t border-gray-200">
                  <h4 className="text-md font-medium text-gray-900 mb-4">
                    Session Information
                  </h4>
                  <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Token Type
                      </label>
                      <div className="mt-1 text-sm text-gray-900 bg-gray-50 rounded-md px-3 py-2">
                        {session.tokenType}
                      </div>
                    </div>
                    
                    <div>
                      <label className="block text-sm font-medium text-gray-700">
                        Expires At
                      </label>
                      <div className="mt-1 text-sm text-gray-900 bg-gray-50 rounded-md px-3 py-2">
                        {new Date(session.expiresAt).toLocaleString()}
                      </div>
                    </div>
                  </div>
                  
                  {session.error && (
                    <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md">
                      <div className="text-yellow-800">
                        <strong>Session Warning:</strong> {session.error}
                      </div>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </AppLayout>
  )
}
