import { withAuth } from "next-auth/middleware"

export default withAuth(
  function middleware() {
    // Can add additional middleware logic here
    // For example, role-based access control
  },
  {
    callbacks: {
      authorized: ({ token, req }) => {
        // Check if user is authenticated
        if (!token) return false

        // Check for specific routes that require specific roles
        const { pathname } = req.nextUrl

        // Admin routes require admin role
        if (pathname.startsWith('/admin')) {
          return token.roles?.includes('admin') || false
        }

        // Dashboard and other protected routes just need authentication
        if (pathname.startsWith('/dashboard') || pathname.startsWith('/workflows')) {
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
    '/workflows/:path*',
    '/admin/:path*',
    '/profile/:path*',
  ]
}
