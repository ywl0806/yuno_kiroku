import { cn } from '@/lib/utils'
import { FC } from 'react'

type Props = React.HTMLAttributes<HTMLDivElement>
export const HeaderContainer: FC<Props> = ({ className, ...props }) => {
  return <div className={cn(`py-[0.5rem] shadow-2xl shadow-darkGray`, className)} {...props} />
}
