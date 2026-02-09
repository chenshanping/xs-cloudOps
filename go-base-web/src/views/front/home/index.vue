<template>
  <div class="front-home-page">
    <!-- Hero 区域 -->
    <section class="hero-section">
      <div class="hero-content">
        <h1>智能 AI 对话助手</h1>
        <p class="hero-desc">基于先进的大语言模型，为您提供智能、高效的对话体验</p>
        <div class="hero-actions">
          <a-button type="primary" size="large" @click="handleStartChat">
            <template #icon><MessageOutlined /></template>
            开始对话
          </a-button>
          <a-button size="large" v-if="!userStore.token" @click="router.push('/login')">
            立即登录
          </a-button>
        </div>
      </div>
      <div class="hero-illustration">
        <RobotOutlined class="robot-icon" />
      </div>
    </section>
    
    <!-- 功能特性 -->
    <section class="features-section">
      <h2 class="section-title">核心功能</h2>
      <a-row :gutter="24">
        <a-col :xs="24" :sm="12" :md="6" v-for="feature in features" :key="feature.title">
          <div class="feature-card">
            <div class="feature-icon" :style="{ background: feature.color }">
              <component :is="feature.icon" />
            </div>
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.desc }}</p>
          </div>
        </a-col>
      </a-row>
    </section>
    
    <!-- 使用场景 -->
    <section class="scenarios-section">
      <h2 class="section-title">使用场景</h2>
      <a-row :gutter="24">
        <a-col :xs="24" :md="8" v-for="scenario in scenarios" :key="scenario.title">
          <div class="scenario-card">
            <component :is="scenario.icon" class="scenario-icon" />
            <h3>{{ scenario.title }}</h3>
            <p>{{ scenario.desc }}</p>
          </div>
        </a-col>
      </a-row>
    </section>
    
    <!-- 开始使用 -->
    <section class="cta-section">
      <div class="cta-content">
        <h2>准备好开始了吗？</h2>
        <p>立即体验智能 AI 对话，让工作更高效</p>
        <a-button type="primary" size="large" @click="handleStartChat">
          免费开始使用
        </a-button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  MessageOutlined,
  RobotOutlined,
  ThunderboltOutlined,
  GlobalOutlined,
  SafetyOutlined,
  ClockCircleOutlined,
  CodeOutlined,
  EditOutlined,
  SearchOutlined
} from '@ant-design/icons-vue'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

const features = [
  {
    icon: ThunderboltOutlined,
    title: '智能对话',
    desc: '基于先进AI模型，理解上下文，给出精准回复',
    color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  },
  {
    icon: GlobalOutlined,
    title: '联网搜索',
    desc: '实时联网获取最新信息，回答更加准确全面',
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)'
  },
  {
    icon: SafetyOutlined,
    title: '安全可靠',
    desc: '数据加密传输，对话内容安全有保障',
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
  },
  {
    icon: ClockCircleOutlined,
    title: '历史记录',
    desc: '自动保存对话历史，随时回顾查看',
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)'
  }
]

const scenarios = [
  {
    icon: CodeOutlined,
    title: '编程开发',
    desc: '代码编写、Bug调试、技术问答，提升开发效率'
  },
  {
    icon: EditOutlined,
    title: '内容创作',
    desc: '文章撰写、文案润色、创意激发，助力内容产出'
  },
  {
    icon: SearchOutlined,
    title: '知识问答',
    desc: '专业知识查询、学习辅导、问题解答，随时随地学习'
  }
]

const handleStartChat = () => {
  if (!userStore.token) {
    message.info('请先登录后使用')
    router.push('/login')
    return
  }
  router.push('/front/ai')
}
</script>

<style scoped lang="less">
.front-home-page {
  .hero-section {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 60px 0;
    
    .hero-content {
      max-width: 500px;
      
      h1 {
        font-size: 42px;
        font-weight: 700;
        color: #1a1a1a;
        margin-bottom: 16px;
        line-height: 1.2;
      }
      
      .hero-desc {
        font-size: 18px;
        color: #666;
        margin-bottom: 32px;
        line-height: 1.6;
      }
      
      .hero-actions {
        display: flex;
        gap: 16px;
      }
    }
    
    .hero-illustration {
      .robot-icon {
        font-size: 200px;
        color: #1890ff;
        opacity: 0.8;
      }
    }
  }
  
  .section-title {
    text-align: center;
    font-size: 28px;
    font-weight: 600;
    color: #1a1a1a;
    margin-bottom: 40px;
  }
  
  .features-section {
    padding: 60px 0;
    background: #fff;
    margin: 0 -24px;
    padding-left: 24px;
    padding-right: 24px;
    border-radius: 16px;
    
    .feature-card {
      text-align: center;
      padding: 32px 24px;
      border-radius: 12px;
      transition: all 0.3s;
      
      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
      }
      
      .feature-icon {
        width: 64px;
        height: 64px;
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 20px;
        
        :deep(.anticon) {
          font-size: 28px;
          color: #fff;
        }
      }
      
      h3 {
        font-size: 18px;
        font-weight: 600;
        color: #1a1a1a;
        margin-bottom: 8px;
      }
      
      p {
        font-size: 14px;
        color: #666;
        margin: 0;
        line-height: 1.6;
      }
    }
  }
  
  .scenarios-section {
    padding: 60px 0;
    
    .scenario-card {
      background: #fff;
      padding: 32px;
      border-radius: 12px;
      text-align: center;
      height: 100%;
      transition: all 0.3s;
      
      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
      }
      
      .scenario-icon {
        font-size: 48px;
        color: #1890ff;
        margin-bottom: 20px;
      }
      
      h3 {
        font-size: 18px;
        font-weight: 600;
        color: #1a1a1a;
        margin-bottom: 12px;
      }
      
      p {
        font-size: 14px;
        color: #666;
        margin: 0;
        line-height: 1.6;
      }
    }
  }
  
  .cta-section {
    padding: 80px 0;
    text-align: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    margin: 0 -24px -24px;
    border-radius: 16px 16px 0 0;
    
    .cta-content {
      h2 {
        font-size: 32px;
        font-weight: 600;
        color: #fff;
        margin-bottom: 12px;
      }
      
      p {
        font-size: 16px;
        color: rgba(255, 255, 255, 0.85);
        margin-bottom: 32px;
      }
      
      .ant-btn {
        height: 48px;
        padding: 0 32px;
        font-size: 16px;
      }
    }
  }
}

@media (max-width: 768px) {
  .front-home-page {
    .hero-section {
      flex-direction: column;
      text-align: center;
      padding: 40px 0;
      
      .hero-content {
        h1 {
          font-size: 28px;
        }
        
        .hero-desc {
          font-size: 16px;
        }
        
        .hero-actions {
          justify-content: center;
        }
      }
      
      .hero-illustration {
        margin-top: 40px;
        
        .robot-icon {
          font-size: 120px;
        }
      }
    }
    
    .features-section,
    .scenarios-section {
      padding: 40px 0;
      
      .ant-col {
        margin-bottom: 24px;
      }
    }
  }
}
</style>
