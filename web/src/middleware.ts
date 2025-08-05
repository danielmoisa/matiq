import { withAuth } from "next-auth/middleware"
import { NextResponse } from "next/server"

export default withAuth(
  function middleware(req) {
    const token = req.nextauth.token
    
    // If there's a token error (like refresh failure), redirect to login
    if (token?.error === "RefreshAccessTokenError") {
      const loginUrl = new URL('/auth/login', req.url)
      loginUrl.searchParams.set('error', 'SessionExpired')
      return NextResponse.redirect(loginUrl)
    }

    // Check if token has expired
    const expiresAt = token?.expiresAt as number | undefined
    if (expiresAt && Date.now() > expiresAt) {
      const loginUrl = new URL('/auth/login', req.url)
      loginUrl.searchParams.set('error', 'SessionExpired')
      return NextResponse.redirect(loginUrl)
    }

    // Can add additional middleware logic here
    // For example, role-based access control
  },
  {
    callbacks: {
      authorized: ({ token, req }) => {
        // Check if user is authenticated
        if (!token) return false

        // Check for token errors (like refresh failures)
        if (token.error === "RefreshAccessTokenError") {
          return false
        }

        // Check if token has expired
        const expiresAt = token.expiresAt as number | undefined
        if (expiresAt && Date.now() > expiresAt) {
          return false
        }

        // Check for specific routes that require specific roles
        const { pathname } = req.nextUrl

        // Admin routes require admin role
        if (pathname.startsWith('/admin')) {
          return token.roles?.includes('admin') || false
        }

        // Dashboard and other protected routes just need authentication
        if (pathname.startsWith('/dashboard') || pathname.startsWith('/flows')) {
          return true
        }

        return true
      },
    },
  }
)

export const config = {
  matcher: [
    // Protect these routes
    '/dashboard/:path*',
    '/flows/:path*',
    '/admin/:path*',
    '/profile/:path*',
  ]
}
