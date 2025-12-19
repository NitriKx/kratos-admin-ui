import axios, { type AxiosInstance, type AxiosError } from 'axios'
import type { Identity, Session, IdentitySchema, Stats, PaginatedResponse, LoginResponse } from '@/types'

// Runtime config from window.__RUNTIME_CONFIG__ (injected by config.js)
// Falls back to VITE_API_URL for development, then to empty string (relative URLs)
declare global {
  interface Window {
    __RUNTIME_CONFIG__?: {
      apiUrl?: string
    }
  }
}

const API_URL = window.__RUNTIME_CONFIG__?.apiUrl ?? import.meta.env.VITE_API_URL ?? ''

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    // Add auth interceptor
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('auth_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    })

    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('auth_token')
          window.location.href = '/login'
        }
        return Promise.reject(error)
      }
    )
  }

  // Auth
  async login(password: string): Promise<LoginResponse> {
    const response = await this.client.post<LoginResponse>('/api/auth/login', { password })
    return response.data
  }

  // Identities
  async getIdentities(page = 1, perPage = 20): Promise<PaginatedResponse<Identity>> {
    const response = await this.client.get<PaginatedResponse<Identity>>('/api/identities', {
      params: { page, per_page: perPage }
    })
    return response.data
  }

  async getIdentity(id: string): Promise<Identity> {
    const response = await this.client.get<Identity>(`/api/identities/${id}`)
    return response.data
  }

  async createIdentity(data: { schema_id: string; traits: Record<string, unknown>; state?: string }): Promise<Identity> {
    const response = await this.client.post<Identity>('/api/identities', data)
    return response.data
  }

  async updateIdentity(id: string, data: { schema_id: string; traits: Record<string, unknown>; state?: string }): Promise<Identity> {
    const response = await this.client.put<Identity>(`/api/identities/${id}`, data)
    return response.data
  }

  async deleteIdentity(id: string): Promise<void> {
    await this.client.delete(`/api/identities/${id}`)
  }

  async getIdentitySessions(id: string): Promise<{ data: Session[] }> {
    const response = await this.client.get<{ data: Session[] }>(`/api/identities/${id}/sessions`)
    return response.data
  }

  // Sessions
  async getSessions(page = 1, perPage = 20): Promise<PaginatedResponse<Session>> {
    const response = await this.client.get<PaginatedResponse<Session>>('/api/sessions', {
      params: { page, per_page: perPage }
    })
    return response.data
  }

  async revokeSession(id: string): Promise<void> {
    await this.client.delete(`/api/sessions/${id}`)
  }

  // Schemas
  async getSchemas(): Promise<{ data: IdentitySchema[] }> {
    const response = await this.client.get<{ data: IdentitySchema[] }>('/api/schemas')
    return response.data
  }

  // Stats
  async getStats(): Promise<Stats> {
    const response = await this.client.get<Stats>('/api/stats')
    return response.data
  }
}

export const api = new ApiClient()




