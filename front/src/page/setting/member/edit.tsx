
import { MemberEditContainer } from '@/feature/settings/components/member-edit-container'
import { useGetGroups } from '@/feature/settings/hooks/use-get-groups'
import { useGetMember } from '@/feature/settings/hooks/use-get-member'
import { Suspense } from 'react'
import { useParams } from 'react-router-dom'

export const SettingsMemberEditPage = () => {

  const { memberId } = useParams<{ memberId: string }>()

  const { data: member } = useGetMember(Number(memberId))
  const { data: groups } = useGetGroups()

  return (
    <>
      {member && groups && (
        <MemberEditContainer memberId={Number(memberId)} member={member} groups={groups} />
      )}
    </>
  )
}
