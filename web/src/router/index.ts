import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import { useConfigStore } from '@/store/config';
import { useUserStore } from '@/store/user';
import { message } from 'ant-design-vue';
import type { Menu } from '@/types';
import { resolveViewModulePath } from './view-resolver';


// 动态导入视图组件的映射
const viewModules = import.meta.glob('../views/**/*.vue')

// 前台路由（用户端）
const frontRoutes: RouteRecordRaw[] = [
  {
    path: '/front',
    name: 'FrontLayout',
    component: () => import('@/layouts/FrontLayout.vue'),
    redirect: '/front/home',
    children: [
      {
        path: 'home',
        name: 'FrontHome',
        component: () => import('@/views/front/home/index.vue'),
        meta: { title: '首页', icon: 'HomeOutlined', showInNav: true }
      },
      {
        path: 'ai',
        name: 'FrontAI',
        component: () => import('@/views/front/ai/index.vue'),
        meta: { title: 'AI对话', icon: 'svg:aiChat', showInNav: true, requiresAuth: true }
      },
      {
        path: 'test',
        name: 'FrontTest',
        component: () => import('@/views/front/test/index.vue'),
        meta: { title: '测试页面', icon: 'MessageOutlined', showInNav: true, requiresAuth: true }
      },
      {
        path: 'profile',
        name: 'FrontProfile',
        component: () => import('@/views/front/profile/index.vue'),
        meta: { title: '个人中心', requiresAuth: true, showInNav: false }
      }
    ]
  }
]

// 后台路由（管理员端）
const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/login/index.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/register/index.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/forgot-password/index.vue'),
    meta: { title: '忘记密码' }
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/views/auth/reset-password/index.vue'),
    meta: { title: '重置密码' }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layouts/BasicLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/admin/dashboard/index.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/admin/profile/index.vue'),
        meta: { title: '个人中心' }
      },
      {
        path: 'ai',
        name: 'AIChat',
        component: () => import('@/views/admin/ai/index.vue'),
        meta: { title: 'AI对话' }
      },
      {
        path: 'system/storage',
        name: 'StorageRedirect',
        redirect: () => ({
          path: '/system/config',
          query: { tab: 'storage' }
        }),
        meta: { title: '存储设置' }
      }
    ]
  }
]

// 404路由（需要在动态路由添加后再添加）
const notFoundRoute: RouteRecordRaw = {
  path: '/:pathMatch(.*)*',
  name: 'NotFound',
  component: () => import('@/views/error/404.vue'),
  meta: { title: '404' }
}

const router = createRouter({
  history: createWebHistory(),
  routes: [...frontRoutes, ...constantRoutes]
})

// 已添加的动态路由名称（用于避免重复添加）
let dynamicRoutesAdded = false

// 根据菜单生成路由
function generateRoutes(menus: Menu[]): RouteRecordRaw[] {
  const routes: RouteRecordRaw[] = []
  
  const processMenu = (menu: Menu) => {
    // 只处理页面类型的菜单（type=2）
    if (menu.type === 2 && menu.component && menu.path) {
      const componentPath = resolveViewModulePath(menu.component)
      
      // 查找组件
      const component = componentPath ? viewModules[componentPath] : undefined
      
      if (component) {
        routes.push({
          path: menu.path.startsWith('/') ? menu.path.slice(1) : menu.path,
          name: `Route_${menu.id}`,
          component: component,
          meta: {
            title: menu.name,
            icon: menu.icon,
            permission: menu.permission
          }
        })
      }
    }
    
    // 递归处理子菜单
    if (menu.children?.length) {
      menu.children.forEach(processMenu)
    }
  }
  
  menus.forEach(processMenu)
  return routes
}

