<template>
  <div class="er-diagram-container" ref="containerRef">
    <div class="er-toolbar">
      <a-space>
        <a-button size="small" @click="zoomIn"><ZoomInOutlined /></a-button>
        <a-button size="small" @click="zoomOut"><ZoomOutOutlined /></a-button>
        <a-button size="small" @click="resetZoom"><ExpandOutlined /> 重置</a-button>
        <a-button size="small" @click="autoLayout"><ReloadOutlined /> 自动布局</a-button>
        <a-divider type="vertical" />
        <span style="color: #666; font-size: 12px">拖拽实体可调整位置</span>
        <a-divider type="vertical" />
        <a-button size="small" type="primary" @click="exportPNG"><DownloadOutlined /> 导出PNG</a-button>
        <a-button size="small" @click="exportSVG"><DownloadOutlined /> 导出SVG</a-button>
      </a-space>
    </div>
    <div class="er-canvas-wrapper" ref="wrapperRef" @wheel="handleWheel" @mousedown="handleCanvasMouseDown" @mousemove="handleCanvasMouseMove" @mouseup="handleCanvasMouseUp" @mouseleave="handleCanvasMouseUp">
      <svg 
        ref="svgRef" 
        :width="canvasWidth" 
        :height="canvasHeight" 
        :viewBox="`${viewBox.x} ${viewBox.y} ${viewBox.width} ${viewBox.height}`"
        class="er-svg"
      >
        <defs>
          <!-- 渐变定义 -->
          <linearGradient id="entityGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" style="stop-color:#4A90D9;stop-opacity:1" />
            <stop offset="100%" style="stop-color:#357ABD;stop-opacity:1" />
          </linearGradient>
          <linearGradient id="attrGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" style="stop-color:#5BA0E0;stop-opacity:1" />
            <stop offset="100%" style="stop-color:#4A90D9;stop-opacity:1" />
          </linearGradient>
          <linearGradient id="relationGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" style="stop-color:#4A90D9;stop-opacity:1" />
            <stop offset="100%" style="stop-color:#357ABD;stop-opacity:1" />
          </linearGradient>
          <!-- 阴影 -->
          <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
            <feDropShadow dx="2" dy="2" stdDeviation="2" flood-opacity="0.3"/>
          </filter>
        </defs>
        
        <!-- 连线层 -->
        <g class="links-layer">
          <!-- 实体到属性的连线 -->
          <line 
            v-for="(link, idx) in attrLinks" 
            :key="'attr-link-' + idx"
            :x1="link.x1" 
            :y1="link.y1" 
            :x2="link.x2" 
            :y2="link.y2"
            stroke="#2C5F8D"
            stroke-width="2"
          />
          <!-- 实体到关系的连线 -->
          <g v-for="(link, idx) in relationLinks" :key="'rel-link-' + idx">
            <line 
              :x1="link.x1" 
              :y1="link.y1" 
              :x2="link.x2" 
              :y2="link.y2"
              stroke="#2C5F8D"
              stroke-width="2"
            />
            <!-- 基数标签 -->
            <text 
              :x="link.labelX" 
              :y="link.labelY" 
              fill="#333"
              font-size="14"
              font-weight="bold"
              text-anchor="middle"
            >{{ link.cardinality }}</text>
          </g>
        </g>
        
        <!-- 实体层 -->
        <g class="entities-layer">
          <g 
            v-for="entity in entities" 
            :key="'entity-' + entity.id" 
            class="entity-group draggable"
            :class="{ dragging: draggingEntity === entity.id }"
            @mousedown.stop="startDragEntity($event, entity)"
            style="cursor: move"
          >
            <rect 
              :x="entity.x - entity.width/2" 
              :y="entity.y - entity.height/2"
              :width="entity.width"
              :height="entity.height"
              fill="url(#entityGradient)"
              stroke="#2C5F8D"
              :stroke-width="draggingEntity === entity.id ? 3 : 2"
              rx="3"
              filter="url(#shadow)"
            />
            <text 
              :x="entity.x" 
              :y="entity.y + 5"
              fill="white"
              font-size="16"
              font-weight="bold"
              text-anchor="middle"
              style="pointer-events: none"
            >{{ entity.name }}</text>
          </g>
        </g>
        
        <!-- 属性层 -->
        <g class="attributes-layer">
          <g v-for="attr in attributes" :key="'attr-' + attr.id" class="attr-group">
            <ellipse 
              :cx="attr.x" 
              :cy="attr.y"
              :rx="attr.width/2"
              :ry="attr.height/2"
              fill="url(#attrGradient)"
              stroke="#2C5F8D"
              stroke-width="2"
              filter="url(#shadow)"
            />
            <text 
              :x="attr.x" 
              :y="attr.y + 4"
              fill="white"
              font-size="12"
              font-weight="500"
              text-anchor="middle"
              :text-decoration="attr.isPrimary ? 'underline' : 'none'"
            >{{ attr.name }}</text>
          </g>
        </g>
        
        <!-- 关系层 -->
        <g class="relations-layer">
          <g v-for="rel in relations" :key="'rel-' + rel.id" class="relation-group">
            <polygon 
              :points="getDiamondPoints(rel)"
              fill="url(#relationGradient)"
              stroke="#2C5F8D"
              stroke-width="2"
              filter="url(#shadow)"
            />
            <text 
              :x="rel.x" 
              :y="rel.y + 5"
              fill="white"
              font-size="14"
              font-weight="bold"
              text-anchor="middle"
            >{{ rel.name }}</text>
          </g>
        </g>
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { ZoomInOutlined, ZoomOutOutlined, ExpandOutlined, DownloadOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'

// Props
interface EntityConfig {
  name: string
  comment: string
  columns: {
    name: string
    comment: string
    isPrimary?: boolean
  }[]
}

interface RelationConfig {
  name: string
  from: string
  to: string
  fromCardinality: string  // '1', 'n', 'm'
  toCardinality: string
}

interface Props {
  entities: EntityConfig[]
  relations: RelationConfig[]
}

const props = defineProps<Props>()

// Refs
const containerRef = ref<HTMLElement>()
const wrapperRef = ref<HTMLElement>()
const svgRef = ref<SVGSVGElement>()

// 画布尺寸
const canvasWidth = ref(1200)
const canvasHeight = ref(800)

// 视图控制
const viewBox = ref({ x: 0, y: 0, width: 1200, height: 800 })
const scale = ref(1)
const isPanning = ref(false)
const panStart = ref({ x: 0, y: 0 })

// 布局后的元素
interface LayoutEntity {
  id: string
  name: string
  x: number
  y: number
  width: number
  height: number
}

interface LayoutAttribute {
  id: string
  name: string
  entityId: string
  x: number
  y: number
  width: number
  height: number
  isPrimary: boolean
}

interface LayoutRelation {
  id: string
  name: string
  x: number
  y: number
  width: number
  height: number
  from: string
  to: string
  fromCardinality: string
  toCardinality: string
}

interface AttrLink {
  x1: number
  y1: number
  x2: number
  y2: number
}

interface RelationLink {
  x1: number
  y1: number
  x2: number
  y2: number
  labelX: number
  labelY: number
  cardinality: string
}

const entities = ref<LayoutEntity[]>([])
const attributes = ref<LayoutAttribute[]>([])
const relations = ref<LayoutRelation[]>([])
const attrLinks = ref<AttrLink[]>([])
const relationLinks = ref<RelationLink[]>([])

// 拖拽状态
const draggingEntity = ref<string | null>(null)
const dragStartPos = ref({ x: 0, y: 0 })
const entityStartPos = ref({ x: 0, y: 0 })

// 实体位置缓存（用于保持用户调整后的位置）
const entityPositions = ref<Map<string, { x: number; y: number }>>(new Map())

// 计算布局
const calculateLayout = () => {
  const entityList: LayoutEntity[] = []
  const attrList: LayoutAttribute[] = []
  const relList: LayoutRelation[] = []
  const attrLinkList: AttrLink[] = []
  const relLinkList: RelationLink[] = []

  if (!props.entities || props.entities.length === 0) return

  // 实体布局参数
  const entityWidth = 80
  const entityHeight = 40
  const attrWidth = 70
  const attrHeight = 30
  const relWidth = 70
  const relHeight = 50
  const entitySpacingX = 300
  const entitySpacingY = 250
  const attrRadius = 80

  // 计算实体位置 - 网格布局
  const cols = Math.ceil(Math.sqrt(props.entities.length))
  const startX = 200
  const startY = 150

  props.entities.forEach((entity, idx) => {
    const col = idx % cols
    const row = Math.floor(idx / cols)
    const x = startX + col * entitySpacingX
    const y = startY + row * entitySpacingY

    const entityId = `entity-${idx}`
    entityList.push({
      id: entityId,
      name: entity.comment || entity.name,
      x,
      y,
      width: entityWidth,
      height: entityHeight
    })

    // 计算属性位置 - 围绕实体分布
    const attrCount = entity.columns.length
    const angleStep = (Math.PI * 2) / Math.max(attrCount, 1)
    const startAngle = -Math.PI / 2

    entity.columns.forEach((col, colIdx) => {
      const angle = startAngle + colIdx * angleStep
      const attrX = x + Math.cos(angle) * attrRadius
      const attrY = y + Math.sin(angle) * attrRadius

      attrList.push({
        id: `attr-${idx}-${colIdx}`,
        name: col.comment || col.name,
        entityId,
        x: attrX,
        y: attrY,
        width: attrWidth,
        height: attrHeight,
        isPrimary: col.isPrimary || col.name === 'id'
      })

      // 添加实体到属性的连线
      attrLinkList.push({
        x1: x,
        y1: y,
        x2: attrX,
        y2: attrY
      })
    })
  })

  // 计算关系位置 - 放在两个实体之间
  props.relations.forEach((rel, idx) => {
    const fromEntity = entityList.find(e => 
      props.entities.find(pe => (pe.comment || pe.name) === e.name)?.name === rel.from ||
      props.entities.find(pe => (pe.comment || pe.name) === e.name)?.comment === rel.from
    )
    const toEntity = entityList.find(e => 
      props.entities.find(pe => (pe.comment || pe.name) === e.name)?.name === rel.to ||
      props.entities.find(pe => (pe.comment || pe.name) === e.name)?.comment === rel.to
    )

    if (fromEntity && toEntity) {
      const relX = (fromEntity.x + toEntity.x) / 2
      const relY = (fromEntity.y + toEntity.y) / 2

      relList.push({
        id: `rel-${idx}`,
        name: rel.name,
        x: relX,
        y: relY,
        width: relWidth,
        height: relHeight,
        from: fromEntity.id,
        to: toEntity.id,
        fromCardinality: rel.fromCardinality,
        toCardinality: rel.toCardinality
      })

      // 添加实体到关系的连线
      const fromMidX = (fromEntity.x + relX) / 2
      const fromMidY = (fromEntity.y + relY) / 2
      const toMidX = (toEntity.x + relX) / 2
      const toMidY = (toEntity.y + relY) / 2

      relLinkList.push({
        x1: fromEntity.x,
        y1: fromEntity.y,
        x2: relX,
        y2: relY,
        labelX: fromMidX,
        labelY: fromMidY - 10,
        cardinality: rel.fromCardinality
      })

      relLinkList.push({
        x1: relX,
        y1: relY,
        x2: toEntity.x,
        y2: toEntity.y,
        labelX: toMidX,
        labelY: toMidY - 10,
        cardinality: rel.toCardinality
      })
    }
  })

  // 更新画布大小
  const allX = [...entityList.map(e => e.x), ...attrList.map(a => a.x)]
  const allY = [...entityList.map(e => e.y), ...attrList.map(a => a.y)]
  const maxX = Math.max(...allX) + 150
  const maxY = Math.max(...allY) + 100
  canvasWidth.value = Math.max(maxX, 800)
  canvasHeight.value = Math.max(maxY, 600)
  viewBox.value = { x: 0, y: 0, width: canvasWidth.value, height: canvasHeight.value }

  entities.value = entityList
  attributes.value = attrList
  relations.value = relList
  attrLinks.value = attrLinkList
  relationLinks.value = relLinkList
}

// 获取菱形顶点
const getDiamondPoints = (rel: LayoutRelation): string => {
  const { x, y, width, height } = rel
  const hw = width / 2
  const hh = height / 2
  return `${x},${y - hh} ${x + hw},${y} ${x},${y + hh} ${x - hw},${y}`
}

// 缩放控制
const zoomIn = () => {
  scale.value = Math.min(scale.value * 1.2, 3)
  updateViewBox()
}

const zoomOut = () => {
  scale.value = Math.max(scale.value / 1.2, 0.3)
  updateViewBox()
}

const resetZoom = () => {
  scale.value = 1
  viewBox.value = { x: 0, y: 0, width: canvasWidth.value, height: canvasHeight.value }
}

const updateViewBox = () => {
  const newWidth = canvasWidth.value / scale.value
  const newHeight = canvasHeight.value / scale.value
  viewBox.value.width = newWidth
  viewBox.value.height = newHeight
}

const handleWheel = (e: WheelEvent) => {
  e.preventDefault()
  if (e.deltaY < 0) {
    zoomIn()
  } else {
    zoomOut()
  }
}

// 实体拖拽
const startDragEntity = (e: MouseEvent, entity: LayoutEntity) => {
  draggingEntity.value = entity.id
  const svgRect = svgRef.value?.getBoundingClientRect()
  if (!svgRect) return
  
  dragStartPos.value = { x: e.clientX, y: e.clientY }
  entityStartPos.value = { x: entity.x, y: entity.y }
}

const updateEntityPosition = (entityId: string, newX: number, newY: number) => {
  // 更新实体位置
  const entity = entities.value.find(e => e.id === entityId)
  if (!entity) return
  
  entity.x = newX
  entity.y = newY
  
  // 更新属性位置（属性围绕实体）
  const entityAttrs = attributes.value.filter(a => a.entityId === entityId)
  const attrCount = entityAttrs.length
  const angleStep = (Math.PI * 2) / Math.max(attrCount, 1)
  const startAngle = -Math.PI / 2
  const attrRadius = 80
  
  entityAttrs.forEach((attr, idx) => {
    const angle = startAngle + idx * angleStep
    attr.x = newX + Math.cos(angle) * attrRadius
    attr.y = newY + Math.sin(angle) * attrRadius
  })
  
  // 更新连线
  updateLinks()
}

const updateLinks = () => {
  // 更新属性连线
  const newAttrLinks: AttrLink[] = []
  attributes.value.forEach(attr => {
    const entity = entities.value.find(e => e.id === attr.entityId)
    if (entity) {
      newAttrLinks.push({
        x1: entity.x,
        y1: entity.y,
        x2: attr.x,
        y2: attr.y
      })
    }
  })
  attrLinks.value = newAttrLinks
  
  // 更新关系连线
  const newRelLinks: RelationLink[] = []
  relations.value.forEach(rel => {
    const fromEntity = entities.value.find(e => e.id === rel.from)
    const toEntity = entities.value.find(e => e.id === rel.to)
    
    if (fromEntity && toEntity) {
      // 更新关系位置到两个实体中间
      rel.x = (fromEntity.x + toEntity.x) / 2
      rel.y = (fromEntity.y + toEntity.y) / 2
      
      const fromMidX = (fromEntity.x + rel.x) / 2
      const fromMidY = (fromEntity.y + rel.y) / 2
      const toMidX = (toEntity.x + rel.x) / 2
      const toMidY = (toEntity.y + rel.y) / 2
      
      newRelLinks.push({
        x1: fromEntity.x,
        y1: fromEntity.y,
        x2: rel.x,
        y2: rel.y,
        labelX: fromMidX,
        labelY: fromMidY - 10,
        cardinality: rel.fromCardinality
      })
      
      newRelLinks.push({
        x1: rel.x,
        y1: rel.y,
        x2: toEntity.x,
        y2: toEntity.y,
        labelX: toMidX,
        labelY: toMidY - 10,
        cardinality: rel.toCardinality
      })
    }
  })
  relationLinks.value = newRelLinks
}

// 画布交互
const handleCanvasMouseDown = (e: MouseEvent) => {
  if (draggingEntity.value) return
  isPanning.value = true
  panStart.value = { x: e.clientX, y: e.clientY }
}

const handleCanvasMouseMove = (e: MouseEvent) => {
  if (draggingEntity.value) {
    // 拖拽实体
    const dx = (e.clientX - dragStartPos.value.x) / scale.value
    const dy = (e.clientY - dragStartPos.value.y) / scale.value
    const newX = entityStartPos.value.x + dx
    const newY = entityStartPos.value.y + dy
    updateEntityPosition(draggingEntity.value, newX, newY)
  } else if (isPanning.value) {
    // 平移画布
    const dx = (e.clientX - panStart.value.x) / scale.value
    const dy = (e.clientY - panStart.value.y) / scale.value
    viewBox.value.x -= dx
    viewBox.value.y -= dy
    panStart.value = { x: e.clientX, y: e.clientY }
  }
}

const handleCanvasMouseUp = () => {
  if (draggingEntity.value) {
    // 保存实体位置
    const entity = entities.value.find(e => e.id === draggingEntity.value)
    if (entity) {
      entityPositions.value.set(entity.id, { x: entity.x, y: entity.y })
    }
  }
  draggingEntity.value = null
  isPanning.value = false
}

// 自动布局（重置位置）
const autoLayout = () => {
  entityPositions.value.clear()
  calculateLayout()
  message.success('已重新布局')
}

// 导出PNG
const exportPNG = async () => {
  if (!svgRef.value) return
  
  try {
    const svgData = new XMLSerializer().serializeToString(svgRef.value)
    const svgBlob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' })
    const url = URL.createObjectURL(svgBlob)
    
    const img = new Image()
    img.onload = () => {
      const canvas = document.createElement('canvas')
      canvas.width = canvasWidth.value * 2  // 2x for better quality
      canvas.height = canvasHeight.value * 2
      const ctx = canvas.getContext('2d')!
      ctx.fillStyle = '#f5f5f5'
      ctx.fillRect(0, 0, canvas.width, canvas.height)
      ctx.scale(2, 2)
      ctx.drawImage(img, 0, 0)
      
      canvas.toBlob((blob) => {
        if (blob) {
          const link = document.createElement('a')
          link.download = 'ER-Diagram.png'
          link.href = URL.createObjectURL(blob)
          link.click()
          URL.revokeObjectURL(link.href)
        }
      }, 'image/png')
      
      URL.revokeObjectURL(url)
    }
    img.src = url
    message.success('正在导出PNG...')
  } catch (error) {
    message.error('导出失败')
    console.error(error)
  }
}

// 导出SVG
const exportSVG = () => {
  if (!svgRef.value) return
  
  try {
    const svgData = new XMLSerializer().serializeToString(svgRef.value)
    const blob = new Blob([svgData], { type: 'image/svg+xml;charset=utf-8' })
    const link = document.createElement('a')
    link.download = 'ER-Diagram.svg'
    link.href = URL.createObjectURL(blob)
    link.click()
    URL.revokeObjectURL(link.href)
    message.success('SVG导出成功')
  } catch (error) {
    message.error('导出失败')
    console.error(error)
  }
}

// 监听数据变化
watch(() => [props.entities, props.relations], () => {
  nextTick(() => {
    calculateLayout()
  })
}, { deep: true, immediate: true })

onMounted(() => {
  calculateLayout()
})
</script>

<style scoped>
.er-diagram-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f0f2f5;
  border-radius: 8px;
  overflow: hidden;
}

.er-toolbar {
  padding: 12px 16px;
  background: white;
  border-bottom: 1px solid #e8e8e8;
}

.er-canvas-wrapper {
  flex: 1;
  overflow: hidden;
  cursor: grab;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background: #e8eaed;
}

.er-canvas-wrapper:active {
  cursor: grabbing;
}

.er-svg {
  display: block;
}

.entity-group,
.attr-group,
.relation-group {
  transition: filter 0.2s;
}

.entity-group:hover {
  filter: brightness(1.1);
}

.entity-group.draggable {
  cursor: move;
}

.entity-group.dragging {
  filter: brightness(1.2) drop-shadow(0 0 8px rgba(74, 144, 217, 0.6));
}

.attr-group:hover,
.relation-group:hover {
  filter: brightness(1.05);
}
</style>
