import { getConfigsByKeys, batchUpdateConfigs } from '@/api/config';
import { defineStore } from 'pinia';
import { ref } from 'vue';


// 系统配置键
export const CONFIG_KEYS = [
  'sys_name',
  'sys_logo',
  'register_logo',
  'menu_bg_color',
  'menu_text_color',
  'menu_active_text_color',
  'menu_active_bg_color',
  'header_bg_color',
  'header_text_color',
  // 登录页配置
  'login_bg_image',
  'login_title',
  'login_subtitle',
  'login_bg_color',
  'login_slogan',
  'login_desc',
  'login_features',
  'login_features_max',
  'login_images',
  'login_images_max',
  // 注册相关配置
  'enable_register',
  // 邮箱配置
  'email_smtp_host',
  'email_smtp_port',
  'email_username',
  'email_password',
  'email_from_name',
  // 安全配置
  'login_captcha_enabled',
  'register_captcha_enabled',
  'register_email_verify',
  'frontend_url',
  // 前台模式配置
  'front_mode'  // 'full': 完整前台, 'profile': 仅个人中心(用于身份认证)
] as const

// 默认配置
const DEFAULT_CONFIG: Record<string, string> = {
  sys_name: 'Go RBAC Admin',
  sys_logo: '/src/assets/logo.svg',
  menu_bg_color: '#001529',
  menu_text_color: 'rgba(255, 255, 255, 0.65)',
  menu_active_text_color: '#ffffff',
  menu_active_bg_color: '#1890ff',
  header_bg_color: '#ffffff',
  header_text_color: '#333333',
  // 登录页默认配置
  login_bg_image: '',
  login_title: '欢迎回来',
  login_subtitle: '',
  login_bg_color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  login_slogan: '智能化企业管理平台',
  login_desc: '实时数据联动，智能分析策略，多渠道展示整合，让管理更简单，让体验更美好。',
  login_features: JSON.stringify([
    { icon: 'CheckCircleOutlined', title: '智能分析', desc: '动态调整业务策略' },
    { icon: 'SafetyOutlined', title: '安全合规', desc: '多层级权限保障' },
    { icon: 'LineChartOutlined', title: '数据洞察', desc: '全方位数据分析' },
    { icon: 'ThunderboltOutlined', title: '高效管理', desc: '简化日常运营流程' }
  ]),
  login_features_max: '4',
  login_images: JSON.stringify([]),
  login_images_max: '4',
  // 注册相关配置
  enable_register: 'false',
  // 邮箱配置
  email_smtp_host: '',
  email_smtp_port: '587',
  email_username: '',
  email_password: '',
  email_from_name: '',
  // 安全配置
  login_captcha_enabled: '0',
  register_captcha_enabled: '0',
  register_email_verify: '0',
  frontend_url: 'http://localhost:5173',
  // 前台模式: 'full' = 完整前台, 'profile' = 仅个人中心
  front_mode: 'full'
}

export type ConfigKey = typeof CONFIG_KEYS[number]

export const useConfigStore = defineStore('config', () => {
  const configs = ref<Record<string, string>>({ ...DEFAULT_CONFIG })
  const loading = ref(false)
  const loaded = ref(false) // 标记是否已加载

  // 获取单个配置值
  const get = (key: string): string => {
    return configs.value[key] || DEFAULT_CONFIG[key] || ''
  }

  // 加载配置（支持强制刷新）
  const loadConfigs = async (force = false) => {
    // 已加载且非强制刷新，直接返回
    if (loaded.value && !force) {
      return
    }
    loading.value = true
    try {
      const res = await getConfigsByKeys([...CONFIG_KEYS])
      const data = res.data
      for (const key of CONFIG_KEYS) {
        if (data[key]) {
          configs.value[key] = data[key].value
        }
      }
      loaded.value = true
    } catch (e) {
      console.error('加载配置失败', e)
    } finally {
      loading.value = false
    }
  }

  // 更新配置
  const updateConfigs = async (newConfigs: Record<string, string>) => {
    await batchUpdateConfigs(newConfigs)
    // 更新本地状态
    for (const key in newConfigs) {
      configs.value[key] = newConfigs[key]
    }
  }

  // 设置单个配置（本地）
  const set = (key: ConfigKey, value: string) => {
    configs.value[key] = value
  }

  return {
    configs,
    loading,
    loaded,
    get,
    set,
    loadConfigs,
    updateConfigs
  }
})
