import { Button } from '@/components/ui/button'
import { FormField } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { RadioList } from '@/components/ui/radio-list'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { SettingsSubPageLayout } from '@/feature/settings/components/settings-sub-page-layout'
import { useUpdateMember } from '@/feature/settings/hooks/use-update-member'
import { FAMILY_TITLE_OPTIONS, Group, Member } from '@/types'
import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import { z } from 'zod'

interface Props {
    memberId: number
    member: Member
    groups: Group[]
}

const formSchema = z.object({
    familyTitle: z.enum(FAMILY_TITLE_OPTIONS).optional(),
    customFamilyTitle: z.string().optional().nullable(),
})

type FormValues = z.infer<typeof formSchema>

export const MemberEditContainer = ({ memberId, member, groups }: Props) => {
    const { t } = useTranslation()
    const navigate = useNavigate()

    const { mutate: updateMember, isPending: isUpdating } = useUpdateMember()
    const [selectedGroupId, setSelectedGroupId] = useState<number>(member?.group_id ?? 0)
    const form = useForm<FormValues>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            familyTitle: member?.family_title ?? undefined,
            customFamilyTitle: member?.custom_family_title ?? null,
        },
    })

    const familyTitle = form.watch('familyTitle')

    const handleSubmit = (data: FormValues) => {
        updateMember(
            { id: memberId, groupId: selectedGroupId, familyTitle: data.familyTitle ?? '', customFamilyTitle: data.customFamilyTitle ?? '' },
            { onSuccess: () => navigate(-1) },
        )
    }

    return (
        <SettingsSubPageLayout title={t('settings.member.editTitle')}>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
                <div className="px-4 pt-3 space-y-3">
                    <div className="rounded-2xl bg-white shadow-sm ring-1 ring-black/[0.04] divide-y divide-stone-100">
                        <div className="px-4 py-4">
                            <p className="text-sm font-semibold text-stone-800">{member?.name || member?.username}</p>
                        </div>

                        <div className="px-4 py-4 space-y-1.5">
                            <FormField
                                control={form.control}
                                name="familyTitle"
                                render={({ field }) => (
                                    <>
                                        <Label>{t('settings.member.familyTitleLabel')}</Label>
                                        <Select value={field.value} onValueChange={field.onChange} disabled={isUpdating}>
                                            <SelectTrigger className="border-stone-200">
                                                <SelectValue placeholder={t('settings.member.inviteTitlePlaceholder')} />
                                            </SelectTrigger>
                                            <SelectContent>
                                                {FAMILY_TITLE_OPTIONS.map((option) => (
                                                    <SelectItem key={option} value={option}>
                                                        {t(`settings.member.familyTitle.${option}`)}
                                                    </SelectItem>
                                                ))}
                                            </SelectContent>
                                        </Select>
                                    </>
                                )}
                            />
                        </div>

                        {familyTitle === 'custom' && (
                            <div className="px-4 py-4 space-y-1.5">
                                <Label>{t('settings.member.customTitleLabel')}</Label>
                                <FormField
                                    control={form.control}
                                    name="customFamilyTitle"
                                    render={({ field }) => (
                                        <Input
                                            value={field.value ?? ''}
                                            onChange={field.onChange}
                                            placeholder={t('settings.member.customTitlePlaceholder')}
                                            disabled={isUpdating}
                                            className="border-stone-200 focus-visible:ring-stone-400"
                                        />
                                    )}
                                />
                            </div>
                        )}

                        <div className="px-4 py-4 space-y-1.5">
                            <Label>{t('settings.member.groupLabel')}</Label>
                            <RadioList
                                options={groups?.map((group) => ({ value: group.id, label: group.name })) ?? []}
                                value={selectedGroupId}
                                onChange={(val) => setSelectedGroupId(val as number)}
                                disabled={isUpdating}
                            />
                        </div>
                    </div>

                    <Button
                        type="submit"
                        className="w-full h-11 rounded-xl bg-stone-900 hover:bg-stone-800 font-medium mt-2"
                        disabled={isUpdating}
                    >
                        {t('common.save')}
                    </Button>
                </div>
            </form>
        </SettingsSubPageLayout>
    )
}
