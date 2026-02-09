/**
 * 心理题目 AI 提示词工具
 * 
 * 使用说明：
 * 1. 导入需要的函数
 * 2. 传入参数生成提示词
 * 3. 调用 AI 接口获取结果
 * 4. 使用解析函数处理返回数据
 */

// ==================== 类型定义 ====================

/** 题目类型 */
export type QuestionType = 1 | 2 | 3  // 1:单选 2:多选 3:量表题

/** 选项结构 */
export interface OptionItem {
  label: string      // 选项标签 A/B/C/D
  content: string    // 选项内容
  score?: number     // 分数（量表题使用）
}

/** 单题生成结果 */
export interface SingleQuestionResult {
  options: OptionItem[]
  answer: string | string[]  // 单选为字符串，多选为数组，量表题为空
  analysis: string
}

/** 完整题目结构（批量生成使用） */
export interface FullQuestionResult {
  content: string           // 题干
  options: OptionItem[]
  answer: string | string[]
  analysis: string
}

/** 单题生成参数 */
export interface SinglePromptParams {
  type: QuestionType        // 题目类型
  categoryName: string      // 分类名称
  content: string           // 题目内容（题干）
}

/** 批量生成参数 */
export interface BatchPromptParams {
  type: QuestionType        // 题目类型
  categoryName: string      // 分类名称
  count: number             // 生成数量
  topic?: string            // 主题描述（可选）
  existingQuestions?: string[]  // 已有题目列表（避免重复）
}

// ==================== 类型映射 ====================

/** 题目类型名称映射 */
export const questionTypeNames: Record<QuestionType, string> = {
  1: '单选题',
  2: '多选题',
  3: '量表题（心理测评）'
}

/** 题目类型颜色映射 */
export const questionTypeColors: Record<QuestionType, string> = {
  1: 'blue',
  2: 'purple',
  3: 'green'
}

// ==================== 提示词模板 ====================

/**
 * 量表题选项模板说明
 * 用于提示 AI 生成符合心理测评规范的选项
 */
const SCALE_OPTIONS_GUIDE = `
常见的心理测评量表选项模式：
- 频率类：几乎每天(10分)、大部分日子(7分)、很少(4分)、完全没有(0分)
- 程度类：非常符合(10分)、比较符合(7分)、有点符合(4分)、完全不符合(0分)
- 同意程度：非常同意(10分)、同意(7分)、不同意(4分)、非常不同意(0分)
- 频率类2：总是(10分)、经常(7分)、有时(4分)、很少(2分)、从不(0分)`

/**
 * JSON 输出格式模板 - 单题（单选/多选）
 */
const SINGLE_CHOICE_JSON_TEMPLATE = `{
  "options": [
    {"label": "A", "content": "选项A内容"},
    {"label": "B", "content": "选项B内容"},
    {"label": "C", "content": "选项C内容"},
    {"label": "D", "content": "选项D内容"}
  ],
  "answer": "A",
  "analysis": "详细的答案解析，支持Markdown格式"
}`

/**
 * JSON 输出格式模板 - 单题（多选）
 */
const MULTI_CHOICE_JSON_TEMPLATE = `{
  "options": [
    {"label": "A", "content": "选项A内容"},
    {"label": "B", "content": "选项B内容"},
    {"label": "C", "content": "选项C内容"},
    {"label": "D", "content": "选项D内容"}
  ],
  "answer": ["A", "B"],
  "analysis": "详细的答案解析，支持Markdown格式"
}`

/**
 * JSON 输出格式模板 - 单题（量表题）
 */
const SCALE_JSON_TEMPLATE = `{
  "options": [
    {"label": "A", "content": "几乎每天", "score": 10},
    {"label": "B", "content": "大部分日子", "score": 7},
    {"label": "C", "content": "很少", "score": 4},
    {"label": "D", "content": "完全没有", "score": 0}
  ],
  "analysis": "评分说明，包括各分数段含义和建议"
}`

