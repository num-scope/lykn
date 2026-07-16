import { useEffect, useState } from 'react'
import { FolderKanban, KeyRound, ShieldCheck, ShieldX } from 'lucide-react'
import { toast } from 'sonner'
import { fetchDashboardSummary } from '@/api/lykn'
import { getErrorMessage } from '@/lib/lykn'
import type { DashboardSummary } from '@/types/api'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

const emptySummary: DashboardSummary = {
  project_count: 0,
  license_count: 0,
  active_license_count: 0,
  expired_license_count: 0,
}

export function Dashboard() {
  const [loading, setLoading] = useState(true)
  const [summary, setSummary] = useState<DashboardSummary>(emptySummary)

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true)
      try {
        const data = await fetchDashboardSummary()
        if (!cancelled) setSummary(data)
      } catch (error) {
        if (!cancelled) {
          toast.error(getErrorMessage(error, '加载仪表盘失败'))
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => {
      cancelled = true
    }
  }, [])

  const cards = [
    {
      title: '项目数',
      value: summary.project_count,
      description: '已创建的授权项目',
      icon: FolderKanban,
    },
    {
      title: 'License 总数',
      value: summary.license_count,
      description: '历史签发总量',
      icon: KeyRound,
    },
    {
      title: '生效中',
      value: summary.active_license_count,
      description: '当前有效 License',
      icon: ShieldCheck,
    },
    {
      title: '已过期',
      value: summary.expired_license_count,
      description: '已超过有效期',
      icon: ShieldX,
    },
  ]

  return (
    <>
      <Header>
        <div className='ms-auto flex items-center space-x-2'>
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      <Main>
        <div className='mb-6 space-y-1'>
          <h1 className='text-2xl font-bold tracking-tight'>仪表盘</h1>
          <p className='text-muted-foreground'>授权业务概览</p>
        </div>

        <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
          {cards.map((card) => (
            <Card key={card.title}>
              <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
                <CardTitle className='text-sm font-medium'>{card.title}</CardTitle>
                <card.icon className='h-4 w-4 text-muted-foreground' />
              </CardHeader>
              <CardContent>
                {loading ? (
                  <Skeleton className='h-8 w-20' />
                ) : (
                  <div className='text-2xl font-bold'>{card.value}</div>
                )}
                <CardDescription className='mt-1'>{card.description}</CardDescription>
              </CardContent>
            </Card>
          ))}
        </div>
      </Main>
    </>
  )
}
