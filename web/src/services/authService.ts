import { getSession, signOut } from 'next-auth/react'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL

export class AuthenticatedApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public response?: Record<string, unknown>
  ) {
    super(message)
    this.name = 'AuthenticatedApiError'
  }
}

// Helper function to make authenticated API calls
export async function authenticatedFetch(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  const session = await getSession()

  if (!session?.accessToken || session?.error) {
    // Handle session errors or missing tokens
    if (session?.error === 'RefreshAccessTokenError') {
      // Token refresh failed, redirect to login
      await signOut({ redirect: false })
      window.location.href = '/auth/login?error=SessionExpired'
      throw new AuthenticatedApiError('Session expired', 401)
    }
    throw new AuthenticatedApiError('No access token available', 401)
  }

  // Check if token is expired
  if (session.expiresAt && Date.now() > session.expiresAt) {
    await signOut({ redirect: false })
    window.location.href = '/auth/login?error=SessionExpired'
    throw new AuthenticatedApiError('Session expired', 401)
  }

  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `${session.tokenType} ${session.accessToken}`,
    ...options.headers,
  }

  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers,
  })

  // Handle 401 responses (unauthorized)
  if (response.status === 401) {
    await signOut({ redirect: false })
    window.location.href = '/auth/login?error=SessionExpired'
    throw new AuthenticatedApiError('Session expired', 401)
  }

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}))
    throw new AuthenticatedApiError(
      errorData.message || `HTTP ${response.status}`,
      response.status,
      errorData
    )
  }

  return response
}

// API service functions
export const authService = {
  // Get user profile
  async getProfile() {
    const response = await authenticatedFetch('/api/v1/auth/profile')
    return response.json()
  },

  // Validate current token
  async validateToken() {
    const response = await authenticatedFetch('/api/v1/auth/validate')
    return response.json()
  },

  // Generic authenticated GET request
  async get(url: string) {
    const response = await authenticatedFetch(url)
    return response.json()
  },

  // Generic authenticated POST request
  async post(url: string, data: Record<string, unknown>) {
    const response = await authenticatedFetch(url, {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  // Generic authenticated PUT request
  async put(url: string, data: Record<string, unknown>) {
    const response = await authenticatedFetch(url, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  // Generic authenticated DELETE request
  async delete(url: string) {
    const response = await authenticatedFetch(url, {
      method: 'DELETE',
    })
    return response.json()
  },
}

// Workflow service functions
export const workflowService = {
  // Create a new workflow
  async createWorkflow(teamId: string, workflowData: Record<string, unknown>) {
    return authService.post(`/api/v1/teams/${teamId}/workflow`, workflowData)
  },

  // Get a specific workflow
  async getWorkflow(teamId: string, workflowId: string) {
    return authService.get(`/api/v1/teams/${teamId}/workflow/${workflowId}`)
  },

  // Update a workflow
  async updateWorkflow(teamId: string, workflowId: string, workflowData: Record<string, unknown>) {
    return authService.put(`/api/v1/teams/${teamId}/workflow/${workflowId}`, workflowData)
  },

  // Delete a workflow
  async deleteWorkflow(teamId: string, workflowId: string) {
    return authService.delete(`/api/v1/teams/${teamId}/workflow/${workflowId}`)
  },
}
