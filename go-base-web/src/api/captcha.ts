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
    login_captcha_type: string
    register_captcha_enabled: boolean
    register_email_verify: boolean
    slider_captcha_bg?: string
  }>('/captcha/config')
}

// 获取滑动验证码
export function getSliderCaptcha() {
  return request.get<{
    captcha_id: string
    bg_width: number
    bg_height: number
    slider_y: number
    target_x: number
  }>('/captcha/slider')
}

// 验证滑动验证码
export function verifySliderCaptcha(captchaId: string, x: number) {
  return request.post<{ success: boolean }>('/captcha/slider/verify', { captcha_id: captchaId, x })
}
