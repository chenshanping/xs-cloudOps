import { login, getUserInfo, type LoginParams, logout } from '@/api/auth';
import type { User, Menu } from '@/types';
import { resetRouter } from '@/router';
import { defineStore } from 'pinia';
import { ref } from 'vue';


export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)
  const menus = ref<Menu[]>([])
  const permissions = ref<string[]>([])

  // 检查是否有某个权限
  function hasPermission(permission: string): boolean {
    // 超级管理员拥有所有权限
    if (permissions.value.includes('*')) {
      return true
    }
    return permissions.value.includes(permission)
  }

  // 检查是否有多个权限中的任意一个
  function hasAnyPermission(perms: string[]): boolean {
    if (permissions.value.includes('*')) {
      return true
    }
    return perms.some(p => permissions.value.includes(p))
  }

  // 登录
  async function loginAction(data: LoginParams) {
    const res = await login(data)
    token.value = res.data.token
    localStorage.setItem('token', res.data.token)
    return res
  }

  // 获取用户信息
  async function getUserInfoAction() {
    const res = await getUserInfo()
    user.value = res.data.user
    menus.value = res.data.menus
    permissions.value = res.data.permissions || []
    return res
  }

  // 登出
  async function logoutAction() {
    try {
      await logout()  // 先调用 API（此时 token 还在）
    } catch (e) {
      console.error(e)
    }
    token.value = ''
    user.value = null
    menus.value = []
    permissions.value = []
   
    localStorage.removeItem('token')
    resetRouter()
    
  }

  return {
    token,
    user,
    menus,
    permissions,
    hasPermission,
    hasAnyPermission,
    loginAction,
    getUserInfoAction,
    logoutAction,
  }
})
