import type { UserInfo } from '@vben/types';

import { requestClient } from '#/api/request';

/**
 * 获取用户信息
 */
export async function getUserInfoApi() {
  return requestClient.get<UserInfo>('/user/info');
}

/**
 * 发送验证码
 */
export async function sendCodeApi(email: string) {
  return requestClient.post('/user/send-code', { email });
}

/**
 * 注册并激活
 * @example
 * {
 * "code": "string",
 * "confirm_password": "string",
 * "email": "string",
 * "password": "string",
 * "username": "string"
}
 */
export type RegisterAndActivateParams = {
  email: string;
  password: string;
  code: string;
  confirmPassword: string;
  username: string;
};
export async function registerAndActivateApi(
  params: RegisterAndActivateParams,
) {
  return requestClient.post('/user/register', JSON.stringify(params));
}

/**
 * 登录
 */
export type LoginParams = {
  username: string;
  password: string;
};
export async function loginApi(params: LoginParams) {
  return requestClient.post('/user/login', params);
}

/**
 * 退出登录
 */
export async function logoutApi() {
  return requestClient.post('/user/logout');
}

// 忘记密码-发送验证码
export async function forgotPasswordSendCodeApi(email: string) {
  return requestClient.post('/user/forgot-password/send-code', { email });
}

// 忘记密码-重置密码
export type ForgotPasswordResetParams = {
  code: 'string';
  confirmPassword: 'string';
  email: 'string';
  newPassword: 'string';
};
export async function forgotPasswordResetApi(
  params: ForgotPasswordResetParams,
) {
  return requestClient.post('/user/forgot-password/reset', params);
}
