import { FC, ReactNode, useEffect, useRef, useState } from 'react'

type Props = {
  children: ReactNode
  threshold?: number
  transitionDuration?: string
}

// 가장 가까운 스크롤 가능한 부모 요소 찾기
// const findScrollableParent = (element: HTMLElement | null): HTMLElement | null => {
//   if (!element) return null

//   let parent = element.parentElement
//   while (parent) {
//     const style = window.getComputedStyle(parent)
//     const overflowY = style.overflowY || style.overflow
//     const hasScrollableContent = parent.scrollHeight > parent.clientHeight

//     if ((overflowY === 'auto' || overflowY === 'scroll') && hasScrollableContent) {
//       return parent
//     }

//     parent = parent.parentElement
//   }

//   return null
// }

export const ScrollHideWrapper: FC<Props> = ({ children, threshold = 50, transitionDuration = '0.1s' }) => {
  const lastScrollY = useRef(0)
  const [isVisible, setIsVisible] = useState(true)
  const wrapperRef = useRef<HTMLDivElement>(null)
  const scrollContainerRef = useRef<HTMLElement | null>(null)

  useEffect(() => {
    if (!wrapperRef.current) return

    // 클래스로 스크롤 컨테이너 찾기 (모든 요소)
    const scrollContainers = Array.from(document.getElementsByClassName('scroll-container')) as HTMLElement[]
    scrollContainerRef.current = scrollContainers[0] || null

    const handleScroll = (target: HTMLElement | Window) => {
      const currentScrollY = target === window ? window.scrollY : (target as HTMLElement).scrollTop

      // 스크롤이 아래로 내려가면 숨기기, 위로 올라가면 보이기
      if (currentScrollY > lastScrollY.current && currentScrollY > threshold) {
        // 아래로 스크롤하고 threshold 이상 스크롤했을 때
        setIsVisible(false)
      } else if (currentScrollY < lastScrollY.current) {
        // 위로 스크롤할 때
        setIsVisible(true)
      }

      lastScrollY.current = currentScrollY
    }

    // 모든 스크롤 컨테이너에 이벤트 리스너 추가
    const cleanupFunctions: (() => void)[] = []

    scrollContainers.forEach((container) => {
      const handler = () => handleScroll(container)
      container.addEventListener('scroll', handler, { passive: true })
      cleanupFunctions.push(() => container.removeEventListener('scroll', handler))
    })

    // window에도 이벤트 리스너 추가 (fallback)
    const windowHandler = () => handleScroll(window)
    window.addEventListener('scroll', windowHandler, { passive: true })
    cleanupFunctions.push(() => window.removeEventListener('scroll', windowHandler))

    // cleanup: 모든 이벤트 리스너 제거
    return () => {
      cleanupFunctions.forEach((cleanup) => cleanup())
    }
  }, [threshold])

  return (
    <div
      ref={wrapperRef}
      className="z-98 bg-background"
      style={{
        transform: isVisible ? 'translateY(0)' : 'translateY(-100%)',
        transition: `transform ${transitionDuration} ease-in-out`,
      }}
    >
      {children}
    </div>
  )
}
