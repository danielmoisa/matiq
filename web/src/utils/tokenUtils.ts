// Utility functions for testing token expiration
export function isTokenExpired(expiresAt: number): boolean {
  return Date.now() > expiresAt
}

export function getTimeUntilExpiry(expiresAt: number): number {
  return Math.max(0, expiresAt - Date.now())
}

export function formatTimeUntilExpiry(expiresAt: number): string {
  const timeLeft = getTimeUntilExpiry(expiresAt)
  
  if (timeLeft <= 0) {
    return 'Expired'
  }
  
  const minutes = Math.floor(timeLeft / (1000 * 60))
  const seconds = Math.floor((timeLeft % (1000 * 60)) / 1000)
  
  if (minutes > 0) {
    return `${minutes}m ${seconds}s`
  } else {
    return `${seconds}s`
  }
}
