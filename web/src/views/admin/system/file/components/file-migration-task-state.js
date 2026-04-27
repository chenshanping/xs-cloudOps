export function shouldRestoreMigrationTaskOnOpen(task) {
  if (!task) {
    return false
  }

  return task.status === 'SCANNING' || task.status === 'RUNNING'
}
