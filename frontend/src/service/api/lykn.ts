import { request } from '../request';

function getFilename(disposition = '') {
  const utf8Match = disposition.match(/filename\*=UTF-8''([^;]+)/i);
  if (utf8Match?.[1]) return decodeURIComponent(utf8Match[1]);

  const normalMatch = disposition.match(/filename="?([^";]+)"?/i);
  return normalMatch?.[1] || 'download';
}

async function download(url: string): Promise<Api.Lykn.DownloadFile> {
  const { data, response, error } = await request<Blob, 'blob'>({ url, responseType: 'blob' });

  if (error || !data) {
    throw error;
  }

  return {
    blob: data,
    filename: getFilename(response.headers['content-disposition'])
  };
}

export function fetchDashboardSummary() {
  return request<Api.Lykn.DashboardSummary>({ url: '/dashboard/summary' });
}

export function fetchProjects() {
  return request<Api.Lykn.ProjectRecord[]>({ url: '/projects' });
}

export function fetchProject(id: number) {
  return request<Api.Lykn.ProjectRecord>({ url: `/projects/${id}` });
}

export function createProject(data: Api.Lykn.CreateProjectPayload) {
  return request<Api.Lykn.ProjectRecord>({ url: '/projects', method: 'post', data });
}

export function updateProject(id: number, data: Api.Lykn.UpdateProjectPayload) {
  return request<Api.Lykn.ProjectRecord>({ url: `/projects/${id}`, method: 'put', data });
}

export function deleteProject(id: number) {
  return request<{ deleted: boolean }>({ url: `/projects/${id}`, method: 'delete' });
}

export function downloadProjectPublicKey(id: number) {
  return download(`/projects/${id}/public-key`);
}

export function fetchFeatures() {
  return request<Api.Lykn.FeatureRecord[]>({ url: '/features' });
}

export function createFeature(data: Api.Lykn.FeaturePayload) {
  return request<Api.Lykn.FeatureRecord>({ url: '/features', method: 'post', data });
}

export function updateFeature(id: number, data: Api.Lykn.FeaturePayload) {
  return request<Api.Lykn.FeatureRecord>({ url: `/features/${id}`, method: 'put', data });
}

export function deleteFeature(id: number) {
  return request<{ deleted: boolean }>({ url: `/features/${id}`, method: 'delete' });
}

export function fetchPlans() {
  return request<Api.Lykn.PlanRecord[]>({ url: '/plans' });
}

export function createPlan(data: Api.Lykn.PlanPayload) {
  return request<Api.Lykn.PlanRecord>({ url: '/plans', method: 'post', data });
}

export function updatePlan(id: number, data: Api.Lykn.PlanPayload) {
  return request<Api.Lykn.PlanRecord>({ url: `/plans/${id}`, method: 'put', data });
}

export function deletePlan(id: number) {
  return request<{ deleted: boolean }>({ url: `/plans/${id}`, method: 'delete' });
}

export function fetchProjectLicenses(projectId: number) {
  return request<Api.Lykn.LicenseRecord[]>({ url: `/projects/${projectId}/licenses` });
}

export function issueLicense(projectId: number, data: Api.Lykn.IssueLicensePayload) {
  return request<Api.Lykn.LicenseRecord>({ url: `/projects/${projectId}/licenses`, method: 'post', data });
}

export function fetchLicense(id: number) {
  return request<Api.Lykn.LicenseRecord>({ url: `/licenses/${id}` });
}

export function downloadLicense(id: number) {
  return download(`/licenses/${id}/download`);
}
