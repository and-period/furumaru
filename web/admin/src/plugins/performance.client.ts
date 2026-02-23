import { usePerformanceMetrics } from '~/composables/usePerformanceMetrics'

export default defineNuxtPlugin(() => {
  if (import.meta.dev) {
    const { measureLCP, measureCLS } = usePerformanceMetrics()

    measureLCP((value) => {
      console.log(`[Performance] LCP: ${value.toFixed(0)}ms`)
    })

    measureCLS((value) => {
      console.log(`[Performance] CLS: ${value.toFixed(4)}`)
    })
  }
})
