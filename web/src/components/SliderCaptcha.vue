<template>
  <div class="slider-captcha">
    <div class="slider-captcha-container" :style="{ width: bgWidth + 'px', height: bgHeight + 'px' }">
      <!-- 背景图 -->
      <canvas ref="bgCanvas" class="bg-canvas" :width="bgWidth" :height="bgHeight"></canvas>
      <!-- 滑块 -->
      <canvas 
        ref="sliderCanvas" 
        class="slider-canvas" 
        :width="sliderSize" 
        :height="sliderSize"
        :style="{ top: sliderY + 'px', left: sliderX + 'px' }"
      ></canvas>
      <!-- 验证状态 -->
      <div v-if="status === 'success'" class="status-mask success">
        <CheckCircleOutlined /> 验证成功
      </div>
      <div v-if="status === 'error'" class="status-mask error">
        <CloseCircleOutlined /> 验证失败
      </div>
      <!-- 刷新按钮 -->
      <div class="refresh-btn" @click="refresh" v-if="status !== 'success'">
        <ReloadOutlined />
      </div>
    </div>
    <!-- 滑动条 -->
    <div class="slider-bar" :class="{ 'is-moving': isMoving }">
      <div class="slider-track" :style="{ width: sliderX + 'px' }"></div>
      <div 
        class="slider-handle" 
        :style="{ left: sliderX + 'px' }"
        @mousedown="handleMouseDown"
        @touchstart="handleTouchStart"
      >
        <template v-if="status === 'success'">
          <CheckOutlined style="color: #52c41a" />
        </template>
        <template v-else-if="status === 'error'">
          <CloseOutlined style="color: #ff4d4f" />
        </template>
        <template v-else>
          <DoubleRightOutlined />
        </template>
      </div>
      <span v-if="sliderX === 0 && status === 'pending'" class="slider-tip">向右拖动滑块完成验证</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { 
  ReloadOutlined, 
  CheckCircleOutlined, 
  CloseCircleOutlined,
  DoubleRightOutlined,
  CheckOutlined,
  CloseOutlined
} from '@ant-design/icons-vue'
import { getSliderCaptcha, verifySliderCaptcha } from '@/api/captcha'

const props = defineProps<{
  backgroundImage?: string // 自定义背景图片URL
}>()

const emit = defineEmits<{
  (e: 'success', captchaId: string): void
  (e: 'fail'): void
}>()

const bgWidth = 280
const bgHeight = 160
const sliderSize = 40
const bgImageLoaded = ref<HTMLImageElement | null>(null)

const bgCanvas = ref<HTMLCanvasElement | null>(null)
const sliderCanvas = ref<HTMLCanvasElement | null>(null)

const captchaId = ref('')
const targetX = ref(0)
const sliderY = ref(0)
const sliderX = ref(0)
const isMoving = ref(false)
const startX = ref(0)
const status = ref<'pending' | 'success' | 'error'>('pending')

// 生成随机背景图案
const drawBackground = (ctx: CanvasRenderingContext2D, targetXPos: number, targetYPos: number) => {
  // 如果有自定义背景图片且已加载
  if (bgImageLoaded.value) {
    ctx.drawImage(bgImageLoaded.value, 0, 0, bgWidth, bgHeight)
  } else {
    // 生成渐变背景
    const gradient = ctx.createLinearGradient(0, 0, bgWidth, bgHeight)
    gradient.addColorStop(0, '#667eea')
    gradient.addColorStop(1, '#764ba2')
    ctx.fillStyle = gradient
    ctx.fillRect(0, 0, bgWidth, bgHeight)
    
    // 添加随机圆形装饰
    for (let i = 0; i < 15; i++) {
      ctx.beginPath()
      const x = Math.random() * bgWidth
      const y = Math.random() * bgHeight
      const r = Math.random() * 20 + 5
      ctx.arc(x, y, r, 0, Math.PI * 2)
      ctx.fillStyle = `rgba(255, 255, 255, ${Math.random() * 0.15 + 0.05})`
      ctx.fill()
    }
  }
  
  // 绘制缺口
  ctx.save()
  ctx.globalCompositeOperation = 'destination-out'
  drawSliderPath(ctx, targetXPos, targetYPos)
  ctx.fill()
  ctx.restore()
  
  // 绘制缺口边框
  ctx.save()
  ctx.strokeStyle = 'rgba(255, 255, 255, 0.8)'
  ctx.lineWidth = 2
  drawSliderPath(ctx, targetXPos, targetYPos)
  ctx.stroke()
  ctx.restore()
}

