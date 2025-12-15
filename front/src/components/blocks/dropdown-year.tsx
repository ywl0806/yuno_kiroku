import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { FC } from 'react'

type Props = {
  years: number[]
  selectedYear: number
  onChange: (year: number) => void
}
export const DropdownYear: FC<Props> = ({ years, selectedYear, onChange }) => {
  return (
    <div className="">
      <Select
        value={selectedYear.toString()}
        onValueChange={(value) => {
          onChange(parseInt(value))
        }}
      >
        <SelectTrigger className=" shadow-none ring-0 focus:ring-0">
          <SelectValue placeholder="YYYY" />
        </SelectTrigger>
        <SelectContent className="z-999">
          {years.map((year) => (
            <SelectItem key={year} value={year.toString()}>
              {year}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  )
}
