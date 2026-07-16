import { apiDownload, apiRequest } from '@/lib/api-client'
import type {
  CreateProjectPayload,
  DashboardSummary,
  FeaturePayload,
  FeatureRecord,
  IssueLicensePayload,
  LicenseRecord,
  PlanPayload,
  PlanRecord,
  ProjectRecord,
  UpdateProjectPayload,
} from '@/types/api'

export function fetchDashboardSummary() {
  return apiRequest<DashboardSummary>('/dashboard/summary')
}

export function fetchProjects() {
  return apiRequest<ProjectRecord[]>('/projects')
}

export function fetchProject(id: number) {
  return apiRequest<ProjectRecord>(`/projects/${id}`)
}

export function createProject(data: CreateProjectPayload) {
  return apiRequest<ProjectRecord>('/projects', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export function updateProject(id: number, data: UpdateProjectPayload) {
  return apiRequest<ProjectRecord>(`/projects/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export function deleteProject(id: number) {
  return apiRequest<{ deleted: boolean }>(`/projects/${id}`, {
    method: 'DELETE',
  })
}

export function downloadProjectPublicKey(id: number) {
  return apiDownload(`/projects/${id}/public-key`)
}

export function fetchFeatures() {
  return apiRequest<FeatureRecord[]>('/features')
}

export function createFeature(data: FeaturePayload) {
  return apiRequest<FeatureRecord>('/features', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export function updateFeature(id: number, data: FeaturePayload) {
  return apiRequest<FeatureRecord>(`/features/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export function deleteFeature(id: number) {
  return apiRequest<{ deleted: boolean }>(`/features/${id}`, {
    method: 'DELETE',
  })
}

export function fetchPlans() {
  return apiRequest<PlanRecord[]>('/plans')
}

export function createPlan(data: PlanPayload) {
  return apiRequest<PlanRecord>('/plans', {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export function updatePlan(id: number, data: PlanPayload) {
  return apiRequest<PlanRecord>(`/plans/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  })
}

export function deletePlan(id: number) {
  return apiRequest<{ deleted: boolean }>(`/plans/${id}`, {
    method: 'DELETE',
  })
}

export function fetchProjectLicenses(projectId: number) {
  return apiRequest<LicenseRecord[]>(`/projects/${projectId}/licenses`)
}

export function issueLicense(projectId: number, data: IssueLicensePayload) {
  return apiRequest<LicenseRecord>(`/projects/${projectId}/licenses`, {
    method: 'POST',
    body: JSON.stringify(data),
  })
}

export function fetchLicense(id: number) {
  return apiRequest<LicenseRecord>(`/licenses/${id}`)
}

export function downloadLicense(id: number) {
  return apiDownload(`/licenses/${id}/download`)
}
