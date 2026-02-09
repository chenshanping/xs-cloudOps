<template>
  <a-modal
    v-model:open="visible"
    :title="title"
    width="520px"
    @ok="handleSubmit"
    @cancel="handleCancel"
    :okButtonProps="{ disabled: !formData.audit_status }"
  >
    <a-form
      ref="formRef"
      :model="formData"
      layout="vertical"
    >
      <!-- 审批结果卡片选择 -->
      <a-form-item 
        label="审批结果" 
        name="audit_status" 
        :rules="[{ required: true, message: '请选择审批结果' }]"
      >
        <div class="audit-cards">
          <div 
            class="audit-card approve" 
            :class="{ active: formData.audit_status === 1 }"
            @click="selectAuditStatus(1)"
          >
            <CheckCircleFilled class="card-icon" />
            <span class="card-title">通过</span>
            <span class="card-desc">审核通过，允许该申请</span>
          </div>
          <div 
            class="audit-card reject" 
            :class="{ active: formData.audit_status === 2 }"
            @click="selectAuditStatus(2)"
          >
            <CloseCircleFilled class="card-icon" />
            <span class="card-title">拒绝</span>
            <span class="card-desc">审核不通过，需要补充资料</span>
          </div>
        </div>
      </a-form-item>

      <!-- 快捷选择 -->
      <a-form-item v-if="formData.audit_status === 1" label="常用审批意见">
        <div class="quick-reasons">
          <a-tag 
            v-for="reason in approveReasons" 
            :key="reason"
            class="reason-tag approve"
            @click="selectReason(reason)"
          >
            {{ reason }}
          </a-tag>
        </div>
      </a-form-item>
      <a-form-item v-if="formData.audit_status === 2" label="常用拒绝原因">
        <div class="quick-reasons">
          <a-tag 
            v-for="reason in rejectReasons" 
            :key="reason"
            class="reason-tag reject"
            @click="selectReason(reason)"
          >
            {{ reason }}
          </a-tag>
        </div>
      </a-form-item>
      
      <a-form-item 
        :label="formData.audit_status === 2 ? '拒绝原因' : '审批备注'" 
        name="audit_remark" 
        :rules="remarkRules"
      >
        <a-textarea
          v-model:value="formData.audit_remark"
          :rows="3"
:placeholder="formData.audit_status === 2 ? '请输入拒绝原因，便于申请人了解具体问题' : '请填写审批备注'"
          :maxlength="500"
          show-count
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { CheckCircleFilled, CloseCircleFilled } from '@ant-design/icons-vue'
import type { FormInstance } from 'ant-design-vue'

interface Props {
  open: boolean
  title?: string
}

interface Emits {
  (e: 'update:open', value: boolean): void
  (e: 'confirm', data: { audit_status: number; audit_remark: string }): void
}

const props = withDefaults(defineProps<Props>(), {
  title: '审批'
})

const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()
const visible = ref(false)

const formData = reactive({
  audit_status: undefined as number | undefined,
  audit_remark: ''
})

// 常用通过原因
const approveReasons = [
  '资料齐全，审核通过',
  '符合要求，同意通过',
  '信息无误，予以通过',
  '审核无异常'
]

// 常用拒绝原因
const rejectReasons = [
  '资料不完整',
  '证件照片不清晰',
  '信息填写有误',
  '证件已过期',
  '不符合申请条件'
]

// 审批备注验证规则：都必填
const remarkRules = computed(() => {
  if (formData.audit_status === 2) {
    return [{ required: true, message: '拒绝时必须填写原因' }]
  }
  return [{ required: true, message: '请填写审批备注' }]
})

// 选择审批状态
const selectAuditStatus = (status: number) => {
  formData.audit_status = status
  // 切换时清空备注
  formData.audit_remark = ''
}

// 选择快捷拒绝原因
const selectReason = (reason: string) => {
  if (formData.audit_remark) {
    formData.audit_remark += '；' + reason
  } else {
    formData.audit_remark = reason
  }
}

// 监听 open 变化
watch(() => props.open, (val) => {
  visible.value = val
  if (val) {
    // 打开时重置表单
    formData.audit_status = undefined
    formData.audit_remark = ''
    formRef.value?.clearValidate()
  }
})

// 监听 visible 变化，同步回父组件
watch(visible, (val) => {
  emit('update:open', val)
})

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    emit('confirm', {
      audit_status: formData.audit_status!,
      audit_remark: formData.audit_remark
    })
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const handleCancel = () => {
  visible.value = false
}
</script>

<style scoped lang="less">
.audit-cards {
  display: flex;
  gap: 16px;
  
  .audit-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px 16px;
    border: 2px solid #e8e8e8;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;
    
    &:hover {
      border-color: #d9d9d9;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
    }
    
    .card-icon {
      font-size: 36px;
      margin-bottom: 8px;
    }
    
    .card-title {
      font-size: 16px;
      font-weight: 600;
      margin-bottom: 4px;
    }
    
    .card-desc {
      font-size: 12px;
      color: #999;
      text-align: center;
    }
    
    &.approve {
      .card-icon {
        color: #d9d9d9;
      }
      .card-title {
        color: #666;
      }
      
      &.active {
        border-color: #52c41a;
        background: #f6ffed;
        
        .card-icon {
          color: #52c41a;
        }
        .card-title {
          color: #52c41a;
        }
      }
    }
    
    &.reject {
      .card-icon {
        color: #d9d9d9;
      }
      .card-title {
        color: #666;
      }
      
      &.active {
        border-color: #ff4d4f;
        background: #fff2f0;
        
        .card-icon {
          color: #ff4d4f;
        }
        .card-title {
          color: #ff4d4f;
        }
      }
    }
  }
}

.quick-reasons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  
  .reason-tag {
    cursor: pointer;
    border-style: dashed;
    
    &.approve:hover {
      color: #52c41a;
      border-color: #52c41a;
    }
    
    &.reject:hover {
      color: #ff4d4f;
      border-color: #ff4d4f;
    }
  }
}
</style>
