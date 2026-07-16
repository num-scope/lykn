import {
  FolderKanban,
  KeyRound,
  LayoutDashboard,
  Package,
  Puzzle,
} from 'lucide-react'
import { type SidebarData } from '../types'

export const sidebarData: SidebarData = {
  user: {
    name: '管理员',
    email: 'admin',
    avatar: '/avatars/shadcn.jpg',
  },
  navGroups: [
    {
      title: '业务',
      items: [
        {
          title: '仪表盘',
          url: '/',
          icon: LayoutDashboard,
        },
        {
          title: '项目',
          url: '/projects',
          icon: FolderKanban,
        },
        {
          title: '功能',
          url: '/features',
          icon: Puzzle,
        },
        {
          title: '套餐',
          url: '/plans',
          icon: Package,
        },
        {
          title: 'License',
          url: '/licenses',
          icon: KeyRound,
        },
      ],
    },
  ],
}