// 绘制滑块
const drawSlider = (ctx: CanvasRenderingContext2D, bgCtx: CanvasRenderingContext2D, targetXPos: number, targetYPos: number) => {
  // 从背景图中截取滑块区域
  ctx.clearRect(0, 0, sliderSize, sliderSize)
  
  // 创建临时canvas获取原始背景
  const tempCanvas = document.createElement('canvas')
  tempCanvas.width = bgWidth
  tempCanvas.height = bgHeight
  const tempCtx = tempCanvas.getContext('2d')!
  
  // 重新绘制不带缺口的背景
  if (bgImageLoaded.value) {
    tempCtx.drawImage(bgImageLoaded.value, 0, 0, bgWidth, bgHeight)
  } else {
    const gradient = tempCtx.createLinearGradient(0, 0, bgWidth, bgHeight)
    gradient.addColorStop(0, '#667eea')
    gradient.addColorStop(1, '#764ba2')
    tempCtx.fillStyle = gradient
    tempCtx.fillRect(0, 0, bgWidth, bgHeight)
    
    for (let i = 0; i < 15; i++) {
      tempCtx.beginPath()
      const x = Math.random() * bgWidth
      const y = Math.random() * bgHeight
      const r = Math.random() * 20 + 5
      tempCtx.arc(x, y, r, 0, Math.PI * 2)
      tempCtx.fillStyle = `rgba(255, 255, 255, ${Math.random() * 0.15 + 0.05})`
      tempCtx.fill()
    }
  }
  
  // 裁剪滑块形状
  ctx.save()
  drawSliderPath(ctx, 0, 0)
  ctx.clip()
  ctx.drawImage(tempCanvas, targetXPos, targetYPos, sliderSize, sliderSize, 0, 0, sliderSize, sliderSize)
  ctx.restore()
  
  // 绘制滑块边框
  ctx.strokeStyle = 'rgba(255, 255, 255, 0.9)'
  ctx.lineWidth = 2
  drawSliderPath(ctx, 0, 0)
  ctx.stroke()
  
  // 添加阴影效果
  ctx.shadowColor = 'rgba(0, 0, 0, 0.3)'
  ctx.shadowBlur = 5
  ctx.shadowOffsetX = 2
  ctx.shadowOffsetY = 2
}

// 绘制滑块路径（带凸起）
const drawSliderPath = (ctx: CanvasRenderingContext2D, x: number, y: number) => {
  const r = 5
  const size = sliderSize
  const bumpSize = 8
  
  ctx.beginPath()
  ctx.moveTo(x + r, y)
  ctx.lineTo(x + size / 2 - bumpSize, y)
  ctx.arc(x + size / 2, y, bumpSize, Math.PI, 0, false) // 顶部凸起
  ctx.lineTo(x + size - r, y)
  ctx.arc(x + size - r, y + r, r, -Math.PI / 2, 0, false)
  ctx.lineTo(x + size, y + size / 2 - bumpSize)
  ctx.arc(x + size, y + size / 2, bumpSize, -Math.PI / 2, Math.PI / 2, false) // 右侧凸起
  ctx.lineTo(x + size, y + size - r)
  ctx.arc(x + size - r, y + size - r, r, 0, Math.PI / 2, false)
  ctx.lineTo(x + r, y + size)
  ctx.arc(x + r, y + size - r, r, Math.PI / 2, Math.PI, false)
  ctx.lineTo(x, y + r)
  ctx.arc(x + r, y + r, r, Math.PI, -Math.PI / 2, false)
  ctx.closePath()
}

// 预加载背景图片
const preloadBgImage = (): Promise<void> => {
  return new Promise((resolve) => {
    console.log('[SliderCaptcha] backgroundImage prop:', props.backgroundImage)
    if (!props.backgroundImage) {
      bgImageLoaded.value = null
      resolve()
      return
    }
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => {
      console.log('[SliderCaptcha] 背景图加载成功')
      bgImageLoaded.value = img
      resolve()
    }
    img.onerror = (e) => {
      console.error('[SliderCaptcha] 背景图加载失败:', e)
      bgImageLoaded.value = null
      resolve()
    }
    img.src = props.backgroundImage
  })
}