/**
 * JSON 输出格式模板 - 批量（单选/多选）
 */
const BATCH_CHOICE_JSON_TEMPLATE = `[
  {
    "content": "题目内容1",
    "options": [
      {"label": "A", "content": "选项内容"},
      {"label": "B", "content": "选项内容"},
      {"label": "C", "content": "选项内容"},
      {"label": "D", "content": "选项内容"}
    ],
    "answer": "A",
    "analysis": "答案解析"
  }
]`

/**
 * JSON 输出格式模板 - 批量（量表题）
 */
const BATCH_SCALE_JSON_TEMPLATE = `[
  {
    "content": "题目内容1",
    "options": [
      {"label": "A", "content": "几乎每天", "score": 10},
      {"label": "B", "content": "大部分日子", "score": 7},
      {"label": "C", "content": "很少", "score": 4},
      {"label": "D", "content": "完全没有", "score": 0}
    ],
    "analysis": "评分说明"
  }
]`

// ==================== 提示词生成函数 ====================

/**
 * 生成单题提示词（用于已有题干，生成选项和答案）
 * 
 * @param params - 生成参数
 * @returns 提示词字符串
 * 
 * @example
 * ```ts
 * const prompt = buildSingleQuestionPrompt({
 *   type: 1,
 *   categoryName: '认知心理学',
 *   content: '以下哪个不属于记忆的三个阶段？'
 * })
 * ```
 */
export function buildSingleQuestionPrompt(params: SinglePromptParams): string {
  const { type, categoryName, content } = params
  const typeName = questionTypeNames[type]

  if (type === 3) {
    // 量表题
    return `你是一个专业的心理测评题目设计专家。请根据以下题目内容，生成适合心理测评量表的选项和评分说明。

题目分类：${categoryName || '心理测评'}
题目内容：${content}

要求：
1. 量表题没有正确答案，每个选项都有对应的分数
2. 选项通常表示频率或程度
3. 分数通常从高到低或从低到高，根据题目语义决定
4. 生成4个选项
5. 评分说明要包含各分数段的含义解读
${SCALE_OPTIONS_GUIDE}

请直接返回JSON格式，不要包含任何其他文字或markdown标记：
${SCALE_JSON_TEMPLATE}`
  }

  // 单选/多选题
  const optionCount = type === 2 ? '4-6个' : '4个'
  const answerDesc = type === 2 ? '数组格式，包含所有正确选项' : '单个字母'
  const jsonTemplate = type === 2 ? MULTI_CHOICE_JSON_TEMPLATE : SINGLE_CHOICE_JSON_TEMPLATE

  return `你是一个专业的心理学知识考试题目出题专家。请根据以下信息生成题目的选项、答案和解析。

题目分类：${categoryName || '心理学'}
题目类型：${typeName}
题目内容：${content}

要求：
1. 生成${optionCount}选项，标签为A、B、C、D等
2. 选项内容要专业、准确、有迷惑性
3. 正确答案要准确
4. 答案解析要详细、有教育意义，包括：
   - 正确答案的解释
   - 其他选项为何错误
   - 相关的心理学理论或概念

请直接返回JSON格式，不要包含任何其他文字或markdown标记：
${jsonTemplate}`
}

/**
 * 生成批量题目提示词（AI生成完整题目）
 * 
 * @param params - 生成参数
 * @returns 提示词字符串
 * 
 * @example
 * ```ts
 * const prompt = buildBatchQuestionPrompt({
 *   type: 3,
 *   categoryName: '抑郁量表',
 *   count: 5,
 *   topic: 'PHQ-9抑郁筛查',
 *   existingQuestions: ['在过去两周内，你是否感到情绪低落？']
 * })
 * ```
 */
