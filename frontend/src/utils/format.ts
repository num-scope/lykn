import type { DownloadFile, LicenseRecord } from "../api/api";

const dateTimeFormatter = new Intl.DateTimeFormat("zh-CN", {
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
  hour: "2-digit",
  minute: "2-digit",
});

const dateFormatter = new Intl.DateTimeFormat("zh-CN", {
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
});

export const formatDateTime = (value?: string) => {
  if (!value) return "-";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "-" : dateTimeFormatter.format(date);
};

export const formatDate = (value?: string) => {
  if (!value) return "-";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "-" : dateFormatter.format(date);
};

export const formatPeriod = (start: string, end: string) =>
  `${formatDate(start)} - ${formatDate(end)}`;

export const getLicenseState = (license: LicenseRecord) => {
  const now = Date.now();
  const start = new Date(license.not_before).getTime();
  const end = new Date(license.not_after).getTime();
  if (Number.isNaN(start) || Number.isNaN(end)) return "unknown";
  if (now < start) return "upcoming";
  if (now > end) return "expired";
  return "active";
};

export const downloadFile = ({ blob, filename }: DownloadFile) => {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
};

export const safeJsonStringify = (value: unknown) => JSON.stringify(value ?? {}, null, 2);

export const getErrorMessage = (error: unknown, fallback = "操作失败") =>
  error instanceof Error && error.message ? error.message : fallback;
