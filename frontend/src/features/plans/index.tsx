import { useCallback, useEffect, useMemo, useState } from 'react'
import { Loader2, Plus, RefreshCw, Search } from 'lucide-react'
import { toast } from 'sonner'
import {
  createPlan,
  deletePlan,
  fetchFeatures,
  fetchPlans,
  updatePlan,
} from '@/api/lykn'
import { formatDateTime, getErrorMessage } from '@/lib/lykn'
import type { FeatureRecord, PlanPayload, PlanRecord } from '@/types/api'
import { ConfirmDialog } from '@/components/confirm-dialog'
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
import { Switch } from '@/components/ui/switch'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Textarea } from '@/components/ui/textarea'

const emptyForm: PlanPayload = {
  code: '',
  name: '',
  description: '',
  feature_ids: [],
  max_users: 0,
  max_devices: 1,
  enabled: true,
}

export function PlansPage() {
  const [loading, setLoading] = useState(true)
  const [rows, setRows] = useState<PlanRecord[]>([])
  const [features, setFeatures] = useState<FeatureRecord[]>([])
  const [keyword, setKeyword] = useState('')
  const [enabledFilter, setEnabledFilter] = useState<string>('all')

  const [drawerOpen, setDrawerOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [editing, setEditing] = useState<PlanRecord | null>(null)
  const [form, setForm] = useState<PlanPayload>(emptyForm)

  const [deleteTarget, setDeleteTarget] = useState<PlanRecord | null>(null)
  const [deleting, setDeleting] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const [planRows, featureRows] = await Promise.all([fetchPlans(), fetchFeatures()])
      setRows(planRows)
      setFeatures(featureRows)
    } catch (error) {
      toast.error(getErrorMessage(error, '加载套餐失败'))
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  const filtered = useMemo(() => {
    const q = keyword.trim().toLowerCase()
    return rows.filter((item) => {
      const matchKeyword =
        !q ||
        item.code.toLowerCase().includes(q) ||
        item.name.toLowerCase().includes(q) ||
        item.description.toLowerCase().includes(q)
      const matchEnabled =
        enabledFilter === 'all' || String(item.enabled) === enabledFilter
      return matchKeyword && matchEnabled
    })
  }, [rows, keyword, enabledFilter])

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setDrawerOpen(true)
  }

  function openEdit(record: PlanRecord) {
    setEditing(record)
    setForm({
      code: record.code,
      name: record.name,
      description: record.description,
      feature_ids: record.features.map((item) => item.id),
      max_users: record.max_users,
      max_devices: record.max_devices,
      enabled: record.enabled,
    })
    setDrawerOpen(true)
  }

  function toggleFeature(id: number, checked: boolean) {
    setForm((s) => ({
      ...s,
      feature_ids: checked
        ? [...s.feature_ids, id]
        : s.feature_ids.filter((item) => item !== id),
    }))
  }

  async function submitForm() {
    if (!form.code.trim() || !form.name.trim()) {
      toast.error('请填写套餐编码和名称')
      return
    }
    setSubmitting(true)
    try {
      const payload: PlanPayload = {
        ...form,
        code: form.code.trim(),
        name: form.name.trim(),
        description: form.description.trim(),
        max_users: Number(form.max_users) || 0,
        max_devices: Number(form.max_devices) || 0,
      }
      if (editing) {
        await updatePlan(editing.id, payload)
        toast.success('套餐已更新')
      } else {
        await createPlan(payload)
        toast.success('套餐已创建')
      }
      setDrawerOpen(false)
      await load()
    } catch (error) {
      toast.error(getErrorMessage(error))
    } finally {
      setSubmitting(false)
    }
  }

  async function confirmDelete() {
    if (!deleteTarget) return
    setDeleting(true)
    try {
      await deletePlan(deleteTarget.id)
      toast.success('套餐已删除')
      setDeleteTarget(null)
      await load()
    } catch (error) {
      toast.error(getErrorMessage(error, '删除套餐失败'))
    } finally {
      setDeleting(false)
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
            <h1 className='text-2xl font-bold tracking-tight'>套餐管理</h1>
            <p className='text-muted-foreground'>组合功能与额度，供 License 签发使用</p>
          </div>
          <div className='flex items-center gap-2'>
            <Button variant='outline' size='sm' onClick={() => void load()}>
              <RefreshCw className={loading ? 'animate-spin' : ''} />
              刷新
            </Button>
            <Button size='sm' onClick={openCreate}>
              <Plus />
              新建套餐
            </Button>
          </div>
        </div>

        <Card className='mb-4'>
          <CardContent className='flex flex-wrap gap-3 pt-6'>
            <div className='relative min-w-[220px] flex-1'>
              <Search className='absolute start-2.5 top-2.5 size-4 text-muted-foreground' />
              <Input
                className='ps-8'
                placeholder='搜索编码 / 名称 / 描述'
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
              />
            </div>
            <Select value={enabledFilter} onValueChange={setEnabledFilter}>
              <SelectTrigger className='w-[140px]'>
                <SelectValue placeholder='状态' />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='all'>全部状态</SelectItem>
                <SelectItem value='true'>启用</SelectItem>
                <SelectItem value='false'>停用</SelectItem>
              </SelectContent>
            </Select>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>套餐列表</CardTitle>
            <CardDescription>共 {filtered.length} 个套餐</CardDescription>
          </CardHeader>
          <CardContent className='overflow-x-auto'>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>套餐</TableHead>
                  <TableHead>编码</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>功能</TableHead>
                  <TableHead>额度</TableHead>
                  <TableHead>更新时间</TableHead>
                  <TableHead className='text-end'>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={7} className='h-24 text-center'>
                      <Loader2 className='mx-auto animate-spin' />
                    </TableCell>
                  </TableRow>
                ) : filtered.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className='h-24 text-center text-muted-foreground'>
                      暂无套餐
                    </TableCell>
                  </TableRow>
                ) : (
                  filtered.map((record) => (
                    <TableRow key={record.id}>
                      <TableCell className='font-medium'>{record.name}</TableCell>
                      <TableCell className='font-mono text-sm'>{record.code}</TableCell>
                      <TableCell>
                        <Badge variant={record.enabled ? 'default' : 'secondary'}>
                          {record.enabled ? '启用' : '停用'}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className='flex max-w-[260px] flex-wrap gap-1'>
                          {record.features.length ? (
                            record.features.map((item) => (
                              <Badge key={item.id} variant='outline'>
                                {item.name}
                              </Badge>
                            ))
                          ) : (
                            <span className='text-muted-foreground'>未绑定</span>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        用户 {record.max_users || '不限'} / 设备 {record.max_devices || '-'}
                      </TableCell>
                      <TableCell>{formatDateTime(record.updated_at)}</TableCell>
                      <TableCell className='text-end'>
                        <div className='flex justify-end gap-1'>
                          <Button size='sm' variant='outline' onClick={() => openEdit(record)}>
                            编辑
                          </Button>
                          <Button
                            size='sm'
                            variant='destructive'
                            onClick={() => setDeleteTarget(record)}
                          >
                            删除
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </Main>

      <Sheet open={drawerOpen} onOpenChange={setDrawerOpen}>
        <SheetContent className='w-full overflow-y-auto sm:max-w-lg'>
          <SheetHeader>
            <SheetTitle>{editing ? '编辑套餐' : '新建套餐'}</SheetTitle>
          </SheetHeader>
          <div className='grid gap-4 px-4 py-2'>
            <div className='grid gap-2'>
              <Label>套餐编码</Label>
              <Input
                value={form.code}
                onChange={(e) => setForm((s) => ({ ...s, code: e.target.value }))}
                placeholder='例如：pro'
              />
            </div>
            <div className='grid gap-2'>
              <Label>套餐名称</Label>
              <Input
                value={form.name}
                onChange={(e) => setForm((s) => ({ ...s, name: e.target.value }))}
                placeholder='请输入套餐名称'
              />
            </div>
            <div className='grid gap-2'>
              <Label>套餐描述</Label>
              <Textarea
                value={form.description}
                onChange={(e) => setForm((s) => ({ ...s, description: e.target.value }))}
                placeholder='请输入套餐描述'
              />
            </div>
            <div className='grid gap-2'>
              <Label>包含功能</Label>
              <div className='max-h-48 space-y-2 overflow-y-auto rounded-md border p-3'>
                {features.length === 0 ? (
                  <p className='text-sm text-muted-foreground'>暂无功能，请先创建功能</p>
                ) : (
                  features.map((item) => (
                    <label key={item.id} className='flex items-center gap-2 text-sm'>
                      <Checkbox
                        checked={form.feature_ids.includes(item.id)}
                        disabled={!item.enabled && !form.feature_ids.includes(item.id)}
                        onCheckedChange={(checked) =>
                          toggleFeature(item.id, checked === true)
                        }
                      />
                      <span>
                        {item.name}
                        <span className='ms-1 text-muted-foreground'>· {item.code}</span>
                        {!item.enabled && (
                          <span className='ms-1 text-xs text-muted-foreground'>(已停用)</span>
                        )}
                      </span>
                    </label>
                  ))
                )}
              </div>
            </div>
            <div className='grid grid-cols-2 gap-3'>
              <div className='grid gap-2'>
                <Label>用户额度</Label>
                <Input
                  type='number'
                  min={0}
                  value={form.max_users}
                  onChange={(e) =>
                    setForm((s) => ({ ...s, max_users: Number(e.target.value) }))
                  }
                />
                <p className='text-xs text-muted-foreground'>0 表示不限制</p>
              </div>
              <div className='grid gap-2'>
                <Label>设备额度</Label>
                <Input
                  type='number'
                  min={1}
                  value={form.max_devices}
                  onChange={(e) =>
                    setForm((s) => ({ ...s, max_devices: Number(e.target.value) }))
                  }
                />
              </div>
            </div>
            <div className='flex items-center justify-between rounded-md border p-3'>
              <div>
                <div className='text-sm font-medium'>启用状态</div>
                <div className='text-xs text-muted-foreground'>停用后不可用于签发</div>
              </div>
              <Switch
                checked={form.enabled}
                onCheckedChange={(enabled) => setForm((s) => ({ ...s, enabled }))}
              />
            </div>
          </div>
          <SheetFooter>
            <Button variant='outline' onClick={() => setDrawerOpen(false)}>
              取消
            </Button>
            <Button disabled={submitting} onClick={() => void submitForm()}>
              {submitting && <Loader2 className='animate-spin' />}
              保存
            </Button>
          </SheetFooter>
        </SheetContent>
      </Sheet>

      <ConfirmDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title='删除套餐'
        desc='确定删除该套餐？'
        confirmText={deleting ? '删除中...' : '删除'}
        destructive
        isLoading={deleting}
        handleConfirm={() => void confirmDelete()}
      />
    </>
  )
}