export function buildBatchQuestionPrompt(params: BatchPromptParams): string {
  const { type, categoryName, count, topic, existingQuestions = [] } = params
  const typeName = questionTypeNames[type]

  // 构建已有题目提示（避免重复）
  let existingTip = ''
  if (existingQuestions.length > 0) {
    const maxShow = 20
    const existingList = existingQuestions
      .slice(0, maxShow)
      .map((q, i) => `${i + 1}. ${q}`)
      .join('\n')
    const moreText = existingQuestions.length > maxShow 
      ? `\n...等共${existingQuestions.length}道题目` 
      : ''
    existingTip = `

【重要】以下是该分类下已有的题目，请勿生成相同或相似的题目：
${existingList}${moreText}`
  }

  if (type === 3) {
    // 量表题批量生成
    return `你是一个专业的心理测评题目设计专家。请生成${count}道心理测评量表题目。

分类：${categoryName || '心理测评'}
主题：${topic || '心理健康评估'}${existingTip}

要求：
1. 每道题目都是独立的心理测评问题
2. 题目内容要专业、科学、适合心理测评
3. 量表题没有正确答案，每个选项有对应分数
4. 选项表示频率或程度
5. 每题4个选项，分数从高到低或从低到高
6. 不要生成与已有题目相同或类似的内容
${SCALE_OPTIONS_GUIDE}

请直接返回JSON数组格式，不要包含任何其他文字：
${BATCH_SCALE_JSON_TEMPLATE}`
  }

  // 单选/多选题批量生成
  const answerFormat = type === 2 ? '["A", "B"]数组格式' : '"A"单个字母'

  return `你是一个专业的心理学考试出题专家。请生成${count}道${typeName}。

分类：${categoryName || '心理学'}
主题：${topic || '心理学知识'}${existingTip}

要求：
1. 每道题目都是独立的、有价值的考试题目
2. 题目内容要专业、准确、有区分度
3. 选项要专业、准确、有迷惑性
4. 每题4个选项(A/B/C/D)
5. 答案格式：${answerFormat}
6. 解析要简明扼要
7. 不要生成与已有题目相同或类似的内容

请直接返回JSON数组格式，不要包含任何其他文字：
${BATCH_CHOICE_JSON_TEMPLATE}`
}

// ==================== 响应解析函数 ====================

/**
 * 从 AI 响应中提取 JSON 字符串
 * 
 * @param content - AI 返回的原始内容
 * @returns 提取的 JSON 字符串
 */
function extractJSON(content: string): string {
  // 尝试提取 markdown 代码块中的 JSON
  const jsonMatch = content.match(/```(?:json)?\s*([\s\S]*?)```/)
  if (jsonMatch) {
    return jsonMatch[1].trim()
  }

  // 去除前后空白，找到第一个非空字符
  const trimmed = content.trim()
  if (!trimmed) return content
  
  const firstChar = trimmed[0]
  
  // 如果以 [ 开头，提取数组
  if (firstChar === '[') {
    const arrayStart = trimmed.indexOf('[')
    const arrayEnd = trimmed.lastIndexOf(']')
    if (arrayStart !== -1 && arrayEnd !== -1 && arrayEnd > arrayStart) {
      return trimmed.slice(arrayStart, arrayEnd + 1)
    }
  }
  
  // 如果以 { 开头，提取对象
  if (firstChar === '{') {
    const objectStart = trimmed.indexOf('{')
    const objectEnd = trimmed.lastIndexOf('}')
    if (objectStart !== -1 && objectEnd !== -1 && objectEnd > objectStart) {
      return trimmed.slice(objectStart, objectEnd + 1)
    }
  }

  return content
}

/**
 * 解析单题生成结果
 * 
 * @param content - AI 返回的原始内容
 * @param type - 题目类型
 * @returns 解析后的结果
 * @throws 解析失败时抛出错误
 * 
 * @example
 * ```ts
 * try {
 *   const result = parseSingleQuestionResponse(aiResponse, 1)
 *   console.log(result.options, result.answer, result.analysis)
 * } catch (e) {
 *   console.error('解析失败:', e)
 * }
 * ```
 */
