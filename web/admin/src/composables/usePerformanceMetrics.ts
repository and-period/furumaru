export function usePerformanceMetrics() {
  function measureLCP(callback: (value: number) => void): void {
    if (typeof PerformanceObserver === 'undefined') {
      return
    }
    new PerformanceObserver((list) => {
      const entries = list.getEntries()
      const last = entries[entries.length - 1]
      if (last) {
        callback(last.startTime)
      }
    }).observe({ type: 'largest-contentful-paint', buffered: true })
  }

  function measureCLS(callback: (value: number) => void): void {
    if (typeof PerformanceObserver === 'undefined') {
      return
    }
    let clsValue = 0
    new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        if (!(entry as any).hadRecentInput) {
          clsValue += (entry as any).value
        }
      }
      callback(clsValue)
    }).observe({ type: 'layout-shift', buffered: true })
  }

  return { measureLCP, measureCLS }
}
