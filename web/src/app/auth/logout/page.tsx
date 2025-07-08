"use client"

import { useEffect } from 'react'
import { signOut, useSession } from 'next-auth/react'
import { useRouter } from 'next/navigation'

export default function LogoutPage() {
  const { data: session } = useSession()
  const router = useRouter()

  useEffect(() => {
    const handleLogout = async () => {
      if (session?.refreshToken) {
        try {
          // Call the backend logout endpoint to invalidate tokens
          await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/logout`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              refresh_token: session.refreshToken,
            }),
          })
        } catch (error) {
          console.error('Backend logout error:', error)
          // Continue with NextAuth logout even if backend logout fails
        }
      }

      // Sign out from NextAuth
      await signOut({
        redirect: false,
      })

      // Redirect to login page
      router.push('/auth/login')
    }

    handleLogout()
  }, [session, router])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h2 className="mt-6 text-3xl font-extrabold text-gray-900">
            Signing out...
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            Please wait while we sign you out safely.
          </p>
        </div>
      </div>
    </div>
  )
}
