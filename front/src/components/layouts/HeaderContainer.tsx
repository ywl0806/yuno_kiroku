import { FC } from 'react'

type Props = React.HTMLAttributes<HTMLDivElement>
export const HeaderContainer: FC<Props> = ({ className, ...props }) => {
  return <div className={`shadow-darkGray py-[0.5rem] shadow-[0px_0px_5px_5px] ${className}`} {...props} />
}
