import { AdminRole } from '../api'

export interface SettingMenu {
  text: string;
  color?: string;
  action: () => void;
  roles: AdminRole[];
}
