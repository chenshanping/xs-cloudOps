import type { LoginForm, ApiResponse } from '@/types';
import request from '@/utils/request';


// 登录请求参数
export interface LoginParams {
  username: string
  password: string
  captcha_id?: string
  captcha?: string
}

// 注册请求参数
export interface RegisterParams {
  username: string
  password: string
  email: string
  email_code?: string
  captcha_id?: string
  captcha?: string
  captcha_code?: string
}

// 发送邮箱验证码参数
export interface SendEmailCodeParams {
  email: string
  type?: string  // 'register' | 'reset_password'
  captcha_id?: string
  captcha?: string
}

// 重置密码参数（通过token）
export interface ResetPasswordParams {
  token: string
  new_password: string
}

// 重置密码参数（通过邮箱验证码）
export interface ResetPasswordByEmailParams {
  email: string
  email_code: string
  new_password: string
}

// 重置密码参数（通过用户名 + 图形验证码）
export interface ResetPasswordByUsernameParams {
  username: string
  new_password: string
  captcha_id: string
  captcha: string
}

// 登录
export function login(data: LoginParams) {
  return request.post<any, ApiResponse>('/auth/login', data)
}

// 注册
export function register(data: RegisterParams) {
  return request.post<any, ApiResponse>('/auth/register', data)
}

// 发送邮箱验证码
export function sendEmailCode(data: SendEmailCodeParams) {
  return request.post<any, ApiResponse>('/auth/send-email-code', data)
}

// 重置密码（通过token）
export function resetPassword(data: ResetPasswordParams) {
  return request.post<any, ApiResponse>('/auth/reset-password', data)
}

// 重置密码（通过邮箱验证码）
export function resetPasswordByEmail(data: ResetPasswordByEmailParams) {
  return request.post<any, ApiResponse>('/auth/reset-password-by-email', data)
}

// 重置密码（通过用户名 + 图形验证码）
export function resetPasswordByUsername(data: ResetPasswordByUsernameParams) {
  return request.post<any, ApiResponse>('/auth/reset-password-by-username', data)
}

// 获取用户信息
export function getUserInfo() {
  return request.get<any, ApiResponse>('/auth/userinfo')
}

// 登出
export function logout() {
  return request.post<any, ApiResponse>('/auth/logout', {}, { silent: true })
}

// 刷新Token
export function refreshToken() {
  return request.post<any, ApiResponse>('/auth/refresh')
}
