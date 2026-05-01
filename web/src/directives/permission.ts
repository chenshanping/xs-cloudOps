import type { Directive, DirectiveBinding, WatchStopHandle } from 'vue'
import { watch } from 'vue'
import { useUserStore } from '@/store/user'

/**
 * 权限指令
 * 使用方式：
 * v-permission="'system:user:edit'"  - 单个权限
 * v-permission="['system:user:edit', 'system:user:delete']" - 多个权限（满足任意一个即可）
 */
type PermissionValue = string | string[]

type PermissionElement = HTMLElement & {
  __permissionOriginalDisplay__?: string
  __permissionWatcherStop__?: WatchStopHandle
}

const resolvePermissionVisible = (value: PermissionValue | undefined, userStore: ReturnType<typeof useUserStore>) => {
  if (!value) {
    return true
  }
  if (Array.isArray(value)) {
    return userStore.hasAnyPermission(value)
  }
  return userStore.hasPermission(value)
}

const applyPermissionVisibility = (
  el: PermissionElement,
  binding: DirectiveBinding<PermissionValue>,
  userStore: ReturnType<typeof useUserStore>
) => {
  const hasPermission = resolvePermissionVisible(binding.value, userStore)
  if (hasPermission) {
    el.style.display = el.__permissionOriginalDisplay__ ?? ''
    el.removeAttribute('aria-hidden')
    return
  }

  el.style.display = 'none'
  el.setAttribute('aria-hidden', 'true')
}

export const permission: Directive = {
  mounted(el: PermissionElement, binding: DirectiveBinding<PermissionValue>) {
    const userStore = useUserStore()
    el.__permissionOriginalDisplay__ = el.style.display
    el.__permissionWatcherStop__ = watch(
      () => userStore.permissions.slice(),
      () => applyPermissionVisibility(el, binding, userStore),
      { immediate: true }
    )
  },
  updated(el: PermissionElement, binding: DirectiveBinding<PermissionValue>) {
    const userStore = useUserStore()
    applyPermissionVisibility(el, binding, userStore)
  },
  unmounted(el: PermissionElement) {
    el.__permissionWatcherStop__?.()
    delete el.__permissionWatcherStop__
    delete el.__permissionOriginalDisplay__
  }
}

export default permission
