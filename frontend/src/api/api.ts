export interface ApiEnvelope<T> {
  code: number;
  message: string;
  data?: T;
}

export interface UserAccount {
  id: number;
  username: string;
}

export interface LoginResponse {
  access_token: string;
  token_type: string;
  expires_at: string;
  user: UserAccount;
}

export interface ProjectRecord {
  id: number;
  name: string;
  description: string;
  public_key?: string;
  key_bits: number;
  created_at: string;
  updated_at: string;
}

export interface FeatureRecord {
  id: number;
  code: string;
  name: string;
  description: string;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface PlanRecord {
  id: number;
  code: string;
  name: string;
  description: string;
  features: FeatureRecord[];
  max_users: number;
  max_devices: number;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface LicenseLimits {
  max_users: number;
  max_devices: number;
}

export interface LicenseRecord {
  id: number;
  uuid: string;
  project_id: number;
  subject_name: string;
  subject_email: string;
  subject_org: string;
  plan_id?: number;
  plan_name: string;
  plan: string;
  not_before: string;
  not_after: string;
  features: string[];
  limits: LicenseLimits;
  metadata: Record<string, unknown>;
  created_at: string;
}

export interface CreateProjectPayload {
  name: string;
  description: string;
  key_bits: number;
}

export interface UpdateProjectPayload {
  name: string;
  description: string;
}

export interface CreateFeaturePayload {
  code: string;
  name: string;
  description: string;
  enabled: boolean;
}

export type UpdateFeaturePayload = CreateFeaturePayload;

export interface CreatePlanPayload {
  code: string;
  name: string;
  description: string;
  feature_ids: number[];
  max_users: number;
  max_devices: number;
  enabled: boolean;
}

export type UpdatePlanPayload = CreatePlanPayload;

export interface LicenseHardwarePayload {
  hostname?: string;
  cpu_id?: string;
  disk_serial?: string;
  mac_addresses?: string[];
}

export interface IssueLicensePayload {
  subject: {
    name: string;
    email: string;
    organization: string;
  };
  plan_id: number;
  not_before: string;
  not_after: string;
  hardware: LicenseHardwarePayload;
}

export interface StoredSession {
  token: string;
  tokenType: string;
  expiresAt: string;
  user: UserAccount;
}

export interface DownloadFile {
  blob: Blob;
  filename: string;
}

export class ApiError extends Error {
  readonly status: number;
  readonly code?: number;

  constructor(message: string, status: number, code?: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

const requireEnv = (key: string, value?: string) => {
  if (!value) {
    throw new Error(`Missing required environment variable: ${key}`);
  }
  return value;
};

const API_BASE = requireEnv("VITE_API_BASE", import.meta.env.VITE_API_BASE);
const SESSION_KEY = "lykn.user.session";

const getSessionToken = () => getStoredSession()?.token || "";

const toHeaders = (body?: BodyInit | null): HeadersInit => {
  const headers: Record<string, string> = {
    Accept: "application/json",
  };
  if (body !== undefined && body !== null) {
    headers["Content-Type"] = "application/json";
  }
  const token = getSessionToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return headers;
};

const parseEnvelope = async <T>(response: Response): Promise<ApiEnvelope<T>> => {
  try {
    return (await response.json()) as ApiEnvelope<T>;
  } catch {
    throw new ApiError(response.statusText || "响应格式错误", response.status);
  }
};

const request = async <T>(path: string, init: RequestInit = {}): Promise<T> => {
  const body = init.body ?? undefined;
  const headers = new Headers(toHeaders(body));
  new Headers(init.headers).forEach((value, key) => headers.set(key, value));

  const response = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers,
  });
  const envelope = await parseEnvelope<T>(response);
  if (!response.ok || envelope.code !== 0) {
    throw new ApiError(
      envelope.message || response.statusText || "请求失败",
      response.status,
      envelope.code,
    );
  }
  return envelope.data as T;
};

const requestRaw = async (path: string): Promise<DownloadFile> => {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: toHeaders(),
  });
  if (!response.ok) {
    const envelope = await parseEnvelope<never>(response);
    throw new ApiError(
      envelope.message || response.statusText || "下载失败",
      response.status,
      envelope.code,
    );
  }
  const disposition = response.headers.get("Content-Disposition") || "";
  const filename = extractFilename(disposition) || "download";
  return {
    blob: await response.blob(),
    filename,
  };
};

