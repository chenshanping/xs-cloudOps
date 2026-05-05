<template>
  <div class="auth-container">
    <!-- 左侧背景图区域 -->
    <div class="showcase-section" :style="showcaseStyle">
      <div class="showcase-overlay">
        <div class="showcase-content">
          <img v-if="loginConfig.logo" :src="loginConfig.logo" alt="Logo" class="logo" />
          <h1 class="system-name">{{ loginConfig.system_name || '管理系统' }}</h1>
          <p class="slogan">{{ loginConfig.slogan || '欢迎使用' }}</p>
          <p v-if="loginConfig.login_desc" class="login-desc">{{ loginConfig.login_desc }}</p>
          
          <div v-if="loginConfig.features?.length" class="features">
            <div v-for="(feature, index) in loginConfig.features.slice(0, 4)" :key="index" class="feature-item">
              <component :is="getIcon(feature.icon)" class="feature-icon" />
              <div class="feature-text">
                <h3>{{ feature.title }}</h3>
                <p>{{ feature.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧表单区 -->
    <div class="form-section">
      <div class="form-box">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, provide } from 'vue'
import { 
  CheckCircleOutlined,
  SafetyOutlined,
  LineChartOutlined,
  ThunderboltOutlined,
  RocketOutlined,
  SettingOutlined,
  CloudOutlined,
  TeamOutlined,
  GlobalOutlined,
  DashboardOutlined,
  DatabaseOutlined,
  ApiOutlined,
  BulbOutlined,
  BookOutlined,
  FileProtectOutlined,
  SolutionOutlined,
  ExperimentOutlined,
  FundOutlined,
  ApartmentOutlined,
  ScheduleOutlined,
  MobileOutlined,
  LikeOutlined,
  StarOutlined,
  CrownOutlined
} from '@ant-design/icons-vue'
import { getCaptchaConfig } from '@/api/captcha'
import { useConfigStore } from '@/store/config'
const configStore = useConfigStore()
interface LoginConfig {
  logo?: string
  system_name?: string
  slogan?: string
  features?: Array<{
    icon: string
    title: string
    description: string
  }>
  login_desc?: ''
  showcase_images?: string[]
  login_bg_image?: string
  login_bg_color?: string
}

interface CaptchaConfig {
  login_captcha_enabled: boolean
  register_captcha_enabled: boolean
  register_email_verify: boolean
}

const loginConfig = ref<LoginConfig>({})
const captchaConfig = ref<CaptchaConfig>({
  login_captcha_enabled: false,
  register_captcha_enabled: false,
  register_email_verify: false,
})

// 提供给子组件
provide('captchaConfig', captchaConfig)

const iconMap: Record<string, any> = {
  CheckCircleOutlined,
  SafetyOutlined,
  LineChartOutlined,
  ThunderboltOutlined,
  RocketOutlined,
  SettingOutlined,
  CloudOutlined,
  TeamOutlined,
  GlobalOutlined,
  DashboardOutlined,
  DatabaseOutlined,
  ApiOutlined,
  BulbOutlined,
  BookOutlined,
  FileProtectOutlined,
  SolutionOutlined,
  ExperimentOutlined,
  FundOutlined,
  ApartmentOutlined,
  ScheduleOutlined,
  MobileOutlined,
  LikeOutlined,
  StarOutlined,
  CrownOutlined
}

const getIcon = (iconName: string) => {
  return iconMap[iconName] 
}

const showcaseStyle = computed(() => {
  const styles: Record<string, string> = {}
  if (loginConfig.value.login_bg_image) {
    styles.backgroundImage = `url(${loginConfig.value.login_bg_image})`
    styles.backgroundSize = 'cover'
    styles.backgroundPosition = 'center'
  } else if (loginConfig.value.login_bg_color) {
    styles.background = loginConfig.value.login_bg_color
  } else {
    // 默认渐变背景
    styles.background = 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  }
  return styles
})

const loadConfig = async () => {
  try {
    console.log(configStore.get('sys_logo'))
    // 并行加载配置
    const [captchaRes] = await Promise.all([
      getCaptchaConfig(),
    ])

    captchaConfig.value = {
      login_captcha_enabled: captchaRes.data.login_captcha_enabled,
      register_captcha_enabled: captchaRes.data.register_captcha_enabled,
      register_email_verify: captchaRes.data.register_email_verify,
    }

    loginConfig.value = {
      logo: configStore.get('sys_logo'),
      system_name: configStore.get('sys_name'),
      slogan: configStore.get('login_slogan'),
      login_desc: configStore.get('login_desc'),
      features: JSON.parse(configStore.get('login_features')),
      showcase_images: JSON.parse(configStore.get('login_images')),
      login_bg_image: configStore.get('login_bg_image'),
      login_bg_color: configStore.get('login_bg_color'),
    }
    console.log(loginConfig.value)
  } catch (error) {
    console.error('加载配置失败:', error)
  }
}

onMounted(() => {
  loadConfig()
})
</script>

<style scoped lang="scss">
.auth-container {
  min-height: 100vh;
  display: flex;
  width: 100%;
}

.showcase-section {
  flex: 1;
  position: relative;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;

  .showcase-overlay {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .showcase-content {
    text-align: center;
    color: #fff;
    padding: 40px;

    .logo {
      width: 80px;
      height: 80px;
      margin-bottom: 24px;
      border-radius: 16px;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
    }

    .system-name {
      font-size: 36px;
      font-weight: 600;
      margin: 0 0 16px 0;
      text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
    }

    .slogan {
      font-size: 18px;
      margin: 0 0 12px 0;
      opacity: 0.9;
    }

    .login-desc {
      font-size: 14px;
      margin: 0 0 32px 0;
      opacity: 0.8;
      max-width: 400px;
      line-height: 1.6;
    }

    .features {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 16px;
      max-width: 500px;
      text-align: left;

      .feature-item {
        display: flex;
        align-items: flex-start;
        gap: 12px;
        padding: 16px;
        background: rgba(255, 255, 255, 0.15);
        backdrop-filter: blur(10px);
        border-radius: 8px;
        border: 1px solid rgba(255, 255, 255, 0.2);
        transition: all 0.3s ease;

        &:hover {
          background: rgba(255, 255, 255, 0.25);
          transform: translateY(-2px);
        }

        .feature-icon {
          font-size: 24px;
          color: #fff;
          flex-shrink: 0;
        }

        .feature-text {
          flex: 1;
          min-width: 0;

          h3 {
            font-size: 14px;
            font-weight: 600;
            margin: 0 0 4px 0;
            color: #fff;
          }

          p {
            font-size: 12px;
            margin: 0;
            opacity: 0.8;
            line-height: 1.4;
            color: #fff;
          }
        }
      }
    }
  }
}

.form-section {
  width: 480px;
  min-width: 480px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 60px;
  background: #fff;
  box-shadow: -4px 0 20px rgba(0, 0, 0, 0.1);
}

.form-box {
  width: 100%;
  max-width: 360px;
}

@media (max-width: 900px) {
  .auth-container {
    flex-direction: column;
  }

  .showcase-section {
    min-height: 200px;

    .showcase-content {
      padding: 30px;

      .logo {
        width: 60px;
        height: 60px;
        margin-bottom: 16px;
      }

      .system-name {
        font-size: 28px;
        margin-bottom: 8px;
      }

      .slogan {
        font-size: 16px;
      }

      .login-desc {
        display: none;
      }
    }
  }

  .form-section {
    width: 100%;
    min-width: auto;
    flex: 1;
    padding: 40px 30px;
  }
}

@media (max-width: 600px) {
  .showcase-section {
    min-height: 160px;

    .showcase-content {
      padding: 20px;

      .logo {
        width: 48px;
        height: 48px;
        margin-bottom: 12px;
      }

      .system-name {
        font-size: 24px;
      }

      .slogan {
        font-size: 14px;
      }
    }
  }

  .form-section {
    padding: 30px 20px;
  }

  .form-box {
    max-width: 100%;
  }
}
</style>
