const STORAGE_KEY = 'furumaru_recently_viewed'
const MAX_ITEMS = 20

export function useRecentlyViewed() {
  function getItems(): string[] {
    if (import.meta.server) {
      return []
    }
    try {
      const raw = localStorage.getItem(STORAGE_KEY)
      if (!raw) {
        return []
      }
      const parsed = JSON.parse(raw)
      if (!Array.isArray(parsed)) {
        return []
      }
      return parsed.filter((item: unknown) => typeof item === 'string')
    }
    catch {
      return []
    }
  }

  function addItem(id: string): void {
    if (import.meta.server) {
      return
    }
    const items = getItems().filter(item => item !== id)
    items.unshift(id)
    const trimmed = items.slice(0, MAX_ITEMS)
    localStorage.setItem(STORAGE_KEY, JSON.stringify(trimmed))
  }

  function clearItems(): void {
    if (import.meta.server) {
      return
    }
    localStorage.removeItem(STORAGE_KEY)
  }

  return {
    addItem,
    getItems,
    clearItems,
  }
}
