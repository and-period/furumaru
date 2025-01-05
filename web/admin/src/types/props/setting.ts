import type { AdminType } from '../api'

export interface SettingMenu {
  text: string
  color?: string
  action: () => void
  adminTypes: AdminType[]
}
