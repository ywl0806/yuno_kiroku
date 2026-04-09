import { DropdownYear } from '@/components/blocks/dropdown-year'
import { cn } from '@/lib/utils'
import { MediaItemRange } from '@/types'
import { FC, useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'

type SelectedDate = {
  year: number
  month: number
}

type Props = {
  date: string
  range: MediaItemRange[]
}
export const HomeHeader: FC<Props> = ({ date, range }) => {
  const nav = useNavigate()

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

  useEffect(() => {
    if (!date) return
    const [y, m] = date.split('-')
    setSelectedDate({ year: parseInt(y), month: parseInt(m) })
  }, [date])

  return (
    <div className="pt-5">
      <div className="flex items-center justify-center gap-2 pt-2">
        <DropdownYear
          years={years}
          selectedYear={selectedDate.year}
          onChange={(changeYear) => {
            setSelectedDate({ ...selectedDate, year: changeYear })
            const newD = range.find((y) => y.year === changeYear)
            if (newD) {
              const newDStr = `${newD.year}-${newD.month}`
              nav(`/${newDStr}`)
            }
          }}
        />
      </div>

      <div className="z-99 h-full bg-background">
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
                    'relative h-10 w-[20%] shrink-0 border-b-2 px-0 text-base transition-colors',
                    isSelected
                      ? 'border-primary font-medium text-primary'
                      : 'border-transparent text-muted-foreground',
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
