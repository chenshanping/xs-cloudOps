import type { Directive, DirectiveBinding, WatchStopHandle } from 'vue'
import { watch } from 'vue'
import { useUserStore } from '@/store/user'
import type { PermissionCheckMode, PermissionDirectiveValue } from '@/utils/permission'

/**
 * 权限指令
 * 使用方式：
 * v-permission="'system:user:edit'"  - 单个权限
 * v-permission="['system:user:edit', 'system:user:delete']" - 多个权限（满足任意一个即可）
 */
type PermissionValue = PermissionDirectiveValue

type PermissionElement = HTMLElement & {
  __permissionOriginalCursor__?: string
  __permissionOriginalDisplay__?: string
  __permissionOriginalOpacity__?: string
  __permissionOriginalPointerEvents__?: string
  __permissionOriginalTabIndex__?: string | null
  __permissionOriginalTitle__?: string | null
  __permissionWatcherStop__?: WatchStopHandle
}

const resolvePermissionRule = (value: PermissionValue | undefined) => {
  if (!value) {
    return {
      allOf: [] as string[],
      anyOf: [] as string[],
      mode: 'hide' as PermissionCheckMode,
      reason: '',
    }
  }

  if (typeof value === 'string') {
    return {
      allOf: [] as string[],
      anyOf: [value],
      mode: 'hide' as PermissionCheckMode,
      reason: '',
    }
  }

  if (Array.isArray(value)) {
    return {
      allOf: [] as string[],
      anyOf: value,
      mode: 'hide' as PermissionCheckMode,
      reason: '',
    }
  }

  const inlinePermissions = value.permission
    ? Array.isArray(value.permission)
      ? value.permission
      : [value.permission]
    : []

  return {
    allOf: value.allOf || [],
    anyOf: value.anyOf?.length ? value.anyOf : inlinePermissions,
    mode: value.mode || 'hide',
    reason: value.reason || '',
  }
}

const resolvePermissionVisible = (value: PermissionValue | undefined, userStore: ReturnType<typeof useUserStore>) => {
  const rule = resolvePermissionRule(value)
  const allMatched = rule.allOf.every(permission => userStore.hasPermission(permission))
  const anyMatched = rule.anyOf.length === 0 || userStore.hasAnyPermission(rule.anyOf)

  return {
    mode: rule.mode,
    reason: rule.reason,
    visible: allMatched && anyMatched,
  }
}

const restoreDisabledState = (el: PermissionElement) => {
  el.style.pointerEvents = el.__permissionOriginalPointerEvents__ ?? ''
  el.style.opacity = el.__permissionOriginalOpacity__ ?? ''
  el.style.cursor = el.__permissionOriginalCursor__ ?? ''
  el.removeAttribute('aria-disabled')

  if (el.__permissionOriginalTabIndex__ === null) {
    el.removeAttribute('tabindex')
  } else if (typeof el.__permissionOriginalTabIndex__ === 'string') {
    el.setAttribute('tabindex', el.__permissionOriginalTabIndex__)
  }

  if (el.__permissionOriginalTitle__ === null) {
    el.removeAttribute('title')
  } else if (typeof el.__permissionOriginalTitle__ === 'string') {
    el.setAttribute('title', el.__permissionOriginalTitle__)
  }
}

const applyDisabledState = (el: PermissionElement, reason: string) => {
  el.style.display = el.__permissionOriginalDisplay__ ?? ''
  el.removeAttribute('aria-hidden')
  el.setAttribute('aria-disabled', 'true')
  el.style.pointerEvents = 'none'
  el.style.opacity = '0.56'
  el.style.cursor = 'not-allowed'
  el.setAttribute('tabindex', '-1')

  if (reason) {
    el.setAttribute('title', reason)
  }
}

const applyPermissionVisibility = (
  el: PermissionElement,
  binding: DirectiveBinding<PermissionValue>,
  userStore: ReturnType<typeof useUserStore>
) => {
  const permissionState = resolvePermissionVisible(binding.value, userStore)
  if (permissionState.visible) {
    el.style.display = el.__permissionOriginalDisplay__ ?? ''
    el.removeAttribute('aria-hidden')
    restoreDisabledState(el)
    return
  }

  if (permissionState.mode === 'disable') {
    applyDisabledState(el, permissionState.reason)
    return
  }

  restoreDisabledState(el)
  el.style.display = 'none'
  el.setAttribute('aria-hidden', 'true')
}

export const permission: Directive = {
  mounted(el: PermissionElement, binding: DirectiveBinding<PermissionValue>) {
    const userStore = useUserStore()
    el.__permissionOriginalCursor__ = el.style.cursor
    el.__permissionOriginalDisplay__ = el.style.display
    el.__permissionOriginalOpacity__ = el.style.opacity
    el.__permissionOriginalPointerEvents__ = el.style.pointerEvents
    el.__permissionOriginalTabIndex__ = el.getAttribute('tabindex')
    el.__permissionOriginalTitle__ = el.getAttribute('title')
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
    delete el.__permissionOriginalCursor__
    delete el.__permissionWatcherStop__
    delete el.__permissionOriginalDisplay__
    delete el.__permissionOriginalOpacity__
    delete el.__permissionOriginalPointerEvents__
    delete el.__permissionOriginalTabIndex__
    delete el.__permissionOriginalTitle__
  }
}

export default permission
