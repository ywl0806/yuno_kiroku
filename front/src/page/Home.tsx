import colors from '@/colors'
import { DropdownYear } from '@/components/blocks/DropdownYear'
import { PhotoGridContainer } from '@/components/blocks/PhotoGridContainer'
import { HeaderContainer } from '@/components/layouts/HeaderContainer'
import { usePhotosRange } from '@/hooks/usePhotosRange'
import { Tab, Tabs } from '@mui/material'
import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { Swiper as SwiperCore } from 'swiper'
import 'swiper/css'
import { Controller } from 'swiper/modules'
import { Swiper, SwiperClass, SwiperSlide } from 'swiper/react'

SwiperCore.use([Controller])

type SelectedDate = {
  year: number
  month: number
}

export const HomePage = () => {
  const { date } = useParams<{ date: string }>()

  const [swiper, setSwiper] = useState<SwiperClass | null>(null)

  const nav = useNavigate()
  const { range, isFetched } = usePhotosRange()

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
    setSelectedDate({
      year: parseInt(y),
      month: parseInt(m),
    })

    if (swiper) {
      const findIndex = range.findIndex((r) => r.year === parseInt(y) && r.month === parseInt(m))
      swiper.slideTo(findIndex)
    }
  }, [date, range])

  return (
    <>
      <HeaderContainer>
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

        <div>
          {isFetched && (
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
      </HeaderContainer>
      <Swiper
        spaceBetween={10}
        slidesPerView={1}
        onSlideChange={(swiper) => {
          nav(`/${range[swiper.activeIndex].year}-${range[swiper.activeIndex].month}`)
        }}
        controller={{ control: swiper }}
        onSwiper={(swiper) => setSwiper(swiper)}
        className="h-full w-full py-[0.6rem]"
      >
        {range.map((ran) => {
          return (
            <SwiperSlide key={`${ran.year}-${ran.month}`} className="min-h-[80vh] w-full">
              <PhotoGridContainer year={ran.year} month={ran.month} />
            </SwiperSlide>
          )
        })}
      </Swiper>
    </>
  )
}
