"use client"

import { useSession, signOut } from 'next-auth/react'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'

export function useAuth(requireAuth = true) {
  const { data: session, status } = useSession()
  const router = useRouter()

  useEffect(() => {
    // Handle session errors (like token refresh failures)
    if (session?.error === 'RefreshAccessTokenError') {
      signOut({ 
        redirect: false 
      }).then(() => {
        router.push('/auth/login?error=SessionExpired')
      })
      return
    }

    // Handle token expiration
    if (session?.expiresAt && Date.now() > session.expiresAt) {
      signOut({ 
        redirect: false 
      }).then(() => {
        router.push('/auth/login?error=SessionExpired')
      })
      return
    }

    if (requireAuth && status === 'unauthenticated') {
      router.push('/auth/login')
    }
  }, [requireAuth, status, router, session])

  return {
    session,
    status,
    isLoading: status === 'loading',
    isAuthenticated: status === 'authenticated' && !session?.error,
    user: session?.user,
    accessToken: session?.accessToken,
    refreshToken: session?.refreshToken,
    error: session?.error,
  }
}

export function useRequireAuth() {
  return useAuth(true)
}

export function useOptionalAuth() {
  return useAuth(false)
}

// Hook to check if user has specific role
export function useHasRole(requiredRole: string) {
  const { user } = useAuth()
  return user?.roles?.includes(requiredRole) || false
}

// Hook for admin-only access
export function useRequireAdmin() {
  const { session, status } = useAuth(true)
  const router = useRouter()
  const isAdmin = session?.user?.roles?.includes('admin') || false

  useEffect(() => {
    if (status === 'authenticated' && !isAdmin) {
      router.push('/dashboard') // Redirect non-admin users
    }
  }, [status, isAdmin, router])

  return {
    session,
    status,
    isAdmin,
    isLoading: status === 'loading',
  }
}
