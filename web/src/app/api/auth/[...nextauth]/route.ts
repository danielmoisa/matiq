import NextAuth from "next-auth"
import CredentialsProvider from "next-auth/providers/credentials"
import type { JWT } from "next-auth/jwt"
import type { LoginResponse } from "@/types/auth"

const handler = NextAuth({
  providers: [
    CredentialsProvider({
      name: "credentials",
      credentials: {
        username: { label: "Username", type: "text" },
        password: { label: "Password", type: "password" }
      },
      async authorize(credentials) {
        if (!credentials?.username || !credentials?.password) {
          return null
        }

        try {
          const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              username: credentials.username,
              password: credentials.password,
            }),
          })

          const data: LoginResponse = await response.json()

          if (response.ok && data.success && data.data) {
            return {
              id: data.data.user.id,
              name: data.data.user.username,
              email: data.data.user.email,
              username: data.data.user.username,
              enabled: data.data.user.enabled,
              roles: data.data.user.roles,
              accessToken: data.data.access_token,
              refreshToken: data.data.refresh_token,
              tokenType: data.data.token_type,
              expiresIn: data.data.expires_in,
              expiresAt: new Date(data.data.expires_at).getTime(),
            }
          }
          
          return null
        } catch (error) {
          console.error('Login error:', error)
          return null
        }
      }
    })
  ],
  session: {
    strategy: "jwt",
    maxAge: 24 * 60 * 60, // 24 hours
  },
  jwt: {
    maxAge: 24 * 60 * 60, // 24 hours
  },
  callbacks: {
    async jwt({ token, user, account }): Promise<JWT> {
      // Initial sign in
      if (account && user) {
        return {
          ...token,
          id: user.id,
          username: user.username,
          email: user.email || '',
          enabled: user.enabled,
          roles: user.roles,
          accessToken: user.accessToken,
          refreshToken: user.refreshToken,
          tokenType: user.tokenType,
          expiresIn: user.expiresIn,
          expiresAt: user.expiresAt,
        } as JWT
      }

      // Return previous token if the access token has not expired yet
      if (Date.now() < (token.expiresAt as number)) {
        return token
      }

      // Access token has expired, try to refresh it
      return await refreshAccessToken(token)
    },
    async session({ session, token }) {
      if (token) {
        session.user = {
          id: token.id as string,
          username: token.username as string,
          email: token.email as string,
          enabled: token.enabled as boolean,
          roles: token.roles as string[],
          name: token.username as string,
          image: null,
        }
        session.accessToken = token.accessToken as string
        session.refreshToken = token.refreshToken as string
        session.tokenType = token.tokenType as string
        session.expiresAt = token.expiresAt as number
        session.error = token.error
      }
      return session
    },
  },
  pages: {
    signIn: '/auth/login',
    signOut: '/auth/logout',
  },
  debug: process.env.NODE_ENV === 'development',
})

async function refreshAccessToken(token: JWT): Promise<JWT> {
  try {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        refresh_token: token.refreshToken,
      }),
    })

    const data: LoginResponse = await response.json()

    if (response.ok && data.success && data.data) {
      return {
        ...token,
        accessToken: data.data.access_token,
        refreshToken: data.data.refresh_token,
        expiresAt: new Date(data.data.expires_at).getTime(),
        error: undefined,
      }
    }

    return {
      ...token,
      error: "RefreshAccessTokenError",
    }
  } catch (error) {
    console.error('Refresh token error:', error)
    return {
      ...token,
      error: "RefreshAccessTokenError",
    }
  }
}

export { handler as GET, handler as POST }
