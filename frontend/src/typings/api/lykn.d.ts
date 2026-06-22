declare namespace Api {
  namespace Lykn {
    interface DashboardSummary {
      project_count: number;
      license_count: number;
      active_license_count: number;
      expired_license_count: number;
    }

    interface ProjectRecord {
      id: number;
      name: string;
      description: string;
      public_key?: string;
      key_bits: number;
      created_at: string;
      updated_at: string;
    }

    interface CreateProjectPayload {
      name: string;
      description: string;
      key_bits: number;
    }

    interface UpdateProjectPayload {
      name: string;
      description: string;
    }

    interface FeatureRecord {
      id: number;
      code: string;
      name: string;
      description: string;
      enabled: boolean;
      created_at: string;
      updated_at: string;
    }

    interface FeaturePayload {
      code: string;
      name: string;
      description: string;
      enabled: boolean;
    }

    interface PlanRecord {
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

    interface PlanPayload {
      code: string;
      name: string;
      description: string;
      feature_ids: number[];
      max_users: number;
      max_devices: number;
      enabled: boolean;
    }

    interface LicenseLimits {
      max_users: number;
      max_devices: number;
    }

    interface LicenseRecord {
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

    interface LicenseHardwarePayload {
      hostname?: string;
      cpu_id?: string;
      disk_serial?: string;
      mac_addresses?: string[];
      ip_addresses?: string[];
    }

    interface IssueLicensePayload {
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

    interface DownloadFile {
      blob: Blob;
      filename: string;
    }
  }
}
