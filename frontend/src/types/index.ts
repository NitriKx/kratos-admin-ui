// Identity credential types
export interface IdentityCredential {
  type: string
  identifiers?: string[]
  created_at?: string
  updated_at?: string
}

// Identity types
export interface Identity {
  id: string
  schema_id: string
  schema_url?: string
  state: IdentityState
  state_changed_at?: string
  traits: Record<string, unknown>
  verifiable_addresses?: VerifiableAddress[]
  recovery_addresses?: RecoveryAddress[]
  metadata_public?: Record<string, unknown>
  metadata_admin?: Record<string, unknown>
  credentials?: Record<string, IdentityCredential>
  created_at: string
  updated_at: string
}

export type IdentityState = 'active' | 'inactive'

export interface VerifiableAddress {
  id: string
  value: string
  verified: boolean
  via: string
  status: string
  verified_at?: string
  created_at: string
  updated_at: string
}

export interface RecoveryAddress {
  id: string
  value: string
  via: string
  created_at: string
  updated_at: string
}

// Session types
export interface Session {
  id: string
  active: boolean
  expires_at: string
  authenticated_at: string
  authenticator_assurance_level: string
  authentication_methods?: AuthenticationMethod[]
  issued_at: string
  identity: Identity
  devices?: SessionDevice[]
}

export interface AuthenticationMethod {
  method: string
  aal: string
  completed_at: string
}

export interface SessionDevice {
  id: string
  ip_address?: string
  user_agent?: string
  location?: string
}

// Schema types
export interface IdentitySchema {
  id: string
  schema: Record<string, unknown>
}

// Stats types
export interface Stats {
  active_identities: number
  active_sessions: number
}

// API response types
export interface PaginatedResponse<T> {
  data: T[]
  page: number
  per_page: number
  total?: number
}

export interface LoginResponse {
  token: string
  expires_at: number
}

// UI helper types
export interface TableColumn {
  key: string
  label: string
  sortable?: boolean
  width?: string
}




