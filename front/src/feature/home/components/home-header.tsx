import colors from '@/colors'
import { DropdownYear } from '@/components/blocks/dropdown-year'
// import { ScrollHideWrapper } from '@/components/layouts/scroll-hide-wrapper'
import { PhotoRange } from '@/service/get-photos-range'
import { Tab, Tabs } from '@mui/material'
import { FC, useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'

type SelectedDate = {
  year: number
  month: number
}

type Props = {
  date: string
  range: PhotoRange[]
  rangeFetched: boolean
}
export const HomeHeader: FC<Props> = ({ date, range, rangeFetched }) => {
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
    <div>
      {/* <ScrollHideWrapper> */}
      <div className="flex justify-center pt-2">
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
      {/* </ScrollHideWrapper> */}

      <div className="z-99 bg-background h-full">
        {rangeFetched && (
          <Tabs
            sx={{
              '&.MuiTabs-root': {
                minHeight: '0rem',
              },
            }}
            variant="scrollable"
            value={currentYearMonthStr}
            onChange={(_, v) => {
              nav(`/${v}`)
            }}
            scrollButtons
            TabIndicatorProps={{
              style: {
                backgroundColor: colors.amethyst,
              },
            }}
          >
            {range.length > 0 ? (
              range.map((ran) => {
                return (
                  <Tab
                    sx={{
                      '&.MuiButtonBase-root': {
                        fontSize: '1rem',
                        minWidth: '0rem',
                        width: '20%',

                        minHeight: '0',
                        height: '2.5rem',
                        paddingBottom: '0',
                        paddingTop: '0',
                        '&.Mui-selected': {
                          color: colors.amethyst,
                        },
                      },
                    }}
                    key={`${ran.year}-${ran.month}`}
                    label={ran.month}
                    value={`${ran.year}-${ran.month}`}
                    onClick={() => {
                      nav(`/${ran.year}-${ran.month}`)
                    }}
                  ></Tab>
                )
              })
            ) : (
              <Tab label={currentYearMonthStr.split('-')[1]} value={currentYearMonthStr} />
            )}
          </Tabs>
        )}
      </div>
    </div>
  )
}
