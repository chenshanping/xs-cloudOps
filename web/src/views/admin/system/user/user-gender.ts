export interface GenderDictOptionInput {
  label: string
  value: string | number
  tag_type?: string
}

export interface GenderOption {
  label: string
  value: number
  tag_type?: string
}

export function normalizeGenderDictOptions(options: GenderDictOptionInput[]): GenderOption[] {
  return options
    .map(option => ({
      label: option.label,
      value: Number(option.value),
      tag_type: option.tag_type,
    }))
    .filter(option => Number.isFinite(option.value))
}

export function resolveGenderOption(options: GenderOption[], value?: number) {
  return options.find(option => option.value === value)
}

export function resolveGenderLabel(options: GenderOption[], value?: number) {
  return resolveGenderOption(options, value)?.label ?? '-'
}
