import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { LANGUAGE_OPTIONS } from '@/feature/auth/consts'
import i18n from '@/i18n'

function currentLanguageValue(): string {
  const lang = i18n.language ?? ''
  if (lang.startsWith('ja') || lang === 'jp') return 'jp'
  return 'kr'
}

export function LoginLanguageSelect() {
  return (
    <Select value={currentLanguageValue()} onValueChange={(value) => i18n.changeLanguage(value)}>
      <SelectTrigger className="w-[140px]">
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        {LANGUAGE_OPTIONS.map((opt) => (
          <SelectItem key={opt.value} value={opt.value}>
            {opt.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