const extractFilename = (disposition: string) => {
  const utf8Match = disposition.match(/filename\*=UTF-8''([^;]+)/i);
  if (utf8Match?.[1]) return decodeURIComponent(utf8Match[1]);
  const normalMatch = disposition.match(/filename="?([^";]+)"?/i);
  return normalMatch?.[1] || "";
};

export const getStoredSession = (): StoredSession | null => {
  const raw = window.localStorage.getItem(SESSION_KEY);
  if (!raw) return null;
  try {
    const session = JSON.parse(raw) as StoredSession;
    if (!session.token || new Date(session.expiresAt).getTime() <= Date.now()) {
      clearStoredSession();
      return null;
    }
    return session;
  } catch {
    clearStoredSession();
    return null;
  }
};

export const setStoredSession = (login: LoginResponse) => {
  const session: StoredSession = {
    token: login.access_token,
    tokenType: login.token_type,
    expiresAt: login.expires_at,
    user: login.user,
  };
  window.localStorage.setItem(SESSION_KEY, JSON.stringify(session));
  return session;
};

export const clearStoredSession = () => {
  window.localStorage.removeItem(SESSION_KEY);
};

export const api = {
  login(payload: { username: string; password: string }) {
    return request<LoginResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },
  listProjects() {
    return request<ProjectRecord[]>("/projects");
  },
  createProject(payload: CreateProjectPayload) {
    return request<ProjectRecord>("/projects", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },
  getProject(id: number) {
    return request<ProjectRecord>(`/projects/${id}`);
  },
  updateProject(id: number, payload: UpdateProjectPayload) {
    return request<ProjectRecord>(`/projects/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
    });
  },
  deleteProject(id: number) {
    return request<{ deleted: boolean }>(`/projects/${id}`, {
      method: "DELETE",
    });
  },
  downloadProjectPublicKey(id: number) {
    return requestRaw(`/projects/${id}/public-key`);
  },
  listFeatures() {
    return request<FeatureRecord[]>("/features");
  },
  createFeature(payload: CreateFeaturePayload) {
    return request<FeatureRecord>("/features", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },
  updateFeature(id: number, payload: UpdateFeaturePayload) {
    return request<FeatureRecord>(`/features/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
    });
  },
  deleteFeature(id: number) {
    return request<{ deleted: boolean }>(`/features/${id}`, {
      method: "DELETE",
    });
  },
  listPlans() {
    return request<PlanRecord[]>("/plans");
  },
  createPlan(payload: CreatePlanPayload) {
    return request<PlanRecord>("/plans", {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },
  updatePlan(id: number, payload: UpdatePlanPayload) {
    return request<PlanRecord>(`/plans/${id}`, {
      method: "PUT",
      body: JSON.stringify(payload),
    });
  },
  deletePlan(id: number) {
    return request<{ deleted: boolean }>(`/plans/${id}`, {
      method: "DELETE",
    });
  },
  listProjectLicenses(projectId: number) {
    return request<LicenseRecord[]>(`/projects/${projectId}/licenses`);
  },
  issueLicense(projectId: number, payload: IssueLicensePayload) {
    return request<LicenseRecord>(`/projects/${projectId}/licenses`, {
      method: "POST",
      body: JSON.stringify(payload),
    });
  },
  getLicense(id: number) {
    return request<LicenseRecord>(`/licenses/${id}`);
  },
  downloadLicense(id: number) {
    return requestRaw(`/licenses/${id}/download`);
  },
};