// 添加动态路由
export function addDynamicRoutes(menus: Menu[] | null) {
  if (dynamicRoutesAdded) return
  
  const dynamicRoutes = generateRoutes(menus || [])
  
  // 确保每个路由都能正确添加
  try {
    dynamicRoutes.forEach(route => {
      // 检查路由是否已存在，避免重复添加
      if (!router.hasRoute(route.name as string)) {
        router.addRoute('Layout', route)
      }
    })
    
    // 最后添加404路由（确保它在所有路由之后）
    if (!router.hasRoute('NotFound')) {
      router.addRoute(notFoundRoute)
    }
    
    dynamicRoutesAdded = true
  } catch (error) {
    console.error('添加动态路由失败', error)
  }
}

// 重置路由（登出时调用）
export function resetRouter() {
  dynamicRoutesAdded = false
  // 获取所有路由，移除动态添加的
  const routes = router.getRoutes()
  routes.forEach(route => {
    if (route.name && (String(route.name).startsWith('Route_') || route.name === 'NotFound')) {
      router.removeRoute(route.name)
    }
  })
}

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const configStore = useConfigStore()
  
  // 判断是否是前台路由
  const isFrontRoute = to.path.startsWith('/front')
  
  // 公开页面直接放行
  const publicPaths = ['/login', '/register', '/forgot-password', '/reset-password']
  if (publicPaths.includes(to.path)) {
    next()
    return
  }
  
  // 前台路由处理
  if (isFrontRoute) {
    // 检查是否需要登录
    if (to.meta?.requiresAuth && !userStore.token) {
      next('/login')
      return
    }
    
    // 已登录但未加载用户信息
    if (userStore.token && !userStore.user) {
      try {
        await userStore.getUserInfoAction()
        await configStore.loadConfigs()
      } catch (error) {
        // 前台页面加载失败不强制跳转，继续访问
        console.error('加载用户信息失败', error)
      }
    }
    
    // 检查前台模式：如果是 profile 模式，只允许访问个人中心
    const frontMode = configStore.get('front_mode')
    if (frontMode === 'profile' && to.name !== 'FrontProfile') {
      // profile 模式下，其他前台页面重定向到个人中心
      next('/front/profile')
      return
    }
    
    next()
    return
  }

  // 后台路由处理
  if (!userStore.token) {
    next('/login')
    return
  }

  if (!userStore.user) {
    try {
      await userStore.getUserInfoAction()
      // 加载系统配置
      await configStore.loadConfigs()
      // 添加动态路由
      addDynamicRoutes(userStore.menus)
      
      // 检查是否有后台菜单权限（普通用户没有菜单，重定向到前台）
      if (!userStore.menus || userStore.menus.length === 0) {
        const frontMode = configStore.get('front_mode')
        // 根据前台模式决定跳转目标
        next(frontMode === 'profile' ? '/front/profile' : '/front')
        return
      }
      
      // 重新导航到目标路由（确保动态路由已添加）
      // 使用 nextTick 确保路由已完全注册
      next({ ...to, replace: true })
      return
    } catch (error) {
      console.error('加载用户信息失败', error)
      await userStore.logoutAction()
      next('/login')
      return
    }
  }
  
  // 确保动态路由已添加（防止某些情况下路由丢失）
  if (!dynamicRoutesAdded && userStore.menus && userStore.menus.length > 0) {
    addDynamicRoutes(userStore.menus)
    next({ ...to, replace: true })
    return
  }

  // 已加载用户信息，检查后台权限
  if (!userStore.menus || userStore.menus.length === 0) {
    const frontMode = configStore.get('front_mode')
    next(frontMode === 'profile' ? '/front/profile' : '/front')
    return
  }

  next()
})

// 路由后置守卫 - 更新页面标题
router.afterEach((to) => {
  // 延迟执行，确保配置已加载
  setTimeout(() => {
    const configStore = useConfigStore()
    const sysName = configStore.get('sys_name')
    const pageTitle = to.meta?.title as string
    document.title = pageTitle ? `${pageTitle} - ${sysName}` : sysName
  }, 0)
})

export default router
