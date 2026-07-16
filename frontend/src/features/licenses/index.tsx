import { useCallback, useEffect, useMemo, useState } from 'react'
import { Loader2, Plus, RefreshCw, Search } from 'lucide-react'
import { toast } from 'sonner'
import {
  downloadLicense,
  fetchPlans,
  fetchProjectLicenses,
  issueLicense,
} from '@/api/lykn'
import {
  formatDateTime,
  formatPeriod,
  getErrorMessage,
  getLicenseState,
  licenseStateText,
  saveBlobFile,
  toDatetimeLocalValue,
} from '@/lib/lykn'
import { useAuthStore } from '@/stores/auth-store'
import type {
  LicenseHardwarePayload,
  LicenseRecord,
  LicenseState,
  PlanRecord,
} from '@/types/api'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { ThemeSwitch } from '@/components/theme-switch'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Sheet,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

type HardwareKey = keyof LicenseHardwarePayload

type IssueForm = {
  subjectName: string
  subjectEmail: string
  subjectOrg: string
  planId?: number
  notBefore: string
  notAfter: string
  hardwareKeys: HardwareKey[]
  hostname: string
  cpu_id: string
  disk_serial: string
  mac_addresses: string
}

const hardwareOptions: { label: string; value: HardwareKey }[] = [
  { label: '主机名', value: 'hostname' },
  { label: 'CPU ID', value: 'cpu_id' },
  { label: '磁盘序列号', value: 'disk_serial' },
  { label: 'MAC 地址', value: 'mac_addresses' },
]

function defaultIssueForm(planId?: number): IssueForm {
  const now = new Date()
  const end = new Date(now)
  end.setFullYear(end.getFullYear() + 1)
  return {
    subjectName: '',
    subjectEmail: '',
    subjectOrg: '',
    planId,
    notBefore: toDatetimeLocalValue(now),
    notAfter: toDatetimeLocalValue(end),
    hardwareKeys: [],
    hostname: '',
    cpu_id: '',
    disk_serial: '',
    mac_addresses: '',
  }
}

function stateVariant(state: LicenseState): 'default' | 'secondary' | 'destructive' | 'outline' {
  if (state === 'active') return 'default'
  if (state === 'upcoming') return 'secondary'
  if (state === 'expired') return 'destructive'
  return 'outline'
}

