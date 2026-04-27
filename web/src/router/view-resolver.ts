const ADMIN_VIEW_PREFIXES = ['system/', 'monitor/', 'ai/'] as const
const ADMIN_ROOT_VIEWS = new Set(['dashboard', 'profile', 'ai'])
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

  if (normalized.startsWith('front/')) {
    return `../views/${appendVueSuffix(normalized)}`
  }

  if (AUTH_ROOT_VIEWS.has(normalized)) {
    return resolveAuthViewPath(normalized)
  }

  if (normalized.endsWith('/index')) {
    const rootView = normalized.slice(0, -'/index'.length)
    if (AUTH_ROOT_VIEWS.has(rootView)) {
      return resolveAuthViewPath(rootView)
    }
    if (ADMIN_ROOT_VIEWS.has(rootView)) {
      return resolveAdminViewPath(rootView)
    }
  }

  if (ADMIN_ROOT_VIEWS.has(normalized)) {
    return resolveAdminViewPath(normalized)
  }

  if (ADMIN_VIEW_PREFIXES.some((prefix) => normalized.startsWith(prefix))) {
    return resolveAdminViewPath(normalized)
  }

  return `../views/${appendVueSuffix(normalized)}`
}
