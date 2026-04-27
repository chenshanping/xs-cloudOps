export interface DictTypeLike {
  id: number
  name: string
  type: string
}

const normalizeKeyword = (value: string) => value.trim().toLowerCase()

export function filterDictTypes<T extends DictTypeLike>(dictTypes: T[], keyword: string) {
  const normalizedKeyword = normalizeKeyword(keyword)

  if (!normalizedKeyword) {
    return dictTypes
  }

  return dictTypes.filter(item => {
    const name = item.name.toLowerCase()
    const type = item.type.toLowerCase()
    return name.includes(normalizedKeyword) || type.includes(normalizedKeyword)
  })
}

export function reconcileSelectedType<T extends DictTypeLike>(dictTypes: T[], selectedType: T | null) {
  if (!selectedType) {
    return null
  }

  return dictTypes.find(item => item.id === selectedType.id) ?? selectedType
}

export function canApplyDictDataResponse(requestType: string, currentType?: string | null) {
  return Boolean(requestType) && requestType === currentType
}
