import { DropdownYear } from '@/components/blocks/dropdown-year'
import { cn } from '@/lib/utils'
import { MediaItemRange } from '@/types'
import { RefreshCcw } from 'lucide-react'
import { FC, useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQueryClient } from '@tanstack/react-query'

type SelectedDate = {
  year: number
  month: number
}

type Props = {
  date: string
  range: MediaItemRange[]
  hidden: boolean

}
export const HomeHeader: FC<Props> = ({ date, range, hidden }) => {
  const nav = useNavigate()
  const queryClient = useQueryClient()
  const [selectedDate, setSelectedDate] = useState<SelectedDate>({
    year: parseInt(date?.split('-')[0] ?? '0') ?? new Date().getFullYear(),
    month: parseInt(date?.split('-')[1] ?? '0') ?? new Date().getMonth() + 1,
  })

  const currentYearMonthStr = useMemo(() => {
    return `${selectedDate.year}-${selectedDate.month}`
  }, [selectedDate.month, selectedDate.year])

  const years = useMemo(() => {
    const set = new Set<number>()
    const thisYear = new Date().getFullYear()
    set.add(thisYear)
    range.forEach((r) => {
      set.add(r.year)
    })

    return Array.from(set)
  }, [range])

  const refresh = () => {
    queryClient.refetchQueries({ queryKey: ['mediaItems', selectedDate.year, selectedDate.month] })
    queryClient.refetchQueries({ queryKey: ['mediaItemsRange'] })
  }

  useEffect(() => {
    if (!date) return
    const [y, m] = date.split('-')
    setSelectedDate({ year: parseInt(y), month: parseInt(m) })
  }, [date])

  return (
    <div>
      {/* Year dropdown: CSS grid trick for smooth height animation */}
      <div className={cn('grid transition-all duration-300 ease-in-out', hidden ? 'grid-rows-[0fr]' : 'grid-rows-[1fr]')}>
        <div className="overflow-hidden">
          <div className={cn('flex items-center justify-center gap-2 pt-2 pb-1 transition-opacity duration-300', hidden ? 'opacity-0' : 'opacity-100')}>
            <div className="w-1/3" />
            <div className="w-1/3 flex justify-center">
              <div className='w-20'>
                <DropdownYear
                  years={years}
                  selectedYear={selectedDate.year}
                  onChange={(changeYear) => {
                    setSelectedDate({ ...selectedDate, year: changeYear })
                    const newD = range.find((y) => y.year === changeYear)
                    if (newD) {
                      nav(`/${newD.year}-${newD.month}`)
                    }
                  }}
                />
              </div>
            </div>
            <button type="button" onClick={refresh} className="w-1/3 flex justify-end flex-row pr-5">
              <RefreshCcw className="size-5 text-muted-foreground hover:text-foreground transition-all active:scale-95 active:rotate-180 duration-300" />
            </button>
          </div>
        </div>
      </div>

      {/* Month tabs */}
      <div className="bg-background">
        {range.length > 0 && (
          <div className="flex overflow-x-auto">
            {range.map((ran) => {
              const value = `${ran.year}-${ran.month}`
              const isSelected = value === currentYearMonthStr
              return (
                <button
                  key={value}
                  type="button"
                  onClick={() => nav(`/${value}`)}
                  className={cn(
                    'w-[20%] shrink-0 border-b-3 px-0 text-base transition-all duration-300',
                    isSelected ? 'border-primary font-medium text-primary' : 'border-transparent text-muted-foreground',
                    hidden ? 'h-8 text-sm' : 'h-10 text-base',
                  )}
                >
                  {ran.month}
                </button>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}
