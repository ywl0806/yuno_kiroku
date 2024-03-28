import { FormControl, MenuItem, Select } from '@mui/material'
import { FC } from 'react'

type Props = {
  years: number[]
  selectedYear: number
  onChange: (year: number) => void
}
export const DropdownYear: FC<Props> = ({ years, selectedYear, onChange }) => {
  return (
    <div className="flex justify-center">
      <FormControl variant="standard" sx={{ m: 1, minWidth: '6rem' }}>
        <Select
          value={selectedYear}
          sx={{ fontSize: '1rem' }}
          onChange={(e) => {
            onChange(e.target.value as number)
          }}
        >
          {years.map((year) => (
            <MenuItem key={year} value={year}>
              {year}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </div>
  )
}
