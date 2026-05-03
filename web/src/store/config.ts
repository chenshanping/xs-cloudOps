import { getConfigList, getConfigsByKeys, batchUpdateConfigs } from '@/api/config';
import { defineStore } from 'pinia';
import { ref } from 'vue';


// 公开配置键：未登录阶段允许加载
export const PUBLIC_CONFIG_KEYS = [
  'sys_name',
  'sys_logo',
  'register_logo',
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
  // 前台模式配置
  'front_mode',  // 'full': 完整前台, 'profile': 仅个人中心, 'none': 无前台(纯后台模式)
  'user_profile_button_visible',
] as const

// 后台配置键：登录后后台按需补充加载
export const ADMIN_CONFIG_KEYS = [
  // 邮箱配置
  'email_smtp_host',
  'email_smtp_port',
  'email_username',
  'email_password',
  'email_from_name',
  // 安全配置
  'login_captcha_enabled',
  'login_captcha_type',
  'login_max_retry',
  'login_lock_time',
  'register_email_verify',
  'frontend_url',
  'slider_captcha_bg',
  'user_default_password',
  'public_config_keys',
  'file_delete_mode',
  'storage_type',
  'storage_local_config',
  'storage_aliyun_config',
  'storage_tencent_config',
  'storage_minio_config'
] as const

export const CONFIG_KEYS = [...PUBLIC_CONFIG_KEYS, ...ADMIN_CONFIG_KEYS] as const

// 默认配置
const DEFAULT_CONFIG: Record<string, string> = {
  sys_name: 'Go RBAC Admin',
  sys_logo: '/src/assets/logo.svg',
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
  login_captcha_type: 'digit',
  login_max_retry: '5',
  login_lock_time: '15',
  register_email_verify: '0',
  frontend_url: 'http://localhost:5173',
  slider_captcha_bg: '',
  // 前台模式: 'full' = 完整前台, 'profile' = 仅个人中心, 'none' = 无前台(纯后台)
  front_mode: 'full',
  user_profile_button_visible: 'false',
  user_default_password: '123456',
  public_config_keys: JSON.stringify(PUBLIC_CONFIG_KEYS),
  file_delete_mode: 'logical',
  storage_type: 'local',
  storage_local_config: JSON.stringify({
    base_path: 'uploads',
    base_url: '/api/v1/upload'
  }),
  storage_aliyun_config: JSON.stringify({
    endpoint: '',
    access_key_id: '',
    access_key_secret: '',
    bucket_name: '',
    region: ''
  }),
  storage_tencent_config: JSON.stringify({
    region: '',
    secret_id: '',
    secret_key: '',
    bucket: '',
    app_id: ''
  }),
  storage_minio_config: JSON.stringify({
    endpoint: '',
    access_key_id: '',
    secret_access_key: '',
    bucket_name: '',
    use_ssl: false
  })
}

export type ConfigKey = typeof CONFIG_KEYS[number]
type ConfigScope = 'public' | 'all'

export const useConfigStore = defineStore('config', () => {
  const configs = ref<Record<string, string>>({ ...DEFAULT_CONFIG })
  const loading = ref(false)
  const publicLoaded = ref(false)
  const adminLoaded = ref(false)

  // 获取单个配置值
  const get = (key: string): string => {
    return configs.value[key] || DEFAULT_CONFIG[key] || ''
  }

  const applyConfigs = (keys: readonly string[], data: Record<string, any>) => {
    for (const key of keys) {
      if (data[key]) {
        configs.value[key] = data[key].value
      }
    }
  }

  const applyConfigList = (items: Array<{ key: string; value: string }>) => {
    for (const item of items) {
      configs.value[item.key] = item.value
    }
  }

  // 加载配置（支持按公开/后台作用域加载）
  const loadConfigs = async (force = false, scope: ConfigScope = 'public') => {
    const needPublic = force || !publicLoaded.value
    const needAdmin = scope === 'all' && (force || !adminLoaded.value)

    if (!needPublic && !needAdmin) {
      return
    }

    loading.value = true
    try {
      if (needAdmin) {
        const res = await getConfigList()
        const items = Array.isArray(res.data) ? res.data : []
        applyConfigList(items)
        publicLoaded.value = true
        adminLoaded.value = true
        return
      }

      if (needPublic) {
        const res = await getConfigsByKeys([...PUBLIC_CONFIG_KEYS])
        const data = res.data
        applyConfigs(PUBLIC_CONFIG_KEYS, data)
        publicLoaded.value = true
      }
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
    publicLoaded,
    adminLoaded,
    get,
    set,
    loadConfigs,
    updateConfigs
  }
})
