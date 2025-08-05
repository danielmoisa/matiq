"use client"

import Link from 'next/link'
import { useSession } from 'next-auth/react'
import { useRouter } from 'next/navigation'
import { ReactNode } from 'react'

interface AppLayoutProps {
  children: ReactNode
  showAuthButtons?: boolean
}

export function AppLayout({ children, showAuthButtons = true }: AppLayoutProps) {
  const { data: session } = useSession()
  const router = useRouter()

  const handleLogout = () => {
    router.push('/auth/logout')
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      {/* Header */}
      <div className="bg-white shadow-sm border-b">
        <div className="mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div>
              <Link href={session ? "/dashboard" : "/"}>
                <h1 className="text-2xl font-bold text-gray-900 cursor-pointer hover:text-gray-700 transition-colors">
                  Matiq
                </h1>
              </Link>
            </div>

            {/* Navigation and Auth Buttons */}
            <div className="flex items-center space-x-4">
              {session ? (
                <>
                  {/* Authenticated Navigation */}
                  <nav className="hidden md:flex space-x-6">
                    <Link
                      href="/dashboard"
                      className="text-gray-700 hover:text-gray-900 font-medium transition-colors"
                    >
                      Dashboard
                    </Link>
                    <Link
                      href="/flows"
                      className="text-gray-700 hover:text-gray-900 font-medium transition-colors"
                    >
                      Flows
                    </Link>
                    <Link
                      href="/profile"
                      className="text-gray-700 hover:text-gray-900 font-medium transition-colors"
                    >
                      Profile
                    </Link>
                  </nav>

                  {/* User Info and Logout */}
                  <div className="flex items-center space-x-4 border-l border-gray-200 pl-4">
                    <span className="text-gray-700 hidden sm:block">
                    {session.user.username}
                    </span>
                    <button
                      onClick={handleLogout}
                      className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
                    >
                      Sign Out
                    </button>
                  </div>
                </>
              ) : (
                showAuthButtons && (
                  <>
                    {/* Unauthenticated Buttons */}
                    <Link
                      href="/auth/login"
                      className="border border-indigo-600 text-indigo-600 hover:bg-indigo-600 hover:text-white px-6 py-2 rounded-lg font-medium transition-colors"
                    >
                      Sign In
                    </Link>
                    <Link
                      href="/auth/register"
                      className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
                    >
                      Get Started
                    </Link>
                  </>
                )
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <main className="flex-1">
        {children}
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <div className="text-gray-600 text-sm">
              © 2025 Matiq. All rights reserved.
            </div>
            <div className="flex space-x-6 mt-4 md:mt-0">
              <Link href="/docs" className="text-gray-600 hover:text-gray-900 text-sm">
                Documentation
              </Link>
              <Link href="/support" className="text-gray-600 hover:text-gray-900 text-sm">
                Support
              </Link>
              <Link href="/privacy" className="text-gray-600 hover:text-gray-900 text-sm">
                Privacy
              </Link>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