// 加载验证码
const loadCaptcha = async () => {
  status.value = 'pending'
  sliderX.value = 0
  
  try {
    // 先预加载背景图
    await preloadBgImage()
    
    const res = await getSliderCaptcha()
    captchaId.value = res.data.captcha_id
    targetX.value = res.data.target_x
    sliderY.value = res.data.slider_y
    
    // 绘制背景和滑块
    if (bgCanvas.value && sliderCanvas.value) {
      const bgCtx = bgCanvas.value.getContext('2d')!
      const sliderCtx = sliderCanvas.value.getContext('2d')!
      
      drawBackground(bgCtx, targetX.value, sliderY.value)
      drawSlider(sliderCtx, bgCtx, targetX.value, sliderY.value)
    }
  } catch (e) {
    console.error('获取滑动验证码失败', e)
  }
}

// 鼠标/触摸事件处理
const handleMouseDown = (e: MouseEvent) => {
  if (status.value === 'success') return
  isMoving.value = true
  startX.value = e.clientX - sliderX.value
  
  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
}

const handleTouchStart = (e: TouchEvent) => {
  if (status.value === 'success') return
  isMoving.value = true
  startX.value = e.touches[0].clientX - sliderX.value
  
  document.addEventListener('touchmove', handleTouchMove)
  document.addEventListener('touchend', handleTouchEnd)
}

const handleMouseMove = (e: MouseEvent) => {
  if (!isMoving.value) return
  updateSliderPosition(e.clientX)
}

const handleTouchMove = (e: TouchEvent) => {
  if (!isMoving.value) return
  updateSliderPosition(e.touches[0].clientX)
}

const updateSliderPosition = (clientX: number) => {
  let newX = clientX - startX.value
  newX = Math.max(0, Math.min(newX, bgWidth - sliderSize))
  sliderX.value = newX
}

const handleMouseUp = async () => {
  if (!isMoving.value) return
  isMoving.value = false
  
  document.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseup', handleMouseUp)
  
  await verify()
}

const handleTouchEnd = async () => {
  if (!isMoving.value) return
  isMoving.value = false
  
  document.removeEventListener('touchmove', handleTouchMove)
  document.removeEventListener('touchend', handleTouchEnd)
  
  await verify()
}

// 验证
const verify = async () => {
  try {
    const res = await verifySliderCaptcha(captchaId.value, sliderX.value)
    if (res.data?.success) {
      status.value = 'success'
      emit('success', captchaId.value)
    } else {
      status.value = 'error'
      emit('fail')
      setTimeout(() => {
        loadCaptcha()
      }, 1000)
    }
  } catch {
    status.value = 'error'
    emit('fail')
    setTimeout(() => {
      loadCaptcha()
    }, 1000)
  }
}

// 刷新
const refresh = () => {
  loadCaptcha()
}

// 暴露方法
defineExpose({
  refresh,
  getCaptchaId: () => captchaId.value
})

onMounted(() => {
  loadCaptcha()
})

onUnmounted(() => {
  document.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseup', handleMouseUp)
  document.removeEventListener('touchmove', handleTouchMove)
  document.removeEventListener('touchend', handleTouchEnd)
})
</script>

<style scoped>
.slider-captcha {
  width: 280px;
  user-select: none;
}

.slider-captcha-container {
  position: relative;
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.bg-canvas {
  display: block;
}

.slider-canvas {
  position: absolute;
  cursor: pointer;
}

.status-mask {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  gap: 8px;
}

.status-mask.success {
  background: rgba(82, 196, 26, 0.9);
  color: #fff;
}

.status-mask.error {
  background: rgba(255, 77, 79, 0.9);
  color: #fff;
}

.refresh-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 28px;
  height: 28px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: all 0.3s;
}

.refresh-btn:hover {
  background: #fff;
  color: #1890ff;
}

.slider-bar {
  position: relative;
  height: 40px;
  background: #f5f5f5;
  border-radius: 4px;
  margin-top: 12px;
  border: 1px solid #e8e8e8;
}

.slider-bar.is-moving {
  border-color: #1890ff;
}

.slider-track {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px 0 0 4px;
}

.slider-handle {
  position: absolute;
  top: 0;
  width: 40px;
  height: 38px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #667eea;
  font-size: 16px;
  transition: background 0.2s;
}

.slider-handle:hover {
  background: #f0f0f0;
}

.slider-tip {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #999;
  font-size: 13px;
  white-space: nowrap;
  pointer-events: none;
}
</style>
