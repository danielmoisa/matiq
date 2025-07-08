import { DefaultSession, DefaultUser } from "next-auth"
import { DefaultJWT } from "next-auth/jwt"

declare module "next-auth" {
  interface Session {
    user: {
      id: string
      username: string
      email: string
      enabled: boolean
      roles: string[]
    } & DefaultSession["user"]
    accessToken: string
    refreshToken: string
    tokenType: string
    expiresAt: number
    error?: string
  }

  interface User extends DefaultUser {
    username: string
    enabled: boolean
    roles: string[]
    accessToken: string
    refreshToken: string
    tokenType: string
    expiresIn: number
    expiresAt: number
  }
}

declare module "next-auth/jwt" {
  interface JWT extends DefaultJWT {
    id: string
    username: string
    email: string
    enabled: boolean
    roles: string[]
    accessToken: string
    refreshToken: string
    tokenType: string
    expiresIn: number
    expiresAt: number
    error?: string
  }
}

export interface LoginResponse {
  success: boolean
  message: string
  data?: {
    user: {
      id: string
      username: string
      email: string
      enabled: boolean
      roles: string[]
    }
    access_token: string
    refresh_token: string
    token_type: string
    expires_in: number
    expires_at: string
  }
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface LogoutRequest {
  refresh_token: string
}

export interface ApiError {
  success: false
  message: string
}
