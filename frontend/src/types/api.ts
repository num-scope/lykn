export type ApiResponse<T> = {
  code: number
  message: string
  data?: T
}

export type UserAccount = {
  id: number
  username: string
}

export type LoginResponse = {
  access_token: string
  token_type: string
  expires_at: string
  user: UserAccount
}

export type DashboardSummary = {
  project_count: number
  license_count: number
  active_license_count: number
  expired_license_count: number
}

export type ProjectRecord = {
  id: number
  name: string
  description: string
  public_key?: string
  key_bits: number
  created_at: string
  updated_at: string
}

export type CreateProjectPayload = {
  name: string
  description: string
  key_bits: number
}

export type UpdateProjectPayload = {
  name: string
  description: string
}

export type FeatureRecord = {
  id: number
  code: string
  name: string
  description: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export type FeaturePayload = {
  code: string
  name: string
  description: string
  enabled: boolean
}

export type PlanRecord = {
  id: number
  code: string
  name: string
  description: string
  features: FeatureRecord[]
  max_users: number
  max_devices: number
  enabled: boolean
  created_at: string
  updated_at: string
}

export type PlanPayload = {
  code: string
  name: string
  description: string
  feature_ids: number[]
  max_users: number
  max_devices: number
  enabled: boolean
}

export type LicenseLimits = {
  max_users: number
  max_devices: number
}

export type LicenseRecord = {
  id: number
  uuid: string
  project_id: number
  subject_name: string
  subject_email: string
  subject_org: string
  plan_id?: number
  plan_name: string
  plan: string
  not_before: string
  not_after: string
  features: string[]
  limits: LicenseLimits
  metadata: Record<string, unknown>
  created_at: string
}

export type LicenseHardwarePayload = {
  hostname?: string
  cpu_id?: string
  disk_serial?: string
  mac_addresses?: string[]
  ip_addresses?: string[]
}

export type IssueLicensePayload = {
  subject: {
    name: string
    email: string
    organization: string
  }
  plan_id: number
  not_before: string
  not_after: string
  hardware: LicenseHardwarePayload
}

export type DownloadFile = {
  blob: Blob
  filename: string
}

export type LicenseState = 'active' | 'upcoming' | 'expired' | 'unknown'
