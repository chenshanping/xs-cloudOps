import request from '@/utils/request';


// 获取图形验证码
export function getCaptcha() {
  return request.get<{
    captcha_id: string
    captcha_image: string
  }>('/captcha')
}

// 获取验证码配置
export function getCaptchaConfig():any {
  return request.get<{
    login_captcha_enabled: boolean
    register_captcha_enabled: boolean
    register_email_verify: boolean
  }>('/captcha/config')
}