export function parseSingleQuestionResponse(
  content: string, 
  type: QuestionType
): SingleQuestionResult {
  const jsonStr = extractJSON(content)
  const parsed = JSON.parse(jsonStr)

  const result: SingleQuestionResult = {
    options: [],
    answer: type === 2 ? [] : '',
    analysis: parsed.analysis || ''
  }

  // 处理选项
  if (Array.isArray(parsed.options)) {
    result.options = parsed.options.map((opt: any, index: number) => ({
      label: opt.label || String.fromCharCode(65 + index),
      content: opt.content || '',
      score: opt.score ?? 0
    }))
  }

  // 处理答案（量表题无答案）
  if (type === 1) {
    result.answer = parsed.answer || ''
  } else if (type === 2) {
    if (Array.isArray(parsed.answer)) {
      result.answer = parsed.answer
    } else if (typeof parsed.answer === 'string') {
      result.answer = parsed.answer.split(',').map((s: string) => s.trim())
    }
  }
  // type === 3 量表题无答案，保持空

  return result
}

/**
 * 解析批量生成结果
 * 
 * @param content - AI 返回的原始内容
 * @param type - 题目类型
 * @returns 解析后的题目数组
 * @throws 解析失败时抛出错误
 * 
 * @example
 * ```ts
 * try {
 *   const questions = parseBatchQuestionResponse(aiResponse, 3)
 *   questions.forEach(q => {
 *     console.log(q.content, q.options)
 *   })
 * } catch (e) {
 *   console.error('解析失败:', e)
 * }
 * ```
 */
export function parseBatchQuestionResponse(
  content: string, 
  type: QuestionType
): FullQuestionResult[] {
  const jsonStr = extractJSON(content)
  const parsed = JSON.parse(jsonStr)

  if (!Array.isArray(parsed)) {
    throw new Error('返回格式错误：期望数组')
  }

  return parsed.map((item: any) => ({
    content: item.content || '',
    options: (item.options || []).map((opt: any, idx: number) => ({
      label: opt.label || String.fromCharCode(65 + idx),
      content: opt.content || '',
      score: opt.score ?? 0
    })),
    answer: type === 3 ? '' : (item.answer || ''),
    analysis: item.analysis || ''
  }))
}

// ==================== 工具函数 ====================

/**
 * 获取默认量表题选项
 * 用于初始化表单或重置
 */
export function getDefaultScaleOptions(): OptionItem[] {
  return [
    { label: 'A', content: '几乎每天', score: 10 },
    { label: 'B', content: '大部分日子', score: 7 },
    { label: 'C', content: '很少', score: 4 },
    { label: 'D', content: '完全没有', score: 0 },
  ]
}

/**
 * 获取默认普通选项
 * 用于初始化表单或重置
 */
export function getDefaultOptions(): OptionItem[] {
  return [
    { label: 'A', content: '' },
    { label: 'B', content: '' },
    { label: 'C', content: '' },
    { label: 'D', content: '' },
  ]
}

/**
 * 选项标签列表
 */
export const OPTION_LABELS = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H']

/**
 * 格式化答案显示
 * 
 * @param answer - 答案（字符串或数组）
 * @returns 格式化后的字符串
 */
export function formatAnswer(answer: string | string[] | undefined): string {
  if (!answer) return '-'
  if (Array.isArray(answer)) {
    return answer.join(', ')
  }
  return answer
}

/**
 * 判断选项是否为正确答案
 * 
 * @param answer - 答案
 * @param label - 选项标签
 * @returns 是否为正确答案
 */
export function isCorrectAnswer(answer: string | string[] | undefined, label: string): boolean {
  if (!answer) return false
  if (Array.isArray(answer)) {
    return answer.includes(label)
  }
  return answer === label
}
