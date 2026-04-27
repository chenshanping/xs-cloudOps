export function computeDynamicRouteSyncPlan(currentNames, nextNames) {
  const current = new Set(currentNames || [])
  const next = new Set(nextNames || [])

  const removeNames = []
  const replaceNames = []
  const addNames = []

  for (const name of current) {
    if (!next.has(name)) {
      removeNames.push(name)
      continue
    }
    replaceNames.push(name)
  }

  for (const name of next) {
    if (!current.has(name)) {
      addNames.push(name)
    }
  }

  return {
    removeNames,
    replaceNames,
    addNames,
  }
}
