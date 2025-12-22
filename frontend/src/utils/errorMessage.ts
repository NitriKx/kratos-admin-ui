import type { AxiosError } from 'axios'

interface ApiErrorResponse {
  error?: string
  details?: string
  message?: string
}

export function getErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    const axiosError = error as AxiosError<ApiErrorResponse>
    if (axiosError.response?.data) {
      const data = axiosError.response.data
      return data.details || data.error || data.message || error.message
    }
    return error.message
  }
  return 'An unexpected error occurred'
}



