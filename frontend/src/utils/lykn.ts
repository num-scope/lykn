import type { AxiosError } from 'axios';
import type { FlatResponseData } from '@sa/axios';

export function getServiceErrorMessage(error: unknown, fallback = '操作失败') {
  const axiosError = error as AxiosError<App.Service.Response>;
  const responseMessage = axiosError.response?.data?.message || axiosError.response?.data?.msg;

  if (responseMessage) return responseMessage;
  if (error instanceof Error && error.message) return error.message;

  return fallback;
}

export function formatDateTime(value?: string) {
  if (!value) return '-';

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;

  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date);
}

export function formatPeriod(start?: string, end?: string) {
  return `${formatDateTime(start)} 至 ${formatDateTime(end)}`;
}

export function getLicenseState(record: Api.Lykn.LicenseRecord) {
  const now = Date.now();
  const start = new Date(record.not_before).getTime();
  const end = new Date(record.not_after).getTime();

  if (Number.isNaN(start) || Number.isNaN(end)) return 'unknown';
  if (now < start) return 'upcoming';
  if (now > end) return 'expired';
  return 'active';
}

export function saveBlobFile(file: Api.Lykn.DownloadFile) {
  const url = URL.createObjectURL(file.blob);
  const link = document.createElement('a');

  link.href = url;
  link.download = file.filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

export function unwrapFlatData<T>(response: FlatResponseData<App.Service.Response<T>, T>, fallback = '请求失败') {
  const { data, error } = response;

  if (error || data === null) {
    throw new Error(getServiceErrorMessage(error, fallback));
  }

  return data;
}
