import { request } from '../request';

export function fetchLogin(userName: string, password: string) {
  return request<Api.Auth.LoginResponse>({
    url: '/auth/login',
    method: 'post',
    data: {
      username: userName,
      password
    }
  });
}

export function fetchGetUserInfo() {
  return request<Api.Auth.UserAccount>({ url: '/auth/me' });
}

export function fetchRefreshToken(refreshToken: string) {
  return request<Api.Auth.LoginToken>({
    url: '/auth/refreshToken',
    method: 'post',
    data: {
      refreshToken
    }
  });
}

export function fetchCustomBackendError(code: string, msg: string) {
  return request({ url: '/auth/error', params: { code, msg } });
}
