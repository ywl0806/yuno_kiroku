export const calculateImageSize = (width: number, height: number, maxWidth: number, maxHeigh: number) => {
  let newWidth = width
  let newHeight = height
  if (width > maxWidth) {
    newWidth = maxWidth
    newHeight = Math.round((maxWidth / width) * height)
  }
  if (newHeight > maxHeigh) {
    newHeight = maxHeigh
    newWidth = Math.round((maxHeigh / height) * width)
  }
  return { width: newWidth, height: newHeight }
}
export const resizeImageView = (
  width: number,
  height: number,
  maxWidth: number,
  maxHeigh: number,
  orientation: number = 1,
): {
  width: number
  height: number
} => {
  if (orientation === 6 || orientation === 8) {
    return calculateImageSize(height, width, maxWidth, maxHeigh)
  }
  return calculateImageSize(width, height, maxWidth, maxHeigh)
}
