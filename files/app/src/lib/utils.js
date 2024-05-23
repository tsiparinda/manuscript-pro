export function isElementInViewport(el) {
  const rect = el.getBoundingClientRect()
  const windowHeight =
    window.innerHeight || document.documentElement.clientHeight
  const windowWidth = window.innerWidth || document.documentElement.clientWidth

  return (
    (rect.top >= 0 && rect.top <= windowHeight) || // Top edge is in viewport
    (rect.bottom <= windowHeight && rect.bottom >= 0) || // Bottom edge is in viewport
    (rect.top < 0 &&
      rect.bottom > windowHeight && // Middle section is in viewport
      rect.left >= 0 &&
      rect.right <= windowWidth) // Whole width is in viewport
  )
}
