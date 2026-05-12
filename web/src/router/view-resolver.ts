const AUTH_ROOT_VIEWS = new Set([
  'login',
  'register',
  'forgot-password',
  'reset-password',
])

function trimSlashes(value: string) {
  return value.replace(/^\/+|\/+$/g, '')
}

function appendVueSuffix(componentPath: string) {
  return componentPath.endsWith('/index')
    ? `${componentPath}.vue`
    : `${componentPath}/index.vue`
}

function resolveAuthViewPath(componentPath: string) {
  return `../views/auth/${appendVueSuffix(componentPath)}`
}

function resolveAdminViewPath(componentPath: string) {
  return `../views/admin/${appendVueSuffix(componentPath)}`
}

export function resolveViewModulePath(componentPath: string) {
  const normalized = trimSlashes(componentPath)

  if (!normalized) {
    return null
  }

  // front/ 前缀 → views/front/
  if (normalized.startsWith('front/')) {
    return `../views/${appendVueSuffix(normalized)}`
  }

  // auth 页面 → views/auth/
  if (AUTH_ROOT_VIEWS.has(normalized)) {
    return resolveAuthViewPath(normalized)
  }
  if (normalized.endsWith('/index')) {
    const rootView = normalized.slice(0, -'/index'.length)
    if (AUTH_ROOT_VIEWS.has(rootView)) {
      return resolveAuthViewPath(rootView)
    }
  }

  // 其他所有组件 → views/admin/ (后台管理页面默认路径)
  return resolveAdminViewPath(normalized)
}
