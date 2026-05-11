import { AlbumGroupPermission, Group } from '@/types'
import { useTranslation } from 'react-i18next'

interface Props {
  groups: Group[]
  value: AlbumGroupPermission[]
  onChange: (permissions: AlbumGroupPermission[]) => void
}

export const AlbumGroupPermissionSelector = ({ groups, value, onChange }: Props) => {
  const { t } = useTranslation()

  const hasPermission = (groupId: number, permission: 'R' | 'W') =>
    value.some((p) => p.group_id === groupId && p.permission === permission)

  const togglePermission = (groupId: number, permission: 'R' | 'W') => {
    if (hasPermission(groupId, permission)) {
      onChange(value.filter((p) => !(p.group_id === groupId && p.permission === permission)))
    } else {
      onChange([...value, { group_id: groupId, permission }])
    }
  }

  return (
    <div className="space-y-2">
      {groups.map((group) => (
        <div key={group.id} className="flex items-center justify-between rounded-md border px-3 py-2">
          <span className="text-sm">{group.is_admin ? t('settings.group.adminName') : group.name}</span>
          <div className="flex items-center gap-4">
            <label className="flex cursor-pointer items-center gap-1.5 text-sm">
              <input
                type="checkbox"
                checked={hasPermission(group.id, 'R')}
                onChange={() => togglePermission(group.id, 'R')}
                className="accent-primary"
              />
              {t('settings.album.permRead')}
            </label>
            <label className="flex cursor-pointer items-center gap-1.5 text-sm">
              <input
                type="checkbox"
                checked={hasPermission(group.id, 'W')}
                onChange={() => togglePermission(group.id, 'W')}
                className="accent-primary"
              />
              {t('settings.album.permWrite')}
            </label>
          </div>
        </div>
      ))}
    </div>
  )
}
