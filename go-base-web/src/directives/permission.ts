import type { Directive, DirectiveBinding } from 'vue'
import { useUserStore } from '@/store/user'

/**
 * 权限指令
 * 使用方式：
 * v-permission="'system:user:edit'"  - 单个权限
 * v-permission="['system:user:edit', 'system:user:delete']" - 多个权限（满足任意一个即可）
 */
export const permission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const userStore = useUserStore()
    const value = binding.value

    if (!value) return

    let hasPermission = false

    if (Array.isArray(value)) {
      hasPermission = userStore.hasAnyPermission(value)
    } else {
      hasPermission = userStore.hasPermission(value)
    }

    if (!hasPermission) {
      // 移除元素
      el.parentNode?.removeChild(el)
    }
  }
}

export default permission
