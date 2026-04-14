import { Button } from '@/components/ui/button'
import { FormField } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { RadioList } from '@/components/ui/radio-list'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useUpdateMember } from '@/feature/settings/hooks/use-update-member'
import { FAMILY_TITLE_OPTIONS, Group, Member } from '@/types'
import { zodResolver } from '@hookform/resolvers/zod'
import { ChevronLeft } from 'lucide-react'
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
        updateMember({ id: memberId, groupId: selectedGroupId, familyTitle: data.familyTitle ?? '', customFamilyTitle: data.customFamilyTitle ?? '' }, { onSuccess: () => navigate(-1) })
    }

    return (
        <div className="mx-auto h-full max-w-[50rem] pt-4">
            <div className="flex items-center gap-2 px-4 py-2">
                <button onClick={() => navigate(-1)} className="text-muted-foreground hover:text-foreground">
                    <ChevronLeft className="size-5" />
                </button>
                <span className="text-sm font-medium">{t('settings.member.editTitle')}</span>
            </div>
            <form onSubmit={form.handleSubmit(handleSubmit)}>
                <div className="space-y-5 px-4 pt-6">
                    <div className="space-y-1.5">
                        <p className="text-sm font-medium">{member?.name || member?.username}</p>
                    </div>

                    <div className="space-y-1.5">
                        <FormField
                            control={form.control}
                            name="familyTitle"
                            render={({ field }) => (
                                <div className="space-y-1.5">
                                    <Label>{t('settings.member.familyTitleLabel')}</Label>
                                    <Select value={field.value} onValueChange={field.onChange} disabled={isUpdating}>
                                        <SelectTrigger>
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
                                </div>
                            )}
                        />
                    </div>

                    {familyTitle === 'custom' && (
                        <div className="space-y-1.5">
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
                                    />
                                )}
                            />
                        </div>
                    )}

                    <div className="space-y-1.5">
                        <Label>{t('settings.member.groupLabel')}</Label>
                        <RadioList
                            options={groups?.map((group) => ({ value: group.id, label: group.name })) ?? []}
                            value={selectedGroupId}
                            onChange={(val) => setSelectedGroupId(val as number)}
                            disabled={isUpdating}
                        />
                    </div>
                </div>

                <div className="px-4 pt-6 flex justify-end mt-4">
                    <Button type="submit" disabled={isUpdating}>{t('common.save')}</Button>
                </div>
            </form>
        </div>
    )
}
