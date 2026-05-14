export function isEditablePasteTarget(target: EventTarget | null) {
  if (!(target instanceof HTMLElement)) {
    return false
  }
  return Boolean(target.closest('input, textarea, [contenteditable], .w-e-text-container'))
}

export function normalizePastedFile(file: File) {
  if (file.name && file.name !== 'image.png') {
    return file
  }

  const extension = file.type.split('/')[1] || 'png'
  return new File([file], `pasted-${Date.now()}.${extension}`, { type: file.type })
}

export function getClipboardFiles(
  event: ClipboardEvent,
  options?: {
    multiple?: boolean
    normalizeUnnamedImage?: boolean
  }
) {
  const multiple = options?.multiple ?? true
  const normalizeUnnamedImage = options?.normalizeUnnamedImage ?? true
  const items = Array.from(event.clipboardData?.items || [])
  const fallbackFiles = Array.from(event.clipboardData?.files || [])
  const seen = new Set<string>()
  const files: File[] = []

  const pushFile = (file: File | null) => {
    if (!file) {
      return
    }

    const normalized = normalizeUnnamedImage ? normalizePastedFile(file) : file
    const key = [normalized.name, normalized.size, normalized.type, normalized.lastModified].join('::')
    if (seen.has(key)) {
      return
    }

    seen.add(key)
    files.push(normalized)
  }

  items.forEach((item) => {
    if (item.kind === 'file') {
      pushFile(item.getAsFile())
    }
  })

  fallbackFiles.forEach((file) => {
    pushFile(file)
  })

  return multiple ? files : files.slice(0, 1)
}
