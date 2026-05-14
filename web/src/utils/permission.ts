import { computed, type ComputedRef } from 'vue'
import { useUserStore } from '@/store/user'

export type PermissionCheckMode = 'disable' | 'hide'

export interface PermissionDirectiveConfig {
  allOf?: string[]
  anyOf?: string[]
  mode?: PermissionCheckMode
  permission?: string | string[]
  reason?: string
}

export type PermissionDirectiveValue = string | string[] | PermissionDirectiveConfig

export const withDisabledPermission = (
  permission: string | string[],
  reason = '当前账号无权限执行该操作'
): PermissionDirectiveConfig => ({
  permission,
  mode: 'disable',
  reason,
})

export const withHiddenPermission = (permission: string | string[]): PermissionDirectiveConfig => ({
  permission,
  mode: 'hide',
})

/**
 * 权限相关工具函数
 */
export function usePermission() {
  const userStore = useUserStore()

  /**
   * 检查是否有某个权限
   */
  const hasPermission = (permission: string): boolean => {
    return userStore.hasPermission(permission)
  }

  /**
   * 检查是否有多个权限中的任意一个
   */
  const hasAnyPermission = (permissions: string[]): boolean => {
    return userStore.hasAnyPermission(permissions)
  }

  /**
   * 检查是否有全部权限
   */
  const hasAllPermissions = (permissions: string[]): boolean => {
    return permissions.every(p => userStore.hasPermission(p))
  }

  return {
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
  }
}

/**
 * 表格列配置接口
 */
export interface TableColumn {
  title: string
  dataIndex?: string
  key: string
  width?: number
  [key: string]: any
}

/**
 * 动态表格列 - 根据权限显示/隐藏操作列
 * @param baseColumns 基础列配置（不含操作列）
 * @param actionColumn 操作列配置
 * @param actionPermissions 操作列需要的权限（有任意一个就显示操作列）
 */
export function useTableColumns(
  baseColumns: TableColumn[],
  actionColumn: TableColumn | null,
  actionPermissions: string[] = []
): ComputedRef<TableColumn[]> {
  const userStore = useUserStore()

  return computed(() => {
    const columns = [...baseColumns]
    
    // 如果有操作列配置，且有任意操作权限，则添加操作列
    if (actionColumn && actionPermissions.length > 0) {
      if (userStore.hasAnyPermission(actionPermissions)) {
        columns.push(actionColumn)
      }
    } else if (actionColumn) {
      // 没有权限限制，直接添加操作列
      columns.push(actionColumn)
    }
    
    return columns
  })
}
