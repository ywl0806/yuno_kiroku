import { DropdownYear } from '@/components/blocks/DropdownYear'
import { PhotoGrid } from '@/components/blocks/PhotoGrid'
import { usePhotosRange } from '@/hooks/usePhotosRange'
import { Tab, Tabs } from '@mui/material'
import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

export const HomePage = () => {
  const { date } = useParams<{ date: string }>()

  const nav = useNavigate()
  const { range, isFetched } = usePhotosRange()
  const [year, setYear] = useState<number>(parseInt(date?.split('-')[0] ?? '') ?? new Date().getFullYear())
  const [currentRange, setCurrentRange] = useState(date ?? `${new Date().getFullYear()}-${new Date().getMonth() + 1}`)

  const years = useMemo(() => {
    const set = new Set<number>()
    range.forEach((r) => {
      set.add(r.year)
    })
    return Array.from(set)
  }, [range])

  useEffect(() => {
    if (!date) return
    const [y] = date.split('-')
    setYear(parseInt(y))
  }, [date])

  return (
    <div>
      <div>
        <DropdownYear
          years={years}
          selectedYear={year}
          onChange={(changeYear) => {
            setYear(changeYear)
            const newD = range.find((y) => y.year === changeYear)
            if (newD) {
              const newDStr = `${newD.year}-${newD.month}`
              nav(`/${newDStr}`)
              setCurrentRange(newDStr)
            }
          }}
        />
      </div>

      <div>
        {isFetched && (
          <Tabs
            variant="scrollable"
            value={currentRange}
            onChange={(_, v) => setCurrentRange(v)}
            scrollButtons
            allowScrollButtonsMobile
          >
            {range.length > 0 ? (
              range.map((ran) => {
                return (
                  <Tab
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
              <Tab label={currentRange.split('-')[1]} value={currentRange} />
            )}
          </Tabs>
        )}
      </div>

      <div>
        <PhotoGrid />
      </div>
    </div>
  )
}
