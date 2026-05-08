import { SCROLL_THRESHOLD } from '@/config/config'
import { useCallback, useRef, useState } from 'react'

export function useScrollHide(threshold = SCROLL_THRESHOLD) {
  const [hidden, setHidden] = useState(false)
  const prevScrollY = useRef(0)

  const handleScroll = useCallback(
    (scrollTop: number) => {
      const clamped = Math.max(0, scrollTop)
      const diff = clamped - prevScrollY.current

      if (Math.abs(diff) < threshold) return
      if (diff > 0 && clamped > threshold) {
        setHidden(true)
      } else if (diff < 0) {
        setHidden(false)
      }
      prevScrollY.current = clamped
    },
    [threshold],
  )

  const reset = useCallback(() => {
    setHidden(false)
    prevScrollY.current = 0
  }, [])

  return { hidden, handleScroll, reset }
}