export function LicensesPage() {
  const projects = useAuthStore((s) => s.auth.projects)
  const selectedProjectId = useAuthStore((s) => s.auth.selectedProjectId)
  const loadingProjects = useAuthStore((s) => s.auth.loadingProjects)
  const selectProject = useAuthStore((s) => s.auth.selectProject)
  const loadProjects = useAuthStore((s) => s.auth.loadProjects)

  const [loading, setLoading] = useState(false)
  const [rows, setRows] = useState<LicenseRecord[]>([])
  const [plans, setPlans] = useState<PlanRecord[]>([])
  const [keyword, setKeyword] = useState('')
  const [stateFilter, setStateFilter] = useState<string>('all')

  const [drawerOpen, setDrawerOpen] = useState(false)
  const [issuing, setIssuing] = useState(false)
  const [form, setForm] = useState<IssueForm>(defaultIssueForm())

  const [detailOpen, setDetailOpen] = useState(false)
  const [detailRecord, setDetailRecord] = useState<LicenseRecord | null>(null)
  const [downloadingId, setDownloadingId] = useState<number>()

  const enabledPlans = useMemo(() => plans.filter((item) => item.enabled), [plans])
  const selectedPlan = useMemo(
    () => plans.find((item) => item.id === form.planId),
    [plans, form.planId]
  )

  const loadLicenses = useCallback(async () => {
    if (!selectedProjectId) {
      setRows([])
      return
    }
    setLoading(true)
    try {
      setRows(await fetchProjectLicenses(selectedProjectId))
    } catch (error) {
      toast.error(getErrorMessage(error, '加载 License 失败'))
    } finally {
      setLoading(false)
    }
  }, [selectedProjectId])

  const loadPlans = useCallback(async () => {
    try {
      setPlans(await fetchPlans())
    } catch (error) {
      toast.error(getErrorMessage(error, '加载套餐失败'))
    }
  }, [])

  useEffect(() => {
    void loadProjects()
    void loadPlans()
  }, [loadProjects, loadPlans])

  useEffect(() => {
    void loadLicenses()
  }, [loadLicenses])

  const filtered = useMemo(() => {
    const q = keyword.trim().toLowerCase()
    return rows.filter((item) => {
      const matchKeyword =
        !q ||
        item.uuid.toLowerCase().includes(q) ||
        item.subject_name.toLowerCase().includes(q) ||
        item.subject_email.toLowerCase().includes(q) ||
        item.subject_org.toLowerCase().includes(q) ||
        item.plan.toLowerCase().includes(q) ||
        item.plan_name.toLowerCase().includes(q)
      const state = getLicenseState(item)
      const matchState = stateFilter === 'all' || state === stateFilter
      return matchKeyword && matchState
    })
  }, [rows, keyword, stateFilter])

  function openIssueDrawer() {
    if (!selectedProjectId) {
      toast.warning('请先选择项目')
      return
    }
    if (!enabledPlans.length) {
      toast.warning('请先创建并启用套餐')
      return
    }
    setForm(defaultIssueForm(enabledPlans[0]?.id))
    setDrawerOpen(true)
  }

  function hasHardwareKey(key: HardwareKey) {
    return form.hardwareKeys.includes(key)
  }

  function toggleHardwareKey(key: HardwareKey, checked: boolean) {
    setForm((s) => ({
      ...s,
      hardwareKeys: checked
        ? [...s.hardwareKeys, key]
        : s.hardwareKeys.filter((item) => item !== key),
    }))
  }

  function getHardwarePayload(): LicenseHardwarePayload {
    const payload: LicenseHardwarePayload = {}
    if (hasHardwareKey('hostname')) payload.hostname = form.hostname.trim()
    if (hasHardwareKey('cpu_id')) payload.cpu_id = form.cpu_id.trim()
    if (hasHardwareKey('disk_serial')) payload.disk_serial = form.disk_serial.trim()
    if (hasHardwareKey('mac_addresses')) {
      payload.mac_addresses = form.mac_addresses
        .split(',')
        .map((item) => item.trim())
        .filter(Boolean)
    }
    return payload
  }

  async function submitIssue() {
    if (!selectedProjectId || !form.planId) return
    if (!form.subjectName.trim()) {
      toast.error('请输入授权对象')
      return
    }
    if (!form.notBefore || !form.notAfter) {
      toast.error('请选择有效期')
      return
    }

    setIssuing(true)
    try {
      await issueLicense(selectedProjectId, {
        subject: {
          name: form.subjectName.trim(),
          email: form.subjectEmail.trim(),
          organization: form.subjectOrg.trim(),
        },
        plan_id: form.planId,
        not_before: new Date(form.notBefore).toISOString(),
        not_after: new Date(form.notAfter).toISOString(),
        hardware: getHardwarePayload(),
      })
      toast.success('License 已签发')
      setDrawerOpen(false)
      await loadLicenses()
    } catch (error) {
      toast.error(getErrorMessage(error, '签发 License 失败'))
    } finally {
      setIssuing(false)
    }
  }

  async function handleDownload(record: LicenseRecord) {
    setDownloadingId(record.id)
    try {
      saveBlobFile(await downloadLicense(record.id))
    } catch (error) {
      toast.error(getErrorMessage(error, '下载 License 失败'))
    } finally {
      setDownloadingId(undefined)
    }
  }

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
        <div className='mb-4 flex flex-wrap items-center justify-between gap-3'>
          <div>
            <h1 className='text-2xl font-bold tracking-tight'>License 管理</h1>
            <p className='text-muted-foreground'>按项目签发、查看与下载 License</p>
          </div>
          <div className='flex items-center gap-2'>
            <Button variant='outline' size='sm' onClick={() => void loadLicenses()}>
              <RefreshCw className={loading ? 'animate-spin' : ''} />
              刷新
            </Button>
            <Button size='sm' onClick={openIssueDrawer}>
              <Plus />
              签发 License
            </Button>
          </div>
        </div>

        <Card className='mb-4'>
          <CardContent className='flex flex-wrap gap-3 pt-6'>
            <Select
              value={selectedProjectId ? String(selectedProjectId) : undefined}
              onValueChange={(value) => selectProject(Number(value))}
              disabled={loadingProjects || projects.length === 0}
            >
              <SelectTrigger className='w-[220px]'>
                <SelectValue placeholder='选择项目' />
              </SelectTrigger>
              <SelectContent>
                {projects.map((project) => (
                  <SelectItem key={project.id} value={String(project.id)}>
                    {project.name} · {project.key_bits} bit
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <div className='relative min-w-[220px] flex-1'>
              <Search className='absolute start-2.5 top-2.5 size-4 text-muted-foreground' />
              <Input
                className='ps-8'
                placeholder='搜索客户 / 邮箱 / UUID / 套餐'
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
              />
            </div>
            <Select value={stateFilter} onValueChange={setStateFilter}>
              <SelectTrigger className='w-[140px]'>
                <SelectValue placeholder='有效性' />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='all'>全部状态</SelectItem>
                <SelectItem value='active'>生效中</SelectItem>
                <SelectItem value='upcoming'>未开始</SelectItem>
                <SelectItem value='expired'>已过期</SelectItem>
              </SelectContent>
            </Select>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>License 列表</CardTitle>
            <CardDescription>
              {selectedProjectId
                ? `共 ${filtered.length} 条`
                : '请先选择项目'}
            </CardDescription>
          </CardHeader>
          <CardContent className='overflow-x-auto'>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>授权对象</TableHead>
                  <TableHead>UUID</TableHead>
                  <TableHead>套餐</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>有效期</TableHead>
                  <TableHead>功能</TableHead>
                  <TableHead>签发时间</TableHead>
                  <TableHead className='text-end'>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {!selectedProjectId ? (
                  <TableRow>
                    <TableCell colSpan={8} className='h-24 text-center text-muted-foreground'>
                      请先选择项目
                    </TableCell>
                  </TableRow>
                ) : loading ? (
                  <TableRow>
                    <TableCell colSpan={8} className='h-24 text-center'>
                      <Loader2 className='mx-auto animate-spin' />
                    </TableCell>
                  </TableRow>
                ) : filtered.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} className='h-24 text-center text-muted-foreground'>
                      暂无 License
                    </TableCell>
                  </TableRow>
                ) : (
                  filtered.map((record) => {
                    const state = getLicenseState(record)
                    return (
                      <TableRow key={record.id}>
                        <TableCell>
                          <div className='font-medium'>{record.subject_name}</div>
                          <div className='text-xs text-muted-foreground'>
                            {record.subject_email || '-'}
                          </div>
                        </TableCell>
                        <TableCell className='max-w-[180px] truncate font-mono text-xs'>
                          {record.uuid}
                        </TableCell>
                        <TableCell>{record.plan_name || record.plan || '-'}</TableCell>
                        <TableCell>
                          <Badge variant={stateVariant(state)}>{licenseStateText(state)}</Badge>
                        </TableCell>
                        <TableCell className='text-xs'>
                          {formatPeriod(record.not_before, record.not_after)}
                        </TableCell>
                        <TableCell>
                          <div className='flex max-w-[200px] flex-wrap gap-1'>
                            {record.features.length ? (
                              record.features.map((item) => (
                                <Badge key={item} variant='outline'>
                                  {item}
                                </Badge>
                              ))
                            ) : (
                              <span className='text-muted-foreground'>未配置</span>
                            )}
                          </div>
                        </TableCell>
                        <TableCell>{formatDateTime(record.created_at)}</TableCell>
                        <TableCell className='text-end'>
                          <div className='flex justify-end gap-1'>
                            <Button
                              size='sm'
                              variant='outline'
                              onClick={() => {
                                setDetailRecord(record)
                                setDetailOpen(true)
                              }}
                            >
                              详情
                            </Button>
                            <Button
                              size='sm'
                              variant='outline'
                              disabled={downloadingId === record.id}
                              onClick={() => void handleDownload(record)}
                            >
                              下载
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    )
                  })
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </Main>

      <Sheet open={drawerOpen} onOpenChange={setDrawerOpen}>
        <SheetContent className='w-full overflow-y-auto sm:max-w-xl'>
          <SheetHeader>
            <SheetTitle>签发 License</SheetTitle>
          </SheetHeader>
          <div className='grid gap-4 px-4 py-2'>
            <div className='grid gap-2'>
              <Label>授权对象</Label>
              <Input
                value={form.subjectName}
                onChange={(e) => setForm((s) => ({ ...s, subjectName: e.target.value }))}
                placeholder='请输入客户或组织名称'
              />
            </div>
            <div className='grid gap-2'>
              <Label>邮箱</Label>
              <Input
                value={form.subjectEmail}
                onChange={(e) => setForm((s) => ({ ...s, subjectEmail: e.target.value }))}
                placeholder='请输入邮箱'
              />
            </div>
            <div className='grid gap-2'>
              <Label>组织</Label>
              <Input
                value={form.subjectOrg}
                onChange={(e) => setForm((s) => ({ ...s, subjectOrg: e.target.value }))}
                placeholder='请输入组织'
              />
            </div>
            <div className='grid gap-2'>
              <Label>套餐</Label>
              <Select
                value={form.planId ? String(form.planId) : undefined}
                onValueChange={(value) =>
                  setForm((s) => ({ ...s, planId: Number(value) }))
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder='请选择套餐' />
                </SelectTrigger>
                <SelectContent>
                  {enabledPlans.map((item) => (
                    <SelectItem key={item.id} value={String(item.id)}>
                      {item.name} · {item.code}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {selectedPlan && (
                <p className='text-xs text-muted-foreground'>
                  功能 {selectedPlan.features.length} 个，用户{' '}
                  {selectedPlan.max_users || '不限'}，设备 {selectedPlan.max_devices || '-'}
                </p>
              )}
            </div>
            <div className='grid grid-cols-2 gap-3'>
              <div className='grid gap-2'>
                <Label>生效时间</Label>
                <Input
                  type='datetime-local'
                  value={form.notBefore}
                  onChange={(e) => setForm((s) => ({ ...s, notBefore: e.target.value }))}
                />
              </div>
              <div className='grid gap-2'>
                <Label>过期时间</Label>
                <Input
                  type='datetime-local'
                  value={form.notAfter}
                  onChange={(e) => setForm((s) => ({ ...s, notAfter: e.target.value }))}
                />
              </div>
            </div>
            <div className='grid gap-2'>
              <Label>硬件绑定</Label>
              <div className='flex flex-wrap gap-3'>
                {hardwareOptions.map((item) => (
                  <label key={item.value} className='flex items-center gap-2 text-sm'>
                    <Checkbox
                      checked={hasHardwareKey(item.value)}
                      onCheckedChange={(checked) =>
                        toggleHardwareKey(item.value, checked === true)
                      }
                    />
                    {item.label}
                  </label>
                ))}
              </div>
            </div>
            {hasHardwareKey('hostname') && (
              <div className='grid gap-2'>
                <Label>主机名</Label>
                <Input
                  value={form.hostname}
                  onChange={(e) => setForm((s) => ({ ...s, hostname: e.target.value }))}
                  placeholder='例如：office-mac'
                />
              </div>
            )}
            {hasHardwareKey('cpu_id') && (
              <div className='grid gap-2'>
                <Label>CPU ID</Label>
                <Input
                  value={form.cpu_id}
                  onChange={(e) => setForm((s) => ({ ...s, cpu_id: e.target.value }))}
                  placeholder='例如：BFEBFBFF000906EA'
                />
              </div>
            )}
            {hasHardwareKey('disk_serial') && (
              <div className='grid gap-2'>
                <Label>磁盘序列号</Label>
                <Input
                  value={form.disk_serial}
                  onChange={(e) => setForm((s) => ({ ...s, disk_serial: e.target.value }))}
                  placeholder='例如：DISK-SERIAL-001'
                />
              </div>
            )}
            {hasHardwareKey('mac_addresses') && (
              <div className='grid gap-2'>
                <Label>MAC 地址</Label>
                <Input
                  value={form.mac_addresses}
                  onChange={(e) =>
                    setForm((s) => ({ ...s, mac_addresses: e.target.value }))
                  }
                  placeholder='多个地址用英文逗号分隔'
                />
              </div>
            )}
          </div>
          <SheetFooter>
            <Button variant='outline' onClick={() => setDrawerOpen(false)}>
              取消
            </Button>
            <Button disabled={issuing} onClick={() => void submitIssue()}>
              {issuing && <Loader2 className='animate-spin' />}
              签发
            </Button>
          </SheetFooter>
        </SheetContent>
      </Sheet>

      <Sheet open={detailOpen} onOpenChange={setDetailOpen}>
        <SheetContent className='w-full overflow-y-auto sm:max-w-lg'>
          <SheetHeader>
            <SheetTitle>License 详情</SheetTitle>
          </SheetHeader>
          {detailRecord ? (
            <div className='space-y-3 px-4 py-2 text-sm'>
              <DetailItem label='UUID' value={detailRecord.uuid} mono />
              <DetailItem label='授权对象' value={detailRecord.subject_name} />
              <DetailItem label='邮箱' value={detailRecord.subject_email || '-'} />
              <DetailItem label='组织' value={detailRecord.subject_org || '-'} />
              <DetailItem
                label='套餐'
                value={detailRecord.plan_name || detailRecord.plan || '-'}
              />
              <DetailItem
                label='有效期'
                value={formatPeriod(detailRecord.not_before, detailRecord.not_after)}
              />
              <DetailItem
                label='额度'
                value={`用户 ${detailRecord.limits?.max_users || '不限'} / 设备 ${detailRecord.limits?.max_devices || '-'}`}
              />
              <div>
                <div className='mb-1 text-muted-foreground'>功能</div>
                <div className='flex flex-wrap gap-1'>
                  {detailRecord.features.length ? (
                    detailRecord.features.map((item) => (
                      <Badge key={item} variant='outline'>
                        {item}
                      </Badge>
                    ))
                  ) : (
                    <span className='text-muted-foreground'>未配置</span>
                  )}
                </div>
              </div>
            </div>
          ) : (
            <p className='px-4 text-sm text-muted-foreground'>暂无详情</p>
          )}
        </SheetContent>
      </Sheet>
    </>
  )
}

function DetailItem({
  label,
  value,
  mono,
}: {
  label: string
  value: string
  mono?: boolean
}) {
  return (
    <div>
      <div className='mb-1 text-muted-foreground'>{label}</div>
      <div className={mono ? 'break-all font-mono text-xs' : ''}>{value}</div>
    </div>
  )
}
